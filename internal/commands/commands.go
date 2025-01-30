package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"pokedexcli/internal/models"
	"pokedexcli/internal/pokeapi"
	"pokedexcli/internal/pokecache"
)

type cliCommand struct {
	name        string
	description string
	callback    func(config *cmdConfig, inputs []string) error
}

type cmdConfig struct {
	Next     string
	Previous string
}

const (
	baseURL         = "https://pokeapi.co/api/v2"
	locationAreaURL = baseURL + "/location-area"
	pokemonURL      = baseURL + "/pokemon"
)

var (
	pokedex  = make(map[string]models.Pokemon)
	Commands map[string]cliCommand
	config   = cmdConfig{
		Next:     locationAreaURL + "?offset=0&limit=20",
		Previous: "",
	}
	cache     = pokecache.NewCache(5 * time.Second)
	apiClient = pokeapi.NewClient(cache)
)

func init() {
	Commands = map[string]cliCommand{
		"map": {
			name:        "map",
			description: "Displays the map of the current location",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Go back to the previous page",
			callback:    commandMapBack,
		},
		"explore": {
			name:        "explore",
			description: "Explore the current location",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Inspect a Pokemon",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "List all caught pokemon",
			callback:    commandPokedex,
		},
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
	}
}

func commandMapBack(config *cmdConfig, inputs []string) error {
	if config.Previous == "" {
		fmt.Println("No previous page")
		return nil
	}

	config.Next = config.Previous

	return commandMap(config, inputs)
}

func commandMap(config *cmdConfig, inputs []string) error {
	body, err := apiClient.GetWithCache(config.Next)
	if err != nil {
		return err
	}

	var locationAreasResponse models.LocationAreas
	if err := json.Unmarshal(body, &locationAreasResponse); err != nil {
		return err
	}

	config.Next = locationAreasResponse.Next
	config.Previous = locationAreasResponse.Previous

	for _, locationArea := range locationAreasResponse.Results {
		fmt.Println(locationArea.Name)
	}

	return nil
}

func commandExplore(config *cmdConfig, inputs []string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("No location area provided")
	}

	area := inputs[0]
	fmt.Printf("Exploring %s...\n", area)

	url := locationAreaURL + "/" + area

	body, err := apiClient.GetWithCache(url)
	if err != nil {
		return err
	}

	var locationAreaResponse models.LocationArea
	if err := json.Unmarshal(body, &locationAreaResponse); err != nil {
		return err
	}

	fmt.Println("Found Pokemon:")
	for _, encounter := range locationAreaResponse.PokemonEncounters {
		fmt.Printf(" - %s\n", encounter.Pokemon.Name)
	}

	return nil
}

func commandCatch(_ *cmdConfig, inputs []string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("No Pokemon provided")
	}

	pokemonName := inputs[0]
	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	url := pokemonURL + "/" + pokemonName

	body, err := apiClient.Get(url)
	if err != nil {
		return err
	}

	var pokemon models.Pokemon
	if err := json.Unmarshal(body, &pokemon); err != nil {
		return err
	}

	catchRate := 0.7 - float64(pokemon.BaseExperience)/1000

	if rand.Float64() <= catchRate {
		fmt.Printf("%s was caught!\n", pokemonName)
		pokedex[pokemonName] = pokemon
		return nil
	}

	fmt.Printf("%s escaped!\n", pokemonName)

	return nil
}

func commandInspect(_ *cmdConfig, inputs []string) error {
	if len(inputs) == 0 {
		return fmt.Errorf("No Pokemon provided")
	}

	pokemonName := inputs[0]
	pokemon, ok := pokedex[pokemonName]
	if !ok {
		return fmt.Errorf("you have not caught that pokemon")
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf(" - %s: %d\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, pkType := range pokemon.Types {
		fmt.Printf(" - %s\n", pkType.Type.Name)
	}

	return nil
}

func commandPokedex(_ *cmdConfig, _ []string) error {
	if len(pokedex) == 0 {
		fmt.Println("No pokemon caught")
		return nil
	}

	fmt.Println("Your Pokedex:")
	for name := range pokedex {
		fmt.Println(" -", name)
	}

	return nil
}

func commandExit(_ *cmdConfig, _ []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)

	return nil
}

func commandHelp(_ *cmdConfig, _ []string) error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println()

	for _, command := range Commands {
		fmt.Printf("%s: %s\n", command.name, command.description)
	}

	return nil
}

func validateCommand(input []string) (cliCommand, error) {
	if len(input) == 0 {
		return cliCommand{}, fmt.Errorf("No command provided")
	}

	command, ok := Commands[input[0]]
	if !ok {
		return cliCommand{}, fmt.Errorf("Unknown command")
	}

	return command, nil
}

func ExecuteCommand(input []string) error {
	command, err := validateCommand(input)
	if err != nil {
		return err
	}

	return command.callback(&config, input[1:])
}
