package event

import (
	"database/sql"
	"errors"
	"log"
	"web-app/app/services/core"
)

type Event struct {
	ID        int64  `json:"id"`
	Name      string `form:"name" json:"name" xml:"name" binding:"required"`
	Date      string `form:"date" json:"date" xml:"date" binding:"required"`
	CreatedAt string `json:"created_at"`
	UserId    int64  `json:"user_id"`
}

func NewEventModel() *Event {
	return &Event{}
}

// Create inserts a new event into the database
func (e *Event) Create() error {
	db, _ := core.NewPostgresService()

	query := `
        INSERT INTO events (name, date, user_id)
        VALUES ($1, $2, $3)
        RETURNING id, created_at`

	result, err := db.Create(query, e.Name, e.Date, e.UserId)
	if err != nil {
		log.Printf("error creating event: %v", err)
		return err
	}

	defer db.Close()

	// Get the ID and CreatedAt from the result
	err = result.Scan(&e.ID, &e.CreatedAt)
	if err != nil {
		log.Printf("error scanning event: %v", err)
	}

	return nil
}

// Find retrieves an event by its ID
func (e *Event) Find() error {
	db, _ := core.NewPostgresService()

	if e.ID == 0 {
		return errors.New("id is required")
	}

	query := `
        SELECT id, name, date, created_at, user_id 
        FROM events 
        WHERE id = $1`

	rows, err := db.Read(query, e.ID)
	if err != nil {
		log.Printf("error finding event: %v", err)
		return err
	}
	defer db.Close()

	if rows.Next() {
		err := rows.Scan(&e.ID, &e.Name, &e.Date, &e.CreatedAt, &e.UserId)
		if err != nil {
			log.Printf("error scanning event: %v", err)
			return err
		}
		return nil
	}

	return sql.ErrNoRows
}

// Update modifies an existing event in the database
func (e *Event) Update() error {
	db, _ := core.NewPostgresService()

	if e.ID == 0 {
		return errors.New("id is required")
	}

	defer db.Close()

	query := `
        UPDATE events 
        SET name = $1, date = $2, user_id = $3
        WHERE id = $4`

	_, err := db.Update(query, e.Name, e.Date, e.UserId, e.ID)
	if err != nil {
		log.Printf("error updating event: %v", err)
		return err
	}

	return nil
}

// Delete removes an event from the database
func (e *Event) Delete() error {
	db, _ := core.NewPostgresService()

	if e.ID == 0 {
		return errors.New("id is required")
	}

	defer db.Close()

	query := `
        DELETE FROM events 
        WHERE id = $1`

	_, err := db.Delete(query, e.ID)
	if err != nil {
		log.Printf("error deleting event: %v", err)
		return err
	}

	return nil
}

// Paginate retrieves a paginated list of events
func (e *Event) Paginate(limit, page int) ([]Event, error) {
	db, _ := core.NewPostgresService()

	// Set default values
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	query := `
        SELECT id, name, date, created_at, user_id
        FROM events
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2`

	rows, err := db.Read(query, limit, offset)
	if err != nil {
		log.Printf("error paginating events: %v", err)
		return nil, err
	}
	defer db.Close()

	events := make([]Event, 0, limit)
	for rows.Next() {
		var event Event
		err := rows.Scan(&event.ID, &event.Name, &event.Date, &event.CreatedAt, &event.UserId)
		if err != nil {
			log.Printf("error scanning event: %v", err)
			return nil, err
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		log.Printf("error iterating over rows: %v", err)
		return nil, err
	}

	return events, nil
}
