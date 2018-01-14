package minesweeper

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func TestDefaultUI_initSymbols(t *testing.T) {
	width := 12
	height := 800
	renderer := &defaultUI{}

	renderer.initSymbols(width, height)

	if len(renderer.xSymbols) != width {
		t.Fatalf("Unexpected amount of symbols for x axis is set: %d", len(renderer.xSymbols))
	}

	if len(renderer.ySymbols) != height {
		t.Fatalf("Unexpected amount of symbols for y axis is set: %d.", len(renderer.ySymbols))
	}

	firstX := renderer.xSymbols[0]
	if firstX != 1 {
		t.Errorf("Unexpected symbol is returned: %d", firstX)
	}

	lastX := renderer.xSymbols[width-1]
	if lastX != width {
		t.Errorf("Unexpected symbol is returned: %d", lastX)
	}

	firstY := renderer.ySymbols[0]
	if firstY != "a" {
		t.Errorf("Unexpected symbol is returned: %s", firstY)
	}

	lastY := renderer.ySymbols[height-1]
	if lastY != "adt" {
		t.Errorf("Unexpected symbol is returned: %s", lastY)
	}
}

func Test_dispState(t *testing.T) {
	tests := []struct {
		state    CellState
		expected string
	}{
		{
			state:    Closed,
			expected: " ",
		},
		{
			state:    Opened,
			expected: "-",
		},
		{
			state:    Flagged,
			expected: "F",
		},
		{
			state:    Exploded,
			expected: "X",
		},
		{
			state: 999,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if test.expected != "" {
						t.Fatal("Panicked unexpectedly.")
					}
				}
			}()

			result := dispState(test.state)

			if result != test.expected {
				t.Errorf(`Expected "%s" but "%s" was returned.`, test.expected, result)
			}
		})
	}
}

func TestDefaultUI_Render(t *testing.T) {
	field := &Field{
		Width:  2,
		Height: 2,
		Cells: [][]Cell{
			{
				&cell{state: Closed},
				&cell{state: Opened},
			},
			{
				&cell{state: Flagged},
				&cell{state: Exploded},
			},
		},
	}

	w := bytes.NewBuffer([]byte{})
	r := &defaultUI{}
	_, err := r.Render(w, field)

	if err != nil {
		t.Fatalf("Unexpected error is returned: %s.", err.Error())
	}

	str := w.String()
	for _, state := range []CellState{Closed, Opened, Flagged, Exploded} {
		if !strings.Contains(str, dispState(state)) {
			t.Errorf("Expected cell state for %s is not included.", state.String())
		}
	}

	if len(strings.Split(str, "\n")) != 3 {
		fmt.Println(len(strings.Split(str, "\n")))
		t.Errorf("Unexpected number of lines: \n%s", str)
	}
}

func TestDefaultUI_ParseInput(t *testing.T) {
	tests := []struct {
		xSymbols []int
		ySymbols []string
		input    []byte
		opType   OpType
		expected *Coordinate
	}{
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b", "c"},
			input:    []byte("2 c"),
			opType:   Open,
			expected: &Coordinate{X: 1, Y: 2},
		},
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b", "c"},
			input:    []byte("2 b f"),
			opType:   Flag,
			expected: &Coordinate{X: 1, Y: 1},
		},
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b", "c"},
			input:    []byte("2 b flag"),
			opType:   Flag,
			expected: &Coordinate{X: 1, Y: 1},
		},
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b", "c"},
			input:    []byte("2 a u"),
			opType:   Unflag,
			expected: &Coordinate{X: 1, Y: 0},
		},
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b", "c"},
			input:    []byte("2 a unflag"),
			opType:   Unflag,
			expected: &Coordinate{X: 1, Y: 0},
		},
		{
			input: []byte("2 invalid"),
		},
		{
			input: []byte("invalid abc"),
		},
		{
			input: []byte("invalid number of fields"),
		},
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b"},
			input:    []byte("100 a"),
		},
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b"},
			input:    []byte("1 zzz"),
		},
		{
			xSymbols: []int{1, 2},
			ySymbols: []string{"a", "b", "c"},
			input:    []byte("2 a invalid"),
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			ui := &defaultUI{
				xSymbols: test.xSymbols,
				ySymbols: test.ySymbols,
			}

			opType, coord, err := ui.ParseInput(test.input)

			if test.expected == nil {
				if err == nil {
					t.Fatal("Expected error is not returned.")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error is returned: %s.", err.Error())
			}

			if opType != test.opType {
				t.Errorf("Expected OpType to be %d, but was %d.", test.opType, opType)
			}

			if coord.X != test.expected.X {
				t.Errorf("Expected X to be %d, but was %d.", test.expected.X, coord.X)
			}

			if coord.Y != test.expected.Y {
				t.Errorf("Expected Y to be %d, but was %d.", test.expected.Y, coord.Y)
			}
		})
	}
}
