package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
	models "github.com/luongquochai/goBlog/db/models"
)

var Queries *models.Queries
var db *sql.DB

func InitDB(connStr string) *sql.DB {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to open a DB connection: ", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping DB: ", err)
	}

	Queries = models.New(db)

	return db
}

func InitSchema(db *sql.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username TEXT UNIQUE NOT NULL,
        password TEXT NOT NULL
    );
    CREATE TABLE IF NOT EXISTS tasks (
        id SERIAL PRIMARY KEY,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        user_id INTEGER NOT NULL REFERENCES users(id)
    );
    `
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatal("Failed to initialize schema: ", err)
	}
}
