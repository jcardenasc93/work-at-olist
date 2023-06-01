package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jcardenasc93/work-at-olist/app/db"
	"github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
)

var mockDB *db.MockDB = db.NewMockDB()

func populateAuthors() {
	var authors []*models.Author
	for i := 0; i < 10; i++ {
		author := models.NewAuthor(uint64(i+1), fmt.Sprintf("Author %d", i+1))
		authors = append(authors, author)
	}
	mockDB.SetAuthors(authors)
}

func decodeResponseBody(t *testing.T, body io.ReadCloser) ApiResponse {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		t.Error(err)
	}
	var authors ApiResponse
	err = json.Unmarshal(data, &authors)
	if err != nil {
		t.Error(err)
	}
	return authors
}

func TestGetAuthorsAPINoParams(t *testing.T) {
	populateAuthors()
	handler := HTTPHandleFunc(GetAuthors, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP code %d but got %d", http.StatusOK, response.StatusCode)
	}
	apiRes := decodeResponseBody(t, response.Body)
	authors, _ := apiRes.Data.([]interface{})
	if len(authors) != len(mockDB.Authors[:middlewares.DefaultLimit]) {
		t.Errorf("Expected %d authors but got %d", len(mockDB.Authors), len(authors))
	}
	if apiRes.NextPage != middlewares.DefaultLimit {
		t.Errorf("Expected %d next_page value but got %d", middlewares.DefaultLimit, apiRes.NextPage)
	}
}
