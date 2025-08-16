package character

import (
	"errors"
	"fmt"
	"math"

	"gopkg.in/yaml.v3"
)

var (
	ErrUndefinedStat = errors.New("attempted to use undefined stat")
)

type StatName string

func (s *StatName) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}
	switch StatName(str) {
	case StatStrength, StatDexterity, StatConstitution, StatIntelligence, StatWisdom, StatCharisma, StatYourChoice:
		*s = StatName(str)
		return nil
	default:
		return fmt.Errorf("error occurred while parsing stat from yaml: %v", ErrUndefinedStat)
	}
}

const (
	StatStrength     StatName = "Strength"
	StatDexterity    StatName = "Dexterity"
	StatConstitution StatName = "Constitution"
	StatIntelligence StatName = "Intelligence"
	StatWisdom       StatName = "Wisdom"
	StatCharisma     StatName = "Charisma"
	StatYourChoice   StatName = "YourChoice" // Should be edited in the YAML file
)

type StatBlock struct {
	Strength     int `yaml:"strength"`
	Dexterity    int `yaml:"dexterity"`
	Constitution int `yaml:"constitution"`
	Intelligence int `yaml:"intelligence"`
	Wisdom       int `yaml:"wisdom"`
	Charisma     int `yaml:"charisma"`
}

func (s *StatBlock) GetAbilityScore(stat StatName) (int, error) {
	var statValue int
	switch stat {
	case StatStrength:
		statValue = s.Strength
	case StatDexterity:
		statValue = s.Dexterity
	case StatConstitution:
		statValue = s.Constitution
	case StatIntelligence:
		statValue = s.Intelligence
	case StatWisdom:
		statValue = s.Wisdom
	case StatCharisma:
		statValue = s.Charisma
	case StatYourChoice:
	default:
		return 0, ErrUndefinedStat
	}

	abilityScore := float64(statValue-10) / float64(2)
	return int(math.Floor(abilityScore)), nil
}

func (s *StatBlock) SetStat(stat StatName, value int) error {
	switch stat {
	case StatStrength:
		s.Strength = value
	case StatDexterity:
		s.Dexterity = value
	case StatConstitution:
		s.Constitution = value
	case StatIntelligence:
		s.Intelligence = value
	case StatWisdom:
		s.Wisdom = value
	case StatCharisma:
		s.Charisma = value
	default:
		return ErrUndefinedStat
	}

	return nil
}
