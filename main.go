package main

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/alex-c96/pokedex-cli/internal/pokeapi"
)

type config struct {
	pokeapiClient       pokeapi.Client
	nextLocationAreaURL *string
	prevLocationAreaURL *string
	caughtPokemon       map[string]pokeapi.Pokemon
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

func commandCatch(cfg *config, args []string) error {
	if len(args) < 2 {
		return errors.New("please provide a pokemon to catch")
	}
	pokemonName := args[1]
	pokemonResp, err := cfg.pokeapiClient.GetPokemon(pokemonName)
	if err != nil {
		return err
	}
	baseExp := pokemonResp.BaseExperience

	randNum := rand.Intn(baseExp)
	const threshold = 50
	fmt.Printf("You throw a pokeball at %s\n", pokemonName)
	fmt.Print(".\n")
	time.Sleep(time.Second * 1)
	fmt.Print("..\n")
	time.Sleep(time.Second * 1)
	fmt.Print("...\n")
	time.Sleep(time.Second * 1)
	if randNum > threshold {
		return fmt.Errorf("failed to catch %s\n", pokemonName)
	}
	fmt.Printf("You caught %s!\n", pokemonName)
	cfg.caughtPokemon[pokemonName] = pokemonResp
	return nil
}

func commandInspect(cfg *config, args []string) error {
	if len(args) < 2 {
		return errors.New("please provide a pokemon to inspect")
	}
	pokemonToInspect := args[1]
	pokemonDetails, ok := cfg.caughtPokemon[pokemonToInspect]
	if !ok {
		return errors.New("you have not cauch that pokemon")
	}
	fmt.Printf("Name: %s\n", pokemonDetails.Name)
	fmt.Printf("Height: %v\n", pokemonDetails.Height)
	fmt.Printf("Weight: %v\n", pokemonDetails.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemonDetails.Stats {
		fmt.Printf("\t-%s: %v\n", stat.Stat.Name, stat.BaseStat)
	}
	fmt.Println("Types:")
	for _, types := range pokemonDetails.Types {
		fmt.Printf("\t- %v\n", types.Type.Name)
	}
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

func commandPokedex(cfg *config, args []string) error {
	if len(cfg.caughtPokemon) < 1 {
		return errors.New("you have not caught any pokemon")
	}
	fmt.Println("Your Pokedex:")
	for _, pokemon := range cfg.caughtPokemon {
		fmt.Printf("\t- %s\n", pokemon.Name)
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
		"catch": {
			name:        "catch",
			description: "try and catch a pokemon",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "inspect the pokemon's stats",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "display the pokemon in your pokedex",
			callback:    commandPokedex,
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
		caughtPokemon: make(map[string]pokeapi.Pokemon),
	}

	startRepl(&cfg)
}
