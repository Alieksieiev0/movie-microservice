package main

import (
	"log"

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

	s := Server{MovieService{MovieDatabase{pool}}}
	err = s.Start(":3000")
	if err != nil {
		log.Fatal(err)
	}
}
