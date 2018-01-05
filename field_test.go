package minesweeper

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.FieldWidth == 0 {
		t.Errorf("Config.FieldWidth is not set.")
	}

	if config.FieldHeight == 0 {
		t.Errorf("Config.FieldHeight is not set.")
	}

	if config.MineCnt == 0 {
		t.Errorf("Config.MineCnt is not set.")
	}
}

func TestNewField(t *testing.T) {
	var configs = []*Config{
		{
			FieldWidth:  12,
			FieldHeight: 0,
			MineCnt:     9,
		},
		{
			FieldWidth:  0,
			FieldHeight: 12,
			MineCnt:     9,
		},
		{
			FieldWidth:  12,
			FieldHeight: 12,
			MineCnt:     0,
		},
		{
			FieldWidth:  12,
			FieldHeight: 12,
			MineCnt:     9,
		},
		{
			FieldWidth:  2,
			FieldHeight: 2,
			MineCnt:     10,
		},
	}

	for i, config := range configs {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			field, err := NewField(config)

			if config.FieldWidth == 0 || config.FieldHeight == 0 || config.MineCnt == 0 {
				if err == nil {
					t.Fatal("Error is not returned on invalid *Config.")
				}

				return
			}

			if config.MineCnt >= (config.FieldWidth * config.FieldHeight) {
				if err == nil {
					t.Fatal("Error is not returned on invalid *Config.")
				}

				return
			}

			if field == nil {
				t.Fatal("Field is nil.")
			}

			mineCnt := 0
			for _, row := range field.Cells {
				for _, c := range row {
					if c.hasMine() {
						mineCnt++
					}
				}
			}
			if config.MineCnt != mineCnt {
				t.Errorf("Expected mine count of %d, but was %d.", config.MineCnt, mineCnt)
			}
		})
	}
}

func TestField_Flag(t *testing.T) {
	type test struct {
		field    *Field
		coord    *Coordinate
		expected [][]Cell
	}

	tests := []*test{
		// Only left top corner has a mine and right bottom is opened.
		{
			field: &Field{
				Width:  2,
				Height: 2,
				Cells: [][]Cell{
					{
						&cell{state: Closed},
						&cell{state: Closed},
					},
					{
						&cell{state: Closed},
						&cell{state: Closed},
					},
				},
			},
			coord: &Coordinate{X: 1, Y: 1},
			expected: [][]Cell{
				{
					&cell{state: Closed},
					&cell{state: Closed},
				},
				{
					&cell{state: Closed},
					&cell{state: Flagged},
				},
			},
		},

		// Invalid coordinate is given
		{
			field: &Field{Width: 3, Height: 3},
			coord: &Coordinate{X: 1, Y: 100},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			// See if given coordinate is valid
			if test.coord.X+1 > test.field.Width || test.coord.Y+1 > test.field.Height {
				_, err := test.field.Flag(test.coord)
				if err == nil || err != ErrCoordinateOutOfRange {
					t.Fatalf("Expected error is not returned: %s", err)
				}

				return
			}

			target := test.field.Cells[test.coord.Y][test.coord.X]
			oldStatus := target.State()

			result, err := test.field.Flag(test.coord)

			if oldStatus == Flagged {
				if err == nil {
					t.Fatal("Error should be returned when flagged cell is subject to flag.")
				} else if err != ErrFlaggingFlaggedCell {
					t.Fatal("ErrFlaggingFlaggedCell should be returned when flagged cell is subject to flag.")
				}

				return

			}

			if oldStatus == Closed && result.NewState != Flagged {
				t.Fatalf("Unexpected state is returned: %s", result.NewState)
			}

			for i, row := range test.field.Cells {
				for ii, cell := range row {
					if cell.State() != test.expected[i][ii].State() {
						t.Errorf("Cell with unexpected state is retuned. X: %d, Y: %d. State: %s", i, ii, cell.State())
					}
				}
			}
		})
	}
}

