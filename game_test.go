package minesweeper

import (
	"errors"
	"fmt"
	"testing"
)

type DummyUI struct {
	RenderFunc     func(*Field) string
	ParseInputFunc func(string) (*Coordinate, error)
}

func (ui *DummyUI) Render(field *Field) string {
	return ui.RenderFunc(field)
}

func (ui *DummyUI) ParseInput(str string) (*Coordinate, error) {
	return ui.ParseInputFunc(str)
}

func TestWithUI(t *testing.T) {
	ui := &DummyUI{}

	option := WithUI(ui)

	if option == nil {
		t.Fatal("Expected GameOption is not returned.")
	}

	err := option(&Game{})
	if err != nil {
		t.Fatalf("Unexpected error is returned: %s.", err.Error())
	}
}

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.Field == nil {
		t.Error("Field should be filled with default configuration.")
	}
}

func TestNewGame(t *testing.T) {
	validFieldConfig := &FieldConfig{
		Height:  3,
		Width:   3,
		MineCnt: 1,
	}

	tests := []struct {
		config   *Config
		options  []GameOption
		hasError bool
		ui       UI
	}{
		{
			config: &Config{Field: validFieldConfig},
		},
		{
			config:   &Config{Field: validFieldConfig},
			options:  []GameOption{func(_ *Game) error { return errors.New("dummy") }},
			hasError: true,
		},
		{
			config:   &Config{Field: &FieldConfig{}},
			hasError: true,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			game, err := NewGame(test.config, test.options...)
			if test.hasError {
				if err == nil {
					t.Fatal("Expected error is not returned.")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error is returned: %s.", err.Error())
			}

			if game.ui == nil {
				t.Error("UI should be set even when one is not given via GameOption.")
			}

		})
	}
}

func TestGame_Open(t *testing.T) {
	tests := []struct {
		ui       UI
		field    *Field
		expected *Result
	}{
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (*Coordinate, error) {
					return nil, errors.New("dummy")
				},
			},
		},
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (*Coordinate, error) {
					return &Coordinate{X: 0, Y: 0}, nil
				},
			},
			field: &Field{
				Width:  1,
				Height: 1,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: false, surroundingCnt: 0},
					},
				},
			},
			expected: &Result{
				NewState: Opened,
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			game := &Game{
				ui:    test.ui,
				field: test.field,
			}

			result, err := game.Open("dummy")

			if test.expected == nil {
				if err == nil {
					t.Fatal("Expected error is not returned.")
				}

				return
			}
			if err != nil {
				t.Fatalf("Unexpected error is returned: %s.", err.Error())
			}

			if result.NewState != test.expected.NewState {
				t.Errorf("Expected new state to be %s, but was %s.", test.expected.NewState, result.NewState.String())
			}
		})
	}
}

func TestGame_Render(t *testing.T) {
	str := "dummy"
	ui := &DummyUI{
		RenderFunc: func(_ *Field) string {
			return str
		},
	}
	game := &Game{
		field: &Field{},
		ui:    ui,
	}

	rendered := game.Render()

	if rendered != str {
		t.Errorf("Unexpected output is given: %s.", rendered)
	}
}
