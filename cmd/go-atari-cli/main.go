package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/thewizardplusplus/go-atari-cli/encoding/ascii"
	"github.com/thewizardplusplus/go-atari-cli/encoding/unicode"
	climodels "github.com/thewizardplusplus/go-atari-cli/models"
	models "github.com/thewizardplusplus/go-atari-models"
	"github.com/thewizardplusplus/go-atari-models/encoding/sgf"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders"
	"github.com/thewizardplusplus/go-atari-montecarlo/builders/terminators"
	"github.com/thewizardplusplus/go-atari-montecarlo/searchers"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors"
	"github.com/thewizardplusplus/go-atari-montecarlo/selectors/scorers"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators"
	"github.com/thewizardplusplus/go-atari-montecarlo/simulators/bulky"
	"github.com/thewizardplusplus/go-atari-montecarlo/tree"
)

const (
	ucbFactor = math.Sqrt2
)

var (
	asciiPlaceholders = ascii.Placeholders{
		HorizontalLine: "-",
		VerticalLine:   "|",
		Crosshairs:     "+",
	}
	unicodePlaceholders = ascii.Placeholders{
		HorizontalLine: "\u2500",
		VerticalLine:   "\u2502",
		Crosshairs:     "\u253c",
	}

	wideMargins = ascii.Margins{
		Stone: ascii.StoneMargins{
			HorizontalMargins: ascii.HorizontalMargins{
				Left: 1,
			},
			VerticalMargins: ascii.VerticalMargins{
				Bottom: 1,
			},
		},
		Legend: ascii.LegendMargins{
			Column: ascii.VerticalMargins{
				Top: 1,
			},
			Row: ascii.HorizontalMargins{
				Right: 1,
			},
		},
		Board: ascii.VerticalMargins{
			Top:    1,
			Bottom: 1,
		},
	}
)

type colorCodeGroup map[models.Color]int

func colorize(
	text string,
	color models.Color,
	colorsCodes colorCodeGroup,
) string {
	return setTTYMode(colorsCodes[color]) +
		text +
		setTTYMode(0)
}

func setTTYMode(mode int) string {
	return fmt.Sprintf("\x1b[%dm", mode)
}

type searchSettings struct {
	maximalPass            int
	maximalDuration        time.Duration
	parallelSimulator      bool
	parallelBulkySimulator bool
	parallelBuilder        bool
}

func search(
	board models.Board,
	color models.Color,
	settings searchSettings,
) (models.Move, error) {
	randomSelector :=
		selectors.RandomMoveSelector{}
	generalSelector :=
		selectors.MaximalNodeSelector{
			NodeScorer: scorers.UCBScorer{
				Factor: ucbFactor,
			},
		}

	var simulator simulators.Simulator
	simulator = simulators.RolloutSimulator{
		MoveSelector: randomSelector,
	}
	if settings.parallelSimulator {
		simulator =
			simulators.ParallelSimulator{
				Simulator:   simulator,
				Concurrency: runtime.NumCPU(),
			}
	}

	var bulkySimulator builders.BulkySimulator
	if !settings.parallelBulkySimulator {
		bulkySimulator =
			bulky.FirstNodeSimulator{
				Simulator: simulator,
			}
	} else {
		bulkySimulator =
			bulky.AllNodesSimulator{
				Simulator: simulator,
			}
	}

	var builder builders.Builder
	terminator :=
		terminators.NewGroupTerminator(
			terminators.NewPassTerminator(
				settings.maximalPass,
			),
			terminators.NewTimeTerminator(
				time.Now,
				settings.maximalDuration,
			),
		)
	builder = builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector: generalSelector,
			Simulator:    bulkySimulator,
		},
		Terminator: terminator,
	}
	if settings.parallelBuilder {
		builder = builders.ParallelBuilder{
			Builder:     builder,
			Concurrency: runtime.NumCPU(),
		}
	}

	previousMove :=
		models.NewPreliminaryMove(color)
	root := &tree.Node{
		Move:  previousMove,
		Board: board,
	}
	searcher := searchers.MoveSearcher{
		Builder:      builder,
		NodeSelector: generalSelector,
	}
	node, err := searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}

func check(
	board models.Board,
	color models.Color,
) error {
	previousMove :=
		models.NewPreliminaryMove(color)
	_, err := board.LegalMoves(previousMove)
	return err // don't wrap
}

func writePrompt(
	boardEncoder ascii.BoardEncoder,
	board models.Board,
	color models.Color,
	side climodels.Side,
) error {
	text := boardEncoder.EncodeBoard(board)
	fmt.Println(text)

	err := check(board, color)
	if err != nil {
		return err // don't wrap
	}

	var mark string
	if side == climodels.Searcher {
		mark = "(searching) "
	}

	prompt := makePrompt(color, mark)
	// don't break the line
	fmt.Print(prompt)

	return nil
}

func makePrompt(
	color models.Color,
	data interface{},
) string {
	prompt := ascii.EncodeColor(color)
	return fmt.Sprintf("%s> %v", prompt, data)
}