func TestField_Unflag(t *testing.T) {
	type test struct {
		field    *Field
		coord    *Coordinate
		expected [][]Cell
	}

	tests := []*test{
		{
			field: &Field{
				Width:  2,
				Height: 2,
				Cells: [][]Cell{
					{
						&cell{state: Closed},
						&cell{state: Closed},
					},
					{
						&cell{state: Closed},
						&cell{state: Flagged},
					},
				},
			},
			coord: &Coordinate{X: 1, Y: 1},
			expected: [][]Cell{
				{
					&cell{state: Closed},
					&cell{state: Closed},
				},
				{
					&cell{state: Closed},
					&cell{state: Closed},
				},
			},
		},

		// Invalid coordinate is given
		{
			field: &Field{Width: 3, Height: 3},
			coord: &Coordinate{X: 1, Y: 100},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			// See if given coordinate is valid
			if test.coord.X+1 > test.field.Width || test.coord.Y+1 > test.field.Height {
				_, err := test.field.Unflag(test.coord)
				if err == nil || err != ErrCoordinateOutOfRange {
					t.Fatalf("Expected error is not returned: %s", err)
				}

				return
			}

			target := test.field.Cells[test.coord.Y][test.coord.X]
			oldStatus := target.State()

			result, err := test.field.Unflag(test.coord)

			if oldStatus != Flagged {
				if err == nil {
					t.Fatal("Error should be returned when non-flagged cell is subject to unflag.")
				} else if err != ErrUnflaggingNonFlaggedCell {
					t.Fatal("ErrUnflaggingNonFlaggedCell should be returned when non-flagged cell is subject to unflag.")
				}

				return

			}

			if oldStatus == Flagged && result.NewState != Closed {
				t.Fatalf("Unexpected state is returned: %s", result.NewState)
			}

			for i, row := range test.field.Cells {
				for ii, cell := range row {
					if cell.State() != test.expected[i][ii].State() {
						t.Errorf("Cell with unexpected state is retuned. X: %d, Y: %d. State: %s", i, ii, cell.State())
					}
				}
			}
		})
	}
}

