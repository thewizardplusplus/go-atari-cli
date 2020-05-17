package ascii

import (
	"reflect"
	"strings"
	"testing"

	models "github.com/thewizardplusplus/go-atari-models"
)

func TestBoardEncoder(test *testing.T) {
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
	encoder := NewBoardEncoder(
		EncodeStone,
		placeholders,
		margins,
		2,
	)

	gotEncoder := reflect.
		ValueOf(encoder.encoder).
		Pointer()
	wantEncoder := reflect.
		ValueOf(EncodeStone).
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

func TestBoardEncoderEncodeBoard(
	test *testing.T,
) {
	type fields struct {
		encoder      StoneEncoder
		placeholders Placeholders
		margins      Margins
		stoneWidth   int
	}
	type args struct {
		board models.Board
	}
	type data struct {
		fields fields
		args   args
		want   string
	}

	for _, data := range []data{
		data{
			fields: fields{
				encoder: EncodeStone,
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins:    Margins{},
				stoneWidth: 1,
			},
			args: args{
				board: models.NewBoard(
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
				encoder: EncodeStone,
				placeholders: Placeholders{
					HorizontalLine: "-",
					VerticalLine:   "|",
					Crosshairs:     "+",
				},
				margins:    Margins{},
				stoneWidth: 1,
			},
			args: args{
				board: func() models.Board {
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
			want: "c+o+\n" +
				"bo*o\n" +
				"a+o+\n" +
				" abc",
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
			want: "c +   o   +  \n" +
				"b o   *   o  \n" +
				"a +   o   +  \n" +
				"  a   b   c  ",
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
			want: "c-+---o---+--\n" +
				"b-o---*---o--\n" +
				"a-+---o---+--\n" +
				"  a   b   c  ",
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
				"c+o+\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				"bo*o\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				"a+o+\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4) + "\n" +
				" abc",
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
				"c+o+\n" +
				" |||\n" +
				" |||\n" +
				" |||\n" +
				"bo*o\n" +
				" |||\n" +
				" |||\n" +
				" |||\n" +
				"a+o+\n" +
				" |||\n" +
				" |||\n" +
				" abc",
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
			want: " c  +o+\n" +
				" b  o*o\n" +
				" a  +o+\n" +
				"    abc",
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
			want: "c+o+\n" +
				"bo*o\n" +
				"a+o+\n" +
				strings.Repeat(" ", 4) + "\n" +
				" abc\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4),
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
				" c   +   o   +  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				" b   o   *   o  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				" a   +   o   +  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4) + "\n" +
				"     a   b   c  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4),
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
				" c  -+---o---+--\n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				" b  -o---*---o--\n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				" a  -+---o---+--\n" +
				"     |   |   |  \n" +
				"     |   |   |  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				"     a   b   c  \n" +
				strings.Repeat(" ", 4*4) + "\n" +
				strings.Repeat(" ", 4*4),
		},
		data{
			fields: fields{
				encoder: EncodeStone,
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
				board: func() models.Board {
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
				"c+o+\n" +
				"bo*o\n" +
				"a+o+\n" +
				" abc\n" +
				strings.Repeat(" ", 4) + "\n" +
				strings.Repeat(" ", 4),
		},
	} {
		encoder := BoardEncoder{
			encoder: data.fields.encoder,
			placeholders: data.fields.
				placeholders,
			margins:    data.fields.margins,
			stoneWidth: data.fields.stoneWidth,
		}
		got :=
			encoder.EncodeBoard(data.args.board)

		if got != data.want {
			test.Fail()
		}
	}
}
