package controllers

import (
	m "github.com/jcardenasc93/work-at-olist/app/middlewares"
	"github.com/jcardenasc93/work-at-olist/app/models"
	"net/http"
)

const nameKey string = "name"

func GetAuthors(w http.ResponseWriter, r *http.Request) (*ApiError, *ApiResponse) {
	pageId := r.Context().Value(m.PageIdKey)

	params := r.URL.Query()
	authors, err := models.GetAuthors(pageId.(int), params)
	if err != nil {
		return NewApiError(http.StatusInternalServerError, "Internal Error"), nil
	} else {
		return nil, NewApiResponse(200, authors)
	}
}
