package ascii_test

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
)

func ExampleDecodeColor() {
	color, _ := ascii.DecodeColor("white")
	fmt.Printf("%v\n", color)

	// Output: 1
}

func ExampleEncodeColor() {
	color := ascii.EncodeColor(models.White)
	fmt.Printf("%v\n", color)

	// Output: white
}

func ExampleStoneStorageEncoder_EncodeStoneStorage() {
	stoneEncoder := func(color models.Color) string {
		return string(sgf.EncodeColor(color))
	}
	placeholders := ascii.Placeholders{
		Crosshairs: "+",
	}

	// +-+-+-+-+-+
	// | | | | | |
	// +-+-+-+-+-+
	// | |B|W| | |
	// +-+-+-+-+-+
	// |B|W| |W| |
	// +-+-+-+-+-+
	// | |B|W| | |
	// +-+-+-+-+-+
	// | | | | | |
	// +-+-+-+-+-+
	board := models.NewBoard(models.Size{Width: 5, Height: 5})
	for _, move := range []models.Move{
		{Color: models.Black, Point: models.Point{Column: 1, Row: 1}},
		{Color: models.White, Point: models.Point{Column: 2, Row: 1}},
		{Color: models.Black, Point: models.Point{Column: 0, Row: 2}},
		{Color: models.White, Point: models.Point{Column: 1, Row: 2}},
		{Color: models.White, Point: models.Point{Column: 3, Row: 2}},
		{Color: models.Black, Point: models.Point{Column: 1, Row: 3}},
		{Color: models.White, Point: models.Point{Column: 2, Row: 3}},
	} {
		board = board.ApplyMove(move)
	}

	boardEncoder :=
		ascii.NewStoneStorageEncoder(stoneEncoder, placeholders, ascii.Margins{}, 1)
	fmt.Printf("%v\n", boardEncoder.EncodeStoneStorage(board))

	// Output:
	// e+++++
	// d+BW++
	// cBW+W+
	// b+BW++
	// a+++++
	//  abcde
}

func ExampleStoneStorageEncoder_EncodeStoneStorage_withMargins() {
	stoneEncoder := func(color models.Color) string {
		return string(sgf.EncodeColor(color))
	}
	placeholders := ascii.Placeholders{
		HorizontalLine: " ",
		Crosshairs:     "+",
	}
	margins := ascii.Margins{
		Stone: ascii.StoneMargins{
			HorizontalMargins: ascii.HorizontalMargins{
				Left: 1,
			},
		},
		Legend: ascii.LegendMargins{
			Row: ascii.HorizontalMargins{
				Right: 1,
			},
		},
	}

	// +-+-+-+-+-+
	// | | | | | |
	// +-+-+-+-+-+
	// | |B|W| | |
	// +-+-+-+-+-+
	// |B|W| |W| |
	// +-+-+-+-+-+
	// | |B|W| | |
	// +-+-+-+-+-+
	// | | | | | |
	// +-+-+-+-+-+
	board := models.NewBoard(models.Size{Width: 5, Height: 5})
	for _, move := range []models.Move{
		{Color: models.Black, Point: models.Point{Column: 1, Row: 1}},
		{Color: models.White, Point: models.Point{Column: 2, Row: 1}},
		{Color: models.Black, Point: models.Point{Column: 0, Row: 2}},
		{Color: models.White, Point: models.Point{Column: 1, Row: 2}},
		{Color: models.White, Point: models.Point{Column: 3, Row: 2}},
		{Color: models.Black, Point: models.Point{Column: 1, Row: 3}},
		{Color: models.White, Point: models.Point{Column: 2, Row: 3}},
	} {
		board = board.ApplyMove(move)
	}

	boardEncoder :=
		ascii.NewStoneStorageEncoder(stoneEncoder, placeholders, margins, 1)
	fmt.Printf("%v\n", boardEncoder.EncodeStoneStorage(board))

	// Output:
	// e  + + + + +
	// d  + B W + +
	// c  B W + W +
	// b  + B W + +
	// a  + + + + +
	//    a b c d e
}

func ExampleStoneStorageEncoder_EncodeStoneStorage_withGrid() {
	stoneEncoder := func(color models.Color) string {
		return string(sgf.EncodeColor(color))
	}
	placeholders := ascii.Placeholders{
		HorizontalLine: "-",
		VerticalLine:   "|",
		Crosshairs:     "+",
	}
	margins := ascii.Margins{
		Stone: ascii.StoneMargins{
			HorizontalMargins: ascii.HorizontalMargins{
				Left: 1,
			},
			VerticalMargins: ascii.VerticalMargins{
				Bottom: 1,
			},
		},
		Legend: ascii.LegendMargins{
			Row: ascii.HorizontalMargins{
				Right: 1,
			},
		},
	}

	// +-+-+-+-+-+
	// | | | | | |
	// +-+-+-+-+-+
	// | |B|W| | |
	// +-+-+-+-+-+
	// |B|W| |W| |
	// +-+-+-+-+-+
	// | |B|W| | |
	// +-+-+-+-+-+
	// | | | | | |
	// +-+-+-+-+-+
	board := models.NewBoard(models.Size{Width: 5, Height: 5})
	for _, move := range []models.Move{
		{Color: models.Black, Point: models.Point{Column: 1, Row: 1}},
		{Color: models.White, Point: models.Point{Column: 2, Row: 1}},
		{Color: models.Black, Point: models.Point{Column: 0, Row: 2}},
		{Color: models.White, Point: models.Point{Column: 1, Row: 2}},
		{Color: models.White, Point: models.Point{Column: 3, Row: 2}},
		{Color: models.Black, Point: models.Point{Column: 1, Row: 3}},
		{Color: models.White, Point: models.Point{Column: 2, Row: 3}},
	} {
		board = board.ApplyMove(move)
	}

	boardEncoder :=
		ascii.NewStoneStorageEncoder(stoneEncoder, placeholders, margins, 1)
	fmt.Printf("%v\n", boardEncoder.EncodeStoneStorage(board))

	// Output:
	// e -+-+-+-+-+
	//    | | | | |
	// d -+-B-W-+-+
	//    | | | | |
	// c -B-W-+-W-+
	//    | | | | |
	// b -+-B-W-+-+
	//    | | | | |
	// a -+-+-+-+-+
	//    | | | | |
	//    a b c d e
}
