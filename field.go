package minesweeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"math/rand"
)

var (
	// ErrCoordinateOutOfRange is returned when given *Coordinate points to a non-existing field location.
	ErrCoordinateOutOfRange = errors.New("invalid coordinate is given")
)

// FieldConfig contains some configuration variables for Field.
type FieldConfig struct {
	Width   int `json:"width" yaml:"width"`
	Height  int `json:"height" yaml:"height"`
	MineCnt int `json:"mine_count" yaml:"mine_count"`
}

// NewFieldConfig construct FieldConfig with default values.
// Use json.Unmarshal, yaml.Unmarshal or manual manipulation to override default values.
func NewFieldConfig() *FieldConfig {
	return &FieldConfig{
		Width:   9,
		Height:  9,
		MineCnt: 10,
	}
}

func validateConfig(config *FieldConfig) error {
	if config.Width <= 0 {
		return errors.New("field width is zero")
	}

	if config.Height <= 0 {
		return errors.New("field height is zero")
	}

	if config.MineCnt <= 0 {
		return errors.New("mine count is zero")
	}

	if (config.Width * config.Height) <= config.MineCnt {
		return errors.New("too many mines")
	}

	return nil
}

// Field represents a minefield with given width and height.
// This is merely a representation of minefield, so the state of a game is not part of this.
type Field struct {
	Width  int
	Height int
	Cells  [][]Cell
}

// NewField construct a Field with given configuration.
func NewField(config *FieldConfig) (*Field, error) {
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalild config is given: %s", err.Error())
	}

	grid := func() [][]bool {
		n := config.Width * config.Height
		mines := make([]bool, n)
		for _, v := range rand.Perm(n)[:config.MineCnt] {
			mines[v] = true
		}

		grid := make([][]bool, config.Height)
		for i := 0; i < config.Height; i++ {
			start := i * config.Width
			grid[i] = mines[start : start+config.Width]
		}
		return grid
	}()

	cells := make([][]Cell, config.Height)
	for i, row := range grid {
		cells[i] = make([]Cell, config.Width)

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

				if ii+1 < config.Width && above[ii+1] {
					surroundingCnt++
				}
			}

			if ii > 0 && row[ii-1] {
				surroundingCnt++
			}

			if ii+1 < config.Width && row[ii+1] {
				surroundingCnt++
			}

			if i+1 < config.Height {
				below := grid[i+1]
				if ii > 0 && below[ii-1] {
					surroundingCnt++
				}

				if below[ii] {
					surroundingCnt++
				}

				if ii+1 < config.Width && below[ii+1] {
					surroundingCnt++
				}
			}

			cells[i][ii] = newCell(hasMine, surroundingCnt)
		}
	}

	return &Field{
		Width:  config.Width,
		Height: config.Height,
		Cells:  cells,
	}, nil
}

// Open receives a Coordinate, locate a corresponding cell, and opens it.
// If surrounding cells has no underlying mine, all surrounding cells are recursively opened.
//
// Below errors may be returned:
// - ErrCoordinateOutOfRange ... there is not corresponding cell
// - ErrOpeningOpenedCell ... the target cell is already opened
// - ErrOpeningFlaggedCell ... the target cell is currently flagged and needs to be unflagged before this operation
// - ErrOpeningExplodedCell ... the target cell's underlying mine is already exploded
func (f *Field) Open(coord *Coordinate) (*Result, error) {
	x := coord.X
	y := coord.Y

	if x+1 > f.Width || y+1 > f.Height {
		return nil, ErrCoordinateOutOfRange
	}

	target := f.Cells[y][x]
	result, err := target.open()
	if err != nil {
		return nil, err
	}

	if result.NewState == Exploded {
		return result, nil
	}

	f.openSurroundings(coord)

	return result, nil
}

