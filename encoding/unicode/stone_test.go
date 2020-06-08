package unicode

import (
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestEncodeStone(test *testing.T) {
	type args struct {
		color models.Color
	}
	type data struct {
		args args
		want string
	}

	for _, data := range []data{
		{
			args: args{models.Black},
			want: "\u25cf",
		},
		{
			args: args{models.White},
			want: "\u25cb",
		},
	} {
		got := EncodeStone(data.args.color)

		if got != data.want {
			test.Fail()
		}
	}
}
