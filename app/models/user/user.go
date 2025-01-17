package user

import (
	"database/sql"
	"errors"
	"fmt"
	"web-app/app/services/core"
)

type User struct {
	ID        int64
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	CreatedAt string
	db        *core.PostgresService
}

func NewUserModel() *User {
	db, _ := core.NewPostgresService()
	return &User{
		db: db,
	}
}

// Create implements the Model interface Create method
func (u *User) Create() error {
	query := `
        INSERT INTO users (username, password)
        VALUES ($1, $2)
        RETURNING id`

	result, err := u.db.Create(query, u.Username, u.Password)
	if err != nil {
		return fmt.Errorf("error creating user: %w", err)
	}

	lastInsertId, err := result.LastInsertId()
	if err == nil {
		u.ID = lastInsertId
	}

	return nil
}

// Find implements the Model interface Find method
func (u *User) Find() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	query := `
        SELECT id, username, password, created_at 
        FROM users 
        WHERE username = $1`

	rows, err := u.db.Read(query, u.Username)
	if err != nil {
		return fmt.Errorf("error finding user: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			return fmt.Errorf("error scanning user: %w", err)
		}
		return nil
	}

	return sql.ErrNoRows
}

// FindByUsername implements the Model interface FindByUsername method
func (u *User) FindByUsername() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	query := `
		SELECT id, username, password, created_at
		FROM users
		WHERE username = $1`

	rows, err := u.db.Read(query, u.Username)
	if err != nil {
		return fmt.Errorf("error finding user by username: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.CreatedAt)
		if err != nil {
			return fmt.Errorf("error scanning user: %w", err)
		}
		return nil
	}

	return sql.ErrNoRows
}

// Update implements the Model interface Update method
func (u *User) Update() error {
	if u.ID == 0 {
		return errors.New("id is required")
	}

	query := `
        UPDATE users 
        SET username = $1, password = $2
        WHERE id = $3`

	_, err := u.db.Update(query, u.Username, u.Password, u.ID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}

	return nil
}

// Delete implements the Model interface Delete method
func (u *User) Delete() error {
	if u.ID == 0 {
		return errors.New("id is required")
	}

	query := `
        DELETE FROM users 
        WHERE id = $1`

	_, err := u.db.Delete(query, u.ID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	return nil
}

func (u *User) Paginate(limit, page int) ([]User, error) {
	// Set default values
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	query := `
        SELECT id, username, password, created_at
        FROM users
        ORDER BY id DESC
        LIMIT $1 OFFSET $2`

	rows, err := u.db.Read(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error getting users: %w", err)
	}
	defer rows.Close()

	users := make([]User, 0, limit) // Pre-allocate slice with capacity
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
		users = append(users, user)
	}

	// Check for errors from iterating over rows
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over users: %w", err)
	}

	return users, nil
}
