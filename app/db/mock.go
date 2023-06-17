package db

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

type MockDB struct {
	Authors      []*models.Author
	Books        []*models.Book
	AuthorsBooks map[float64][]float64
}

func NewMockDB() *MockDB { return &MockDB{AuthorsBooks: make(map[float64][]float64)} }

func (m *MockDB) SetAuthors(authors []*models.Author) {
	m.Authors = authors
}

func (m *MockDB) SetBooks(books []*models.Book) {
	m.Books = books
}

func (m *MockDB) Setup() error { return nil }

func (m *MockDB) CreateAuthorTable() error { return nil }

func (m *MockDB) CreateBookTable() error { return nil }

func (m *MockDB) CreateAuthorBookTable() error { return nil }

func (m *MockDB) InsertAuthor(string) error { return nil }

func (m *MockDB) InsertBook(c context.Context, req *models.CreateBookReq) (*models.Book, error) {
	book := models.NewBook(float64(len(m.Books)+1), req.Name, req.Edition, req.PubYear, req.Authors)
	m.Books = append(m.Books, book)
	m.AuthorsBooks[book.Id] = req.Authors
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

func filterBooksByName(books []*models.Book, name string) []*models.Book {
	result := []*models.Book{}
	for _, book := range books {
		if strings.Contains(book.Name, name) {
			result = append(result, book)
		}
	}
	return result

}

func filterBooksByPubYear(books []*models.Book, year string) []*models.Book {
	y, _ := strconv.Atoi(year)
	value := float64(y)
	result := []*models.Book{}

	for _, book := range books {
		if book.PubYear == value {
			result = append(result, book)
		}
	}
	return result

}

func (m *MockDB) FetchBooks(pagination *middlewares.PaginationVals, vals url.Values) ([]*models.Book, error) {
	const nameKey string = "name"
	const pubYearKey string = "publication_year"
	const editionKey string = "edition"
	const authorKey string = "author"
	var books []*models.Book
	limit := pagination.Limit
	pageId := pagination.PageId

	filters := allowedFilters[*models.Book]{
		params: map[string]func([]*models.Book, string) []*models.Book{
			nameKey:    filterBooksByName,
			pubYearKey: filterBooksByPubYear,
		},
	}
	books = applyFilters(m.Books, filters, vals)

	books = books[pageId:]
	if limit < len(books) {
		return books[:limit], nil
	}
	return books, nil
}

func (m *MockDB) FetchAuthorsForBooks([]*models.Book) ([]*models.Book, error) {
	return m.Books, nil
}

func (m *MockDB) applyQueryParams(baseQuery string, q allowedQParams, params url.Values) (query string, paramVals []any) {
	return
}

type modelType interface {
	*models.Author | *models.Book
}

type allowedFilters[T modelType] struct {
	params map[string]func([]T, string) []T
}

func applyFilters[T modelType](data []T, q allowedFilters[T], vals url.Values) []T {
	// result := []T{}
	for key, fun := range q.params {
		if vals.Has(key) {
			keyVal := vals.Get(key)
			data = fun(data, keyVal)
		}
	}

	return data
}

func (m *MockDB) filterByName(name string) (authors []*models.Author) {
	for _, author := range m.Authors {
		if strings.Contains(author.Name, name) {
			authors = append(authors, author)
		}
	}
	return authors
}

func (m *MockDB) sortAndLimit(string) string { return "" }
