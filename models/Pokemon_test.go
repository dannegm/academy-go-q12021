package models_test

import (
	"fmt"
	"testing"

	"github.com/dannegm/academy-go-q12021/models"
)

func TestPokemonToCsvRow(t *testing.T) {
	input := models.Pokemon{
		ID:           1,
		Name:         "bulbasaur",
		TypeOne:      "grass",
		Total:        64,
		HealthPoints: 45,
		Attack:       49,
		Defense:      49,
		Speed:        45,
	}

	expected := "1,bulbasaur,grass,poison,64,45,49,49,45\n"
	result := models.PokemonToCsvRow(input)

	if expected != result {
		t.Error("Expected 1,bulbasaur,grass,poison,64,45,49,49,45, recieved: " + result)
	}
}

func TestPokemonFromString(t *testing.T) {
	input := "1,Bulbasaur,Grass,Poison,318,45,49,49,65,65,45,1,False"

	expected := models.Pokemon{
		ID:           1,
		Name:         "Bulbasaur",
		TypeOne:      "Grass",
		Total:        318,
		HealthPoints: 45,
		Attack:       49,
		Defense:      49,
		Speed:        45,
	}
	result := models.PokemonFromString(input)

	if expected != result {
		t.Error("Expected pokemon structure")
	}
}

func TestStorePokemonInCsv(t *testing.T) {
	input := models.Pokemon{
		ID:           1,
		Name:         "Bulbasaur",
		TypeOne:      "Grass",
		Total:        318,
		HealthPoints: 45,
		Attack:       49,
		Defense:      49,
		Speed:        45,
	}
	expected := models.StorePokemonInCsv(input)
	fmt.Println(expected)

	if expected != nil {
		t.Error("Can't store in a file")
	}
}
