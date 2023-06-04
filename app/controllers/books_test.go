package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
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
