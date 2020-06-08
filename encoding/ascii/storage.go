package ascii

import (
	"strings"

	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
)

// StoneEncoder ...
type StoneEncoder func(color models.Color) string

// Placeholders ...
type Placeholders struct {
	HorizontalLine string
	VerticalLine   string
	Crosshairs     string
}

// StoneStorageEncoder ...
type StoneStorageEncoder struct {
	encoder      StoneEncoder
	placeholders Placeholders
	margins      Margins
	stoneWidth   int
}

// NewStoneStorageEncoder ...
func NewStoneStorageEncoder(
	encoder StoneEncoder,
	placeholders Placeholders,
	margins Margins,
	stoneWidth int,
) StoneStorageEncoder {
	return StoneStorageEncoder{
		encoder:      encoder,
		placeholders: placeholders,
		margins:      margins,
		stoneWidth:   stoneWidth,
	}
}

// EncodeStoneStorage ...
func (encoder StoneStorageEncoder) EncodeStoneStorage(
	storage models.StoneStorage,
) string {
	stoneMargins, legendMargins := encoder.margins.Stone, encoder.margins.Legend

	var rows []string
	var currentRow string
	for _, point := range storage.Size().Points() {
		if len(currentRow) == 0 {
			currentRow += encoder.wrapWithSpaces(
				string(sgf.EncodeAxis(point.Row)),
				legendMargins.Row,
			)
		}

		var encodedStone string
		if color, ok := storage.Stone(point); ok {
			encodedStone = encoder.encoder(color)
		} else {
			encodedStone = encoder.placeholders.Crosshairs
		}
		currentRow += encoder.wrapWithSpaces(
			encodedStone,
			stoneMargins.HorizontalMargins,
			encoder.placeholders.HorizontalLine,
		)

		if lastColumn := storage.Size().Height - 1; point.Column == lastColumn {
			rows = append(rows, currentRow)
			currentRow = ""
		}
	}
	reverse(rows)

	var sparseRows []string
	for _, row := range rows {
		sparseRows = append(sparseRows, encoder.wrapWithEmptyLines(
			[]string{row},
			storage.Size().Width,
			stoneMargins.VerticalMargins,
			encoder.placeholders.VerticalLine,
		)...)
	}

	legendRow := encoder.spaces(legendMargins.Row.Width(1))
	for i := 0; i < storage.Size().Width; i++ {
		legendRow += encoder.wrapWithSpaces(
			string(sgf.EncodeAxis(i)),
			stoneMargins.HorizontalMargins,
		)
	}
	sparseRows = append(sparseRows, encoder.wrapWithEmptyLines(
		[]string{legendRow},
		storage.Size().Width,
		legendMargins.Column,
	)...)

	sparseRows = encoder.wrapWithEmptyLines(
		sparseRows,
		storage.Size().Width,
		encoder.margins.Board,
	)

	return strings.Join(sparseRows, "\n")
}

func (encoder StoneStorageEncoder) wrapWithSpaces(
	text string,
	margins HorizontalMargins,
	optionalSymbol ...string,
) string {
	prefix := encoder.spaces(margins.Left, optionalSymbol...)
	suffix := encoder.spaces(margins.Right, optionalSymbol...)
	return prefix + text + suffix
}

func (encoder StoneStorageEncoder) spaces(
	count int,
	optionalSymbol ...string,
) string {
	var symbol string
	if len(optionalSymbol) != 0 {
		symbol = optionalSymbol[0]
	} else {
		symbol = " "
	}
	return strings.Repeat(symbol, count)
}

func (encoder StoneStorageEncoder) wrapWithEmptyLines(
	lines []string,
	width int,
	margins VerticalMargins,
	optionalSeparator ...string,
) []string {
	var wrappedLines []string
	wrappedLines = append(wrappedLines, encoder.emptyLines(
		margins.Top,
		width,
		optionalSeparator...,
	)...)
	wrappedLines = append(wrappedLines, lines...)
	wrappedLines = append(wrappedLines, encoder.emptyLines(
		margins.Bottom,
		width,
		optionalSeparator...,
	)...)

	return wrappedLines
}

func (encoder StoneStorageEncoder) emptyLines(
	count int,
	width int,
	optionalSeparator ...string,
) []string {
	var lines []string
	for i := 0; i < count; i++ {
		line := encoder.emptyLine(width, optionalSeparator...)
		lines = append(lines, line)
	}

	return lines
}

func (encoder StoneStorageEncoder) emptyLine(
	width int,
	optionalSeparator ...string,
) string {
	stoneMargins, legendMargins := encoder.margins.Stone, encoder.margins.Legend

	var separator string
	if len(optionalSeparator) != 0 {
		separator = optionalSeparator[0]
	} else {
		separator = " "
	}

	line := encoder.spaces(legendMargins.Row.Width(1))
	for i := 0; i < width; i++ {
		line += encoder.spaces(stoneMargins.Left) +
			separator +
			encoder.spaces(stoneMargins.Right)
	}

	return line
}

func reverse(strings []string) {
	left, right := 0, len(strings)-1
	for left < right {
		strings[left], strings[right] = strings[right], strings[left]
		left, right = left+1, right-1
	}
}
