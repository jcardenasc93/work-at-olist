package db

import (
	"database/sql"

	"github.com/jcardenasc93/work-at-olist/app/models"
)

type ApiDB interface {
	CreateAuthorsTable() error
	InsertAuthor(string) error
	FetchAuthors(int, ...any) ([]*models.Author, error)
	execQuery(query string, params ...any) (*sql.Rows, error)
}
