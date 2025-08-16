package repositories

import (
	"database/sql"
	"dndcc/internal/models"
	"errors"
	"fmt"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{
		db: db,
	}
}

func (r *SessionRepository) Create(data *models.Session) (*models.Session, error) {
	query := `
		INSERT INTO sessions (user_id, token, expires_at, ip_address, user_agent)
		VALUES (?, ?, ?, ?, ?);
	`
	result, err := r.db.Exec(
		query,
		data.UserId,
		data.Token,
		data.ExpiresAt,
		data.IpAddress,
		data.UserAgent,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID for session: %w", err)
	}

	// Fetch the newly created record to get all fields, including the ID and timestamps
	createdSession, err := r.Get(int(lastID))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve newly created session (ID: %d): %w", lastID, err)
	}

	return createdSession, nil
}

// Get is a helper method to retrieve a session by its ID.
func (r *SessionRepository) Get(id int) (*models.Session, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, last_activity_at, ip_address, user_agent
		FROM sessions WHERE id = ?;
	`
	var session models.Session

	// Use QueryRow and Scan for a single row result
	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&session.ID,
		&session.UserId,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.IpAddress,
		&session.UserAgent,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to get session by ID %d: %w", id, err)
	}
	return &session, nil
}

// GetByToken retrieves a session by its token.
func (r *SessionRepository) GetByToken(token string) (*models.Session, error) {
	// The username parameter is currently unused as the token is unique
	query := `
		SELECT id, user_id, token, expires_at, created_at, last_activity_at, ip_address, user_agent
		FROM sessions WHERE token = ?;
	`
	var session models.Session
	row := r.db.QueryRow(query, token)
	err := row.Scan(
		&session.ID,
		&session.UserId,
		&session.Token,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.LastActivityAt,
		&session.IpAddress,
		&session.UserAgent,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("session with token not found")
		}
		return nil, fmt.Errorf("failed to get session by token: %w", err)
	}
	return &session, nil
}

// GetAllUserSessions retrieves all active sessions for a given user ID.
func (r *SessionRepository) GetAllUserSessions(userId int) ([]models.Session, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, last_activity_at, ip_address, user_agent
		FROM sessions WHERE user_id = ? ORDER BY last_activity_at DESC;
	`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all sessions for user ID %d: %w", userId, err)
	}
	defer rows.Close()

	var sessions []models.Session
	for rows.Next() {
		var session models.Session
		err := rows.Scan(
			&session.ID,
			&session.UserId,
			&session.Token,
			&session.ExpiresAt,
			&session.CreatedAt,
			&session.LastActivityAt,
			&session.IpAddress,
			&session.UserAgent,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan session row: %w", err)
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return sessions, nil
}

// Update updates an existing session record in the database.
func (r *SessionRepository) Update(data *models.Session, userId int) (*models.Session, error) {
	query := `
		UPDATE sessions
		SET expires_at = ?, last_activity_at = CURRENT_TIMESTAMP, ip_address = ?, user_agent = ?
		WHERE id = ? AND user_id = ?;
	`
	_, err := r.db.Exec(
		query,
		data.ExpiresAt,
		data.IpAddress,
		data.UserAgent,
		data.ID,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update session ID %d: %w", data.ID, err)
	}

	// Retrieve the updated session to get the latest timestamps
	updatedSession, err := r.Get(data.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated session (ID: %d): %w", data.ID, err)
	}

	return updatedSession, nil
}

// Delete removes a session record by its ID.
func (r *SessionRepository) Delete(id int) error {
	query := `DELETE FROM sessions WHERE id = ?;`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete session ID %d: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for session deletion ID %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no session found with ID %d to delete", id)
	}

	return nil
}
