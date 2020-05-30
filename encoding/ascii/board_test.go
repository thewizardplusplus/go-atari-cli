package ascii

import (
	"reflect"
	"strings"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
)

func TestNewStoneStorageEncoder(
	test *testing.T,
) {
	stoneEncoder := func(
		color models.Color,
	) string {
		symbol := sgf.EncodeColor(color)
		return string(symbol)
	}
	placeholders := Placeholders{
		HorizontalLine: "-",
		VerticalLine:   "|",
		Crosshairs:     "+",
	}
	margins := Margins{
		Stone: StoneMargins{
			HorizontalMargins: HorizontalMargins{
				Left:  1,
				Right: 2,
			},
			VerticalMargins: VerticalMargins{
				Top:    3,
				Bottom: 4,
			},
		},
		Legend: LegendMargins{
			Column: VerticalMargins{
				Top:    5,
				Bottom: 6,
			},
			Row: HorizontalMargins{
				Left:  7,
				Right: 8,
			},
		},
	}
	encoder := NewStoneStorageEncoder(
		stoneEncoder,
		placeholders,
		margins,
		2,
	)

	gotEncoder := reflect.
		ValueOf(encoder.encoder).
		Pointer()
	wantEncoder := reflect.
		ValueOf(stoneEncoder).
		Pointer()
	if gotEncoder != wantEncoder {
		test.Fail()
	}

	if !reflect.DeepEqual(
		encoder.placeholders,
		placeholders,
	) {
		test.Fail()
	}

	if !reflect.DeepEqual(
		encoder.margins,
		margins,
	) {
		test.Fail()
	}

	if encoder.stoneWidth != 2 {
		test.Fail()
	}
}

func TestStoneStorageEncoderEncodeStoneStorage(
	test *testing.T,
) {
	type fields struct {
		encoder      StoneEncoder
		placeholders Placeholders
		margins      Margins
		stoneWidth   int
	}
	type args struct {
		storage models.StoneStorage
	}
	type data struct {
		fields fields
		args   args
		want   string
	}

	for _, data := range []data{
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins:    Margins{},
				stoneWidth: 1,
			},
			args: args{
				storage: models.NewBoard(
					models.Size{
						Width:  3,
						Height: 3,
					},
				),
			},
			want: "c+++\n" +
				"b+++\n" +
				"a+++\n" +
				" abc",
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins:    Margins{},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: "c+W+\n" +
				"bWBW\n" +
				"a+W+\n" +
				" abc",
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: " ",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins: Margins{
					Stone: StoneMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: "c +   W   +  \n" +
				"b W   B   W  \n" +
				"a +   W   +  \n" +
				"  a   b   c  ",
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins: Margins{
					Stone: StoneMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: "c-+---W---+--\n" +
				"b-W---B---W--\n" +
				"a-+---W---+--\n" +
				"  a   b   c  ",
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   " ",
					Crosshairs:     "+",
				},
				margins: Margins{
					Stone: StoneMargins{
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: strings.Repeat(" ", 4) + "\n" +
				"c+W+\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				"bWBW\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				"a+W+\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				" abc",
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins: Margins{
					Stone: StoneMargins{
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: " |||\n" +
				"c+W+\n" +
				" |||\n" +
				" |||\n" +
				" |||\n" +
				"bWBW\n" +
				" |||\n" +
				" |||\n" +
				" |||\n" +
				"a+W+\n" +
				" |||\n" +
				" |||\n" +
				" abc",
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins: Margins{
					Legend: LegendMargins{
						Row: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: " c  +W+\n" +
				" b  WBW\n" +
				" a  +W+\n" +
				"    abc",
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins: Margins{
					Legend: LegendMargins{
						Column: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: "c+W+\n" +
				"bWBW\n" +
				"a+W+\n" +
				strings.Repeat(" ", 4) + "\n" +
				" abc\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4),
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: " ",
					VerticalLine:   " ",
					Crosshairs:     "+",
				},
				margins: Margins{
					Stone: StoneMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
					Legend: LegendMargins{
						Column: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
						Row: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: strings.Repeat(" ", 4*4) + "\n" +
				" c   +   W   +  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				" b   W   B   W  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				" a   +   W   +  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				"     a   b   c  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4),
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins: Margins{
					Stone: StoneMargins{
						HorizontalMargins: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
						VerticalMargins: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
					},
					Legend: LegendMargins{
						Column: VerticalMargins{
							Top:    1,
							Bottom: 2,
						},
						Row: HorizontalMargins{
							Left:  1,
							Right: 2,
						},
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: "     |   |   |  \n" +
				" c  -+---W---+--\n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				" b  -W---B---W--\n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				" a  -+---W---+--\n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				"     a   b   c  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4),
		},
		data{
			fields: fields{
				encoder: func(
					color models.Color,
				) string {
					symbol := sgf.EncodeColor(color)
					return string(symbol)
				},
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins: Margins{
					Board: VerticalMargins{
						Top:    1,
						Bottom: 2,
					},
				},
				stoneWidth: 1,
			},
			args: args{
				storage: func() models.StoneStorage {
					board := models.NewBoard(
						models.Size{
							Width:  3,
							Height: 3,
						},
					)

					moves := []models.Move{
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    0,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 0,
								Row:    1,
							},
						},
						models.Move{
							Color: models.Black,
							Point: models.Point{
								Column: 1,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 2,
								Row:    1,
							},
						},
						models.Move{
							Color: models.White,
							Point: models.Point{
								Column: 1,
								Row:    2,
							},
						},
					}
					for _, move := range moves {
						board = board.ApplyMove(move)
					}

					return board
				}(),
			},
			want: strings.Repeat(" ", 4) + "\n" +
				"c+W+\n" +
				"bWBW\n" +
				"a+W+\n" +
				" abc\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4),
		},
	} {
		encoder := StoneStorageEncoder{
			encoder: data.fields.encoder,
			placeholders: data.fields.
				placeholders,
			margins:    data.fields.margins,
			stoneWidth: data.fields.stoneWidth,
		}
		got := encoder.EncodeStoneStorage(
			data.args.storage,
		)

		if got != data.want {
			test.Fail()
		}
	}
}
