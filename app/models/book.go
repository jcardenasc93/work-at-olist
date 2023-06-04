package models

type Book struct {
	Id      float64   `json:"id"`
	Name    string    `json:"name"`
	Edition float64   `json:"edition"`
	PubYear float64   `json:"publication_year"`
	Authors []float64 `json:"authors"`
}

func NewBook(id float64, name string, edition float64, pubYear float64, authors []float64) *Book {
	return &Book{
		Id:      id,
		Name:    name,
		Edition: edition,
		PubYear: pubYear,
		Authors: authors,
	}
}

type CreateBookReq struct {
	Name    string    `json:"name"`
	Edition float64   `json:"edition"`
	PubYear float64   `json:"publication_year"`
	Authors []float64 `json:"authors"`
}
