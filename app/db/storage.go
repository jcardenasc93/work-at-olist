package db

import (
	"context"
	"net/url"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

type allowedQParams struct {
	params map[string]func(string) string
}

type ApiDB interface {
	Setup() error
	CreateAuthorTable() error
	InsertAuthor(string) error
	FetchAuthors(*middlewares.PaginationVals, url.Values) ([]*models.Author, error)
	FetchBooks(*middlewares.PaginationVals, url.Values) ([]*models.Book, error)
	FetchAuthorsForBooks([]*models.Book) ([]*models.Book, error)
	applyQueryParams(string, allowedQParams, url.Values) (string, []any)
	sortAndLimit(string) string
	CreateBookTable() error
	CreateAuthorBookTable() error
	InsertBook(context.Context, *models.CreateBookReq) (*models.Book, error)
}