func readMove(
	reader *bufio.Reader,
	boardEncoder ascii.BoardEncoder,
	board models.Board,
	color models.Color,
	side climodels.Side,
) (models.Move, error) {
	err := writePrompt(
		boardEncoder,
		board,
		color,
		side,
	)
	if err != nil {
		return models.Move{}, err // don't wrap
	}

	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return models.Move{}, fmt.Errorf(
			"unable to read the move: %s",
			err,
		)
	}

	text = strings.TrimSuffix(text, "\n")
	point, err := sgf.DecodePoint(text)
	if err != nil {
		return models.Move{}, fmt.Errorf(
			"unable to decode the point: %s",
			err,
		)
	}

	move := models.Move{
		Color: color,
		Point: point,
	}
	err = board.CheckMove(move)
	if err != nil {
		return models.Move{}, fmt.Errorf(
			"incorrect move: %s",
			err,
		)
	}

	return move, nil
}

func searchMove(
	boardEncoder ascii.BoardEncoder,
	board models.Board,
	color models.Color,
	side climodels.Side,
	settings searchSettings,
) (models.Move, error) {
	err := writePrompt(
		boardEncoder,
		board,
		color,
		side,
	)
	if err != nil {
		return models.Move{}, err // don't wrap
	}

	return search(
		board,
		color,
		settings,
	)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	boardInSGF := flag.String(
		"sgf",
		"",
		"board in SGF "+
			"(default: empty board 5x5)",
	)
	humanColor := flag.String(
		"humanColor",
		"random",
		"human color "+
			"(allowed: random, black, white)",
	)
	pass := flag.Int(
		"pass",
		1000,
		"building pass",
	)
	duration := flag.Duration(
		"duration",
		10*time.Second,
		"building duration (e.g. 72h3m0.5s)",
	)
	parallelSimulator := flag.Bool(
		"parallelSimulator",
		false,
		"parallel simulator",
	)
	parallelBulkySimulator := flag.Bool(
		"parallelBulkySimulator",
		false,
		"parallel bulky simulator",
	)
	parallelBuilder := flag.Bool(
		"parallelBuilder",
		true,
		"parallel builder",
	)
	useUnicode := flag.Bool(
		"unicode",
		false,
		"use Unicode to display stones",
	)
	colorful := flag.Bool(
		"colorful",
		false,
		"use colors to display stones",
	)
	blackColor := flag.Int(
		"blackColor",
		34, // blue
		"SGR parameter "+
			"for ANSI escape sequences "+
			"for setting a color of black stones",
	)
	whiteColor := flag.Int(
		"whiteColor",
		31, // red
		"SGR parameter "+
			"for ANSI escape sequences "+
			"for setting a color of white stones",
	)
	wide := flag.Bool(
		"wide",
		false,
		"display the board wide",
	)
	grid := flag.Bool(
		"grid",
		true,
		"display the board grid",
	)
	flag.Parse()

	board, err := sgf.DecodeBoard(*boardInSGF)
	if err != nil {
		log.Fatal(
			"unable to decode the board: ",
			err,
		)
	}

	parsedHumanColor, err :=
		ascii.DecodeColor(*humanColor)
	switch {
	case err == nil:
	case *humanColor == "random":
		if rand.Intn(2) == 0 {
			parsedHumanColor = models.Black
		} else {
			parsedHumanColor = models.White
		}
	default:
		log.Fatal(
			"unable to decode the color: ",
			err,
		)
	}

	var stoneEncoder ascii.StoneEncoder
	var placeholders ascii.Placeholders
	if *useUnicode {
		stoneEncoder = unicode.EncodeStone
		placeholders = unicodePlaceholders
	} else {
		stoneEncoder = func(
			color models.Color,
		) string {
			symbol := sgf.EncodeColor(color)
			return string(symbol)
		}
		placeholders = asciiPlaceholders
	}
	if *colorful {
		baseStoneEncoder := stoneEncoder
		stoneEncoder = func(
			color models.Color,
		) string {
			text := baseStoneEncoder(color)
			return colorize(
				text,
				color,
				colorCodeGroup{
					models.Black: *blackColor,
					models.White: *whiteColor,
				},
			)
		}
	}
	if !*grid {
		placeholders.HorizontalLine = " "
		placeholders.VerticalLine = " "
	}

	var margins ascii.Margins
	if *wide {
		margins = wideMargins
	}

	side :=
		climodels.NewSide(parsedHumanColor)
	reader := bufio.NewReader(os.Stdin)
	boardEncoder := ascii.NewBoardEncoder(
		stoneEncoder,
		placeholders,
		margins,
		1,
	)
	settings := searchSettings{
		maximalPass:            *pass,
		maximalDuration:        *duration,
		parallelSimulator:      *parallelSimulator,
		parallelBulkySimulator: *parallelBulkySimulator,
		parallelBuilder:        *parallelBuilder,
	}
loop:
	for {
		var currentColor models.Color
		var move models.Move
		var err error
		switch side {
		case climodels.Human:
			currentColor = parsedHumanColor
			move, err = readMove(
				reader,
				boardEncoder,
				board,
				currentColor,
				side,
			)
		case climodels.Searcher:
			currentColor =
				parsedHumanColor.Negative()
			move, err = searchMove(
				boardEncoder,
				board,
				currentColor,
				side,
				settings,
			)
			if err == nil {
				text := sgf.EncodePoint(move.Point)
				fmt.Println(text)
			}
		}
		switch err {
		case nil:
		case models.ErrAlreadyLoss,
			models.ErrAlreadyWin:
			prompt :=
				makePrompt(currentColor, err)
			fmt.Println(prompt)

			break loop
		default:
			log.Print("error: ", err)
			continue loop
		}

		board = board.ApplyMove(move)
		side = side.Invert()
	}
}
