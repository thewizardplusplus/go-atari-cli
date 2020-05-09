package unicode

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
		text = "\u25cf"
	case models.White:
		text = "\u25cb"
	}

	return text
}
