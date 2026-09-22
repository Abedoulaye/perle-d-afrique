package main

import (
	"context"
	"log"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func init(){
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environmental variables")
	}

	pool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}

	db = pool
}