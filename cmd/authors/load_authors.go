package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/jcardenasc93/work-at-olist/app/db"
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
	db, err := db.NewSQLiteDB()
	if err != nil {
		log.Fatal(err)
	}
	err = db.CreateAuthorsTable()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Importing authors from csv file...")
	for {
		author, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if author[0] != "name" {
			err = db.InsertAuthor(author[0])
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	defer file.Close()
	log.Println("Done!")

}
