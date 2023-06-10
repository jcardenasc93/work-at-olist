package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jcardenasc93/work-at-olist/app/db"
	mid "github.com/jcardenasc93/work-at-olist/app/middlewares"
	mod "github.com/jcardenasc93/work-at-olist/app/models"
)

func CreateBook(w http.ResponseWriter, r *http.Request, db db.ApiDB) (*ApiResponse, *ApiError) {
	bookReq := new(mod.CreateBookReq)
	err := json.NewDecoder(r.Body).Decode(bookReq)
	if err != nil {
		return nil, NewApiError(http.StatusBadRequest, "Invalid request body")
	}
	err = checkEmptyVals(bookReq)
	if err != nil {
		return nil, NewApiError(http.StatusBadRequest, err.Error())
	}
	book, err := db.InsertBook(r.Context(), bookReq)
	if err != nil {
		return nil, NewApiError(http.StatusInternalServerError, err.Error())
	}
	return NewApiResponse(http.StatusCreated, book, nil), nil
}

func GetBooks(w http.ResponseWriter, r *http.Request, db db.ApiDB) (*ApiResponse, *ApiError) {
	p, ok := mid.CheckPagination(r.Context())
	if ok == false {
		return nil, NewApiError(http.StatusInternalServerError, "Internal Error")
	}
	params := r.URL.Query()
	books, err := db.FetchBooks(p, params)
	if err != nil {
		return nil, NewApiError(http.StatusBadRequest, "Error fetching books")
	}
	if len(books) > 0 {
		var nextPage int
		if p.PageId == 0 {
			nextPage = p.Limit
		} else {
			nextPage = p.PageId + p.Limit
		}
		return NewApiResponse(200, books, &nextPage), nil
	}
	return NewApiResponse(http.StatusOK, books, nil), nil
}

func checkEmptyVals(bookReq *mod.CreateBookReq) error {
	var nameDef string
	var editionDef float64
	var pubYearDef float64
	if bookReq.Name == nameDef {
		return errors.New("Missing name value")
	}
	if bookReq.Edition == editionDef {
		return errors.New("Missing edition value")
	}
	if bookReq.PubYear == pubYearDef {
		return errors.New("Missing publication_year value")
	}
	if bookReq.Authors == nil {
		return errors.New("Missing authors value")
	}
	return nil
}
