package character_test

import (
	"dndcc/internal/character"
	"testing"
)

func TestGetAbilityIncrease(t *testing.T) {
	tests := []struct {
		Type      character.RaceName
		Subrace   character.SubraceName
		Expected  []character.StatIncrease
		ExpectErr bool
	}{
		// Dwarf
		{character.RaceDwarf, character.SubraceNone, []character.StatIncrease{{character.StatConstitution, 2}}, false},
		{character.RaceDwarf, character.SubraceHillDwarf, []character.StatIncrease{{character.StatWisdom, 1}}, false},
		{character.RaceDwarf, character.SubraceMountainDwarf, []character.StatIncrease{{character.StatStrength, 2}}, false},
		// Elf
		{character.RaceElf, character.SubraceNone, []character.StatIncrease{{character.StatDexterity, 2}}, false},
		{character.RaceElf, character.SubraceHighElf, []character.StatIncrease{{character.StatIntelligence, 1}}, false},
		{character.RaceElf, character.SubraceWoodElf, []character.StatIncrease{{character.StatWisdom, 1}}, false},
		{character.RaceElf, character.SubraceDrow, []character.StatIncrease{{character.StatCharisma, 1}}, false},
		// Halfling
		{character.RaceHalfling, character.SubraceNone, []character.StatIncrease{{character.StatDexterity, 2}}, false},
		{character.RaceHalfling, character.SubraceLightfoot, []character.StatIncrease{{character.StatCharisma, 1}}, false},
		{character.RaceHalfling, character.SubraceStout, []character.StatIncrease{{character.StatConstitution, 1}}, false},
		// Human
		{character.RaceHuman, character.SubraceNone, []character.StatIncrease{
			{character.StatStrength, 1},
			{character.StatCharisma, 1},
			{character.StatConstitution, 1},
			{character.StatDexterity, 1},
			{character.StatIntelligence, 1},
			{character.StatWisdom, 1},
		}, false},
		// Dragonborn
		{character.RaceDragonborn, character.SubraceNone, []character.StatIncrease{
			{character.StatStrength, 2},
			{character.StatCharisma, 1},
		}, false},
		// Gnome
		{character.RaceGnome, character.SubraceNone, []character.StatIncrease{{character.StatIntelligence, 2}}, false},
		{character.RaceGnome, character.SubraceForestGnome, []character.StatIncrease{{character.StatDexterity, 1}}, false},
		{character.RaceGnome, character.SubraceRockGnome, []character.StatIncrease{{character.StatConstitution, 1}}, false},
		// Half-Elf
		{character.RaceHalfElf, character.SubraceNone, []character.StatIncrease{{character.StatCharisma, 2}}, false},
		// Half-Orc
		{character.RaceHalfOrc, character.SubraceNone, []character.StatIncrease{
			{character.StatStrength, 2},
			{character.StatConstitution, 1},
		}, false},
		// Tiefling
		{character.RaceTiefling, character.SubraceNone, []character.StatIncrease{
			{character.StatIntelligence, 1},
			{character.StatCharisma, 2},
		}, false},
		// Error cases
		{character.RaceName("UnknownRace"), character.SubraceNone, nil, true},
		{character.RaceDwarf, character.SubraceName("UnknownSubrace"), nil, true},
	}

	for _, test := range tests {
		race := character.Race{
			Type:    test.Type,
			Subrace: test.Subrace,
		}

		increase, err := race.GetAbilityIncrease()
		if err != nil && !test.ExpectErr {
			t.Fatalf("%s-%s got an unexpected err: %v", test.Type, test.Subrace, err)
		}

		if len(increase) != len(test.Expected) {
			t.Fatalf("%s-%s got incorrect length result: %v; want %v", test.Type, test.Subrace, increase, test.Expected)
		}

		checker := make(map[character.StatName]int)

		for i := range test.Expected {
			checker[test.Expected[i].Stat] = test.Expected[i].Amount
		}

		for i := range increase {
			if checker[increase[i].Stat] != increase[i].Amount {
				t.Fatalf("%s-%s got incorrect increase: %v; want %v", test.Type, test.Subrace, increase, test.Expected)
			}
		}
	}
}

func TestGetMoveSpeed(t *testing.T) {
	tests := []struct {
		Type     character.RaceName
		Subrace  character.SubraceName
		Expected int
	}{
		// Dwarf
		{character.RaceDwarf, character.SubraceNone, 25},
		{character.RaceDwarf, character.SubraceHillDwarf, 25},
		{character.RaceDwarf, character.SubraceMountainDwarf, 25},
		// Elf
		{character.RaceElf, character.SubraceNone, 30},
		{character.RaceElf, character.SubraceHighElf, 30},
		{character.RaceElf, character.SubraceWoodElf, 35},
		{character.RaceElf, character.SubraceDrow, 30},
		// Halfling
		{character.RaceHalfling, character.SubraceNone, 25},
		{character.RaceHalfling, character.SubraceLightfoot, 25},
		{character.RaceHalfling, character.SubraceStout, 25},
		// Human
		{character.RaceHuman, character.SubraceNone, 30},
		// Dragonborn
		{character.RaceDragonborn, character.SubraceNone, 30},
		// Gnome
		{character.RaceGnome, character.SubraceNone, 25},
		{character.RaceGnome, character.SubraceForestGnome, 25},
		{character.RaceGnome, character.SubraceRockGnome, 25},
		// Half-Elf
		{character.RaceHalfElf, character.SubraceNone, 30},
		// Half-Orc
		{character.RaceHalfOrc, character.SubraceNone, 30},
		// Tiefling
		{character.RaceTiefling, character.SubraceNone, 30},
		// Error cases
		{character.RaceName("UnknownRace"), character.SubraceNone, 0},
		{character.RaceDwarf, character.SubraceName("UnknownSubrace"), 25},
	}

	for _, test := range tests {
		race := character.Race{
			Type:    test.Type,
			Subrace: test.Subrace,
		}

		speed := race.GetMoveSpeed()
		if speed != test.Expected {
			t.Fatalf("%s-%s got incorrect move speed: %v; want %v", test.Type, test.Subrace, speed, test.Expected)
		}
	}
}
