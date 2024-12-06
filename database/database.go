package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Connect initializes the database connection and returns *sql.DB and an error
func Connect() (*sql.DB, error) {
	var err error

	// Load database configuration from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Define DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open the database connection
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	// Test the connection
	err = DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("Error pinging the database: %v", err)
	}

	log.Println("Database connected successfully")
	return DB, nil
}
