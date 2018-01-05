package minesweeper

import (
	"errors"
	"fmt"
	"math/rand"
)

var (
	ErrOpeningOpenedCell        = errors.New("opened cell can not be opened")
	ErrOpeningFlaggedCell       = errors.New("flagged cell can not be opened")
	ErrOpeningExplodedCell      = errors.New("exploded cell can not be opened")
	ErrFlaggingOpenedCell       = errors.New("opened cell can not be flagged")
	ErrFlaggingFlaggedCell      = errors.New("flagged cell can not be re-flagged")
	ErrFlaggingExplodedCell     = errors.New("exploded cell can not be flagged")
	ErrUnflaggingNonFlaggedCell = errors.New("non-flagged cell can not be unflagged")
	ErrCoordinateOutOfRange     = errors.New("invalid coordinate is given")
)

type Config struct {
	FieldWidth  int `json:"field_width" yaml:"field_width"`
	FieldHeight int `json:"field_height" yaml:"field_height"`
	MineCnt     int `json:"mine_count" yaml:"mine_count"`
}

func NewConfig() *Config {
	return &Config{
		FieldWidth:  9,
		FieldHeight: 9,
		MineCnt:     10,
	}
}

func validateConfig(config *Config) error {
	if config.FieldWidth <= 0 {
		return errors.New("field width is zero")
	}

	if config.FieldHeight <= 0 {
		return errors.New("field height is zero")
	}

	if config.MineCnt <= 0 {
		return errors.New("mine count is zero")
	}

	if (config.FieldWidth * config.FieldHeight) <= config.MineCnt {
		return errors.New("too many mines")
	}

	return nil
}

type Field struct {
	Width  int
	Height int
	Cells  [][]Cell
}

func NewField(config *Config) (*Field, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalild config is given: %s", err.Error())
	}

	grid := func() [][]bool {
		n := config.FieldWidth * config.FieldHeight
		mines := make([]bool, n)
		for _, v := range rand.Perm(n)[:config.MineCnt] {
			mines[v] = true
		}

		grid := make([][]bool, config.FieldHeight)
		for i := 0; i < config.FieldHeight; i++ {
			start := i * config.FieldWidth
			grid[i] = mines[start : start+config.FieldWidth]
		}
		return grid
	}()

	cells := make([][]Cell, config.FieldHeight)
	for i, row := range grid {
		cells[i] = make([]Cell, config.FieldWidth)

		for ii, hasMine := range row {
			var surroundingCnt int

			if i > 0 {
				above := grid[i-1]
				if ii > 0 && above[ii-1] {
					surroundingCnt++
				}

				if above[ii] {
					surroundingCnt++
				}

				if ii+1 < config.FieldWidth && above[ii+1] {
					surroundingCnt++
				}
			}

			if ii > 0 && row[ii-1] {
				surroundingCnt++
			}

			if ii+1 < config.FieldWidth && row[ii+1] {
				surroundingCnt++
			}

			if i+1 < config.FieldHeight {
				below := grid[i+1]
				if ii > 0 && below[ii-1] {
					surroundingCnt++
				}

				if below[ii] {
					surroundingCnt++
				}

				if ii+1 < config.FieldWidth && below[ii+1] {
					surroundingCnt++
				}
			}

			cells[i][ii] = newCell(hasMine, surroundingCnt)
		}
	}

	return &Field{
		Width:  config.FieldWidth,
		Height: config.FieldHeight,
		Cells:  cells,
	}, nil
}

func (f *Field) Open(coord *Coordinate) (*Result, error) {
	x := coord.X
	y := coord.Y

	if x+1 > f.Width || y+1 > f.Height {
		return nil, ErrCoordinateOutOfRange
	}

	row := f.Cells[y]
	cell := row[x]

	result, err := cell.open()
	if err != nil {
		return nil, err
	}

	if result.NewState == Exploded {
		return result, nil
	}

	if cell.SurroundingCnt() == 0 {
		for _, c := range f.getSurroundingCoordinates(coord) {
			r := f.Cells[c.Y]
			target := r[c.X]
			if target.State() == Closed {
				f.Open(c)
			}
		}
	}

	return result, nil
}

func (f *Field) Flag(coord *Coordinate) (*Result, error) {
	x := coord.X
	y := coord.Y

	if x+1 > f.Width || y+1 > f.Height {
		return nil, ErrCoordinateOutOfRange
	}

	return f.Cells[y][x].flag()
}

func (f *Field) Unflag(coord *Coordinate) (*Result, error) {
	x := coord.X
	y := coord.Y

	if x+1 > f.Width || y+1 > f.Height {
		return nil, ErrCoordinateOutOfRange
	}

	return f.Cells[y][x].unflag()
}

func (f *Field) getSurroundingCoordinates(coord *Coordinate) []*Coordinate {
	x := coord.X
	y := coord.Y

	var coords []*Coordinate
	// Above row
	if y > 0 {
		if x > 1 {
			coords = append(coords, &Coordinate{X: x - 1, Y: y - 1})
		}

		coords = append(coords, &Coordinate{X: x, Y: y - 1})

		if x+1 < f.Width {
			coords = append(coords, &Coordinate{X: x + 1, Y: y - 1})
		}
	}

	// Current row
	if x > 0 {
		coords = append(coords, &Coordinate{X: x - 1, Y: y})
	}

	if x+1 < f.Width {
		coords = append(coords, &Coordinate{X: x + 1, Y: y})
	}

	// Below row
	if y+1 < f.Height {
		if x > 1 {
			coords = append(coords, &Coordinate{X: x - 1, Y: y + 1})
		}

		coords = append(coords, &Coordinate{X: x, Y: y + 1})

		if x+1 < f.Width {
			coords = append(coords, &Coordinate{X: x + 1, Y: y + 1})
		}
	}

	return coords
}

type Coordinate struct {
	X int
	Y int
}

type Result struct {
	NewState State
}
