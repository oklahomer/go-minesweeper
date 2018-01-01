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

func TestCell_State(t *testing.T) {
	state := Exploded
	var c Cell = &cell{state: state}
	if c.State() != state {
		t.Fatalf("Expected state is not returned: %s.", c.(Cell).State())
	}
}

func TestCell_SurroundingCnt(t *testing.T) {
	cnt := 123
	var c Cell = &cell{surroundingCnt: cnt}
	if c.SurroundingCnt() != cnt {
		t.Fatalf("Expected count is not returned: %d.", cnt)
	}
}

func TestCell_open(t *testing.T) {
	tests := []struct {
		cell     *cell
		newState State
		error    error
	}{
		{
			cell:     &cell{state: Closed, mine: false},
			newState: Opened,
		},
		{
			cell:     &cell{state: Closed, mine: true},
			newState: Exploded,
		},
		{
			cell:  &cell{state: Opened},
			error: ErrOpeningOpenedCell,
		},
		{
			cell:  &cell{state: Flagged},
			error: ErrOpeningFlaggedCell,
		},
		{
			cell:  &cell{state: Exploded},
			error: ErrOpeningExplodedCell,
		},
		{
			cell: &cell{state: 123456},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					if test.newState != 0 || test.error != nil {
						// State should be changed or error is expected to be returned
						t.Fatal("Panicked unexpectedly.")
					}
				}
			}()

			result, err := test.cell.open()
			if test.error != err {
				t.Errorf("Unexpected error is returned: %s.", err)
			}

			if test.newState != 0 && test.newState != test.cell.state {
				t.Errorf("Unexpected state: %s.", test.cell.state)
			}

			if test.newState != 0 && test.newState != result.NewState {
				t.Errorf("Unepxected result is returned %+v.", result)
			}
		})
	}
}
