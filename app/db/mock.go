package db

import (
	"net/url"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

type MockDB struct {
	Authors []*models.Author
}

func NewMockDB() *MockDB { return &MockDB{} }

func (m *MockDB) SetAuthors(authors []*models.Author) {
	m.Authors = authors
}

func (m *MockDB) CreateAuthorsTable() error { return nil }

func (m *MockDB) InsertAuthor(string) error { return nil }

func (m *MockDB) FetchAuthors(pagination *middlewares.PaginationVals, vals url.Values) ([]*models.Author, error) {
	limit := pagination.Limit
	pageId := pagination.PageId

	if pageId == 0 {
		return m.Authors[:limit], nil
	}

	return m.Authors, nil
}

func (m *MockDB) filterByName(baseQ string) (query string) { return query }

func (m *MockDB) sortAndLimit(baseQ string) (query string) { return query }
