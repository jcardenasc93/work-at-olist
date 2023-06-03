package db

import (
	"context"
	"net/url"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

type ApiDB interface {
	Setup() error
	CreateAuthorTable() error
	InsertAuthor(string) error
	FetchAuthors(*middlewares.PaginationVals, url.Values) ([]*models.Author, error)
	sortAndLimit(string) string
	CreateBookTable() error
	CreateAuthorBookTable() error
	InsertBook(context.Context, *models.CreateBookReq) (*models.Book, error)
}
