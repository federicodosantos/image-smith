package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/federicodosantos/image-smith/db"
	"github.com/federicodosantos/image-smith/internal/bootstrap"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	var PORT = os.Getenv("APP_PORT")
	if PORT == "" {
		log.Fatalf("port undefined")
	}

	db := db.DBConnection()
	defer db.Close()

	mux := http.NewServeMux()

	b := bootstrap.NewBootstrap(db, mux)

	b.InitApp()

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: mux,
	}

	log.Printf("Running the server on port %s", PORT)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("cannot running the server : ")
	}
}
