package main

import (
	"database/sql"
	"log"
	"os"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Unable to reach database: %v\n", err)
	}

	log.Println("Connected to database successfully")
}

func CloseDB() {
	db.Close()
}
