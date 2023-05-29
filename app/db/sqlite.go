package db

import (
	"database/sql"
	"errors"
	"log"
	"os"
	"path"

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

func (sq *SQLiteDB) CreateAuthorsTable() error {
	log.Println("Creating Authors table...")

	createAuthorsTable := `
    CREATE TABLE IF NOT EXISTS authors (
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

func (sq *SQLiteDB) InsertAuthor(authorName string) error {
	insertAuthorStmt := `
    INSERT INTO authors (name) VALUES (?)
    `
	stmt, err := sq.db.Prepare(insertAuthorStmt)
	if err != nil {
		log.Fatalf("failing prepare: %s", err)
		return err
	}

	_, err = stmt.Exec(authorName)
	if err != nil {
		log.Fatalf("failing execution: %s", err)
		return err
	}
	return nil
}

func (sq *SQLiteDB) FetchAuthors(pageId int, params ...any) ([]*models.Author, error) {
	var authors []*models.Author
	var rows *sql.Rows
	var err error

	query := `SELECT id, name FROM authors
              WHERE id >= ?`

	// if queryParams.Has(nameKey) {
	// 	nameValue := string(queryParams.Get(nameKey))
	// 	// query := filterByName(query)
	// 	rows, err = execQuery(query, nameValue)
	// 	if err != nil {
	// 		return authors, err
	// 	}
	// } else {
	// 	rows, err = execQuery(query, pageId)
	// 	if err != nil {
	// 		return authors, err
	// 	}
	// }
	rows, err = sq.execQuery(query, pageId)
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
		log.Panic(err)
		return nil, err
	}
	return rows, nil
}
