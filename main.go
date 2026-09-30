package main

import (
	"database/sql"
	"log"
	"strings"

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

// 3. Base62 Algorithm
func generateShortCode(id int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const base = 62

	if id == 0 {
		return string(charset[0])
	}
	var shortCodeBuilder strings.Builder
	for id > 0 {
		remainder := id % base
		shortCodeBuilder.WriteByte(charset[remainder])
		id = id / base
	}
	result := shortCodeBuilder.String()
	reversed := ""
	for i := len(result) - 1; i >= 0; i-- {
		reversed += string(result[i])
	}
	return reversed
}

func main() {
	initDB()
	defer db.Close()
}
