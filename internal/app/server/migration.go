package server

import (
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

// Migration -.
func init() {
	pgURL := os.Getenv("DATABASE_DSN")
	if pgURL == "" {
		return
	}
	db, err := goose.OpenDBWithDriver("postgres", pgURL+"?sslmode=disable")
	if err != nil {
		log.Printf("goose: failed to open DB: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := goose.Up(db, "/migrations/"); err != nil {
		log.Printf("goose up: %v", err)
		os.Exit(1)
	}
}
