package main

import (
	"flag"
	"log"

	"github.com/saint0x/file-storage-app/backend/internal/api"
	"github.com/saint0x/file-storage-app/backend/scripts"
)

func main() {
	populateDB := flag.Bool("populate", false, "Populate the database with sample data")
	flag.Parse()

	if *populateDB {
		scripts.PopulateSampleData()
		return
	}

	// Start the server
	if err := api.StartServer(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
