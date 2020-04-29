package ascii

import (
	models "github.com/thewizardplusplus/go-atari-models"
)

// EncodeStone ...
func EncodeStone(
	color models.Color,
) string {
	var text string
	switch color {
	case models.Black:
		text = "*"
	case models.White:
		text = "o"
	}

	return text
}
