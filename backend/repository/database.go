package repository

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Sync() {
	DB = connectDB()
}

func connectDB() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URI"))

	if err != nil {
		log.Fatal("Error while connecting to database.", err)
	}

	log.Println("Database successfully connected.")

	return db
}
