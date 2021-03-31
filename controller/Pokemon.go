package controller

import (
	"math"
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

func GetPokemonListWIthConcurrency(pokedex models.Pokedex) func(*gin.Context) {
	return func(context *gin.Context) {
		messages := []string{}

		items_per_worker, err := strconv.Atoi(context.DefaultQuery("items_per_worker", "10"))
		if err != nil {
			messages = append(messages, "Invalid items_per_worker")
		}

		items, err := strconv.Atoi(context.DefaultQuery("items", "150"))
		if err != nil {
			messages = append(messages, "Invalid items")
		}

		if len(messages) > 0 {
			context.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: messages[0],
			})
		}

		evenOdd := context.DefaultQuery("type", "even")

		// Create the channel for sharing pokemons.
		pokemonChannel := make(chan models.Pokemon)
		// Create a channel "shutdown" to tell goroutines when to terminate.
		shutdown := make(chan struct{})

		// Define the size of the worker pool.
		// chunk total items to process by items_per_worker
		poolSize := int(math.Ceil(float64(items) / float64(items_per_worker)))

		// Create a sync.WaitGroup to monitor the Goroutine pool. Add the count.
		var wg sync.WaitGroup
		wg.Add(poolSize)

		// Limit of pokemons available
		pokemonsAvailable := len(pokedex)

		// Create a fixed size pool of goroutines to retrive pokemons
		for i := 0; i < poolSize; i++ {
			// starts go rutine
			go func(id int) {
				// get the start pokemonID base on the items per worker and the current thread
				pokemonID := (items_per_worker * id)

				// Infinite loop to fetch pokemons
				for {
					// increment the pokemon id to find the next pokemon available
					if pokemonsAvailable > pokemonID {
						pokemonID++
					}

					// get the current pokemon
					poke := pokedex[pokemonID]

					select {
					// inject current pokemon to the cannel
					case pokemonChannel <- poke:

					// exit of gorutine
					case <-shutdown:
						wg.Done()
						return
					}
				}
			}(i)
		}

		// Create a slice to hold the selected pokemon.
		pokedexPool := PokedexList{}
		// iterate channel to generate the pool
		for poke := range pokemonChannel {

			if evenOdd == "even" && poke.ID%2 != 0 {
				pokedexPool = append(pokedexPool, poke)
			}

			if evenOdd == "odd" && poke.ID%2 == 0 {
				pokedexPool = append(pokedexPool, poke)
			}

			// break the loop once we reach the limit of the list
			if len(pokedexPool) == pokemonsAvailable {
				break
			}

			// break the loop once we reach the items requested
			if len(pokedexPool) == items {
				break
			}
		}

		// Send the shutdown signal by closing the channel.
		close(shutdown)

		// Wait for the Goroutines to finish.
		wg.Wait()

		// Send the pokemonChannel as response
		context.JSON(http.StatusOK, PokemonListResponse{
			Pokedex: pokedexPool,
		})
	}
}
