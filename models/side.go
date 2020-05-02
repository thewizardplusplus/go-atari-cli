package models

import (
	models "github.com/thewizardplusplus/go-atari-models"
)

// Side ...
type Side int

// ...
const (
	Searcher Side = iota
	Human
)

// NewSide ...
//
// It detects an initial side
// by a human color.
func NewSide(color models.Color) Side {
	var side Side
	switch color {
	case models.Black:
		side = Human
	case models.White:
		side = Searcher
	}

	return side
}

// Invert ...
func (side Side) Invert() Side {
	if side == Searcher {
		return Human
	}

	return Searcher
}
