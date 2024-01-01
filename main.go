package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

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
	callback    func(*config, []string) error
}

func commandHelp(cfg *config, args []string) error {
	commands := getCommands()
	fmt.Println("Welcome to the Pokedex!")
	fmt.Println("Usage:")
	fmt.Println("")
	for _, val := range commands {
		fmt.Printf("%s: %s\n", val.name, val.description)
	}
	return nil
}

func commandExit(cfg *config, args []string) error {
	fmt.Println("exiting program.")
	os.Exit(0)
	return nil
}

func commandMap(cfg *config, args []string) error {
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

func commandMapb(cfg *config, args []string) error {
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

func commandExplore(cfg *config, args []string) error {
	if len(args) < 2 {
		return errors.New("please provide an area argument")
	}
	exploreResp, err := cfg.pokeapiClient.ExploreLocation(args[1])
	if err != nil {
		return err
	}
	for _, pokemon := range exploreResp.PokemonEncounters {
		fmt.Printf(" - %s\n", pokemon.Pokemon.Name)
	}
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
		"explore": {
			name:        "explore",
			description: "explore an area for Pokemon",
			callback:    commandExplore,
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
			args := strings.Split(line, " ")

			command, ok := commands[args[0]]
			if !ok {
				fmt.Println("invalid command")
				break
			}
			err := command.callback(cfg, args)
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
		pokeapiClient: pokeapi.NewClient(time.Hour),
	}

	startRepl(&cfg)
}
