package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"pokedexcli/internal/commands"
	"pokedexcli/internal/models"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

const (
	baseURL         = "https://pokeapi.co/api/v2"
	locationAreaURL = baseURL + "/location-area"
	pokemonURL      = baseURL + "/pokemon"
)

var (
	pokedex   = make(map[string]models.Pokemon)
	cache     = pokecache.NewCache(5 * time.Second)
	apiClient = pokeapi.NewClient(cache)
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		input := strings.Fields(strings.ToLower(scanner.Text()))

		if len(input) == 0 {
			continue
		}

		if err := commands.ExecuteCommand(input); err != nil {
			fmt.Println(err)
		}
	}
}
