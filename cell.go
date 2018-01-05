package minesweeper

import "fmt"

const (
	Closed State = iota
	Opened
	Flagged
	Exploded
)

type State int

func (s State) String() string {
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

type Cell interface {
	State() State
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
	state          State
	mine           bool
	surroundingCnt int
}

func (c *cell) State() State {
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
