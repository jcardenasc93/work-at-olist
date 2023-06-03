package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jcardenasc93/work-at-olist/app/db"
	m "github.com/jcardenasc93/work-at-olist/app/models"
)

func CreateBook(w http.ResponseWriter, r *http.Request, db db.ApiDB) (*ApiResponse, *ApiError) {
	bookReq := new(m.CreateBookReq)
	err := json.NewDecoder(r.Body).Decode(bookReq)
	if err != nil {
		return nil, NewApiError(http.StatusBadRequest, err.Error())
	}
	book, err := db.InsertBook(r.Context(), bookReq)
	if err != nil {
		return nil, NewApiError(http.StatusInternalServerError, err.Error())
	}
	return NewApiResponse(http.StatusCreated, book, nil), nil
}
