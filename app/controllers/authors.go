package controllers

import (
	"net/http"

	"github.com/jcardenasc93/work-at-olist/app/db"
	m "github.com/jcardenasc93/work-at-olist/app/middlewares"
)

func GetAuthors(w http.ResponseWriter, r *http.Request, db db.ApiDB) (*ApiError, *ApiResponse) {
	pagination := r.Context().Value(m.PaginationKey)
	p, ok := pagination.(*m.PaginationVals)
	if ok == false {
		return NewApiError(http.StatusInternalServerError, "Internal Error"), nil
	}
	params := r.URL.Query()
	authors, err := db.FetchAuthors(p, params)
	if err != nil {
		return NewApiError(http.StatusInternalServerError, "Internal Error"), nil
	} else {
		var nextPage int
		if p.PageId == 0 {
			nextPage = p.Limit
		} else {
			nextPage = p.PageId + p.Limit
		}
		return nil, NewApiResponse(200, authors, nextPage)
	}
}
