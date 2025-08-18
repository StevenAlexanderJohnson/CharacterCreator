package repositories

import (
	"database/sql"
	"dndcc/internal/models"
	"errors"
	"fmt"
)

type CharacterRepository struct {
	db *sql.DB
}

func NewCharacterRepository(db *sql.DB) *CharacterRepository {
	return &CharacterRepository{db}
}

func (r *CharacterRepository) getCharacterProficiencies(tx *sql.Tx, characterId int) ([]string, error) {
	profQuery := `SELECT proficiency FROM character_proficiencies WHERE character_id = ?`

	result, err := tx.Query(profQuery, characterId)
	if err != nil {
		return nil, fmt.Errorf("failed to get character proficiencies: %w", err)
	}
	defer result.Close()

	var proficiencies []string
	for result.Next() {
		var name string
		if err := result.Scan(&name); err != nil {
			return nil, err
		}
		proficiencies = append(proficiencies, name)
	}

	if err := result.Err(); err != nil {
		return nil, err
	}

	return proficiencies, nil
}

func (r *CharacterRepository) Create(data *models.Character) (*models.Character, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	var subraceType sql.NullString
	if data.SubraceType.Valid {
		subraceType = data.SubraceType
	}

	charQuery := `
		INSERT INTO characters (
			owner_id, name, bio, background, class, level, race_type, subrace_type, race_move_speed,
			strength, dexterity, constitution, intelligence, wisdom, charisma, current_health_points
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	result, err := tx.Exec(
		charQuery,
		data.OwnerId, data.Name, data.Bio, data.Background, data.Class, data.Level, data.RaceType,
		subraceType, data.RaceMoveSpeed, data.Strength, data.Dexterity, data.Constitution,
		data.Intelligence, data.Wisdom, data.Charisma, data.CurrentHealthPoints,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to insert character: %w", err)
	}

	charID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert ID for character: %w", err)
	}
	data.ID = int(charID)

	if len(data.BackgroundProficiencies) > 0 {
		profInsertStmt, err := tx.Prepare("INSERT INTO character_proficiencies (character_id, proficiency) VALUES (?, ?);")
		if err != nil {
			return nil, fmt.Errorf("failed to prepare proficiencies insert statement: %w", err)
		}
		defer profInsertStmt.Close()

		for _, profName := range data.BackgroundProficiencies {
			if _, err := profInsertStmt.Exec(charID, profName); err != nil {
				return nil, fmt.Errorf("failed to insert character proficiency for character %d, proficiency %s: %w", charID, profName, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit character creation transaction: %w", err)
	}

	return r.Get(data.ID, data.OwnerId)
}

func (r *CharacterRepository) Get(id int, ownerId int) (*models.Character, error) {
	var character models.Character

	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	charQuery := `
		SELECT
			id, owner_id, name, bio, background, class, level, race_type, subrace_type, race_move_speed,
			strength, dexterity, constitution, intelligence, wisdom, charisma, current_health_points
		FROM characters WHERE id = ? AND owner_id = ?;
	`
	row := tx.QueryRow(charQuery, id, ownerId)
	err = row.Scan(
		&character.ID, &character.OwnerId, &character.Name, &character.Bio, &character.Background, &character.Class,
		&character.Level, &character.RaceType, &character.SubraceType, &character.RaceMoveSpeed,
		&character.Strength, &character.Dexterity, &character.Constitution, &character.Intelligence,
		&character.Wisdom, &character.Charisma, &character.CurrentHealthPoints,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("character with ID %d for owner %d not found", id, ownerId)
		}
		return nil, fmt.Errorf("failed to get character by ID %d: %w", id, err)
	}

	proficiencies, err := r.getCharacterProficiencies(tx, character.ID)
	if err != nil {
		return nil, fmt.Errorf("error getting proficiency for %d: %w", character.ID, err)
	}
	character.BackgroundProficiencies = proficiencies

	return &character, nil
}

func (r *CharacterRepository) GetAll(ownerId int) ([]models.Character, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("error beginning transaction: %w", err)
	}

	query := `
		SELECT
			c.id, c.owner_id, c.name, c.bio, c.background, c.class, c.level, c.race_type, c.subrace_type, c.race_move_speed,
			c.strength, c.dexterity, c.constitution, c.intelligence, c.wisdom, c.charisma, c.current_health_points
		FROM characters c
		WHERE c.owner_id = ?
		ORDER BY c.id;
	`
	rows, err := r.db.Query(query, ownerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all characters for owner %d: %w", ownerId, err)
	}
	defer rows.Close()

	var allCharacters []models.Character

	for rows.Next() {
		var char models.Character

		err := rows.Scan(
			&char.ID, &char.OwnerId, &char.Name, &char.Bio, &char.Background, &char.Class,
			&char.Level, &char.RaceType, &char.SubraceType, &char.RaceMoveSpeed,
			&char.Strength, &char.Dexterity, &char.Constitution, &char.Intelligence,
			&char.Wisdom, &char.Charisma, &char.CurrentHealthPoints,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan character row for owner %d: %w", ownerId, err)
		}

		proficiencies, err := r.getCharacterProficiencies(tx, char.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get character proficiencies: %w", err)
		}
		char.BackgroundProficiencies = proficiencies

		allCharacters = append(allCharacters, char)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during all characters rows iteration for owner %d: %w", ownerId, err)
	}

	return allCharacters, nil
}

func (r *CharacterRepository) Update(data *models.Character, id, ownerId int) (*models.Character, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction for update: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM characters WHERE id = ? AND owner_id = ?)", id, ownerId).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("failed to check character existence: %w", err)
	}
	if !exists {
		return nil, fmt.Errorf("character with ID %d for owner %d not found for update", id, ownerId)
	}

	var subraceType sql.NullString
	if data.SubraceType.Valid {
		subraceType = data.SubraceType
	}

	charUpdateQuery := `
		UPDATE characters SET
			name = ?, bio = ?, background = ?, class = ?, level = ?, race_type = ?, subrace_type = ?, race_move_speed = ?,
			strength = ?, dexterity = ?, constitution = ?, intelligence = ?, wisdom = ?, charisma = ?, current_health_points = ?
		WHERE id = ? AND owner_id = ?;
	`
	_, err = tx.Exec(
		charUpdateQuery,
		data.Name, data.Bio, data.Background, data.Class, data.Level, data.RaceType,
		subraceType, data.RaceMoveSpeed, data.Strength, data.Dexterity, data.Constitution,
		data.Intelligence, data.Wisdom, data.Charisma, data.CurrentHealthPoints,
		id, ownerId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update character ID %d: %w", id, err)
	}

	_, err = tx.Exec("DELETE FROM character_proficiencies WHERE character_id = ?;", id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete existing proficiencies for character ID %d: %w", id, err)
	}

	if len(data.BackgroundProficiencies) > 0 {
		profInsertStmt, err := tx.Prepare("INSERT INTO character_proficiencies (character_id, proficiency) VALUES (?, ?);")
		if err != nil {
			return nil, fmt.Errorf("failed to prepare proficiencies insert statement for update: %w", err)
		}
		defer profInsertStmt.Close()

		for _, profName := range data.BackgroundProficiencies {
			if _, err := profInsertStmt.Exec(id, profName); err != nil {
				return nil, fmt.Errorf("failed to insert character proficiency for character %d, proficiency %s during update: %w", id, profName, err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit character update transaction: %w", err)
	}

	return r.Get(id, ownerId)
}

func (r *CharacterRepository) Delete(id, ownerId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for delete: %w", err)
	}
	defer tx.Rollback()

	result, err := tx.Exec("DELETE FROM characters WHERE id = ? AND owner_id = ?;", id, ownerId)
	if err != nil {
		return fmt.Errorf("failed to delete character ID %d for owner %d: %w", id, ownerId, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected for character deletion ID %d: %w", id, err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no character found with ID %d for owner %d to delete", id, ownerId)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit character deletion transaction: %w", err)
	}

	return nil
}
