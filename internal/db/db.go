// internal/db/db.go

package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	// DB is the database connection pool
	DB *sql.DB
)

// Init initializes the database connection
func Init() {
	var err error

	dbHost := "localhost" //"host.docker.internal" // use localhost for native
	dbPort := "3306"
	dbUser := "root"
	dbPassword := ""
	dbName := "myapp"

	// Create DSN string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open database connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}

	// Ping database to verify connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	log.Println("Database connection established")
}
