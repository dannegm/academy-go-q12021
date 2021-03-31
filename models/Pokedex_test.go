package models_test

import (
	"fmt"
	"testing"

	"github.com/dannegm/academy-go-q12021/models"
)

func TestPokedexFromFile(t *testing.T) {
	result, err := models.PokedexFromFile()

	if err != nil {
		t.Error("Can't read the csv file")
	}

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

	if len(result) != 151 {
		t.Error("Num of received pokemon not match")
	}

	if expected.Name != result[1].Name {
		t.Error("First pokemon does not match")
	}
}

func TestPokedexPool(t *testing.T) {
	pokedex, err := models.PokedexFromFile()
	if err != nil {
		t.Error("Can't read pokemon file")
	}

	result := models.PokedexPool(pokedex, "even", 15, 10)

	if len(result) != 15 {
		t.Error("Worker not return the right amount of pokemon")
	}
}

func TestAppendEvenOdd(t *testing.T) {
	evenList := models.PokedexList{}
	evenPokemon := models.Pokemon{
		ID: 1,
	}
	evenList = models.AppendEvenOdd("odd", evenList, evenPokemon)

	fmt.Print(evenList)

	if len(evenList) != 0 {
		t.Error("Should not appear pokemons")
	}

	evenList = models.AppendEvenOdd("even", evenList, evenPokemon)

	if len(evenList) != 1 {
		t.Error("Should appear at least one pokemon")
	}
}
