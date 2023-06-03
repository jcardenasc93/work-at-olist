package controllers

import (
	"net/http"

	"github.com/jcardenasc93/work-at-olist/app/db"
	m "github.com/jcardenasc93/work-at-olist/app/middlewares"
)

func GetAuthors(w http.ResponseWriter, r *http.Request, db db.ApiDB) (*ApiResponse, *ApiError) {
	pagination := r.Context().Value(m.PaginationKey)
	p, ok := pagination.(*m.PaginationVals)
	if ok == false {
		return nil, NewApiError(http.StatusInternalServerError, "Internal Error")
	}
	params := r.URL.Query()
	authors, err := db.FetchAuthors(p, params)
	if err != nil {
		return nil, NewApiError(http.StatusInternalServerError, "Couldn't fetch authors from database")
	}

	if len(authors) > 0 {
		var nextPage int
		if p.PageId == 0 {
			nextPage = p.Limit
		} else {
			nextPage = p.PageId + p.Limit
		}
		return NewApiResponse(200, authors, &nextPage), nil
	}

	return NewApiResponse(200, authors, nil), nil
}
