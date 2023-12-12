package main

import (
	"bufio"
	"fmt"
	"os"
)

type cliCommand struct {
	name        string
	description string
	callback    func() error
}

func commandHelp() error {
	fmt.Println("help function")
	return nil
}

func commandExit() error {
	fmt.Println("exiting program.")
	os.Exit(0)
	return nil
}

func getCommands() map[string]cliCommand {
	commands := map[string]cliCommand{
		"help": {
			name:        "help",
			description: "display a help message",
			callback:    commandHelp,
		},
		"exit": {
			name:        "exit",
			description: "Used to exit the app",
			callback:    commandExit,
		},
	}
	return commands
}

func main() {

	commands := getCommands()
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("Pokedex > ")

		for scanner.Scan() {
			line := scanner.Text()

			switch line {
			case "help":
				commands["help"].callback()
			case "exit":
				commands["exit"].callback()
			default:
				fmt.Println("invalid command")
			}
			break
		}

		// Check for scanner errors
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading input:", err)
		}

	}
}
