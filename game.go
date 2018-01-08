package minesweeper

import (
	"errors"
	"fmt"
)

var (
	ErrOperatingFinishedGame = errors.New("can not operate on finished game")
)

type GameState int

const (
	_ GameState = iota
	InProgress
	Cleared
	Lost
)

type OpType int

const (
	_ OpType = iota
	Open
	Flag
	Unflag
)

func (s GameState) String() string {
	switch s {
	case InProgress:
		return "InProgress"

	case Cleared:
		return "Cleared"

	case Lost:
		return "Lost"

	default:
		panic(fmt.Sprintf("unknown state is given: %d", s))

	}
}

type GameOption func(*Game) error

func WithUI(ui UI) GameOption {
	return func(g *Game) error {
		g.ui = ui
		return nil
	}
}

type Config struct {
	Field *FieldConfig `json:"field" yaml:"field"`
}

func NewConfig() *Config {
	return &Config{
		Field: NewFieldConfig(),
	}
}

type Game struct {
	field  *Field
	ui     UI
	state  GameState
	quota  int
	opened int
}

func NewGame(config *Config, options ...GameOption) (*Game, error) {
	game := &Game{
		state:  InProgress,
		quota:  config.Field.Width*config.Field.Height - config.Field.MineCnt,
		opened: 0,
	}

	// Apply options
	for _, opt := range options {
		err := opt(game)
		if err != nil {
			return nil, fmt.Errorf("failed to apply GameOption: %s", err.Error())
		}
	}

	// Setup field
	field, err := NewField(config.Field)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize field: %s", err.Error())
	}
	game.field = field

	// Setup ui if not set via GameOption
	if game.ui == nil {
		game.ui = &defaultUI{}
	}

	return game, nil
}

func (g *Game) Operate(str string) (GameState, error) {
	if g.state != InProgress {
		return g.state, ErrOperatingFinishedGame
	}

	opType, coord, err := g.ui.ParseInput(str)
	if err != nil {
		return g.state, fmt.Errorf("failed to parse input: %s", err.Error())
	}

	handleOpenResult := func(r *Result) {
		if r == nil {
			return
		}

		switch r.NewState {
		case Exploded:
			g.state = Lost

		case Opened:
			g.opened++
			if g.quota == g.opened {
				g.state = Cleared
			}

		default:
			panic(fmt.Errorf("invalid operation result is returnd: %s", r.NewState))

		}
	}
	switch opType {
	case Open:
		result, err := g.field.Open(coord)
		handleOpenResult(result)
		return g.state, err

	case Flag:
		_, err := g.field.Flag(coord)
		return g.state, err

	case Unflag:
		_, err := g.field.Unflag(coord)
		return g.state, err

	default:
		panic(fmt.Errorf("invalid OpType is returned: %d", opType))

	}
}

func (g *Game) Render() string {
	return g.ui.Render(g.field)
}
