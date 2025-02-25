package controllers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jcardenasc93/work-at-olist/app/middlewares"
)

func TestGetAuthorsAPI(t *testing.T) {
	t.Run("Success cases", getAuthorsAPISuccess)
	t.Run("Pagination errors", getAuthorsPaginationErr)
}

func getAuthorsAPISuccess(t *testing.T) {
	populateAuthors()
	t.Run("Fetch authors with no params", getAuthorsNoParams)
	t.Run("Fetch authors with params", getAuthorsWithParams)
}

func getAuthorsNoParams(t *testing.T) {
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
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	authors := apiRes.Data.([]interface{})
	if len(authors) != len(mockDB.Authors[:middlewares.DefaultLimit]) {
		t.Errorf("Expected %d authors but got %d", len(mockDB.Authors), len(authors))
	}
	if *apiRes.NextPage != middlewares.DefaultLimit {
		t.Errorf("Expected %d next_page value but got %d", middlewares.DefaultLimit, *apiRes.NextPage)
	}
}

func getAuthorsPaginationErr(t *testing.T) {
	handler := HTTPHandleFunc(GetAuthors, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, "/?page_id=text", nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP code %d but got %d", http.StatusBadRequest, response.StatusCode)
	}

	defer response.Body.Close()

	req = httptest.NewRequest(http.MethodGet, "/?limit=text", nil)
	resRecorder = httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response = resRecorder.Result()
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected HTTP code %d but got %d", http.StatusBadRequest, response.StatusCode)
	}
}

func getAuthorsWithParams(t *testing.T) {
	limit := 5
	resp := getAuthorsWithLimit(t, limit)
	getAuthorsWithPageId(t, limit, *resp.NextPage)
	getAuthorsNameFilter(t, "7", limit)
}

func getAuthorsWithLimit(t *testing.T, limit int) ApiResponse {
	handler := HTTPHandleFunc(GetAuthors, mockDB)
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
	authors := apiRes.Data.([]interface{})
	if len(authors) != limit {
		t.Errorf("Expected %d authors but got %d", limit, len(authors))
	}
	if *apiRes.NextPage != limit {
		t.Errorf("Expected %d next_page value but got %d", middlewares.DefaultLimit, apiRes.NextPage)
	}
	return apiRes
}

func getAuthorsWithPageId(t *testing.T, limit, pageId int) {
	handler := HTTPHandleFunc(GetAuthors, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?limit=%d&page_id=%d", limit, pageId), nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	authors := apiRes.Data.([]interface{})
	author := authors[0].(map[string]any)
	id := author["id"].(float64)
	expetedAuth := mockDB.Authors[limit]
	if id != float64(expetedAuth.Id) {
		t.Errorf("Expected %d author's id but got %v", expetedAuth.Id, id)
	}
}

func getAuthorsNameFilter(t *testing.T, name string, limit int) {
	handler := HTTPHandleFunc(GetAuthors, mockDB)
	testHandler := middlewares.Pagination(handler)
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?limit=%d&name=%s", 3, name), nil)
	resRecorder := httptest.NewRecorder()
	testHandler.ServeHTTP(resRecorder, req)
	response := resRecorder.Result()
	apiRes := decodeResponseBody[ApiResponse](t, response.Body)
	authors := apiRes.Data.([]interface{})
	if len(authors) != 1 {
		t.Errorf("Expected %d authors but got %d", 1, len(authors))
	}
}
