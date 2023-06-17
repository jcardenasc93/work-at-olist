package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
)

const contentType = "application/json"

func TestCreateBookAPI(t *testing.T) {
	t.Run("Success cases", createBookAPISuccess)
	t.Run("Failing cases", createBookAPIErr)
}

func createBookAPISuccess(t *testing.T) {
	populateAuthors()
	server := httptest.NewServer(HTTPHandleFunc(CreateBook, mockDB))
	body := map[string]any{}
	body["name"] = "Testing book"
	body["edition"] = float64(3)
	body["publication_year"] = float64(2002)
	body["authors"] = []float64{1, 2}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Error(err.Error())
	}
	resp, err := http.Post(server.URL, contentType, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Logf("Couldn't make request: %s", err.Error())
	}

	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected %d status code but got %d", http.StatusCreated, resp.StatusCode)
	}

	defer resp.Body.Close()

	apiRes := decodeResponseBody[ApiResponse](t, resp.Body)
	data := apiRes.Data.(map[string]interface{})
	for k := range body {
		checkRespBody(t, body, data, k)
	}

	dbAuthorsBooks := mockDB.AuthorsBooks
	bookId := data["id"].(float64)
	eq := reflect.DeepEqual(dbAuthorsBooks[bookId], body["authors"])
	if eq != true {
		t.Error("Relationship authors_book are wrong")
	}
}

func checkRespBody(t *testing.T, expected map[string]any, apiRes map[string]any, attr string) {
	if attr == "authors" {
		apiAuthors := apiRes[attr].([]any)
		authors := make([]float64, len(apiAuthors))
		for i, v := range apiAuthors {
			authors[i] = v.(float64)
		}
		if reflect.DeepEqual(expected[attr], authors) == false {
			t.Error("Not equal")
		}
		return
	}
	if expected[attr] != apiRes[attr] {
		t.Errorf("Error in attribute %s. Expected %v but got %v", attr, expected[attr], apiRes[attr])
	}
}

func createBookAPIErr(t *testing.T) {
	t.Run("Missing attributes", createBookMissinAtttr)
	t.Run("Wrong values", TestCreateBookWrongVals)
}

func createMissAttrCases() []map[string]any {
	bodyReqs := []map[string]any{}
	body := map[string]any{"name": "Testing book"}
	body1 := map[string]any{"name": "Testing book", "edition": 2}
	body2 := map[string]any{"name": "Testing book", "edition": 2, "publication_year": 2022}
	body3 := map[string]any{"authors": []uint64{1, 2}}
	bodyReqs = append(bodyReqs, body)
	bodyReqs = append(bodyReqs, body1)
	bodyReqs = append(bodyReqs, body2)
	bodyReqs = append(bodyReqs, body3)
	return bodyReqs
}
func createBookMissinAtttr(t *testing.T) {
	server := httptest.NewServer(HTTPHandleFunc(CreateBook, mockDB))
	cases := createMissAttrCases()
	for _, body := range cases {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			t.Error(err.Error())
		}
		resp, err := http.Post(server.URL, contentType, bytes.NewBuffer(jsonBody))
		if err != nil {
			t.Logf("Couldn't make request: %s", err.Error())
		}
		if resp.StatusCode != http.StatusBadRequest {
			t.Errorf("Expected %d status code but got %d", http.StatusBadRequest, resp.StatusCode)
		}
	}
}

func TestCreateBookWrongVals(t *testing.T) {
	server := httptest.NewServer(HTTPHandleFunc(CreateBook, mockDB))
	body := make(map[string]any)
	body["name"] = 12
	body["edition"] = "2012"
	body["publication_year"] = "AAA"
	body["authors"] = 1
	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Error(err.Error())
	}
	resp, err := http.Post(server.URL, contentType, bytes.NewBuffer(jsonBody))
	if err != nil {
		t.Logf("Couldn't make request: %s", err.Error())
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected %d status code but got %d", http.StatusBadRequest, resp.StatusCode)
	}
}

func TestGetBookAPI(t *testing.T) {
	populateAuthors()
	populateBooks()
	t.Run("Success cases", getBookAPISuccess)
	// t.Run("Failing cases", createBookAPIErr)
}

func getBookAPISuccess(t *testing.T) {
	t.Run("Success no params", getBookAPISuccessNoParams)
	t.Run("Success with params", getBookAPIWithParams)
}

func getBookAPISuccessNoParams(t *testing.T) {
	handler := HTTPHandleFunc(GetBooks, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP code %d but got %d", http.StatusOK, response.StatusCode)
	}
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	books := apiRes.Data.([]interface{})
	expectedLen := len(mockDB.Books[:middlewares.DefaultLimit])
	if len(books) != expectedLen {
		t.Errorf("Expected %d books but got %d", expectedLen, len(books))
	}
	if *apiRes.NextPage != middlewares.DefaultLimit {
		t.Errorf("Expected %d next_page value but got %d", middlewares.DefaultLimit, *apiRes.NextPage)
	}
}

func getBookAPIWithParams(t *testing.T) {
	limit := 3
	resp := getBooksWithLimit(t, limit)
	getBooksWithPageId(t, limit, *resp.NextPage)
	getBooksNameFilter(t, "7", limit)
	year := mockDB.Books[0].PubYear
	getBooksPubYearFilter(t, year, limit)
}

func getBooksWithLimit(t *testing.T, limit int) ApiResponse {
	handler := HTTPHandleFunc(GetBooks, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?limit=%d", limit), nil)
	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	req.URL.RawQuery = q.Encode()
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Errorf("Expected HTTP code %d but got %d", http.StatusOK, response.StatusCode)
	}
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	books := apiRes.Data.([]interface{})
	if len(books) != limit {
		t.Errorf("Expected %d books but got %d", limit, len(books))
	}
	if *apiRes.NextPage != limit {
		t.Errorf("Expected %d next_page value but got %d", middlewares.DefaultLimit, apiRes.NextPage)
	}
	return apiRes
}

func getBooksWithPageId(t *testing.T, limit int, pageId int) {
	handler := HTTPHandleFunc(GetBooks, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?limit=%d&page_id=%d", limit, pageId), nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	books := apiRes.Data.([]interface{})
	book := books[0].(map[string]any)
	id := book["id"].(float64)
	expetedBook := mockDB.Authors[limit]
	if id != float64(expetedBook.Id) {
		t.Errorf("Expected %d book's id but got %v", expetedBook.Id, id)
	}
}

func getBooksNameFilter(t *testing.T, name string, limit int) {
	handler := HTTPHandleFunc(GetBooks, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?limit=%d&name=%s", 3, name), nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	books := apiRes.Data.([]interface{})
	if len(books) != 1 {
		t.Errorf("Expected %d books but got %d", 1, len(books))
	}
}

func getBooksPubYearFilter(t *testing.T, year float64, limit int) {
	handler := HTTPHandleFunc(GetBooks, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?limit=%d&year=%v", 3, year), nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	books := apiRes.Data.([]interface{})
	if len(books) == 0 {
		t.Errorf("Expected at least one book")
	}
}
