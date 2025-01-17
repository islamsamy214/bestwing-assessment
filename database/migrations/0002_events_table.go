package migrations

import (
	"log"

	"web-app/app/services/core"
)

type EventTable struct{}

func (*EventTable) Up() {
	log.Println("Creating events table")

	// Initialize the service
	db, err := core.NewPostgresService()
	if err != nil {
		log.Fatalf("Failed to initialize database service: %v", err)
	}
	defer db.Close()

	// Create the table
	_, err = db.Create(`
		CREATE TABLE IF NOT EXISTS events (
			id SERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			date DATE NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			user_id INTEGER NOT NULL,
			CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
	`)

	// Insert into migrations table
	_, err = db.Create(`INSERT INTO migrations (name) VALUES ($1);`, "events")
	if err != nil {
		log.Fatalf("Failed to insert into migrations table: %v", err)
	}

	log.Println("Events table created")
}

func (*EventTable) Down() {
	log.Println("Dropping events table")

	// Initialize the service
	db, err := core.NewPostgresService()
	if err != nil {
		log.Printf("Failed to initialize database service: %v", err)
		return
	}
	defer db.Close()

	// Drop the table
	_, err = db.Delete(`DROP TABLE IF EXISTS events;`)
	if err != nil {
		log.Fatalf("Failed to drop events table: %v", err)
	}

	// Delete from migrations table
	_, err = db.Delete(`DELETE FROM migrations WHERE name = $1;`, "events")
	if err != nil {
		log.Fatalf("Failed to delete from migrations table: %v", err)
	}

	log.Println("Events table dropped")
}
