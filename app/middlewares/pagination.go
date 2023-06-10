package middlewares

import (
	"context"
	"net/http"
	"strconv"
)

const PageIdKey = "page_id"
const LimitKey = "limit"
const DefaultLimit = 2
const PaginationKey = "pagination"

type PaginationVals struct {
	PageId int
	Limit  int
}

func Pagination(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var intPageId int
		var intLimit int
		var err error

		if r.URL.Query().Has(LimitKey) {
			limit := r.URL.Query().Get(LimitKey)
			intLimit, err = strconv.Atoi(limit)
			if err != nil {
				// TODO: Refactor errors to handle trhough ApiError
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}

		} else {
			intLimit = DefaultLimit
		}

		if r.URL.Query().Has(PageIdKey) {
			pageId := r.URL.Query().Get(PageIdKey)
			intPageId, err = strconv.Atoi(pageId)
			if err != nil {
				// TODO: Refactor errors to handle trhough ApiError
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(err.Error()))
				return
			}
		}

		ctx := context.WithValue(r.Context(), PaginationKey, &PaginationVals{
			Limit:  intLimit,
			PageId: intPageId,
		})
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func CheckPagination(context context.Context) (p *PaginationVals, ok bool) {
	pagination := context.Value(PaginationKey)
	p, ok = pagination.(*PaginationVals)
	return p, ok
}
