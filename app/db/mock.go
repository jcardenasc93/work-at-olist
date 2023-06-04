package db

import (
	"context"
	"net/url"
	"strings"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

type MockDB struct {
	Authors      []*models.Author
	Books        []*models.Book
	AuthorsBooks []map[uint64][]uint64
}

func NewMockDB() *MockDB { return &MockDB{} }

func (m *MockDB) SetAuthors(authors []*models.Author) {
	m.Authors = authors
}

func (m *MockDB) Setup() error { return nil }

func (m *MockDB) CreateAuthorTable() error { return nil }

func (m *MockDB) CreateBookTable() error { return nil }

func (m *MockDB) CreateAuthorBookTable() error { return nil }

func (m *MockDB) InsertAuthor(string) error { return nil }

func (m *MockDB) InsertBook(c context.Context, req *models.CreateBookReq) (*models.Book, error) {
	book := models.NewBook(uint64(len(m.Books)+1), req.Name, req.Edition, req.PubYear, req.Authors)
	m.Books = append(m.Books, book)
	m.AuthorsBooks = append(m.AuthorsBooks, map[uint64][]uint64{book.Id: req.Authors})
	return book, nil
}

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
