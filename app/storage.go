package main

import "github.com/jcardenasc93/work-at-olist/app/models"

type ApiDB interface {
	CreateAuthorsTable() error
	InsertAuthor(string) error
	FetchAuthors(int) ([]*models.Author, error)
}
