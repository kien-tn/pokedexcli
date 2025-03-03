package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

type LocArea struct {
	Count    int
	Next     string
	Previous string
	Results  []struct {
		Name string
		Url  string
	}
}

var commands map[string]cliCommand

func main() {
	loc := &LocArea{}
	scanner := bufio.NewScanner(os.Stdin)
	commands = map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "List all available commands",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Display the names of 20 location areas in the Pokemon world",
			callback:    func() error { return commandMap(loc) },
		},
		"mapb": {
			name:        "mapb",
			description: "Display the names of Previous 20 location areas in the Pokemon world",
			callback:    func() error { return commandMapb(loc) },
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
			err := c.callback()
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
func commandMap(loc *LocArea) error {
	var url string
	if loc.Next == "" {
		url = "https://pokeapi.co/api/v2/location-area/"
	} else {
		url = loc.Next
	}
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(loc)
	if err != nil {
		return err
	}
	for _, r := range loc.Results {
		fmt.Println(r.Name)
	}
	return nil
}
func commandMapb(loc *LocArea) error {
	if loc.Previous == "" {
		return nil
	}
	res, err := http.Get(loc.Previous)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&loc)
	if err != nil {
		return err
	}
	for _, r := range loc.Results {
		fmt.Println(r.Name)
	}
	return nil
}
