package db

import (
	"net/url"
	"strings"

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
	const nameKey string = "name"
	var authors []*models.Author
	limit := pagination.Limit
	pageId := pagination.PageId

	if vals.Has(nameKey) {
		nameValue := string(vals.Get(nameKey))
		authors = m.filterByName(nameValue)
	} else {
		authors = m.Authors
	}

	authors = authors[pageId:]
	if limit < len(authors) {
		return authors[:limit], nil
	}
	return authors, nil
}

func (m *MockDB) filterByName(name string) (authors []*models.Author) {
	for _, author := range m.Authors {
		if strings.Contains(author.Name, name) {
			authors = append(authors, author)
		}
	}
	return authors
}

func (m *MockDB) sortAndLimit(baseQ string) (query string) { return query }
