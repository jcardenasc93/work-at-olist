package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

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
