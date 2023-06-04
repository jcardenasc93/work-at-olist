package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jcardenasc93/work-at-olist/app/db"
	m "github.com/jcardenasc93/work-at-olist/app/models"
)

func CreateBook(w http.ResponseWriter, r *http.Request, db db.ApiDB) (*ApiResponse, *ApiError) {
	bookReq := new(m.CreateBookReq)
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

func checkEmptyVals(bookReq *m.CreateBookReq) error {
	var nameDef string
	var editionDef uint16
	var pubYearDef uint32
	if bookReq.Name == nameDef {
		return errors.New("Missing name value")
	}
	if bookReq.Edition == editionDef {
		return errors.New("Missing edition value")
	}
	if uint32(bookReq.PubYear) == pubYearDef {
		return errors.New("Missing publication_year value")
	}
	if bookReq.Authors == nil {
		return errors.New("Missing authors value")
	}
	return nil
}
