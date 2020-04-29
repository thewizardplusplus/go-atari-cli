package ascii

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
		data{
			args: args{models.Black},
			want: "*",
		},
		data{
			args: args{models.White},
			want: "o",
		},
	} {
		got := EncodeStone(data.args.color)

		if got != data.want {
			test.Fail()
		}
	}
}
