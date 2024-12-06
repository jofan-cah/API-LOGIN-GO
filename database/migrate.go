package database

import (
	"database/sql"
	"log"
)

// Migrate akan menjalankan migrasi untuk membuat tabel `users`
func Migrate(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL,
		email VARCHAR(100) UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	);
	`

	// Menjalankan query untuk membuat tabel
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Migration failed: %s", err)
	}

	log.Println("Migration completed")
}
