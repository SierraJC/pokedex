package models

type LocationAreas struct {
	Count    int        `json:"count"`
	Next     string     `json:"next"`
	Previous string     `json:"previous"`
	Results  []Location `json:"results"`
}

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type LocationArea struct {
	Name              string             `json:"name"`
	ID                int                `json:"id"`
	Location          Location           `json:"location"`
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}
