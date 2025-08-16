package repositories

import (
	"database/sql"
	"dndcc/internal/models"
	"errors"
	"fmt"
)

type AuthRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (r *AuthRepository) Create(data *models.Auth) (*models.Auth, error) {
	query := `INSERT INTO auth (username, hashed_password) VALUES (?, ?)`

	result, err := r.db.Exec(query, data.Username, data.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth record: %v", err)
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get the last inserted ID after creating auth")
	}

	createdAuth, err := r.GetId(int(lastId))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve newly created auth record")
	}
	return createdAuth, nil
}

func (r *AuthRepository) GetId(id int) (*models.Auth, error) {
	query := `SELECT id, username, hashed_password, created_at, updated_at FROM auth WHERE id = ?;`

	var auth models.Auth
	err := r.db.QueryRow(query, id).Scan(&auth.ID, &auth.Username, &auth.HashedPassword, &auth.CreatedAt, &auth.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("auth record with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get auth record by id %d: %w", id, err)
	}

	return &auth, nil
}

func (r *AuthRepository) Get(username string) (*models.Auth, error) {
	query := `SELECT id, username, hashed_password, created_at, updated_at FROM auth WHERE username = ?;`

	var auth models.Auth
	err := r.db.QueryRow(query, username).Scan(&auth.ID, &auth.Username, &auth.HashedPassword, &auth.CreatedAt, &auth.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("auth record with username %s not found", username)
		}
		return nil, fmt.Errorf("failed to get auth record by username %s: %w", username, err)
	}

	return &auth, nil
}
