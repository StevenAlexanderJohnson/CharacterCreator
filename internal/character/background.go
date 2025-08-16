package character

import (
	"errors"
)

var (
	ErrUndefinedBackground = errors.New("attempted to use undefined background")
)

type BackgroundName string

const (
	BackgroundAcolyte      BackgroundName = "Acolyte"
	BackgroundCharlatan    BackgroundName = "Charlatan"
	BackgroundCriminal     BackgroundName = "Criminal"
	BackgroundEntertainer  BackgroundName = "Entertainer"
	BackgroundFolkHero     BackgroundName = "Folk Hero"
	BackgroundGuildArtisan BackgroundName = "Guild Artisan"
	BackgroundHermit       BackgroundName = "Hermit"
	BackgroundNoble        BackgroundName = "Noble"
	BackgroundOutlander    BackgroundName = "Outlander"
	BackgroundSage         BackgroundName = "Sage"
	BackgroundSailor       BackgroundName = "Sailor"
	BackgroundSoldier      BackgroundName = "Soldier"
	BackgroundUrchin       BackgroundName = "Urchin"
)

type Background struct {
	Name          BackgroundName `yaml:"name"`
	Proficiencies []SkillName    `yaml:"proficiencies"`
}

func (b *Background) GetProficiencies() []SkillName {
	output := b.Name.getProficiencies()
	if output == nil && b.Proficiencies != nil {
		return b.Proficiencies
	}
	if output == nil {
		output = []SkillName{}
	}
	return output
}

func (b BackgroundName) getProficiencies() []SkillName {
	switch b {
	case BackgroundAcolyte:
		return []SkillName{SkillInsight, SkillReligion}
	case BackgroundCharlatan:
		return []SkillName{SkillDeception, SkillSleightOfHand}
	case BackgroundCriminal:
		return []SkillName{SkillDeception, SkillStealth}
	case BackgroundEntertainer:
		return []SkillName{SkillAcrobatics, SkillPerformance}
	case BackgroundFolkHero:
		return []SkillName{SkillAnimalHandling, SkillSurvival}
	case BackgroundGuildArtisan:
		return []SkillName{SkillInsight, SkillPersuasion}
	case BackgroundHermit:
		return []SkillName{SkillMedicine, SkillReligion}
	case BackgroundNoble:
		return []SkillName{SkillHistory, SkillPersuasion}
	case BackgroundOutlander:
		return []SkillName{SkillAthletics, SkillSurvival}
	case BackgroundSage:
		return []SkillName{SkillArcana, SkillHistory}
	case BackgroundSailor:
		return []SkillName{SkillAthletics, SkillPerception}
	case BackgroundSoldier:
		return []SkillName{SkillAthletics, SkillIntimidation}
	case BackgroundUrchin:
		return []SkillName{SkillSleightOfHand, SkillStealth}
	default:
		return nil
	}
}
