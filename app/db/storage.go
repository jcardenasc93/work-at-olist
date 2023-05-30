package db

import (
	"database/sql"
	"net/url"

	"github.com/jcardenasc93/work-at-olist/app/models"
)

type ApiDB interface {
	CreateAuthorsTable() error
	InsertAuthor(string) error
	FetchAuthors(int, url.Values) ([]*models.Author, error)
	execQuery(string, ...any) (*sql.Rows, error)
	filterByName(string) string
}
