package models

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
