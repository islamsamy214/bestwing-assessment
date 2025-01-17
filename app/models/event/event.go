package event

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"web-app/app/services/core"
)

type Event struct {
	ID        int64
	Name      string `json:"name" binding:"required"`
	Date      string `json:"date"`
	CreatedAt string
	UserId    int64 `json:"user_id"`
	db        *core.PostgresService
}

func NewEventModel() *Event {
	db, _ := core.NewPostgresService()
	return &Event{
		db: db,
	}
}

// Create implements the Model interface Create method
func (e *Event) Create() error {
	e.CreatedAt = time.Now().Format(time.RFC3339)
	if e.Date == "" {
		e.Date = e.CreatedAt
	}

	query := `
        INSERT INTO events (name, date, created_at, user_id)
        VALUES ($1, $2, $3, $4)
        RETURNING id`

	result, err := e.db.Create(query, e.Name, e.Date, e.CreatedAt, e.UserId)
	if err != nil {
		return fmt.Errorf("error creating event: %w", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err == nil {
		e.ID = lastInsertId
	}

	return nil
}

// Find implements the Model interface Find method
func (e *Event) Find() error {
	if e.Name == "" {
		return errors.New("name is required")
	}

	query := `
        SELECT id, name, date, created_at, user_id
        FROM events
        WHERE name = $1`

	rows, err := e.db.Read(query, e.Name)
	if err != nil {
		return fmt.Errorf("error finding event: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&e.ID, &e.Name, &e.Date, &e.CreatedAt, &e.UserId)
		if err != nil {
			return fmt.Errorf("error scanning event: %w", err)
		}
		return nil
	}

	return sql.ErrNoRows
}

// Update implements the Model interface Update method
func (e *Event) Update() error {
	if e.ID == 0 {
		return errors.New("id is required")
	}

	query := `
        UPDATE events
        SET name = $1, date = $2, user_id = $3
        WHERE id = $4`

	_, err := e.db.Update(query, e.Name, e.Date, e.UserId, e.ID)
	if err != nil {
		return fmt.Errorf("error updating event: %w", err)
	}

	return nil
}

// Delete implements the Model interface Delete method
func (e *Event) Delete() error {
	if e.ID == 0 {
		return errors.New("id is required")
	}

	query := `
        DELETE FROM events
        WHERE id = $1`

	_, err := e.db.Delete(query, e.ID)
	if err != nil {
		return fmt.Errorf("error deleting event: %w", err)
	}

	return nil
}

// Paginate implements pagination for events
func (e *Event) Paginate(limit, page int) ([]Event, error) {
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
        ORDER BY id DESC
        LIMIT $1 OFFSET $2`

	rows, err := e.db.Read(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting events: %w", err)
	}
	defer rows.Close()

	events := make([]Event, 0, limit) // Pre-allocate slice with capacity
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Date,
			&event.CreatedAt,
			&event.UserId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning event: %w", err)
		}
		events = append(events, event)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over events: %w", err)
	}

	return events, nil
}

// Additional helper method to get events by user ID
func (e *Event) GetByUserId(userId int64, limit, page int) ([]Event, error) {
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
        WHERE user_id = $1
        ORDER BY date DESC
        LIMIT $2 OFFSET $3`

	rows, err := e.db.Read(query, userId, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting user events: %w", err)
	}
	defer rows.Close()

	events := make([]Event, 0, limit)
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Date,
			&event.CreatedAt,
			&event.UserId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning event: %w", err)
		}
		events = append(events, event)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over user events: %w", err)
	}

	return events, nil
}
