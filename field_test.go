package minesweeper

import (
	"fmt"
	"testing"
)

func TestState_String(t *testing.T) {
	tests := []struct {
		state    State
		expected string
	}{
		{
			state:    Closed,
			expected: "Closed",
		},
		{
			state:    Opened,
			expected: "Opened",
		},
		{
			state:    Flagged,
			expected: "Flagged",
		},
		{
			state:    Exploded,
			expected: "Exploded",
		},
		{
			state: 123,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if test.expected != "" {
						t.Fatalf("Unexpectedly panicked for state: %d", test.state)
					}
				}
			}()

			s := test.state.String()
			if s != test.expected {
				t.Fatalf("Expected %s, but %s was returned.", test.expected, s)
			}
		})
	}
}

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
			for _, row := range field.cells {
				for _, c := range row {
					if c.hasMine {
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

func TestField_Open(t *testing.T) {
	type test struct {
		field    *Field
		coord    *Coordinate
		expected [][]*cell
	}

	tests := []*test{
		// Only left top corner has a mine and right bottom is opened.
		{
			field: &Field{
				Width:  4,
				Height: 4,
				cells: [][]*cell{
					{
						{state: Closed, hasMine: true, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 3, Y: 3},
			expected: [][]*cell{
				{
					{state: Closed, hasMine: true, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
			},
		},

		// Only left top corner has a mine and the cell with index of 2:1 is subject to open
		{
			field: &Field{
				Width:  4,
				Height: 4,
				cells: [][]*cell{
					{
						{state: Closed, hasMine: true, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 2, Y: 1},
			expected: [][]*cell{
				{
					{state: Closed, hasMine: true, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
			},
		},

		// Left top corner has a cell with index of 1:1 have mines and right bottom is opened.
		{
			field: &Field{
				Width:  4,
				Height: 4,
				cells: [][]*cell{
					{
						{state: Closed, hasMine: true, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 2},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 2},
						{state: Closed, hasMine: true, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 0},
					},
				},
			},
			coord: &Coordinate{X: 3, Y: 3},
			expected: [][]*cell{
				{
					{state: Closed, hasMine: true, surroundingCnt: 1},
					{state: Closed, hasMine: false, surroundingCnt: 2},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Closed, hasMine: false, surroundingCnt: 2},
					{state: Closed, hasMine: true, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 1},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
				{
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
					{state: Opened, hasMine: false, surroundingCnt: 0},
				},
			},
		},

		// Center cell has a mine and is subject to open.
		{
			field: &Field{
				Width:  3,
				Height: 3,
				cells: [][]*cell{
					{
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: true, surroundingCnt: 0},
						{state: Closed, hasMine: false, surroundingCnt: 1},
					},
					{
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
						{state: Closed, hasMine: false, surroundingCnt: 1},
					},
				},
			},
			coord: &Coordinate{X: 1, Y: 1},
			expected: [][]*cell{
				{
					{state: Closed, hasMine: false, surroundingCnt: 1},
					{state: Closed, hasMine: false, surroundingCnt: 1},
					{state: Closed, hasMine: false, surroundingCnt: 1},
				},
				{
					{state: Closed, hasMine: false, surroundingCnt: 1},
					{state: Exploded, hasMine: true, surroundingCnt: 0},
					{state: Closed, hasMine: false, surroundingCnt: 1},
				},
				{
					{state: Closed, hasMine: false, surroundingCnt: 1},
					{state: Closed, hasMine: false, surroundingCnt: 1},
					{state: Closed, hasMine: false, surroundingCnt: 1},
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
				cells: [][]*cell{
					{
						{state: Opened, hasMine: false, surroundingCnt: 0},
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
				cells: [][]*cell{
					{
						{state: Flagged, hasMine: true, surroundingCnt: 0},
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

			target := test.field.cells[test.coord.Y][test.coord.X]
			oldStatus := target.state

			result, err := test.field.Open(test.coord)

			if oldStatus == Opened {
				if err == nil {
					t.Fatal("Error should be returned when opened cell is subject to open.")
				} else if err != ErrOpeningOpenedCell {
					t.Fatal("ErrOpeningOpenedCell should be returned when opened cell is subject to open.")
				}

				return

			}

			if target.state == Flagged {
				if err == nil {
					t.Fatal("Error should be returned when flagged cell is subject to open.")
				} else if err != ErrOpeningFlaggedCell {
					t.Fatal("ErrOpeningFlaggedCell should be returned when flagged cell is subject to open.")
				}

				return
			}

			if target.hasMine {
				if result.NewState != Exploded {
					t.Fatalf("State should be exploded when target cell has a mine, but was %s", result.NewState)
				}
			} else if result.NewState != Opened {
				t.Fatalf("Unexpected state is returned: %s", result.NewState)
			}

			for i, row := range test.field.cells {
				for ii, cell := range row {
					if cell.state != test.expected[i][ii].state {
						t.Errorf("Cell has unexpected state is retuned. X: %d, Y: %d. State: %s", i, ii, cell.state)
					}
				}
			}
		})
	}
}
