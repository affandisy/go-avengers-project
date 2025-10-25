package db

import (
	"avenger/internal/domain"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
		log.Fatal("Failed to open database connection:", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database", err)
	}

	slog.Info("PostgreSQL (database/sql) connected successfully",
		slog.String("host", host),
		slog.String("port", port),
		slog.String("database", name),
	)

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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying SQL DB", err)
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	slog.Info("PostgreSQL (GORM) connected successfully",
		slog.String("host", host),
		slog.String("port", port),
		slog.String("database", name),
	)

	if err := db.AutoMigrate(&domain.User{}, &domain.Recipe{}); err != nil {
		log.Fatal("Migration failed")
	}

	slog.Info("Database tables migrated successfully (users, recipes)")

	return db
}
