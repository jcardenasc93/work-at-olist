package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(fileName string) (*sql.DB, error) {
	if fileName == "" {
		return nil, errors.New("DB name couldn't be empty")
	}
	dbConn, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}

	err = dbConn.Ping()
	if err != nil {
		return nil, err
	}

	dbConn.SetMaxOpenConns(5)

	fmt.Println("DB connection success!!!")

	createAuthorsTable := `
    CREATE TABLE IF NOT EXISTS authors (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name VARCHAR(64) NOT NULL
    )
    `

	stmt, err := dbConn.Prepare(createAuthorsTable)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func InsertAuthor(dbConn *sql.DB, authorName string) error {
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
