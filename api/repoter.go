package api

import (
	"github.com/pkmngo-odi/pogo-protos"
)

// Reporter is a common interface for collecting findings
type Reporter interface {
	WildPokemons([]*protos.WildPokemon)
	Forts([]*protos.FortData)
}

// VoidReporter sends any reported entries in to the void
type VoidReporter struct{}

// WildPokemons accepts wild pokémons in to the void
func (r *VoidReporter) WildPokemons([]*protos.WildPokemon) {
	// NOOP
}

// Forts accepts wirld forts in to the void
func (r *VoidReporter) Forts([]*protos.FortData) {
	// NOOP
}
