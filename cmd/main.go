package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/federicodosantos/image-smith/db"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("cannot environment variable : %e", err)
	}

	var PORT = os.Getenv("APP_PORT")
	if PORT == "" {
		log.Fatalf("port undefined")
	}

	db := db.DBConnection()
	defer db.Close()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: mux,
	}

	log.Printf("Running the server on port %s", PORT)
	if err = server.ListenAndServe(); err != nil {
		log.Fatalf("cannot running the server : ")
	}
}
