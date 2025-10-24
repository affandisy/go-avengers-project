package db

import (
	"avenger/internal/domain"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

	return db
}

func InitPostgresGORM() *gorm.DB {

	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	name := os.Getenv("PG_DBNAME")

	if host == "" || port == "" || user == "" || name == "" {
		log.Fatal("Missing environment variables")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pass, name)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	if err := db.AutoMigrate(&domain.User{}, &domain.Recipe{}); err != nil {
		log.Fatal("Migration failed")
	}

	log.Println("Database connected and migrated successfully")

	return db
}
