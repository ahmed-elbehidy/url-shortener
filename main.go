package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
)

var db *sql.DB

type shortenRequest struct {
	URL string `json:"url"`
}

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

// Base62 Algorithm
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

// Building API Endpoints
func shortenURL(c *gin.Context) {
	var req shortenRequest
	if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request."})
		return
	}

	result, err := db.Exec("INSERT INTO urls (original_url, short_code) VALUES (?, ?)", req.URL, "")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	id, _ := result.LastInsertId()
	shortCode := generateShortCode(int(id))

	_, err = db.Exec("UPDATE urls SET short_code = ? WHERE id = ?", shortCode, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update short code"})
		return
	}

	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)
	c.JSON(http.StatusOK, gin.H{"original_url": req.URL, "short_url": shortURL})
}

func redirectURL(c *gin.Context) {
	shortCode := c.Param("code")
	var originalURL string
	var id int

	// The query retrieves the original_url to perform the redirect.
	err := db.QueryRow("SELECT id, original_url FROM urls WHERE short_code = ?", shortCode).Scan(&id, &originalURL)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
		return
	}

	c.Redirect(http.StatusFound, originalURL)
}

func main() {
	initDB()
	defer db.Close()

	r := gin.Default()

	r.POST("/shorten", shortenURL)
	r.GET("/:code", redirectURL)

	fmt.Println("Server is running on http://localhost:8080")
	r.Run(":8080")
}
