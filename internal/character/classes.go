package character

import (
	"encoding/json"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

type HitDie int

var (
	ErrUndefinedClass = errors.New("requested information from undefined class")
)

const (
	HitDieD12 HitDie = 12
	HitDieD10 HitDie = 10
	HitDieD8  HitDie = 8
	HitDieD6  HitDie = 6
)

type ClassName string

const (
	ClassBarbarian ClassName = "Barbarian"
	ClassBard      ClassName = "Bard"
	ClassCleric    ClassName = "Cleric"
	ClassDruid     ClassName = "Druid"
	ClassFighter   ClassName = "Fighter"
	ClassMonk      ClassName = "Monk"
	ClassPaladin   ClassName = "Paladin"
	ClassRanger    ClassName = "Ranger"
	ClassRogue     ClassName = "Rogue"
	ClassSorcerer  ClassName = "Sorcerer"
	ClassWarlock   ClassName = "Warlock"
	ClassWizard    ClassName = "Wizard"
	ClassCommoner  ClassName = "Commoner"
)

func NewClassName(name string) (ClassName, error) {
	class := ClassName(name)
	if !class.IsValid() {
		return "", fmt.Errorf("invalid class name: %s", name)
	}
	return class, nil
}

func (c ClassName) IsValid() bool {
	switch c {
	case ClassBarbarian, ClassBard, ClassCleric, ClassDruid, ClassFighter,
		ClassMonk, ClassPaladin, ClassRanger, ClassRogue, ClassSorcerer,
		ClassWarlock, ClassWizard, ClassCommoner:
		return true
	default:
		return false
	}
}

func (c *ClassName) UnmarshalJSON(data []byte) error {
	var classString string
	if err := json.Unmarshal(data, &classString); err != nil {
		return fmt.Errorf("failed to unmarshal class name: %w", err)
	}
	class := ClassName(classString)
	if !class.IsValid() {
		return fmt.Errorf("invalid class name: %s", classString)
	}
	*c = class
	return nil
}

func (c *ClassName) UnmarshalYAML(value *yaml.Node) error {
	var classString string
	if err := value.Decode(&classString); err != nil {
		return fmt.Errorf("failed to unmarshal class name: %w", err)
	}
	class := ClassName(classString)
	if !class.IsValid() {
		return fmt.Errorf("invalid class name: %s", classString)
	}
	*c = class
	return nil
}

func (c ClassName) GetSavingThrowsProficiencies() []StatName {
	switch c {
	case ClassBarbarian:
		return []StatName{StatStrength, StatConstitution}
	case ClassBard:
		return []StatName{StatDexterity, StatCharisma}
	case ClassCleric:
		return []StatName{StatWisdom, StatCharisma}
	case ClassDruid:
		return []StatName{StatIntelligence, StatWisdom}
	case ClassFighter:
		return []StatName{StatStrength, StatConstitution}
	case ClassMonk:
		return []StatName{StatStrength, StatDexterity}
	case ClassPaladin:
		return []StatName{StatWisdom, StatCharisma}
	case ClassRanger:
		return []StatName{StatStrength, StatDexterity}
	case ClassRogue:
		return []StatName{StatDexterity, StatIntelligence}
	case ClassSorcerer:
		return []StatName{StatConstitution, StatCharisma}
	case ClassWarlock:
		return []StatName{StatWisdom, StatCharisma}
	case ClassWizard:
		return []StatName{StatIntelligence, StatWisdom}
	default:
		return []StatName{}
	}
}

func (c ClassName) GetHitDie() HitDie {
	switch c {
	case ClassBarbarian:
		return HitDieD12
	case ClassFighter, ClassPaladin, ClassRanger:
		return HitDieD10
	case ClassBard, ClassCleric, ClassDruid, ClassMonk, ClassRogue, ClassWarlock:
		return HitDieD8
	case ClassSorcerer, ClassWizard, ClassCommoner:
		return HitDieD6
	default:
		return HitDieD6 // Default to D6 for unknown classes
	}
}
