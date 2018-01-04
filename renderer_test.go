package minesweeper

import (
	"fmt"
	"strings"
	"testing"
)

func TestDefaultRenderer_initSymbols(t *testing.T) {
	width := 12
	height := 800
	renderer := &defaultRenderer{}

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
		state    State
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

func TestDefaultRenderer_Render(t *testing.T) {
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

	r := &defaultRenderer{}
	str := r.Render(field)

	for _, state := range []State{Closed, Opened, Flagged, Exploded} {
		if !strings.Contains(str, dispState(state)) {
			t.Errorf("Expected cell state for %s is not included.", state.String())
		}
	}

	if len(strings.Split(str, "\n")) != 3 {
		fmt.Println(len(strings.Split(str, "\n")))
		t.Errorf("Unexpected number of lines: \n%s", str)
	}
}
