package main

import (
	"log"
	"login-example/api"
	"login-example/db"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("API_PORT")
	db, err := db.NewSqlStorage()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	server := api.NewApiServer(port, db)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
