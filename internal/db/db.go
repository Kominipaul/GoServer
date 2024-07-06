// internal/db/db.go

package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DB is the database connection pool
	DB *sql.DB
)

// Init initializes the database connection
func Init() {
	var err error
	dbUser := os.Getenv("root")
	dbPassword := os.Getenv("")
	dbName := os.Getenv("users")
	dbHost := os.Getenv("localhost")
	dbPort := os.Getenv("3306")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	log.Println("Database connection established")
}

