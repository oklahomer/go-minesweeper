package minesweeper

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

var (
	ErrInvalidInput = errors.New("invalid input is given")
)

type UI interface {
	Render(*Field) string
	ParseInput(string) (*Coordinate, error)
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

func (r *defaultUI) ParseInput(str string) (*Coordinate, error) {
	fields := strings.Fields(str)
	if len(fields) != 2 {
		return nil, ErrInvalidInput
	}

	x, err := strconv.Atoi(fields[0])
	if err != nil {
		return nil, ErrInvalidInput
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
		return nil, ErrInvalidInput
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
		return nil, ErrInvalidInput
	}

	return &Coordinate{X: xCoord, Y: yCoord}, nil
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
