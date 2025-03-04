package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"internal"
	"io"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type cliCommand struct {
	name        string
	description string
	callback    func(args ...string) error
}

type Loc struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

type Pokemon struct {
	Id             int
	Name           string
	BaseExperience int `json:"base_experience"`
}
type LocArea struct {
	ID                   int
	Name                 string
	GameIndex            int
	EncounterMethodRates []struct{} `json:"encounter_method_rates"`
	Location             struct {
		Name string
		URL  string
	}
	Names             []struct{}
	PokemonEncounters []struct {
		Pokemon struct {
			Name string
			URL  string
		}
		VersionDetails []struct{}
	} `json:"pokemon_encounters"`
}

var commands map[string]cliCommand

func main() {
	caughtPokemon := make(map[string]Pokemon)
	loc := &Loc{}
	scanner := bufio.NewScanner(os.Stdin)
	cache := internal.NewCache(time.Duration(time.Second * 5))
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    func(args ...string) error { return commandExit() },
		},
		"help": {
			name:        "help",
			description: "List all available commands",
			callback:    func(args ...string) error { return commandHelp() },
		},
		"map": {
			name:        "map",
			description: "Display the names of 20 location areas in the Pokemon world",
			callback:    func(args ...string) error { return commandMap(loc, cache) },
		},
		"mapb": {
			name:        "mapb",
			description: "Display the names of Previous 20 location areas in the Pokemon world",
			callback:    func(args ...string) error { return commandMapb(loc, cache) },
		},
		"lenc": {
			name:        "len_cache",
			description: "Display the length of the cache",
			callback: func(args ...string) error {
				fmt.Println("Length of cache:---", len(cache.Items))
				return nil
			},
		},
		"explore": {
			name:        "explore",
			description: "Explore pokemons the location area",
			callback: func(args ...string) error {
				if len(args) != 1 {
					return fmt.Errorf("exactly one location area name is required")
				}
				return commandExplore(args[0], cache)
			},
		},
		"catch": {
			name:        "catch",
			description: "Catch a Pokemon",
			callback: func(args ...string) error {
				if len(args) != 1 {
					return fmt.Errorf("exactly one Pokemon name is required")
				}
				return commandCatch(args[0], &caughtPokemon)
			},
		},
	}
	for {
		fmt.Print("Pokedex > ")
		scanner.Scan()
		text := scanner.Text()
		input := cleanInput(text)
		if len(input) == 0 {
			continue
		}
		if c, exists := commands[input[0]]; exists {
			err := c.callback(input[1:]...)
			if err != nil {
				fmt.Println("Error executing command:", err)
			}

		} else {
			fmt.Println("Unknown command.")
		}

	}
}

func commandExit() error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp() error {
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, c := range commands {
		fmt.Printf("%s: %s\n", c.name, c.description)
	}
	return nil
}
func commandMap(loc *Loc, cache *internal.Cache) error {
	var url string
	defer func() {
		for _, r := range loc.Results {
			fmt.Println(r.Name)
		}
	}()
	if loc.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = loc.Next
	}
	if val, exists := cache.Get(url); exists {
		err := json.Unmarshal(val, loc)
		if err != nil {
			return err
		}
		return nil
	}
	// time.Sleep(time.Second * 2)
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	cache.Add(url, body)
	err = json.Unmarshal(body, loc)
	if err != nil {
		return err
	}

	return nil
}
func commandMapb(loc *Loc, cache *internal.Cache) error {
	if loc.Previous == "" {
		return nil
	}
	defer func() {
		for _, r := range loc.Results {
			fmt.Println(r.Name)
		}
	}()
	if val, exists := cache.Get(loc.Previous); exists {
		err := json.Unmarshal(val, loc)
		if err != nil {
			return err
		}
		return nil
	}
	// time.Sleep(time.Second * 2)
	res, err := http.Get(loc.Previous)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	cache.Add(loc.Previous, body)
	err = json.Unmarshal(body, loc)
	if err != nil {
		return err
	}

	return nil
}

func commandExplore(locationName string, cache *internal.Cache) error {

	locArea := &LocArea{}
	defer func() {
		for _, p := range locArea.PokemonEncounters {
			fmt.Println(p.Pokemon.Name)
		}
	}()
	fmt.Println("Exploring location area", locationName, "...")
	if val, exists := cache.Get(locationName); exists {
		err := json.Unmarshal(val, locArea)
		if err != nil {
			return err
		}
		return nil
	}
	res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + locationName)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	cache.Add(locationName, body)
	err = json.Unmarshal(body, locArea)
	if err != nil {
		return err
	}

	return nil
}

func commandCatch(pokemonName string, caughtPokemon *map[string]Pokemon) error {

	res, err := http.Get("https://pokeapi.co/api/v2/pokemon/" + pokemonName)
	if err != nil {
		return fmt.Errorf("cant find that pokemon %v. error fetching pokemon data: %v", pokemonName, err)
	}
	fmt.Println("Throwing a Pokeball at", pokemonName, "...")
	defer res.Body.Close()
	// Generate a random number between 1 and 1000
	chance := rand.Intn(1000) + 1
	pokemon := Pokemon{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("non-OK HTTP status: %s", res.Status)
	}
	err = json.Unmarshal(body, &pokemon)
	if err != nil {
		return err
	}

	if chance < pokemon.BaseExperience {
		fmt.Println(pokemonName, " escaped!")
		return nil
	} else { // for fun, we'll say that the pokemon is caught if the random number is greater than the base experience
		fmt.Println(pokemonName, " was caught!")
	}
	(*caughtPokemon)[pokemonName] = Pokemon{
		Name:           pokemonName,
		Id:             pokemon.Id,
		BaseExperience: pokemon.BaseExperience,
	}
	fmt.Println("Caught Pokemons:---", *caughtPokemon)
	return nil
}
