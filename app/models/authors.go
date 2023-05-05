package models

import (
	"database/sql"
	"fmt"
	"net/url"
)

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func newAuthor(id uint64, name string) Author {
	return Author{
		Id:   id,
		Name: name,
	}
}

const nameKey = "name"

func filterByName(baseQuery string, name string) (query string) {
	filter := `WHERE name LIKE '%'||?||'%'`
	query = fmt.Sprintf("%s %s", baseQuery, filter)
	return
}

func GetAuthors(queryParams url.Values) ([]Author, error) {
	var authors []Author
	var rows *sql.Rows
	var err error

	query := `SELECT id, name FROM authors`

	if queryParams.Has(nameKey) {
		nameValue := string(queryParams.Get(nameKey))
		query := filterByName(query, nameValue)
		rows, err = execQuery(query, nameValue)
		if err != nil {
			return authors, err
		}
	} else {
		rows, err = execQuery(query)
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

		authors = append(authors, newAuthor(id, name))
	}

	return authors, nil
}
