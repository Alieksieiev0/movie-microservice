package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}
	pool, err := ConnectDB("DATABASE_URL")
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	logger := log.New(os.Stdout, "SERVICE INFO: ", log.LstdFlags)
	loggingService := NewLoggingService(logger, NewMovieService(NewMovieDatabase(pool)))
	s := NewServer(loggingService)
	err = s.Start(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
