package page

import "dndcc/internal/character"

type CharacterViewPageData struct {
	ID int
	*character.Character
}

func NewCharacterViewPageData(id int, char *character.Character) *CharacterViewPageData {
	return &CharacterViewPageData{
		ID:        id,
		Character: char,
	}
}
