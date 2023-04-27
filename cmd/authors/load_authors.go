package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jcardenasc93/work-at-olist/app/models"
)

const dbName = "sqlite.db"

func main() {
	os.Remove(dbName)
	var csvFile string
	flag.StringVar(&csvFile, "csv", "", "CSV file path")
	flag.Parse()

	file, err := os.Open(filepath.Base("../../input.csv"))
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(file)
	err = models.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	err = models.CreateAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	for {
		author, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if author[0] != "name" {
			err = models.InsertAuthor(author[0])
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	defer file.Close()

}
