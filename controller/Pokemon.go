package controller

import (
	"net/http"
	"strconv"

	"github.com/dannegm/academy-go-q12021/helpers"
	"github.com/dannegm/academy-go-q12021/models"

	"github.com/gin-gonic/gin"
)

// PokemonListResponse for the response
type PokemonListResponse struct {
	Pokedex models.PokedexList `json:"data"`
}

// PokemonResponse for the response
type PokemonResponse struct {
	Pokemon models.Pokemon `json:"data"`
}

// ErrorResponse for error handling
type ErrorResponse struct {
	Status  int    `json:"satus"`
	Message string `json:"message"`
}

// GetPokemonList Pokes router
func GetPokemonList(pokedex models.Pokedex) func(*gin.Context) {
	return func(context *gin.Context) {
		pokedexList := models.PokedexList{}
		for _, pokemon := range pokedex {
			pokedexList = append(pokedexList, pokemon)
		}

		context.JSON(200, PokemonListResponse{
			Pokedex: pokedexList,
		})
	}
}

// GetPokemonFromAPI get a pokemon from PokeAPI
func GetPokemonFromAPI(pokemonID int) error {
	pokeID := strconv.Itoa(pokemonID)
	URL := "https://pokeapi.co/api/v2/pokemon/" + pokeID

	poke := &models.PokemonApi{}
	err := helpers.GetJson(URL, poke)

	if err != nil {
		return err
	}

	pokemon := models.MapPokemonApi(poke)
	return models.StorePokemonInCsv(pokemon)
}

// GetPokemonByID to get a single Pokemon filter by ID
func GetPokemonByID(pokedex models.Pokedex) func(*gin.Context) {
	return func(context *gin.Context) {
		pokemonID, err := strconv.Atoi(context.Param("pokemonID"))

		errApi := GetPokemonFromAPI(pokemonID)
		if errApi != nil {
			context.JSON(http.StatusNotFound, ErrorResponse{
				Status:  http.StatusNotFound,
				Message: "Pokemon not found",
			})
		}

		if err != nil || pokemonID < 1 {
			context.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid Pokemon ID",
			})
		} else {
			pokemon := pokedex[pokemonID]

			context.JSON(http.StatusOK, PokemonResponse{
				Pokemon: pokemon,
			})
		}

	}
}

// FetchPokemonFromApi get first 150 pokemon from PokeAPI
func FetchPokemonFromApi() func(*gin.Context) {
	return func(context *gin.Context) {
		for pokemonID := 1; pokemonID < 150; pokemonID++ {
			GetPokemonFromAPI(pokemonID)
		}
		context.Abort()
	}
}

// GetPokemonListWIthConcurrency to get a list of pokemon filtered by ID if is odd or even
func GetPokemonListWIthConcurrency(pokedex models.Pokedex) func(*gin.Context) {
	return func(context *gin.Context) {
		messages := []string{}

		items, err := strconv.Atoi(context.DefaultQuery("items", "150"))
		if err != nil {
			messages = append(messages, "Invalid items")
		}

		items_per_worker, err := strconv.Atoi(context.DefaultQuery("items_per_worker", "10"))
		if err != nil {
			messages = append(messages, "Invalid items_per_worker")
		}

		if len(messages) > 0 {
			context.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: messages[0],
			})
		}

		evenOdd := context.DefaultQuery("type", "even")

		// Create a slice to hold the selected pokemon.
		pokedexPool := models.PokedexPool(pokedex, evenOdd, items, items_per_worker)

		// Send the pokemonChannel as response
		context.JSON(http.StatusOK, PokemonListResponse{
			Pokedex: pokedexPool,
		})
	}
}
