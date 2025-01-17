package migrations

import (
	"log"
	"web-app/app/services/core"
)

type UserTable struct{}

func (*UserTable) Up() {
	log.Println("Creating users table")

	// Initialize the service
	db, err := core.NewPostgresService()
	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}
	defer db.Close()

	// Create the users table
	_, err = db.Create(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL,
			password TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			update_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)

	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Insert into migrations table
	_, err = db.Create(`INSERT INTO migrations (name) VALUES ($1);`, "users")
	if err != nil {
		log.Fatalf("Failed to insert into migrations table: %v", err)
	}

	log.Println("Users table created")
}

func (*UserTable) Down() {
	log.Println("Dropping users table")

	// Initialize the service
	db, err := core.NewPostgresService()
	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}
	defer db.Close()

	// Drop the users table
	_, err = db.Delete("DROP TABLE users;")
	if err != nil {
		log.Fatalf("Failed to drop users table: %v", err)
	}

	// Delete from migrations table
	_, err = db.Delete(`DELETE FROM migrations WHERE name = $1;`, "users")
	if err != nil {
		log.Fatalf("Failed to delete from migrations table: %v", err)
	}

	log.Println("Users table dropped")
}
