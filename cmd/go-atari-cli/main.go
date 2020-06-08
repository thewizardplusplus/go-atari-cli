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

// nolint: gochecknoglobals
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

	baseWideMargins = ascii.Margins{
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
	wideStoneMargins = ascii.StoneMargins{
		HorizontalMargins: ascii.HorizontalMargins{
			Left: 1,
		},
		VerticalMargins: ascii.VerticalMargins{
			Bottom: 1,
		},
	}
	extraWideStoneMargins = ascii.StoneMargins{
		HorizontalMargins: ascii.HorizontalMargins{
			Left:  1,
			Right: 1,
		},
		VerticalMargins: ascii.VerticalMargins{
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
	return setTTYMode(colorsCodes[color]) + text + setTTYMode(0)
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
	storage models.StoneStorage,
	color models.Color,
	settings searchSettings,
) (models.Move, error) {
	generator := models.MoveGenerator{}

	randomSelector := selectors.RandomMoveSelector{}
	generalSelector := selectors.MaximalNodeSelector{
		NodeScorer: scorers.UCBScorer{
			Factor: ucbFactor,
		},
	}

	var simulator simulators.Simulator // nolint: staticcheck
	simulator = simulators.RolloutSimulator{
		MoveGenerator: generator,
		MoveSelector:  randomSelector,
	}
	if settings.parallelSimulator {
		simulator = simulators.ParallelSimulator{
			Simulator:   simulator,
			Concurrency: runtime.NumCPU(),
		}
	}

	var bulkySimulator builders.BulkySimulator
	if !settings.parallelBulkySimulator {
		bulkySimulator = bulky.FirstNodeSimulator{
			Simulator: simulator,
		}
	} else {
		bulkySimulator = bulky.AllNodesSimulator{
			Simulator: simulator,
		}
	}

	var builder builders.Builder
	terminator := terminators.NewGroupTerminator(
		terminators.NewPassTerminator(settings.maximalPass),
		terminators.NewTimeTerminator(time.Now, settings.maximalDuration),
	)
	builder = builders.IterativeBuilder{
		Builder: builders.TreeBuilder{
			NodeSelector:  generalSelector,
			MoveGenerator: generator,
			Simulator:     bulkySimulator,
		},
		Terminator: terminator,
	}
	if settings.parallelBuilder {
		builder = builders.ParallelBuilder{
			Builder:     builder,
			Concurrency: runtime.NumCPU(),
		}
	}

	root := &tree.Node{
		Move:    models.NewPreliminaryMove(color),
		Storage: storage,
	}
	searcher := searchers.MoveSearcher{
		MoveGenerator: generator,
		Builder:       builder,
		NodeSelector:  generalSelector,
	}
	node, err := searcher.SearchMove(root)
	if err != nil {
		return models.Move{}, err
	}

	return node.Move, nil
}

func check(storage models.StoneStorage, color models.Color) error {
	generator := models.MoveGenerator{}
	_, err := generator.LegalMoves(storage, models.NewPreliminaryMove(color))
	return err // don't wrap
}

func writePrompt(
	storageEncoder ascii.StoneStorageEncoder,
	storage models.StoneStorage,
	color models.Color,
	side climodels.Side,
) error {
	text := storageEncoder.EncodeStoneStorage(storage)
	fmt.Println(text)

	if err := check(storage, color); err != nil {
		return err // don't wrap
	}

	var mark string
	if side == climodels.Searcher {
		mark = "(searching) "
	}
	prompt := makePrompt(color, mark)
	fmt.Print(prompt) // don't break the line

	return nil
}

func makePrompt(color models.Color, data interface{}) string {
	prompt := ascii.EncodeColor(color)
	return fmt.Sprintf("%s> %v", prompt, data)
}

func readMove(
	reader *bufio.Reader,
	storageEncoder ascii.StoneStorageEncoder,
	storage models.StoneStorage,
	color models.Color,
	side climodels.Side,
) (models.Move, error) {
	if err := writePrompt(storageEncoder, storage, color, side); err != nil {
		return models.Move{}, err // don't wrap
	}

	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return models.Move{}, fmt.Errorf("unable to read the move: %s", err)
	}

	text = strings.TrimSuffix(text, "\n")
	point, err := sgf.DecodePoint(text)
	if err != nil {
		return models.Move{}, fmt.Errorf("unable to decode the point: %s", err)
	}

	move := models.Move{
		Color: color,
		Point: point,
	}
	if err := storage.CheckMove(move); err != nil {
		return models.Move{}, fmt.Errorf("incorrect move: %s", err)
	}

	return move, nil
}

func searchMove(
	storageEncoder ascii.StoneStorageEncoder,
	storage models.StoneStorage,
	color models.Color,
	side climodels.Side,
	settings searchSettings,
) (models.Move, error) {
	if err := writePrompt(storageEncoder, storage, color, side); err != nil {
		return models.Move{}, err // don't wrap
	}

	return search(storage, color, settings)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	storageInSGF := flag.String(
		"sgf",
		"",
		"board in Smart Game Format (default: empty board 5x5)",
	)
	humanColor := flag.String(
		"humanColor",
		"random",
		"human color (allowed: random, black, white)",
	)
	passes := flag.Int("passes", 1000, "building passes")
	duration := flag.Duration(
		"duration",
		10*time.Second,
		"building duration (e.g. 72h3m0.5s)",
	)
	parallelSimulator := flag.Bool(
		"parallelSimulator",
		false,
		"use parallel game simulating of a single node child",
	)
	parallelBulkySimulator := flag.Bool(
		"parallelBulkySimulator",
		false,
		"use parallel game simulating of all node children",
	)
	parallelBuilder := flag.Bool(
		"parallelBuilder",
		true,
		"use parallel tree building",
	)
	useUnicode := flag.Bool("unicode", true, "use Unicode to display stones")
	colorful := flag.Bool("colorful", true, "use colors to display stones")
	blackColor := flag.Int(
		"blackColor",
		34, // blue
		"SGR parameter for ANSI escape sequences for setting a color of black stones",
	)
	whiteColor := flag.Int(
		"whiteColor",
		31, // red
		"SGR parameter for ANSI escape sequences for setting a color of white stones",
	)
	wide := flag.Bool("wide", true, "display the board wide")
	grid := flag.Bool("grid", true, "display the board grid")
	flag.Parse()

	storage, err := sgf.DecodeStoneStorage(*storageInSGF, models.NewBoard)
	if err != nil {
		log.Fatal("unable to decode the board: ", err)
	}

	parsedHumanColor, err := ascii.DecodeColor(*humanColor)
	switch {
	case err == nil:
	case *humanColor == "random":
		if rand.Intn(2) == 0 {
			parsedHumanColor = models.Black
		} else {
			parsedHumanColor = models.White
		}
	default:
		log.Fatal("unable to decode the color: ", err)
	}

	var stoneEncoder ascii.StoneEncoder
	var placeholders ascii.Placeholders
	if *useUnicode {
		stoneEncoder = unicode.EncodeStone
		placeholders = unicodePlaceholders
	} else {
		stoneEncoder = func(color models.Color) string {
			return string(sgf.EncodeColor(color))
		}
		placeholders = asciiPlaceholders
	}
	if *colorful {
		baseStoneEncoder := stoneEncoder
		stoneEncoder = func(color models.Color) string {
			text := baseStoneEncoder(color)
			return colorize(text, color, colorCodeGroup{
				models.Black: *blackColor,
				models.White: *whiteColor,
			})
		}
	}
	if !*grid {
		placeholders.HorizontalLine = " "
		placeholders.VerticalLine = " "
	}

	var margins ascii.Margins
	if *wide {
		margins = baseWideMargins

		if *grid {
			margins.Stone = extraWideStoneMargins
		} else {
			margins.Stone = wideStoneMargins
		}
	}

	side := climodels.NewSide(parsedHumanColor)
	reader := bufio.NewReader(os.Stdin)
	storageEncoder := ascii.NewStoneStorageEncoder(
		stoneEncoder,
		placeholders,
		margins,
		1,
	)
	settings := searchSettings{
		maximalPass:            *passes,
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
			move, err = readMove(reader, storageEncoder, storage, currentColor, side)
		case climodels.Searcher:
			currentColor = parsedHumanColor.Negative()
			move, err = searchMove(storageEncoder, storage, currentColor, side, settings)
			if err == nil {
				text := sgf.EncodePoint(move.Point)
				fmt.Println(text)
			}
		}
		switch err {
		case nil:
		case models.ErrAlreadyLoss, models.ErrAlreadyWin:
			prompt := makePrompt(currentColor, err)
			fmt.Println(prompt)

			break loop
		default:
			log.Print("error: ", err)
			continue loop
		}

		storage = storage.ApplyMove(move)
		side = side.Invert()
	}
}
