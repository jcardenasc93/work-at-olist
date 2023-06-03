package models

type Book struct {
	Id      uint64   `json:"id"`
	Name    string   `json:"name"`
	Edition uint16   `json:"edition"`
	PubYear uint32   `json:"publication_year"`
	Authors []uint64 `json:"authors"`
}

func NewBook(id uint64, name string, edition uint16, pubYear uint32, authors []uint64) *Book {
	return &Book{
		Id:      id,
		Name:    name,
		Edition: edition,
		PubYear: pubYear,
		Authors: authors,
	}
}

type CreateBookReq struct {
	Name    string   `json:"name"`
	Edition uint16   `json:"edition"`
	PubYear uint32   `json:"publication_year"`
	Authors []uint64 `json:"authors"`
}
