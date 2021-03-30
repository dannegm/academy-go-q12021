package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/dannegm/academy-go-q12021/helpers"
	"github.com/dannegm/academy-go-q12021/models"
	"github.com/gin-gonic/gin"
)

// PokedexList stores a list of Pokemon
type PokedexList []models.Pokemon

// PokemonListResponse for the response
type PokemonListResponse struct {
	Pokedex PokedexList `json:"data"`
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
		pokedexList := PokedexList{}
		for _, pokemon := range pokedex {
			pokedexList = append(pokedexList, pokemon)
		}

		context.JSON(200, PokemonListResponse{
			Pokedex: pokedexList,
		})
	}
}

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

func FetchPokemonFromApi() func(*gin.Context) {
	return func(context *gin.Context) {

		for pokemonID := 1; pokemonID < 150; pokemonID++ {
			GetPokemonFromAPI(pokemonID)
		}

		context.Abort()
	}
}

func GetAllPokemonWIthConcurrency(pokedex models.Pokedex) func(*gin.Context) {
	return func(context *gin.Context) {

		messages := []string{}

		items_per_worker, err := strconv.Atoi(context.DefaultQuery("items_per_worker", ""))
		if err != nil {
			messages = append(messages, "Invalid items_per_worker")
		}

		items, err := strconv.Atoi(context.DefaultQuery("items", ""))
		if err != nil {
			messages = append(messages, "Invalid items")
		}

		typeQuery := context.DefaultQuery("type", "")

		// Create the channel for sharing results.
		values := make(chan int)

		// Create a channel "shutdown" to tell goroutines when to terminate.
		shutdown := make(chan struct{})

		// Define the size of the worker pool. Use runtime.GOMAXPROCS(0) to size the pool based on number of processors.
		poolSize := items / items_per_worker

		// Create a sync.WaitGroup to monitor the Goroutine pool. Add the count.
		var wg sync.WaitGroup
		wg.Add(poolSize)

		// Create a fixed size pool of goroutines to generate random numbers.
		for i := 0; i < poolSize; i++ {
			go func(id int) {

				for {

					// Generate a random number up to 1000.
					n := items_per_worker

					// Use a select to either send the number or receive the shutdown signal.
					select {

					// In one case send the random number.
					case values <- n:
						// Lo que ejecuta mi worker

						if typeQuery != "odd" && n%2 != 0 {
							// TODO: Sacar a poke del worker
							poke := pokedex[n]
						}

					// In another case receive from the shutdown channel.
					case <-shutdown:
						fmt.Printf("Worker %d shutting down\n", id)
						wg.Done()
						return
					}
				}

			}(i)
		}

	}
}