func TestField_Open(t *testing.T) {
	type test struct {
		field    *Field
		coord    *Coordinate
		expected [][]Cell
	}

	tests := []*test{
		// Only left top corner has a mine and right bottom is opened.
		{
			field: &Field{
				Width:  4,
				Height: 4,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: true, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 3, Y: 3},
			expected: [][]Cell{
				{
					&cell{state: Closed, mine: true, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
			},
		},

		// Only left top corner has a mine and the cell with index of 2:1 is subject to open
		{
			field: &Field{
				Width:  4,
				Height: 4,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: true, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 2, Y: 1},
			expected: [][]Cell{
				{
					&cell{state: Closed, mine: true, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
			},
		},

		// Left top corner has a cell with index of 1:1 have mines and right bottom is opened.
		{
			field: &Field{
				Width:  4,
				Height: 4,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: true, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 2},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 2},
						&cell{state: Closed, mine: true, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 3, Y: 3},
			expected: [][]Cell{
				{
					&cell{state: Closed, mine: true, surroundingCnt: 1},
					&cell{state: Closed, mine: false, surroundingCnt: 2},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Closed, mine: false, surroundingCnt: 2},
					&cell{state: Closed, mine: true, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
				{
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
					&cell{state: Opened, mine: false, surroundingCnt: 0},
				},
			},
		},

		// Center cell has a mine and is subject to open.
		{
			field: &Field{
				Width:  3,
				Height: 3,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: true, surroundingCnt: 0},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
					},
				},
			},
			coord: &Coordinate{X: 1, Y: 1},
			expected: [][]Cell{
				{
					&cell{state: Closed, mine: false, surroundingCnt: 1},
					&cell{state: Closed, mine: false, surroundingCnt: 1},
					&cell{state: Closed, mine: false, surroundingCnt: 1},
				},
				{
					&cell{state: Closed, mine: false, surroundingCnt: 1},
					&cell{state: Exploded, mine: true, surroundingCnt: 0},
					&cell{state: Closed, mine: false, surroundingCnt: 1},
				},
				{
					&cell{state: Closed, mine: false, surroundingCnt: 1},
					&cell{state: Closed, mine: false, surroundingCnt: 1},
					&cell{state: Closed, mine: false, surroundingCnt: 1},
				},
			},
		},

		// Invalid coordinate is given
		{
			field: &Field{Width: 3, Height: 3},
			coord: &Coordinate{X: 1, Y: 100},
		},
		{
			field: &Field{Width: 3, Height: 3},
			coord: &Coordinate{X: 100, Y: 1},
		},
		{
			field: &Field{Width: 3, Height: 3},
			coord: &Coordinate{X: 100, Y: 100},
		},

		// Open opened cell
		{
			field: &Field{
				Width:  1,
				Height: 1,
				Cells: [][]Cell{
					{
						&cell{state: Opened, mine: false, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 0, Y: 0},
		},

		// Open flagged cell
		{
			field: &Field{
				Width:  1,
				Height: 1,
				Cells: [][]Cell{
					{
						&cell{state: Flagged, mine: true, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 0, Y: 0},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			// See if given coordinate is valid
			if test.coord.X+1 > test.field.Width || test.coord.Y+1 > test.field.Height {
				_, err := test.field.Open(test.coord)
				if err == nil || err != ErrCoordinateOutOfRange {
					t.Fatalf("Expected error is not returned: %s", err)
				}

				return
			}

			target := test.field.Cells[test.coord.Y][test.coord.X]
			oldStatus := target.State()

			result, err := test.field.Open(test.coord)

			if oldStatus == Opened {
				if err == nil {
					t.Fatal("Error should be returned when opened cell is subject to open.")
				} else if err != ErrOpeningOpenedCell {
					t.Fatal("ErrOpeningOpenedCell should be returned when opened cell is subject to open.")
				}

				return

			}

			if target.State() == Flagged {
				if err == nil {
					t.Fatal("Error should be returned when flagged cell is subject to open.")
				} else if err != ErrOpeningFlaggedCell {
					t.Fatal("ErrOpeningFlaggedCell should be returned when flagged cell is subject to open.")
				}

				return
			}

			if target.hasMine() {
				if result.NewState != Exploded {
					t.Fatalf("State should be exploded when target cell has a mine, but was %s", result.NewState)
				}
			} else if result.NewState != Opened {
				t.Fatalf("Unexpected state is returned: %s", result.NewState)
			}

			for i, row := range test.field.Cells {
				for ii, cell := range row {
					if cell.State() != test.expected[i][ii].State() {
						t.Errorf("Cell with unexpected state is retuned. X: %d, Y: %d. State: %s", i, ii, cell.State())
					}
				}
			}
		})
	}
}

func TestField_MarshalJSON(t *testing.T) {
	state := Exploded
	mine := true
	cnt := 2
	field := &Field{
		Width:  1,
		Height: 1,
		Cells: [][]Cell{
			{
				&cell{state: state, mine: mine, surroundingCnt: cnt},
			},
		},
	}

	bytes, err := json.Marshal(field)

	if err != nil {
		t.Fatalf("Unexpected error is returned: %s.", err.Error())
	}

	str := string(bytes)
	if !strings.Contains(str, state.String()) {
		t.Errorf("Expected state value is not included: %s.", str)
	}

	if !strings.Contains(str, fmt.Sprintf("%t", mine)) {
		t.Errorf("Expected has_mine value is not included: %s.", str)
	}

	if !strings.Contains(str, strconv.Itoa(cnt)) {
		t.Errorf("Expected surrounding_count value is not included: %s.", str)
	}
}

func TestField_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		string         string
		hasError       bool
		state          State
		hasMine        bool
		surroundingCnt int
		height         int
		width          int
	}{
		{
			string:         `{"cells":[[{"has_mine":true,"state":"Flagged","surrounding_count":2}]],"height":1,"width":1}`,
			hasError:       false,
			state:          Flagged,
			hasMine:        true,
			surroundingCnt: 2,
			height:         1,
			width:          1,
		},
		{
			string:   `{"cells":[[{"has_mine":true,"state":"Flagged","surrounding_count":2}]],"height":1}`,
			hasError: true,
		},
		{
			string:   `{"cells":[[{"has_mine":true,"state":"Flagged","surrounding_count":2}]],"width":1}`,
			hasError: true,
		},
		{
			string:   `{"height":1,"width":1}`,
			hasError: true,
		},
		{
			string:   `{"cells": "foobar", height":1,"width":1}`,
			hasError: true,
		},
		{
			string:   `{"cells":[[{"has_mine":true,"state":"Flagged"}]],"height":1,"width":1}`,
			hasError: true,
		},
		{
			string:   `{"cells":[[{"has_mine":true,"surrounding_count":2}]],"height":1,"width":1}`,
			hasError: true,
		},
		{
			string:   `{"cells":[[{"state":"Flagged","surrounding_count":2}]],"height":1,"width":1}`,
			hasError: true,
		},
		{
			string:   `{"cells":[[{"has_mine":true,"state":"Dummy","surrounding_count":2}]],"height":1,"width":1}`,
			hasError: true,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			field := &Field{}
			err := json.Unmarshal([]byte(test.string), field)

			if test.hasError {
				if err == nil {
					t.Fatal("Expected error is not returned.")
				}

				return
			}

			if !test.hasError && err != nil {
				t.Fatalf("Unexpected error is returned: %s.", err.Error())
			}

			if field.Width != test.width {
				t.Errorf("Expected width is not set: %d.", field.Width)
			}

			if field.Height != test.height {
				t.Errorf("Expected height is not set: %d.", field.Height)
			}

			cell := field.Cells[0][0]
			if cell.State() != test.state {
				t.Errorf("Expected state is not set: %s.", cell.State().String())
			}

			if cell.hasMine() != test.hasMine {
				t.Errorf("Expected mine is not set: %t.", cell.hasMine())
			}

			if cell.SurroundingCnt() != test.surroundingCnt {
				t.Errorf("Expected surroundingCnt is not set: %d.", cell.SurroundingCnt())
			}
		})
	}
}
