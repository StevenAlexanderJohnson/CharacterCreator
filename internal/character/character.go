package character

import (
	"math"

	"gopkg.in/yaml.v3"
)

type Character struct {
	*StatBlock          `yaml:"stats"`
	Class               ClassName  `yaml:"class"`
	Race                Race       `yaml:"race"`
	Name                string     `yaml:"name"`
	Level               int        `yaml:"level"`
	Background          Background `yaml:"background"`
	Bio                 string     `yaml:"bio"`
	CurrentHealthPoints int        `yaml:"current_hit_points"`
}

func NewCharacter() *Character {
	return &Character{
		StatBlock: &StatBlock{},
		Class:     ClassBarbarian,
		Race: Race{
			Type:         RaceHuman,
			Subrace:      SubraceNone,
			MoveSpeed:    0,
			StatIncrease: []StatIncrease{},
		},
		Name:  "",
		Level: 1,
		Background: Background{
			Name:          BackgroundAcolyte,
			Proficiencies: []SkillName{},
		},
		Bio:                 "",
		CurrentHealthPoints: 0,
	}
}

func CharacterFromYaml(data []byte) (*Character, error) {
	character := Character{}
	if err := yaml.Unmarshal(data, &character); err != nil {
		return nil, err
	}
	if character.CurrentHealthPoints == 0 {
		hp, _ := character.GetMaxHealthPoints()
		character.CurrentHealthPoints = hp
	}
	return &character, nil
}

func (c *Character) SetLevel(level int) *Character {
	c.Level = level
	return c
}

func (c *Character) SetClass(class ClassName) *Character {
	c.Class = class
	return c
}

func (c *Character) AddProficiency(skill SkillName) *Character {
	c.Background.Proficiencies = append(c.Background.Proficiencies, skill)
	return c
}

func (c *Character) GetProficiencyBonus() int {
	return int(math.Floor(float64(c.Level-1)/float64(4))) + 2
}

func (c *Character) GetSavingThrow(stat StatName) (int, error) {
	savingThrow, err := c.StatBlock.GetAbilityScore(stat)
	if err != nil {
		return 0, err
	}
	proficiencies, err := c.Class.GetSavingThrowsProficiencies()
	if err != nil {
		return 0, err
	}
	for _, prof := range proficiencies {
		if prof == stat {
			savingThrow += c.GetProficiencyBonus()
			break
		}
	}
	return savingThrow, nil
}

func (c *Character) GetSkill(skill SkillName) (int, error) {
	var bonus int
	var err error
	switch skill {
	case SkillAthletics:
		bonus, err = c.StatBlock.GetAbilityScore(StatStrength)
	case SkillAcrobatics, SkillSleightOfHand, SkillStealth:
		bonus, err = c.StatBlock.GetAbilityScore(StatDexterity)
	case SkillArcana, SkillHistory, SkillInvestigation, SkillNature, SkillReligion:
		bonus, err = c.StatBlock.GetAbilityScore(StatIntelligence)
	case SkillAnimalHandling, SkillInsight, SkillMedicine, SkillPerception, SkillSurvival:
		bonus, err = c.StatBlock.GetAbilityScore(StatWisdom)
	case SkillDeception, SkillIntimidation, SkillPerformance, SkillPersuasion:
		bonus, err = c.StatBlock.GetAbilityScore(StatCharisma)
	default:
		return 0, ErrUndefinedSkill
	}
	if err != nil {
		return 0, err
	}

	for _, proficiency := range c.Background.Proficiencies {
		if proficiency == skill {
			bonus += c.GetProficiencyBonus()
			break
		}
	}

	return bonus, nil
}

func (c *Character) GetMaxHealthPoints() (int, error) {
	hitDie, err := c.Class.GetHitDie()
	if err != nil {
		return 0, err
	}
	// Discarding err because we are using developer defined stat
	constitution, _ := c.GetAbilityScore(StatConstitution)
	return int(hitDie)*c.Level + constitution, nil
}

func (c *Character) GetArmorClass() int {
	// Discarding err because we are using developer defined stat
	dex, _ := c.GetAbilityScore(StatDexterity)
	return 10 + dex
}

func (c *Character) GetInitiative() int {
	// Discarding err because we are using developer defined stat
	dex, _ := c.GetAbilityScore(StatDexterity)
	return dex
}

func (c *Character) GetMoveSpeed() int {
	return c.Race.GetMoveSpeed()
}
