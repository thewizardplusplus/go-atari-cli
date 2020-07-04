# go-atari-cli

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-atari-cli?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-atari-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-atari-cli)](https://goreportcard.com/report/github.com/thewizardplusplus/go-atari-cli)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-atari-cli.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-atari-cli)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-atari-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-atari-cli)

The [Atari Go](https://senseis.xmp.net/?AtariGo) program with a terminal-based interface.

_**Disclaimer:** this program was written directly on an Android smartphone with the AnGoIde IDE._

## Features

- displaying a board:
  - by symbols (to choose):
    - ASCII;
    - Unicode;
  - by colors (to choose):
    - monochrome;
    - colorful;
  - by size (to choose):
    - terse;
    - wide;
  - by decor (to choose):
    - without the board grid;
    - with the board grid;
  - misc.:
    - marking searching process;
- interacting via text commands (moves in [Smart Game Format](https://senseis.xmp.net/?SGF));
- options:
  - initial position in [Smart Game Format](https://senseis.xmp.net/?SGF);
  - human color (i.e. a computer can move first):
    - support automatic random selecting (optional);
  - move searching restrictions:
    - passes of tree building;
    - duration of tree building;
  - optimization via parallel move searching:
    - parallel game simulating:
      - of a single node child;
      - of all node children;
    - parallel tree building;
  - displaying:
    - switching between ASCII/Unicode modes;
    - switching between monochrome/colorful modes;
    - switching between terse/wide modes;
    - switching between modes without/with the board grid.

## Installation

```
$ go get github.com/thewizardplusplus/go-atari-cli/...
```

## Usage

```
$ go-atari-cli -h | -help | --help
$ go-atari-cli [options]
```

Options:

- `-h`, `-help`, `--help` &mdash; show the help message and exit;
- `-sgf STRING` &mdash; board in [Smart Game Format](https://senseis.xmp.net/?SGF) (default: empty board 5x5);
- `-humanColor {random|black|white}` &mdash; human color (default: `random`);
- `-passes INTEGER` &mdash; building passes (default: `1000`);
- `-duration DURATION` &mdash; building duration (e.g. `72h3m0.5s`; default: `10s`);
- `-parallelSimulator` &mdash; use parallel game simulating of a single node child (default: `false`; for inverting use `-parallelSimulator` or `-parallelSimulator=true`);
- `-parallelBulkySimulator` &mdash; use parallel game simulating of all node children (default: `false`; for inverting use `-parallelBulkySimulator` or `-parallelBulkySimulator=true`);
- `-parallelBuilder` &mdash; use parallel tree building (default: `true`; for inverting use `-parallelBuilder=false`);
- `-unicode` &mdash; use Unicode to display stones (default: `true`; for inverting use `-unicode=false`);
- `-colorful` &mdash; use colors to display stones (default: `true`; for inverting use `-colorful=false`);
- `-blackColor INTEGER` &mdash; SGR parameter for ANSI escape sequences for setting a color of black stones (default: `34`; see for details: https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit);
- `-whiteColor INTEGER` &mdash; SGR parameter for ANSI escape sequences for setting a color of white stones (default: `31`; see for details: https://en.wikipedia.org/wiki/ANSI_escape_code#3/4_bit);
- `-wide` &mdash; display the board wide (default: `true`; for inverting use `-wide=false`);
- `-grid` &mdash; display the board grid (default: `true`; for inverting use `-grid=false`).

## Examples

`ascii.DecodeColor()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/ascii"
)

func main() {
	color, _ := ascii.DecodeColor("white")
	fmt.Printf("%v\n", color)

	// Output: 1
}
```

`ascii.EncodeColor()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-atari-models"
)

func main() {
	color := ascii.EncodeColor(models.White)
	fmt.Printf("%v\n", color)

	// Output: white
}
```

`ascii.StoneStorageEncoder.EncodeStoneStorage()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
)

func main() {
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
```

`ascii.StoneStorageEncoder.EncodeStoneStorage()` with margins:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
)

func main() {
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
```

`ascii.StoneStorageEncoder.EncodeStoneStorage()` with a grid:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/ascii"
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
)

func main() {
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
```

`unicode.EncodeStone()`:

```go
package main

import (
	"fmt"

	"github.com/thewizardplusplus/go-atari-cli/encoding/unicode"
	models "github.com/thewizardplusplus/go-atari-models"
)

func main() {
	stone := unicode.EncodeStone(models.White)
	fmt.Printf("%v\n", stone)

	// Output: ○
}
```

## License

The MIT License (MIT)

Copyright &copy; 2020 thewizardplusplus
