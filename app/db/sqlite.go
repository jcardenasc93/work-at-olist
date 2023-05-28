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

func (sq *SQLiteDB) FetchAuthors(pageId int) ([]*models.Author, error) {
	authors := []*models.Author{}
	return authors, nil
}
