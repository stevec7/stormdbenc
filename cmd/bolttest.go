package main

import (
	"flag"
	"log"
)

func main() {
	filename := flag.String("filename", "", "database file")
	db, err := boltenc.setupDB(filename)
	if err != nil {
		log.Fatalf("Error setting up db, %s", err)
	}

}
