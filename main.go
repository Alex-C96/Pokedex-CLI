package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/alex-c96/pokedex-cli/internal/pokeapi"
)

type config struct {
	pokeapiClient       pokeapi.Client
	nextLocationAreaURL *string
	prevLocationAreaURL *string
}

type cliCommand struct {
	name        string
	description string
	callback    func(*config) error
}

func commandHelp(cfg *config) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, val := range commands {
		fmt.Printf("%s: %s\n", val.name, val.description)
	}
	return nil
}

func commandExit(cfg *config) error {
	fmt.Println("exiting program.")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config) error {
	locationResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.nextLocationAreaURL)
	if err != nil {
		return err
	}
	for _, location := range locationResp.Results {
		fmt.Printf(" - %s\n", location.Name)
	}
	cfg.nextLocationAreaURL = locationResp.Next
	cfg.prevLocationAreaURL = locationResp.Previous
	return nil
}

func commandMapb(cfg *config) error {
	if cfg.prevLocationAreaURL == nil {
		return errors.New("you are on the first page")
	}
	locationResp, err := cfg.pokeapiClient.ListLocationAreas(cfg.prevLocationAreaURL)
	if err != nil {
		return err
	}
	for _, location := range locationResp.Results {
		fmt.Printf(" - %s\n", location.Name)
	}
	cfg.nextLocationAreaURL = locationResp.Next
	cfg.prevLocationAreaURL = locationResp.Previous
	return nil
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "display a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "display a list of locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "display the previous list of locations",
			callback:    commandMapb,
		},
		"exit": {
			name:        "exit",
			description: "Used to exit the Pokedex",
			callback:    commandExit,
		},
	}
	return commands
}

func startRepl(cfg *config) {
	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		for scanner.Scan() {
			line := scanner.Text()

			command, ok := commands[line]
			if !ok {
				fmt.Println("invalid command")
			}
			err := command.callback(cfg)
			if err != nil {
				fmt.Println(err)
			}
			break
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}
	}
}

func main() {

	cfg := config{
		pokeapiClient: pokeapi.NewClient(),
	}

	startRepl(&cfg)
}
