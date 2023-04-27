package models

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var dbConn *sql.DB

func InitDB(fileName string) error {
	err := errors.New("")
	if fileName == "" {
		return errors.New("DB name couldn't be empty")
	}
	dbConn, err = sql.Open("sqlite3", fileName)
	if err != nil {
		return err
	}

	err = dbConn.Ping()
	if err != nil {
		return err
	}

	dbConn.SetMaxOpenConns(5)

	log.Println("DB connection success!!!")

	createAuthorsTable := `
    CREATE TABLE IF NOT EXISTS authors (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(64) NOT NULL
    )
    `

	stmt, err := dbConn.Prepare(createAuthorsTable)
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func execQuery(query string) (*sql.Rows, error) {
	rows, err := dbConn.Query(query)
	if err != nil {
		log.Panic(err)
		return nil, err
	}
	return rows, nil
}

func InsertAuthor(authorName string) error {
	insertAuthorStmt := `
    INSERT INTO authors (name) VALUES (?)
    `
	stmt, err := dbConn.Prepare(insertAuthorStmt)
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
