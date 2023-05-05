package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/jcardenasc93/work-at-olist/app/models"
)

var responseHeaders = map[string]string{
	"Content-Type": "application/json",
}

const nameKey string = "name"

func GetAuthors(w http.ResponseWriter, r *http.Request) {
	for key, value := range responseHeaders {
		w.Header().Set(key, value)
	}

	params := r.URL.Query()
	authors, err := models.GetAuthors(params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		json.NewEncoder(w).Encode(authors)
	}

}
