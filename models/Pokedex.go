package models

import (
	"fmt"
	"math"
	"sync"

	"github.com/dannegm/academy-go-q12021/constants"
	"github.com/dannegm/academy-go-q12021/helpers"
)

// Pokedex stores a map of Pokemon by ID
type Pokedex map[int]Pokemon

// PokedexList stores a list of Pokemon
type PokedexList []Pokemon

// PokedexFromFile return a list of Pokemon
func PokedexFromFile() (pokes Pokedex, err error) {
	rows, err := helpers.ReadFile(fmt.Sprintf("%s/pokemon.csv", constants.AssetsPath))

	if err != nil {
		return
	}

	pokes = make(Pokedex)

	for _, row := range rows[1:167] {
		poke := PokemonFromString(row)

		// Megaevolutions appears as a duplicated key, so, will skipped
		if _, exist := pokes[poke.ID]; !exist {
			pokes[poke.ID] = poke
		}
	}

	return
}

func PokedexPool(pokedex Pokedex, evenOdd string, length int, itemsPerWorker int) PokedexList {
	// Limit of pokemons available
	pokemonsAvailable := len(pokedex)
	// Create the channel for sharing pokemons.
	pokemonChannel := make(chan Pokemon)
	// Create a channel "shutdown" to tell goroutines when to terminate.
	shutdown := make(chan struct{})

	// Define the size of the worker pool.
	// chunk total length to process by itemsPerWorker
	poolSize := int(math.Ceil(float64(length) / float64(itemsPerWorker)))
	// Create a sync.WaitGroup to monitor the Goroutine pool. Add the count.
	var waitGroup sync.WaitGroup
	waitGroup.Add(poolSize)

	// Create a fixed size pool of goroutines to retrive pokemons
	for threadID := 0; threadID < poolSize; threadID++ {
		// starts go rutine
		go func(id int) {
			// get the start pokemonID base on the items per worker and the current thread
			pokemonID := (itemsPerWorker * id)
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
					waitGroup.Done()
					return
				}
			}
		}(threadID)
	}

	// Create a slice to hold the selected pokemon.
	pokedexPool := PokedexList{}
	// iterate channel to generate the pool
	for poke := range pokemonChannel {
		pokedexPool = AppendEvenOdd(evenOdd, pokedexPool, poke)

		// break the loop once we reach the limit of the list
		if len(pokedexPool) == pokemonsAvailable {
			break
		}

		// break the loop once we reach the items requested
		if len(pokedexPool) == length {
			break
		}
	}

	// Send the shutdown signal by closing the channel.
	close(shutdown)
	// Wait for the Goroutines to finish.
	waitGroup.Wait()
	// Return the pool list
	return pokedexPool
}

func AppendEvenOdd(evenOdd string, pokedex PokedexList, poke Pokemon) PokedexList {
	if evenOdd == "even" && poke.ID%2 != 0 {
		return append(pokedex, poke)
	}

	if evenOdd == "odd" && poke.ID%2 == 0 {
		return append(pokedex, poke)
	}

	return pokedex
}
