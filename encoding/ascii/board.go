package ascii

import (
	"strings"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
)

// StoneEncoder ...
type StoneEncoder func(
	color models.Color,
) string

// Placeholders ...
type Placeholders struct {
	HorizontalLine string
	VerticalLine   string
	Crosshairs     string
}

// BoardEncoder ...
type BoardEncoder struct {
	encoder      StoneEncoder
	placeholders Placeholders
	margins      Margins
	stoneWidth   int
}

// NewBoardEncoder ...
func NewBoardEncoder(
	encoder StoneEncoder,
	placeholders Placeholders,
	margins Margins,
	stoneWidth int,
) BoardEncoder {
	return BoardEncoder{
		encoder:      encoder,
		placeholders: placeholders,
		margins:      margins,
		stoneWidth:   stoneWidth,
	}
}

// EncodeBoard ...
func (encoder BoardEncoder) EncodeBoard(
	board models.Board,
) string {
	stoneMargins := encoder.margins.Stone
	legendMargins := encoder.margins.Legend

	var rows []string
	var currentRow string
	points := board.Size().Points()
	for _, point := range points {
		if len(currentRow) == 0 {
			axis := sgf.EncodeAxis(point.Row)
			currentRow += encoder.wrapWithSpaces(
				string(axis),
				legendMargins.Row,
			)
		}

		var encodedStone string
		color, ok := board.Stone(point)
		if ok {
			encodedStone = encoder.encoder(color)
		} else {
			encodedStone =
				encoder.placeholders.Crosshairs
		}
		currentRow += encoder.wrapWithSpaces(
			encodedStone,
			stoneMargins.HorizontalMargins,
		)

		lastColumn := board.Size().Height - 1
		if point.Column == lastColumn {
			rows = append(rows, currentRow)
			currentRow = ""
		}
	}
	reverse(rows)

	var sparseRows []string
	for _, row := range rows {
		sparseRows = append(
			sparseRows,
			encoder.wrapWithEmptyLines(
				[]string{row},
				board.Size().Width,
				stoneMargins.VerticalMargins,
			)...,
		)
	}

	legendRow := encoder.spaces(
		legendMargins.Row.Width(1),
	)
	width := board.Size().Width
	for i := 0; i < width; i++ {
		axis := sgf.EncodeAxis(i)
		legendRow += encoder.wrapWithSpaces(
			string(axis),
			stoneMargins.HorizontalMargins,
		)
	}
	sparseRows = append(
		sparseRows,
		encoder.wrapWithEmptyLines(
			[]string{legendRow},
			board.Size().Width,
			legendMargins.Column,
		)...,
	)

	sparseRows = encoder.wrapWithEmptyLines(
		sparseRows,
		board.Size().Width,
		encoder.margins.Board,
	)

	return strings.Join(sparseRows, "\n")
}

func (encoder BoardEncoder) wrapWithSpaces(
	text string,
	margins HorizontalMargins,
) string {
	prefix := encoder.spaces(margins.Left)
	suffix := encoder.spaces(margins.Right)
	return prefix + text + suffix
}

func (encoder BoardEncoder) spaces(
	length int,
) string {
	return strings.Repeat(" ", length)
}

func (
	encoder BoardEncoder,
) wrapWithEmptyLines(
	lines []string,
	width int,
	margins VerticalMargins,
) []string {
	var wrappedLines []string
	wrappedLines = append(
		wrappedLines,
		encoder.emptyLines(
			margins.Top,
			width,
		)...,
	)
	wrappedLines = append(
		wrappedLines,
		lines...,
	)
	wrappedLines = append(
		wrappedLines,
		encoder.emptyLines(
			margins.Bottom,
			width,
		)...,
	)

	return wrappedLines
}

func (encoder BoardEncoder) emptyLines(
	count int,
	width int,
) []string {
	var lines []string
	for i := 0; i < count; i++ {
		line := encoder.emptyLine(width)
		lines = append(lines, line)
	}

	return lines
}

func (encoder BoardEncoder) emptyLine(
	width int,
) string {
	stoneMargins := encoder.margins.Stone
	legendMargins := encoder.margins.Legend

	line := encoder.spaces(
		legendMargins.Row.Width(1),
	)
	for i := 0; i < width; i++ {
		line += encoder.spaces(
			stoneMargins.
				Width(encoder.stoneWidth),
		)
	}

	return line
}

func reverse(strings []string) {
	left, right := 0, len(strings)-1
	for left < right {
		strings[left], strings[right] =
			strings[right], strings[left]
		left, right = left+1, right-1
	}
}
