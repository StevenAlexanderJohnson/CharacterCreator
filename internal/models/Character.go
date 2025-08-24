package models

import (
	"database/sql"
	"dndcc/internal/character"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var (
	ErrInvalidCharacterName       = errors.New("character name cannot be empty")
	ErrInvalidCharacterBackground = errors.New("character background cannot be empty")
	ErrInvalidCharacterClass      = errors.New("character class cannot be empty")
	ErrInvalidCharacterRace       = errors.New("character race cannot be empty")
	ErrInvalidCharacterSubrace    = errors.New("character subrace cannot be empty if provided")
)

type Character struct {
	ID                      int
	OwnerId                 int
	Name                    string
	Bio                     string
	Background              string
	Class                   string
	Level                   int
	RaceType                string
	SubraceType             sql.NullString
	RaceMoveSpeed           int
	Strength                int
	Dexterity               int
	Constitution            int
	Intelligence            int
	Wisdom                  int
	Charisma                int
	CurrentHealthPoints     int
	BackgroundProficiencies []string
}

func (c *Character) Validate() error {
	if strings.TrimSpace(c.Name) == "" {
		return ErrInvalidCharacterName
	}
	if strings.TrimSpace(c.Background) == "" {
		return ErrInvalidCharacterBackground
	}
	if strings.TrimSpace(c.Class) == "" {
		return ErrInvalidCharacterClass
	}
	if strings.TrimSpace(c.RaceType) == "" {
		return ErrInvalidCharacterRace
	}
	if c.SubraceType.Valid && strings.TrimSpace(c.SubraceType.String) == "" {
		return ErrInvalidCharacterSubrace
	}
	return nil
}

func (c *Character) ToCharacterSheet() *character.Character {
	proficiencies := make([]character.SkillName, len(c.BackgroundProficiencies))
	for i := 0; i < len(c.BackgroundProficiencies); i++ {
		proficiencies[i] = character.SkillName(c.BackgroundProficiencies[i])
	}
	return &character.Character{
		StatBlock: &character.StatBlock{
			Strength:     c.Strength,
			Dexterity:    c.Dexterity,
			Constitution: c.Constitution,
			Intelligence: c.Intelligence,
			Wisdom:       c.Wisdom,
			Charisma:     c.Charisma,
		},
		Class: character.ClassName(c.Class),
		Race: character.Race{
			Type:      character.RaceName(c.RaceType),
			Subrace:   character.SubraceName(c.SubraceType.String),
			MoveSpeed: c.RaceMoveSpeed,
		},
		Name:  c.Name,
		Level: c.Level,
		Background: character.Background{
			Name:          character.BackgroundName(c.Background),
			Proficiencies: proficiencies,
		},
		Bio:                 c.Bio,
		CurrentHealthPoints: c.CurrentHealthPoints,
	}
}

func CharacterFromForm(r *http.Request) (*Character, error) {
	name := r.FormValue("Name")
	if name == "" {
		return nil, fmt.Errorf("name is required to create the character")
	}
	bio := r.FormValue("Bio")
	background := r.FormValue("Background")
	level, err := strconv.Atoi(r.FormValue("Level"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for level: %s", r.FormValue("Level"))
	}
	class := r.FormValue("ClassSelect")
	race := r.FormValue("RaceType")
	subrace := r.FormValue("SubraceType")
	moveSpeed, err := strconv.Atoi(r.FormValue("RaceMoveSpeed"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for move speed: %s", r.FormValue("RaceMoveSpeed"))
	}
	strength, err := strconv.Atoi(r.FormValue("Strength"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for strength: %s", r.FormValue("Strength"))
	}
	dexterity, err := strconv.Atoi(r.FormValue("Dexterity"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for dexterity: %s", r.FormValue("Dexterity"))
	}
	constitution, err := strconv.Atoi(r.FormValue("Constitution"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for constitution: %s", r.FormValue("Constitution"))
	}
	intelligence, err := strconv.Atoi(r.FormValue("Intelligence"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for intelligence: %s", r.FormValue("Intelligence"))
	}
	wisdom, err := strconv.Atoi(r.FormValue("Wisdom"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for wisdom: %s", r.FormValue("Wisdom"))
	}
	charisma, err := strconv.Atoi(r.FormValue("Charisma"))
	if err != nil {
		return nil, fmt.Errorf("invalid value was passed for charisma: %s", r.FormValue("Charisma"))
	}

	return &Character{
		Name:                    name,
		Bio:                     bio,
		Background:              background,
		Level:                   level,
		Class:                   class,
		RaceType:                race,
		SubraceType:             sql.NullString{String: subrace, Valid: true},
		RaceMoveSpeed:           moveSpeed,
		Strength:                strength,
		Dexterity:               dexterity,
		Constitution:            constitution,
		Intelligence:            intelligence,
		Wisdom:                  wisdom,
		Charisma:                charisma,
		BackgroundProficiencies: []string{},
	}, nil
}
