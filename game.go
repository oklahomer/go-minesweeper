package minesweeper

import (
	"fmt"
)

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
	field *Field
	ui    UI
}

func NewGame(config *Config, options ...GameOption) (*Game, error) {
	game := &Game{}

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

func (g *Game) Open(str string) (*Result, error) {
	coord, err := g.ui.ParseInput(str)
	if err != nil {
		return nil, fmt.Errorf("failed to parse input: %s", err.Error())
	}

	return g.field.Open(coord)
}

func (g *Game) Render() string {
	return g.ui.Render(g.field)
}
