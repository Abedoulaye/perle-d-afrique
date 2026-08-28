package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var db *pgxpool.Pool

func init(){
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config, err := pgxpool.ParseConfig("")
	if err != nil {
		log.Fatal(err)
	}

	config.ConnConfig.User = os.Getenv("USERNAME")
	config.ConnConfig.Password = os.Getenv("PASSWORD")
	config.ConnConfig.Host = os.Getenv("DB_HOST")
	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		log.Fatalf("invalid DB_PORT: %v", err)
	}
	config.ConnConfig.Port = uint16(port)
	config.ConnConfig.Database = os.Getenv("DB_NAME")

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	db = pool
}