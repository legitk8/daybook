package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Pool *pgxpool.Pool

func InitDB() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load env file")
	}

	connString := os.Getenv("DATABASE_URL")
	if connString == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	Pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	log.Println("Connected to database")
}
