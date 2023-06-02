package db

import (
	"net/url"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

type ApiDB interface {
	CreateAuthorsTable() error
	InsertAuthor(string) error
	FetchAuthors(*middlewares.PaginationVals, url.Values) ([]*models.Author, error)
	sortAndLimit(string) string
}
