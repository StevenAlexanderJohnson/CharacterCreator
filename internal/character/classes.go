package character

import (
	"errors"
	"fmt"
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

func (c ClassName) GetSavingThrowsProficiencies() ([]StatName, error) {
	switch c {
	case ClassBarbarian:
		return []StatName{StatStrength, StatConstitution}, nil
	case ClassBard:
		return []StatName{StatDexterity, StatCharisma}, nil
	case ClassCleric:
		return []StatName{StatWisdom, StatCharisma}, nil
	case ClassDruid:
		return []StatName{StatIntelligence, StatWisdom}, nil
	case ClassFighter:
		return []StatName{StatStrength, StatConstitution}, nil
	case ClassMonk:
		return []StatName{StatStrength, StatDexterity}, nil
	case ClassPaladin:
		return []StatName{StatWisdom, StatCharisma}, nil
	case ClassRanger:
		return []StatName{StatStrength, StatDexterity}, nil
	case ClassRogue:
		return []StatName{StatDexterity, StatIntelligence}, nil
	case ClassSorcerer:
		return []StatName{StatConstitution, StatCharisma}, nil
	case ClassWarlock:
		return []StatName{StatWisdom, StatCharisma}, nil
	case ClassWizard:
		return []StatName{StatIntelligence, StatWisdom}, nil
	case ClassCommoner:
		return []StatName{}, nil
	default:
		return nil, ErrUndefinedClass
	}
}

func (c ClassName) GetHitDie() (HitDie, error) {
	switch c {
	case ClassBarbarian:
		return HitDieD12, nil
	case ClassFighter, ClassPaladin, ClassRanger:
		return HitDieD10, nil
	case ClassBard, ClassCleric, ClassDruid, ClassMonk, ClassRogue, ClassWarlock:
		return HitDieD8, nil
	case ClassSorcerer, ClassWizard, ClassCommoner:
		return HitDieD6, nil
	default:
		panic(fmt.Sprintf("no hit die defined for: %s", c))
	}
}
