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
		character.CurrentHealthPoints = character.GetMaxHealthPoints()
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

func (c *Character) GetSavingThrow(stat StatName) int {
	savingThrow := c.StatBlock.GetAbilityScore(stat)
	proficiencies := c.Class.GetSavingThrowsProficiencies()
	for _, prof := range proficiencies {
		if prof == stat {
			savingThrow += c.GetProficiencyBonus()
			break
		}
	}
	return savingThrow
}

func (c *Character) GetSkill(skill SkillName) int {
	var bonus int
	switch skill {
	case SkillAthletics:
		bonus = c.StatBlock.GetAbilityScore(StatStrength)
	case SkillAcrobatics, SkillSleightOfHand, SkillStealth:
		bonus = c.StatBlock.GetAbilityScore(StatDexterity)
	case SkillArcana, SkillHistory, SkillInvestigation, SkillNature, SkillReligion:
		bonus = c.StatBlock.GetAbilityScore(StatIntelligence)
	case SkillAnimalHandling, SkillInsight, SkillMedicine, SkillPerception, SkillSurvival:
		bonus = c.StatBlock.GetAbilityScore(StatWisdom)
	case SkillDeception, SkillIntimidation, SkillPerformance, SkillPersuasion:
		bonus = c.StatBlock.GetAbilityScore(StatCharisma)
	default:
		return 0
	}

	for _, proficiency := range c.Background.Proficiencies {
		if proficiency == skill {
			bonus += c.GetProficiencyBonus()
			break
		}
	}

	return bonus
}

func (c *Character) GetMaxHealthPoints() int {
	hitDie := c.Class.GetHitDie()
	constitution := c.GetAbilityScore(StatConstitution)
	return int(hitDie)*c.Level + constitution
}

func (c *Character) GetArmorClass() int {
	dex := c.GetAbilityScore(StatDexterity)
	return 10 + dex
}

func (c *Character) GetInitiative() int {
	dex := c.GetAbilityScore(StatDexterity)
	return dex
}

func (c *Character) GetMoveSpeed() int {
	return c.Race.GetMoveSpeed()
}
