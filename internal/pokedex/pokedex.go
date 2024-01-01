package pokedex

import (
	"github.com/alex-c96/pokedex-cli/internal/pokeapi"
)

type Pokedex struct {
	CaughtPokemon map[string]Pokemon
}
