package character

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

var (
	ErrUndefinedRace    = errors.New("an undefined race was defined")
	ErrUndefinedSubrace = errors.New("an undefined subrace was defined")
)

type RaceName string
type SubraceName string

type StatIncrease struct {
	Stat   StatName
	Amount int
}

const (
	RaceDwarf      RaceName = "Dwarf"
	RaceElf        RaceName = "Elf"
	RaceHalfling   RaceName = "Halfling"
	RaceHuman      RaceName = "Human"
	RaceDragonborn RaceName = "Dragonborn"
	RaceGnome      RaceName = "Gnome"
	RaceHalfElf    RaceName = "Half-Elf"
	RaceHalfOrc    RaceName = "Half-Orc"
	RaceTiefling   RaceName = "Tiefling"
)

func (r *RaceName) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}
	switch RaceName(str) {
	case RaceDwarf, RaceElf,
		RaceHalfling, RaceHuman,
		RaceDragonborn, RaceGnome,
		RaceHalfElf, RaceHalfOrc, RaceTiefling:
		*r = RaceName(str)
		return nil
	default:
		return fmt.Errorf("an error ocurred while parsing race from yaml: %v -> %s", ErrUndefinedRace, str)
	}
}

func (r RaceName) getMoveSpeed() int {
	switch r {
	case RaceElf, RaceHuman, RaceDragonborn, RaceHalfElf, RaceHalfOrc, RaceTiefling:
		return 30
	case RaceDwarf, RaceHalfling, RaceGnome:
		return 25
	default:
		return -1
	}
}

func (r RaceName) getAbilityIncrease() ([]StatIncrease, error) {
	switch r {
	case RaceDwarf:
		return []StatIncrease{{StatConstitution, 2}}, nil
	case RaceElf, RaceHalfling:
		return []StatIncrease{{StatDexterity, 2}}, nil
	case RaceHuman:
		return []StatIncrease{
			{StatStrength, 1},
			{StatCharisma, 1},
			{StatConstitution, 1},
			{StatDexterity, 1},
			{StatIntelligence, 1},
			{StatWisdom, 1},
		}, nil
	case RaceDragonborn:
		return []StatIncrease{
			{StatStrength, 2},
			{StatCharisma, 1},
		}, nil
	case RaceGnome:
		return []StatIncrease{
			{StatIntelligence, 2},
		}, nil
	case RaceHalfElf:
		return []StatIncrease{
			{StatCharisma, 2},
		}, nil
	case RaceHalfOrc:
		return []StatIncrease{
			{StatStrength, 2},
			{StatConstitution, 1},
		}, nil
	case RaceTiefling:
		return []StatIncrease{
			{StatIntelligence, 1},
			{StatCharisma, 2},
		}, nil
	default:
		return []StatIncrease{{StatStrength, 2}}, ErrUndefinedRace
	}
}

const (
	SubraceNone          SubraceName = "None"
	SubraceHillDwarf     SubraceName = "Hill Dwarf"
	SubraceMountainDwarf SubraceName = "Mountain Dwarf"
	SubraceHighElf       SubraceName = "High Elf"
	SubraceWoodElf       SubraceName = "Wood Elf"
	SubraceDrow          SubraceName = "Drow"
	SubraceLightfoot     SubraceName = "Lightfoot"
	SubraceStout         SubraceName = "Stout"
	SubraceForestGnome   SubraceName = "Forest Gnome"
	SubraceRockGnome     SubraceName = "Rock Gnome"
)

func (s SubraceName) getMoveSpeed() int {
	switch s {
	case SubraceWoodElf:
		return 35
	default:
		return -1
	}
}

func (s SubraceName) getAbilityIncrease() ([]StatIncrease, error) {
	switch s {
	case SubraceHillDwarf, SubraceWoodElf:
		return []StatIncrease{{StatWisdom, 1}}, nil
	case SubraceMountainDwarf:
		return []StatIncrease{{StatStrength, 2}}, nil
	case SubraceHighElf:
		return []StatIncrease{{StatIntelligence, 1}}, nil
	case SubraceDrow, SubraceLightfoot:
		return []StatIncrease{{StatCharisma, 1}}, nil
	case SubraceStout, SubraceRockGnome:
		return []StatIncrease{{StatConstitution, 1}}, nil
	case SubraceForestGnome:
		return []StatIncrease{{StatDexterity, 1}}, nil
	case SubraceNone:
		return nil, nil
	default:
		return nil, ErrUndefinedSubrace
	}
}

type Race struct {
	Type         RaceName       `yaml:"type"`
	Subrace      SubraceName    `yaml:"subrace"`
	MoveSpeed    int            `yaml:"move-speed"`
	StatIncrease []StatIncrease `yaml:"stat-increase"`
}

func (r *Race) GetMoveSpeed() int {
	speed := -1
	if r.Subrace != SubraceNone {
		speed = r.Subrace.getMoveSpeed()
	}
	if speed == -1 {
		speed = r.Type.getMoveSpeed()
	}
	if r.MoveSpeed != 0 || speed == -1 {
		speed = r.MoveSpeed
	}

	return speed
}

func (r *Race) GetAbilityIncrease() ([]StatIncrease, error) {
	increase, err := r.Type.getAbilityIncrease()
	if err != nil {
		return nil, err
	}

	if subraceIncrease, err := r.Subrace.getAbilityIncrease(); err != nil {
		return nil, err
	} else if subraceIncrease != nil {
		increase = subraceIncrease
	}

	return increase, nil
}
