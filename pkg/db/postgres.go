package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func InitPostgres() *sql.DB {

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	name := os.Getenv("PG_DBNAME")

	if host == "" || port == "" || user == "" || name == "" {
		log.Fatal("Missing environment variables")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS inventories (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		code VARCHAR(50) UNIQUE NOT NULL,
		stock INTEGER NOT NULL,
		description TEXT,
		status VARCHAR(10) NOT NULL CHECK (status IN ('active', 'broken'))
	);`

	if _, err := db.Exec(createTable); err != nil {
		log.Fatal(err)
	}

	return db
}
