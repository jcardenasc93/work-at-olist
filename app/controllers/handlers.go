package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jcardenasc93/work-at-olist/app/db"
)

type apiFunc func(http.ResponseWriter, *http.Request, db.ApiDB) (*ApiError, *ApiResponse)

type ApiError struct {
	StatusCode int    `json:"status_code"`
	Msg        string `json:"message"`
}

func NewApiError(statusCode int, msg string) *ApiError {
	return &ApiError{
		StatusCode: statusCode,
		Msg:        msg,
	}
}

type ApiResponse struct {
	StatusCode int `json:"status_code"`
	Data       any `json:"data"`
}

func NewApiResponse(statusCode int, data any) *ApiResponse {
	return &ApiResponse{
		StatusCode: statusCode,
		Data:       data,
	}
}

func WriteHttpResponse(w http.ResponseWriter, statusCode int, value any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(value)
}

func HTTPHandleFunc(f apiFunc, db db.ApiDB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err, resp := f(w, r, db); err != nil {
			// Handle Error
			WriteHttpResponse(w, err.StatusCode, err)
			return
		} else {
			WriteHttpResponse(w, resp.StatusCode, resp)
		}
	}
}
