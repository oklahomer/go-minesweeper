package minesweeper

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
)

type DummyUI struct {
	RenderFunc     func(*Field) string
	ParseInputFunc func(string) (OpType, *Coordinate, error)
}

func (ui *DummyUI) Render(field *Field) string {
	return ui.RenderFunc(field)
}

func (ui *DummyUI) ParseInput(str string) (OpType, *Coordinate, error) {
	return ui.ParseInputFunc(str)
}

func TestGameState_String(t *testing.T) {
	tests := []struct {
		state    GameState
		expected string
	}{
		{
			state:    InProgress,
			expected: "InProgress",
		},
		{
			state:    Cleared,
			expected: "Cleared",
		},
		{
			state:    Lost,
			expected: "Lost",
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

			if game.state != InProgress {
				t.Errorf("Unexpected state is set on construction: %d.", game.state)
			}

			if game.quota != test.config.Field.Width*test.config.Field.Height-test.config.Field.MineCnt {
				t.Errorf("Unexpected quota value is set: %d.", game.quota)
			}

			if game.opened != 0 {
				t.Errorf("Unexpected count is set: %d.", game.opened)
			}
		})
	}
}

func TestGame_Operate(t *testing.T) {
	tests := []struct {
		ui             UI
		field          *Field
		resultingState GameState
	}{
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (OpType, *Coordinate, error) {
					return 0, nil, errors.New("dummy")
				},
			},
		},
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (OpType, *Coordinate, error) {
					return Open, &Coordinate{X: 100, Y: 100}, nil
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
		},
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (OpType, *Coordinate, error) {
					return Open, &Coordinate{X: 0, Y: 0}, nil
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
			resultingState: Cleared,
		},
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (OpType, *Coordinate, error) {
					return Open, &Coordinate{X: 0, Y: 0}, nil
				},
			},
			field: &Field{
				Width:  2,
				Height: 2,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: false, surroundingCnt: 1},
					},
					{
						&cell{state: Closed, mine: false, surroundingCnt: 1},
						&cell{state: Closed, mine: true, surroundingCnt: 0},
					},
				},
			},
			resultingState: InProgress,
		},
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (OpType, *Coordinate, error) {
					return Open, &Coordinate{X: 0, Y: 0}, nil
				},
			},
			field: &Field{
				Width:  1,
				Height: 1,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: true, surroundingCnt: 0},
					},
				},
			},
			resultingState: Lost,
		},
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (OpType, *Coordinate, error) {
					return Flag, &Coordinate{X: 0, Y: 0}, nil
				},
			},
			field: &Field{
				Width:  1,
				Height: 1,
				Cells: [][]Cell{
					{
						&cell{state: Closed, mine: true, surroundingCnt: 0},
					},
				},
			},
			resultingState: InProgress,
		},
		{
			ui: &DummyUI{
				ParseInputFunc: func(s string) (OpType, *Coordinate, error) {
					return Unflag, &Coordinate{X: 0, Y: 0}, nil
				},
			},
			field: &Field{
				Width:  1,
				Height: 1,
				Cells: [][]Cell{
					{
						&cell{state: Flagged, mine: true, surroundingCnt: 0},
					},
				},
			},
			resultingState: InProgress,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			quota := 0
			if test.field != nil {
				for _, row := range test.field.Cells {
					for _, c := range row {
						if !c.hasMine() {
							quota++
						}
					}
				}
			}
			game := &Game{
				ui:     test.ui,
				field:  test.field,
				state:  InProgress,
				quota:  quota,
				opened: 0,
			}

			state, err := game.Operate("dummy")

			if test.resultingState == 0 {
				if err == nil {
					t.Fatal("Expected error is not returned.")
				}

				return
			}

			if err != nil {
				t.Fatalf("Unexpected error is returned: %s.", err.Error())
			}

			if state != test.resultingState {
				t.Errorf("Expected new state to be %s, but was %s.", test.resultingState.String(), state.String())
			}

			if state != game.state {
				t.Errorf("Returned state is %s, but stored state is %s.", state.String(), game.state.String())
			}

			if test.resultingState != InProgress {
				state, err = game.Operate("dummy")
				if err == nil {
					t.Error("Error should be returned when operated on finished game.")
				}

				if state != test.resultingState {
					t.Errorf("The state should stay as-is when Game.Operate is called after finished.")
				}
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

func TestGame_Save(t *testing.T) {
	game := &Game{
		field: &Field{
			Width:  2,
			Height: 2,
			Cells: [][]Cell{
				{
					&cell{state: Opened, mine: false, surroundingCnt: 1},
					&cell{state: Closed, mine: false, surroundingCnt: 1},
				},
				{
					&cell{state: Closed, mine: true, surroundingCnt: 0},
					&cell{state: Closed, mine: false, surroundingCnt: 1},
				},
			},
		},
		state:  InProgress,
		quota:  1,
		opened: 1,
	}

	buf := bytes.NewBufferString("")
	i, err := game.Save(buf)

	if err != nil {
		t.Fatalf("Unexpected error is returned: %s.", err.Error())
	}

	if i == 0 {
		t.Error("No byte was written.")
	}

	// {"field":{"cells":[[{"has_mine":false,"state":"Opened","surrounding_count":1},{"has_mine":false,"state":"Closed","surrounding_count":1}],[{"has_mine":true,"state":"Closed","surrounding_count":0},{"has_mine":false,"state":"Closed","surrounding_count":1}]],"height":2,"width":2},"state":"InProgress","quota":1,"opened":1}
	str := buf.String()
	for _, jsonField := range []string{"field", "state", "quota", "opened"} {
		if !strings.Contains(str, jsonField) {
			t.Errorf(`Mandatory field "%s" is not present`, jsonField)
		}
	}
}

func TestRestore(t *testing.T) {
	tests := []struct {
		str      string
		options  []GameOption
		hasError bool
		state    GameState
		quota    int
		opened   int
	}{
		{
			str:    `{"state":"InProgress","quota":1,"opened":2,"field":{"cells":[[{"has_mine":false,"state":"Opened","surrounding_count":1},{"has_mine":false,"state":"Opened","surrounding_count":1}],[{"has_mine":true,"state":"Closed","surrounding_count":0},{"has_mine":false,"state":"Closed","surrounding_count":1}]],"height":2,"width":2}}`,
			state:  InProgress,
			quota:  1,
			opened: 2,
		},
		{
			str:      `{"state":"INVALID_STATE","quota":1,"opened":2,"field":{"cells":[[{"has_mine":false,"state":"Opened","surrounding_count":1},{"has_mine":false,"state":"Opened","surrounding_count":1}],[{"has_mine":true,"state":"Closed","surrounding_count":0},{"has_mine":false,"state":"Closed","surrounding_count":1}]],"height":2,"width":2}}`,
			hasError: true,
		},
		{
			str:      `{"quota":1,"opened":2,"field":{"cells":[[{"has_mine":false,"state":"Opened","surrounding_count":1},{"has_mine":false,"state":"Opened","surrounding_count":1}],[{"has_mine":true,"state":"Closed","surrounding_count":0},{"has_mine":false,"state":"Closed","surrounding_count":1}]],"height":2,"width":2}}`,
			hasError: true,
		},
		{
			str:      `{"state":"InProgress","opened":2,"field":{"cells":[[{"has_mine":false,"state":"Opened","surrounding_count":1},{"has_mine":false,"state":"Opened","surrounding_count":1}],[{"has_mine":true,"state":"Closed","surrounding_count":0},{"has_mine":false,"state":"Closed","surrounding_count":1}]],"height":2,"width":2}}`,
			hasError: true,
		},
		{
			str:      `{"state":"InProgress","quota":1,"field":{"cells":[[{"has_mine":false,"state":"Opened","surrounding_count":1},{"has_mine":false,"state":"Opened","surrounding_count":1}],[{"has_mine":true,"state":"Closed","surrounding_count":0},{"has_mine":false,"state":"Closed","surrounding_count":1}]],"height":2,"width":2}}`,
			hasError: true,
		},
		{
			str:      `{"state":"InProgress","quota":1,"opened":2}`,
			hasError: true,
		},
		{
			str:      `{"state":"InProgress","quota":1,"opened":2,"field":{"width":2}}`,
			hasError: true,
		},
		{
			str:      `{"state":"InProgress","quota":1,"opened":2,"field":{"cells":[[{"has_mine":false,"state":"Opened","surrounding_count":1},{"has_mine":false,"state":"Opened","surrounding_count":1}],[{"has_mine":true,"state":"Closed","surrounding_count":0},{"has_mine":false,"state":"Closed","surrounding_count":1}]],"height":2,"width":2}}`,
			options:  []GameOption{func(_ *Game) error { return errors.New("dummy") }},
			hasError: true,
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			game, err := Restore(strings.NewReader(test.str), test.options...)
			if test.hasError {
				if err == nil {
					t.Fatal("Expected error is not returned.")
				}

				return
			}

			if !test.hasError && err != nil {
				t.Fatalf("Unexpected error is returned: %s.", err.Error())
			}

			if game.ui == nil {
				t.Error("UI must be set.")
			}

			if game.state != test.state {
				t.Errorf("Unexpected state is set: %s.", game.state.String())
			}

			if game.quota != test.quota {
				t.Errorf("Unexpected quota is set: %d.", game.quota)
			}

			if game.opened != test.opened {
				t.Errorf("Unexpected opened is set: %d.", game.opened)
			}
		})
	}
}

func Test_strToGameState(t *testing.T) {
	tests := []struct {
		string string
		state  GameState
	}{
		{
			string: "InProgress",
			state:  InProgress,
		},
		{
			string: "Cleared",
			state:  Cleared,
		},
		{
			string: "Lost",
			state:  Lost,
		},
		{
			string: "INVALID",
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test #%d", i+1), func(t *testing.T) {
			state, err := strToGameState(test.string)

			if test.state == 0 && err == nil {
				t.Fatal("Expected error is not returned.")
			}

			if test.state != 0 && err != nil {
				t.Fatalf("Unexpected error is returned: %s.", err.Error())
			}

			if state != test.state {
				t.Errorf("Unexpected state is returned: %s.", state.String())
			}
		})
	}
}
