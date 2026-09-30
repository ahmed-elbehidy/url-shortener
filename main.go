package main

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite", "./urls.db")
	if err != nil {
		log.Fatal("Error opening the database:", err)
	}

	// Enable WAL mode to prevent database locking under high load
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Println("Failed to enable WAL mode:", err)
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS urls (
		"id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		"original_url" TEXT NOT NULL,
		"short_code" TEXT,
		"clicks" INTEGER DEFAULT 0
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Error creating the table:", err)
	}
}

func main() {
	initDB()
	defer db.Close()
}
