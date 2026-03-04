package main

import (
	"log"
	"os"

	"github.com/UGRORF/price-tracker/internal/api"
)

func main() {
	logger := log.New(os.Stdout, "PRICE-TRACKER: ", log.LstdFlags|log.Lshortfile)

	srv := api.NewServer(logger)

	if err := srv.Start(); err != nil {
		logger.Fatal("Server failed:", err)
	}
}
