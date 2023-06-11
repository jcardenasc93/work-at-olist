package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path"

	m "github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteDB struct {
	db *sql.DB
}

func NewSQLiteDB() (*SQLiteDB, error) {
	err := godotenv.Load(path.Base("../../.env"))
	if err != nil {
		return nil, err
	}

	dbName := os.Getenv("dbName")
	if dbName == "" {
		return nil, errors.New("DB name couldn't be empty")
	}
	dbConn, err := sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(5)
	log.Println("DB connection success!!!")

	return &SQLiteDB{
		db: dbConn,
	}, nil
}

func (sq *SQLiteDB) Setup() error {
	err := sq.CreateAuthorTable()
	if err != nil {
		log.Print(err)
		return err
	}
	err = sq.CreateBookTable()
	if err != nil {
		log.Print(err)
		return err
	}
	err = sq.CreateAuthorBookTable()
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func (sq *SQLiteDB) CreateAuthorTable() error {
	log.Println("Creating author table...")

	createAuthorsTable := `
    CREATE TABLE IF NOT EXISTS author (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(64) NOT NULL
    )
    `

	stmt, err := sq.db.Prepare(createAuthorsTable)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func (sq *SQLiteDB) CreateBookTable() error {
	log.Println("Creating Authors table...")

	createBooksTable := `
    CREATE TABLE IF NOT EXISTS book (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(80) NOT NULL,
        edition INTEGER NOT NULL,
        publication_year INTEGER NOT NULL
    )
    `

	stmt, err := sq.db.Prepare(createBooksTable)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (sq *SQLiteDB) CreateAuthorBookTable() error {
	log.Println("Creating Authors-Books relationship table...")
	createAuthorsBooksTable := `
    CREATE TABLE IF NOT EXISTS author_book (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        author_id INTEGER,
        book_id INTEGER,
        FOREIGN KEY(author_id) REFERENCES author(id)
        ON DELETE NO ACTION,
        FOREIGN KEY(book_id) REFERENCES book(id)
        ON DELETE NO ACTION
    )`
	stmt, err := sq.db.Prepare(createAuthorsBooksTable)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func (sq *SQLiteDB) InsertAuthor(authorName string) error {
	insertAuthorStmt := `
    INSERT INTO author (name) VALUES (?)
    `
	stmt, err := sq.db.Prepare(insertAuthorStmt)
	if err != nil {
		log.Fatalf("failing prepare: %s", err.Error())
		return err
	}

	_, err = stmt.Exec(authorName)
	if err != nil {
		log.Fatalf("failing execution: %s", err.Error())
		return err
	}
	return nil
}

func (sq *SQLiteDB) InsertBook(ctx context.Context, bookData *models.CreateBookReq) (*models.Book, error) {
	insertBookStmt := `INSERT INTO book (name, edition, publication_year)
                       VALUES (?, ?, ?)`
	insertAuthorBookStmt := `INSERT INTO author_book (author_id, book_id)
                             VALUES (?, ?)`

	tx, err := sq.db.BeginTx(ctx, nil)
	defer tx.Rollback()
	if err != nil {
		log.Printf("Failing creating SQL Tx: %s\n", err.Error())
		return nil, err
	}

	bookStmt, err := tx.Prepare(insertBookStmt)
	if err != nil {
		log.Printf("Failing preraring new book statement: %s\n", err.Error())
		return nil, err
	}
	defer bookStmt.Close()
	result, err := bookStmt.ExecContext(ctx, bookData.Name, bookData.Edition, bookData.PubYear)
	if err != nil {
		log.Printf("Failing inserting new book. \nData provided: %v\n%s", bookData, err.Error())
		return nil, err
	}

	bookId, err := result.LastInsertId()
	authorBookStmt, err := tx.Prepare(insertAuthorBookStmt)
	if err != nil {
		log.Printf("Failing preraring new author_book statement:%s\n", err.Error())
		return nil, err
	}
	for _, author := range bookData.Authors {
		_, err = authorBookStmt.ExecContext(ctx, author, bookId)
		if err != nil {
			log.Printf("Failing inserting author_book relationship with author_id: %v, book_id: %v.\n%s", author, bookId, err.Error())
			return nil, err
		}
	}
	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Printf("Failing commit changes in db: %s\n", err.Error())
		return nil, err
	}

	book := models.NewBook(float64(bookId), bookData.Name, bookData.Edition, bookData.PubYear, bookData.Authors)
	return book, nil
}

func (sq *SQLiteDB) filterByName(baseQuery string) (query string) {
	const filter string = `AND name LIKE '%'||?||'%'`
	query = fmt.Sprintf("%s %s", baseQuery, filter)
	return
}

func (sq *SQLiteDB) sortAndLimit(baseQuery string) (query string) {
	const limitStmt string = `ORDER BY id LIMIT ?`
	query = fmt.Sprintf("%s %s", baseQuery, limitStmt)
	return
}

func (sq *SQLiteDB) applyQueryParams(baseQuery string, q allowedQParams, params url.Values) (query string, paramVals []any) {
	query = baseQuery
	for key, fun := range q.params {
		if params.Has(key) {
			query = fun(query)
			keyVal := params.Get(key)
			paramVals = append(paramVals, keyVal)
		}
	}
	return
}

func (sq *SQLiteDB) applySortAndLimit(baseQuery string, pageId int, limit int, paramVals []any) (query string, queryVals []any) {
	query = sq.sortAndLimit(baseQuery)
	queryVals = []any{pageId}
	queryVals = append(queryVals, paramVals...)
	queryVals = append(queryVals, limit)
	return
}

func (sq *SQLiteDB) aggregateAuthorsInBook(books []*models.Book, bookId float64, authorId float64) []*models.Book {
	for _, book := range books {
		if book.Id == bookId {
			book.Authors = append(book.Authors, authorId)
		}
	}
	return books
}

func (sq *SQLiteDB) FetchAuthorsForBooks(books []*models.Book) ([]*models.Book, error) {
	query := `SELECT b.id, ab.author_id  FROM book b
              JOIN author_book ab ON b.id = ab.book_id
              WHERE b.id IN`
	bookIds := []any{}
	if len(books) == 1 {
		query = fmt.Sprintf("%s %s", query, "(?)")
		bookIds = append(bookIds, books[0].Id)
	} else {
		for i, book := range books {
			bookIds = append(bookIds, book.Id)
			if i == 0 {
				query = fmt.Sprintf("%s %s", query, "(?,")
				continue
			}
			if i == len(books)-1 {
				query = fmt.Sprintf("%s %s", query, " ?)")
				continue
			}
			query = fmt.Sprintf("%s %s", query, " ?,")
		}
	}

	rows, err := sq.execQuery(query, bookIds...)
	if err != nil {
		return books, err
	}
	for rows.Next() {
		var bookId float64
		var authorId float64

		err = rows.Scan(&bookId, &authorId)
		if err != nil {
			return books, err
		}

		books = sq.aggregateAuthorsInBook(books, bookId, authorId)
	}

	return books, nil
}

func (sq *SQLiteDB) FetchBooks(pagination *m.PaginationVals, params url.Values) ([]*models.Book, error) {
	const nameKey string = "name"
	const pubYearKey string = "publication_year"
	const editionKey string = "edition"
	const authorKey string = "author"
	var books = []*models.Book{}
	var rows *sql.Rows
	var err error
	pageId := pagination.PageId
	limit := pagination.Limit

	query := `SELECT id, name, edition, publication_year FROM book
              WHERE id > ?`

	allowedParams := allowedQParams{
		params: map[string]func(string) string{
			nameKey: sq.filterByName,
		},
	}

	query, paramVals := sq.applyQueryParams(query, allowedParams, params)
	query, queryVals := sq.applySortAndLimit(query, pageId, limit, paramVals)
	rows, err = sq.execQuery(query, queryVals...)
	if err != nil {
		return books, err
	}

	defer rows.Close()

	for rows.Next() {
		var id float64
		var name string
		var edition float64
		var pubYear float64

		err = rows.Scan(&id, &name, &edition, &pubYear)
		if err != nil {
			return books, err
		}

		books = append(books, models.NewBook(id, name, edition, pubYear, []float64{}))
	}

	if len(books) > 0 {
		books, err = sq.FetchAuthorsForBooks(books)
		if err != nil {
			log.Println(err)
			return books, err
		}
	}

	return books, err
}

func (sq *SQLiteDB) FetchAuthors(pagination *m.PaginationVals, params url.Values) ([]*models.Author, error) {
	const nameKey string = "name"
	var authors = []*models.Author{}
	var rows *sql.Rows
	var err error
	pageId := pagination.PageId
	limit := pagination.Limit

	query := `SELECT id, name FROM author
              WHERE id > ?`

	allowedParams := allowedQParams{
		params: map[string]func(string) string{
			nameKey: sq.filterByName,
		},
	}

	query, paramVals := sq.applyQueryParams(query, allowedParams, params)
	query, queryVals := sq.applySortAndLimit(query, pageId, limit, paramVals)
	rows, err = sq.execQuery(query, queryVals...)
	if err != nil {
		return authors, err
	}

	defer rows.Close()

	for rows.Next() {
		var id uint64
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			return authors, err
		}

		authors = append(authors, models.NewAuthor(id, name))
	}

	return authors, nil
}

func (sq *SQLiteDB) execQuery(query string, params ...any) (*sql.Rows, error) {
	rows, err := sq.db.Query(query, params...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return rows, nil
}
