package minesweeper

import (
	"errors"
	"fmt"
)

var (
	// ErrOpeningOpenedCell is returned when a user tries to open a cell that is already opened.
	ErrOpeningOpenedCell = errors.New("opened cell can not be opened")

	// ErrOpeningFlaggedCell is returned when a user tries to open a flagged cell.
	ErrOpeningFlaggedCell = errors.New("flagged cell can not be opened")

	// ErrOpeningExplodedCell is returned when a user tries to open exploded cell.
	//
	// This is seldom returned since operation via Game.Operate properly handles game state and returns ErrOperatingFinishedGame beforehand;
	// This error can be returned when and only when Field.Open or Cell.Open is directly called.
	ErrOpeningExplodedCell = errors.New("exploded cell can not be opened")

	// ErrFlaggingOpenedCell is returned when a user tries to flag a cell that is currently flagged.
	ErrFlaggingOpenedCell = errors.New("opened cell can not be flagged")

	// ErrFlaggingFlaggedCell is returned when a user tries to flag a cell that is already flagged.
	ErrFlaggingFlaggedCell = errors.New("flagged cell can not be re-flagged")

	// ErrFlaggingExplodedCell is returned when a user tries to flag exploded cell.
	//
	// This is seldom returned since operation via Game.Operate properly handles game state and returns ErrOperatingFinishedGame beforehand;
	// This error can be returned when and only when Field.Open or Cell.Open is directly called.
	ErrFlaggingExplodedCell = errors.New("exploded cell can not be flagged")

	// ErrUnflaggingNonFlaggedCell is returned when a user tries to unflag a cell that is not currently flagged.
	ErrUnflaggingNonFlaggedCell = errors.New("non-flagged cell can not be unflagged")
)

const (
	_ CellState = iota
	Closed
	Opened
	Flagged
	Exploded
)

type CellState int

func (s CellState) String() string {
	switch s {
	case Closed:
		return "Closed"

	case Opened:
		return "Opened"

	case Flagged:
		return "Flagged"

	case Exploded:
		return "Exploded"

	default:
		panic(fmt.Sprintf("unknown state is given: %d", s))

	}
}

func strToCellState(str string) (CellState, error) {
	switch str {
	case "Closed":
		return Closed, nil

	case "Opened":
		return Opened, nil

	case "Flagged":
		return Flagged, nil

	case "Exploded":
		return Exploded, nil

	default:
		return 0, fmt.Errorf("unknown state is given: %s", str)

	}
}

type Cell interface {
	State() CellState
	SurroundingCnt() int
	hasMine() bool
	flag() (*Result, error)
	unflag() (*Result, error)
	open() (*Result, error)
}

func newCell(hasMine bool, surroundingCnt int) Cell {
	return &cell{
		state:          Closed,
		mine:           hasMine,
		surroundingCnt: surroundingCnt,
	}
}

type cell struct {
	state          CellState
	mine           bool
	surroundingCnt int
}

func (c *cell) State() CellState {
	return c.state
}

func (c *cell) SurroundingCnt() int {
	return c.surroundingCnt
}

func (c *cell) hasMine() bool {
	return c.mine
}

func (c *cell) flag() (*Result, error) {
	switch c.state {
	case Closed:
		c.state = Flagged
		return &Result{NewState: Flagged}, nil

	case Opened:
		return nil, ErrFlaggingOpenedCell

	case Flagged:
		return nil, ErrFlaggingFlaggedCell

	case Exploded:
		return nil, ErrFlaggingExplodedCell

	default:
		panic(fmt.Sprintf("unknown state is set: %d", c.state))

	}
}

func (c *cell) unflag() (*Result, error) {
	switch c.state {
	case Closed, Opened, Exploded:
		return nil, ErrUnflaggingNonFlaggedCell

	case Flagged:
		c.state = Closed
		return &Result{NewState: Closed}, nil

	default:
		panic(fmt.Sprintf("unknown state is set: %d", c.state))

	}
}

func (c *cell) open() (*Result, error) {
	switch c.state {
	case Closed:
		if c.hasMine() {
			c.state = Exploded
			return &Result{
				NewState: Exploded,
			}, nil
		} else {
			c.state = Opened
			return &Result{
				NewState: Opened,
			}, nil
		}

	case Opened:
		return nil, ErrOpeningOpenedCell

	case Flagged:
		return nil, ErrOpeningFlaggedCell

	case Exploded:
		return nil, ErrOpeningExplodedCell

	default:
		panic(fmt.Sprintf("unknown state is set: %d", c.state))

	}
}
