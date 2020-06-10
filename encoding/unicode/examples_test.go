package unicode_test

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/unicode"
	models "github.com/thewizardplusplus/go-atari-models"
)

func ExampleEncodeStone() {
	stone := unicode.EncodeStone(models.White)
	fmt.Printf("%v\n", stone)

	// Output: ○
}
