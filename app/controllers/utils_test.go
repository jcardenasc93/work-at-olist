package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/jcardenasc93/work-at-olist/app/db"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

var mockDB *db.MockDB = db.NewMockDB()

func populateAuthors() {
	var authors []*models.Author
	mockDB.Authors = authors
	for i := 0; i < 10; i++ {
		author := models.NewAuthor(uint64(i+1), fmt.Sprintf("Author %d", i+1))
		authors = append(authors, author)
	}
	mockDB.SetAuthors(authors)
}

func populateBooks() {
	var books []*models.Book
	mockDB.Books = books
	for i := 0; i < 10; i++ {
		authors := []float64{}
		if i > 3 {
			authors = append(authors, float64(randIntRange(0, len(mockDB.Authors)-1)))
			authors = append(authors, float64(randIntRange(0, len(mockDB.Authors)-1)))
		}
		book := models.NewBook(float64(i+1), fmt.Sprintf("Book %d", i+1), float64(randIntRange(1, 6)), float64(randIntRange(1990, 2022)), authors)
		books = append(books, book)
	}
	mockDB.SetBooks(books)
}

func randIntRange(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

type respBody interface {
	ApiResponse | ApiError
}

func decodeResponseBody[T respBody](t *testing.T, body io.ReadCloser) T {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}
	var respBody T
	err = json.Unmarshal(data, &respBody)
	if err != nil {
		t.Error(err)
	}
	return respBody
}
