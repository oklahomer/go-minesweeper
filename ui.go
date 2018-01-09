package minesweeper

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	// ErrInvalidInput is returned when user input is invalid.
	ErrInvalidInput = errors.New("invalid input is given")
)

// UI defines an interface to output user friendly representation of a game and receive user input for operation.
type UI interface {
	// Render outputs user friendly representation of a game.
	Render(*Field) string

	// ParseInput receives user input and converts into OpType and Coordinate.
	ParseInput(string) (OpType, *Coordinate, error)
}

type defaultUI struct {
	// [1, 2, 3, 4, ...]
	xSymbols []int

	// [a, b, c, ...., aa, ab, ...]
	ySymbols []string
}

func (r *defaultUI) Render(field *Field) string {
	if len(r.xSymbols) == 0 || len(r.ySymbols) == 0 {
		r.initSymbols(field.Width, field.Height)
	}

	yWidth := len(r.ySymbols[len(r.ySymbols)-1])

	str := ""
	for i := 0; i < yWidth; i++ {
		str += " "
	}

	for _, symbol := range r.xSymbols {
		str += fmt.Sprintf(" %d", symbol)
	}
	str += "\n"

	for i, row := range field.Cells {
		str += r.ySymbols[i]
		for _, cell := range row {
			str += fmt.Sprintf("|%s", dispState(cell.State()))
		}
		if i+1 < field.Height {
			str += "\n"
		}
	}

	return str
}

func (r *defaultUI) ParseInput(str string) (OpType, *Coordinate, error) {
	fields := strings.Fields(str)
	fieldsCnt := len(fields)
	if fieldsCnt != 2 && fieldsCnt != 3 {
		return 0, nil, ErrInvalidInput
	}

	x, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, nil, ErrInvalidInput
	}

	var foundX bool
	xCoord := 0
	for i, v := range r.xSymbols {
		if x == v {
			foundX = true
			xCoord = i
		}
	}
	if !(foundX) {
		return 0, nil, ErrInvalidInput
	}

	var foundY bool
	yCoord := 0
	for i, v := range r.ySymbols {
		if fields[1] == v {
			foundY = true
			yCoord = i
		}
	}
	if !(foundY) {
		return 0, nil, ErrInvalidInput
	}

	coord := &Coordinate{X: xCoord, Y: yCoord}

	if fieldsCnt == 2 {
		return Open, coord, nil
	}

	switch strings.ToLower(fields[2]) {
	case "f", "flag":
		return Flag, coord, nil

	case "u", "unflag":
		return Unflag, coord, nil

	default:
		return 0, nil, ErrInvalidInput

	}
}

func (r *defaultUI) initSymbols(width int, height int) {
	r.xSymbols = make([]int, width)
	for i := 0; i < width; i++ {
		r.xSymbols[i] = i + 1
	}

	r.ySymbols = make([]string, height)
	candidates := [...]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	candidatesN := len(candidates)

	for i := 0; i < height; i++ {
		n := i + 1
		for n > 0 {
			n -= 1
			r.ySymbols[i] = candidates[n%candidatesN] + r.ySymbols[i]
			n = int(math.Floor(float64(n) / float64(candidatesN)))
		}
	}
}

func dispState(s CellState) string {
	switch s {
	case Closed:
		return " "

	case Opened:
		return "-"

	case Flagged:
		return "F"

	case Exploded:
		return "X"

	default:
		panic("invalid state")

	}
}
