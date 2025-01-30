package models

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name           string        `json:"name"`
	URL            string        `json:"url,omitempty"`
	BaseExperience int           `json:"base_experience"`
	Height         int           `json:"height"`
	Weight         int           `json:"weight"`
	Stats          []PokemonStat `json:"stats"`
	Types          []PokemonType `json:"types"`
}

type PokemonStat struct {
	BaseStat int  `json:"base_stat"`
	Stat     Stat `json:"stat"`
}

type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type PokemonType struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}
