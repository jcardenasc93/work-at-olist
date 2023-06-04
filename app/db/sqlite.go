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

func (sq *SQLiteDB) FetchAuthors(pagination *m.PaginationVals, params url.Values) ([]*models.Author, error) {
	const nameKey string = "name"
	var authors = []*models.Author{}
	var rows *sql.Rows
	var err error
	pageId := pagination.PageId
	limit := pagination.Limit

	query := `SELECT id, name FROM author
              WHERE id > ?`

	if params.Has(nameKey) {
		nameValue := string(params.Get(nameKey))
		query = sq.filterByName(query)
		query = sq.sortAndLimit(query)
		rows, err = sq.execQuery(query, pageId, nameValue, limit)
		if err != nil {
			return authors, err
		}
	} else {
		query = sq.sortAndLimit(query)
		rows, err = sq.execQuery(query, pageId, limit)
		if err != nil {
			return authors, err
		}
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
		log.Panic(err)
		return nil, err
	}
	return rows, nil
}
