package page

import (
	"dndcc/internal/character"
	"dndcc/internal/models"
)

type CharacterEditPageData struct {
	Method            string
	Action            string
	Error             string
	Character         *models.Character
	BackgroundOptions []character.BackgroundName
	ClassOptions      []character.ClassName
	RaceOptions       []character.RaceName
	SubraceOptions    []character.SubraceName
}

func NewCharacterEditPageData(method, action, errorMessage string, characterModel *models.Character) *CharacterEditPageData {
	return &CharacterEditPageData{
		Method:    method,
		Action:    action,
		Error:     errorMessage,
		Character: characterModel,
		BackgroundOptions: []character.BackgroundName{
			character.BackgroundAcolyte,
			character.BackgroundCharlatan,
			character.BackgroundCriminal,
			character.BackgroundEntertainer,
			character.BackgroundFolkHero,
			character.BackgroundGuildArtisan,
			character.BackgroundHermit,
			character.BackgroundNoble,
			character.BackgroundOutlander,
			character.BackgroundSage,
			character.BackgroundSailor,
			character.BackgroundSoldier,
			character.BackgroundUrchin,
		},
		ClassOptions: []character.ClassName{
			character.ClassBarbarian,
			character.ClassBard,
			character.ClassCleric,
			character.ClassDruid,
			character.ClassFighter,
			character.ClassMonk,
			character.ClassPaladin,
			character.ClassRanger,
			character.ClassRogue,
			character.ClassSorcerer,
			character.ClassWarlock,
			character.ClassWizard,
			character.ClassCommoner,
		},
		RaceOptions: []character.RaceName{
			character.RaceDwarf,
			character.RaceElf,
			character.RaceHalfling,
			character.RaceHuman,
			character.RaceDragonborn,
			character.RaceGnome,
			character.RaceHalfElf,
			character.RaceHalfOrc,
			character.RaceTiefling,
		},
		SubraceOptions: []character.SubraceName{
			character.SubraceNone,
			character.SubraceHillDwarf,
			character.SubraceMountainDwarf,
			character.SubraceHighElf,
			character.SubraceWoodElf,
			character.SubraceDrow,
			character.SubraceLightfoot,
			character.SubraceStout,
			character.SubraceForestGnome,
			character.SubraceRockGnome,
		},
	}
}
