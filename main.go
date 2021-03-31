package main

import (
	"fmt"

	"github.com/dannegm/academy-go-q12021/controller"
	"github.com/dannegm/academy-go-q12021/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Read pokemons from CSV File
	pokedex, err := models.PokedexFromFile()

	if err != nil {
		fmt.Println("There was an error reading pokedex file")
	}

	// Create server instance
	app := gin.Default()

	// TODO: Implements an Client React APP
	// Serve react static files
	// app.Use(static.Serve("/", static.LocalFile("./client/build", true)))

	// Edpoints
	app.GET("/pokedex", controller.GetPokemonList(pokedex))
	app.GET("/pokemon/:pokemonID", controller.GetPokemonByID(pokedex))
	app.GET("/fetchData", controller.FetchPokemonFromApi())
	app.GET("/pokedex/worker", controller.GetPokemonListWIthConcurrency(pokedex))

	// Mounting server - localhost to avoid networks dialogs
	app.Run("localhost:3000")
}
