package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"io"
	"log"
	"os"

	"github.com/jcardenasc93/work-at-olist/db"
)

const dbName = "sqlite.db"

func main() {
	var csvFile string
	flag.StringVar(&csvFile, "csv", "", "CSV file path")
	flag.Parse()

	file, err := os.Open(csvFile)
	if err != nil {
		log.Fatal(err)
	}
	csvReader := csv.NewReader(file)
	dbConn, err := db.InitDB(dbName)
	if err != nil {
		log.Fatal(err)
	}
	for {
		author, err := csvReader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if author[0] != "name" {
			err = db.InsertAuthor(dbConn, author[0])
			if err != nil {
				log.Fatal(err)
			}
		}

	}
	defer dbConn.Close()
	defer file.Close()

}
