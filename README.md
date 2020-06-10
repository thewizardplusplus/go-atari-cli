# go-atari-cli

[![GoDoc](https://godoc.org/github.com/thewizardplusplus/go-atari-cli?status.svg)](https://godoc.org/github.com/thewizardplusplus/go-atari-cli)
[![Go Report Card](https://goreportcard.com/badge/github.com/thewizardplusplus/go-atari-cli)](https://goreportcard.com/report/github.com/thewizardplusplus/go-atari-cli)
[![Build Status](https://travis-ci.org/thewizardplusplus/go-atari-cli.svg?branch=master)](https://travis-ci.org/thewizardplusplus/go-atari-cli)
[![codecov](https://codecov.io/gh/thewizardplusplus/go-atari-cli/branch/master/graph/badge.svg)](https://codecov.io/gh/thewizardplusplus/go-atari-cli)

The [Atari Go](https://senseis.xmp.net/?AtariGo) program with a terminal-based interface.

_**Disclaimer:** this program was written directly on an Android smartphone with the AnGoIde IDE._

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
