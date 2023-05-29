package models

import (
	"fmt"
)

type Author struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

func NewAuthor(id uint64, name string) *Author {
	return &Author{
		Id:   id,
		Name: name,
	}
}

const nameKey = "name"

func FilterByName(baseQuery string) (query string) {
	filter := `AND name LIKE '%'||?||'%'`
	query = fmt.Sprintf("%s %s", baseQuery, filter)
	return
}
