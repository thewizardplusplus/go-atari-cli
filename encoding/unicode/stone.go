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
		text = "\u25cb"
	case models.White:
		text = "\u25cf"
	}

	return text
}
