package minesweeper

import (
	"errors"
	"fmt"
	"math/rand"
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

	if config.FieldWidth*config.FieldHeight <= config.MineCnt {
		return errors.New("too many mines")
	}

	return nil
}

type Field struct {
	Width  int
	Height int
	cells  [][]*cell
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

	cells := make([][]*cell, config.FieldHeight)
	for i, row := range grid {
		cells[i] = make([]*cell, config.FieldWidth)

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
		cells:  cells,
	}, nil
}

type Coordinate struct {
	X int
	Y int
}

const (
	Closed = iota
	Opened
	Flagged
	Exploded
)

type cell struct {
	state          uint
	hasMine        bool
	surroundingCnt int
}

func newCell(hasMine bool, surroundingCnt int) *cell {
	return &cell{
		state:          Closed,
		hasMine:        hasMine,
		surroundingCnt: surroundingCnt,
	}
}
