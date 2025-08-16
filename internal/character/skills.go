package character

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type SkillName string

var (
	ErrUndefinedSkill = errors.New("attempted to use undefined skill")
)

const (
	// These constants represent all the standard D&D 5e skills.
	SkillAcrobatics     SkillName = "Acrobatics"
	SkillAnimalHandling           = "Animal Handling"
	SkillArcana                   = "Arcana"
	SkillAthletics                = "Athletics"
	SkillDeception                = "Deception"
	SkillHistory                  = "History"
	SkillInsight                  = "Insight"
	SkillIntimidation             = "Intimidation"
	SkillInvestigation            = "Investigation"
	SkillMedicine                 = "Medicine"
	SkillNature                   = "Nature"
	SkillPerception               = "Perception"
	SkillPerformance              = "Performance"
	SkillPersuasion               = "Persuasion"
	SkillReligion                 = "Religion"
	SkillSleightOfHand            = "Sleight of Hand"
	SkillStealth                  = "Stealth"
	SkillSurvival       SkillName = "Survival"
)

func (s *SkillName) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}
	switch SkillName(str) {
	case SkillAcrobatics, SkillAnimalHandling, SkillArcana, SkillAthletics, SkillDeception,
		SkillHistory, SkillInsight, SkillIntimidation, SkillInvestigation, SkillMedicine,
		SkillNature, SkillPerception, SkillPerformance, SkillPersuasion, SkillReligion,
		SkillSleightOfHand, SkillStealth, SkillSurvival:
		*s = SkillName(str)
		return nil
	default:
		return ErrUndefinedSkill
	}
}