func (f *Field) openSurroundings(coord *Coordinate) {
	origin := f.Cells[coord.Y][coord.X]
	if origin.SurroundingCnt() > 0 {
		// At least one surrounding cell has a mine.
		// Do not automatically open all surrounding cells.
		return
	}

	// All surrounding cells are safe to open.
	for _, c := range f.getSurroundingCoordinates(coord) {
		target := f.Cells[c.Y][c.X]

		// Don't open when state is Flagged.
		// And to avoid opening a particular cell multiple times, proceed to open when state is not "Closed."
		if target.State() != Closed {
			continue
		}

		target.open()
		f.openSurroundings(c)
	}
}

// Flag receives a Coordinate, locate a corresponding cell, and flag it to indicate possible underlying mine.
//
// Below errors may be returned:
// - ErrCoordinateOutOfRange ... there is not corresponding cell
// - ErrFlaggingOpenedCell ... the target cell is already opened
// - ErrFlaggingFlaggedCell ... the target cell is already flagged
// - ErrFlaggingExplodedCell ... the target cell's underlying mine is already exploded
func (f *Field) Flag(coord *Coordinate) (*Result, error) {
	x := coord.X
	y := coord.Y

	if x+1 > f.Width || y+1 > f.Height {
		return nil, ErrCoordinateOutOfRange
	}

	return f.Cells[y][x].flag()
}

// Unflag receives a Coordinate, locate a corresponding cell, and flag it to indicate possible underlying mine.
//
// Below errors may be returned:
// - ErrCoordinateOutOfRange ... there is not corresponding cell
// - ErrUnflaggingNonFlaggedCell ... the target cell is not currently flagged
func (f *Field) Unflag(coord *Coordinate) (*Result, error) {
	x := coord.X
	y := coord.Y

	if x+1 > f.Width || y+1 > f.Height {
		return nil, ErrCoordinateOutOfRange
	}

	return f.Cells[y][x].unflag()
}

// MarshalJSON returns JSON representation of Field.
func (f *Field) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	m["width"] = f.Width
	m["height"] = f.Height
	cells := make([][]interface{}, f.Height)
	for i, row := range f.Cells {
		for _, c := range row {
			cells[i] = append(cells[i], map[string]interface{}{
				"state":             c.State().String(),
				"has_mine":          c.hasMine(),
				"surrounding_count": c.SurroundingCnt(),
			})
		}
	}
	m["cells"] = cells
	return json.Marshal(m)
}

// UnmarshalJSON converts given input to Field instance.
func (f *Field) UnmarshalJSON(b []byte) error {
	res := gjson.ParseBytes(b)

	// Set width
	widthValue := res.Get("width")
	if !widthValue.Exists() {
		return errors.New(`"width" field is not given`)
	}
	f.Width = int(widthValue.Int())

	// Set height
	heightValue := res.Get("height")
	if !heightValue.Exists() {
		return errors.New(`"height" field is not given`)
	}
	f.Height = int(heightValue.Int())

	// Set cells
	cellsValue := res.Get("cells")
	if !cellsValue.Exists() {
		return errors.New(`"cells" field is not given`)
	}
	f.Cells = make([][]Cell, f.Height)
	for i, row := range cellsValue.Array() {
		cells := make([]Cell, f.Width)
		for ii, c := range row.Array() {
			stateValue := c.Get("state")
			if !stateValue.Exists() {
				return errors.New(`"state" field is not given`)
			}

			mineValue := c.Get("has_mine")
			if !mineValue.Exists() {
				return errors.New(`"has_mine" field is not given`)
			}

			cntValue := c.Get("surrounding_count")
			if !cntValue.Exists() {
				return errors.New(`"surrounding_count" field is not given`)
			}

			state, err := strToCellState(stateValue.String())
			if err != nil {
				return fmt.Errorf("failed to convert given state value: %s", err.Error())
			}
			cells[ii] = &cell{
				state:          state,
				mine:           mineValue.Bool(),
				surroundingCnt: int(cntValue.Int()),
			}
		}
		f.Cells[i] = cells
	}

	// O.K.
	return nil
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

// Coordinate represents a coordinate of specific location on the Field.
type Coordinate struct {
	X int
	Y int
}

// Result represents a result of given action.
type Result struct {
	NewState CellState
}
