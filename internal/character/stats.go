package character

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"gopkg.in/yaml.v3"
)

var (
	ErrUndefinedStat = errors.New("attempted to use undefined stat")
)

type StatName string

const (
	StatStrength     StatName = "Strength"
	StatDexterity    StatName = "Dexterity"
	StatConstitution StatName = "Constitution"
	StatIntelligence StatName = "Intelligence"
	StatWisdom       StatName = "Wisdom"
	StatCharisma     StatName = "Charisma"
	StatYourChoice   StatName = "YourChoice" // Should be edited in the YAML file
)

func (c StatName) IsValid() bool {
	switch c {
	case StatStrength, StatDexterity, StatConstitution, StatIntelligence, StatWisdom, StatCharisma, StatYourChoice:
		return true
	default:
		return false
	}
}

func (s *StatName) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	stat := StatName(str)
	if !stat.IsValid() {
		return fmt.Errorf("error occurred while parsing stat from json: %v", ErrUndefinedStat)
	}
	*s = stat
	return nil
}

func (s *StatName) UnmarshalYAML(value *yaml.Node) error {
	var str string
	if err := value.Decode(&str); err != nil {
		return err
	}
	stat := StatName(str)
	if !stat.IsValid() {
		return fmt.Errorf("error occurred while parsing stat from yaml: %v", ErrUndefinedStat)
	}
	*s = stat
	return nil
}

type StatBlock struct {
	Strength     int `yaml:"strength"`
	Dexterity    int `yaml:"dexterity"`
	Constitution int `yaml:"constitution"`
	Intelligence int `yaml:"intelligence"`
	Wisdom       int `yaml:"wisdom"`
	Charisma     int `yaml:"charisma"`
}

func (s *StatBlock) GetAbilityScore(stat StatName) int {
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
		return 0
	}

	abilityScore := float64(statValue-10) / float64(2)
	return int(math.Floor(abilityScore))
}

func (s *StatBlock) SetStat(stat StatName, value int) {
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
		return
	}
}
