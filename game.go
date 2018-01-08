package minesweeper

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
)

var (
	// ErrOperatingFinishedGame is returned when a user tries to apply operation to a finished game.
	ErrOperatingFinishedGame = errors.New("can not operate on finished game")
)

// GameState depicts state of the game.
//
// When Cleared or Lost is returned from Game.Operate, the game is finished and no further operation is available.
type GameState int

const (
	_ GameState = iota

	// InProgress represents a state of a game where the game is not finished yet and user operation is available.
	InProgress

	// Cleared represents a state of a game where all safe cells are opened.
	//
	// This state is final so any further Game.Operate call results in returning GameState of Cleared and ErrOperatingFinishedGame.
	Cleared

	// Lost represents a state of a game where non-safe cell was dug and underlying mine has exploded.
	Lost
)

// OpType represents a type of operation a user is applying.
type OpType int

const (
	_ OpType = iota

	// Open represents a kind of operation to open a closed field cell.
	Open

	// Flag represents a kind of operation to flag a closed suspicious field cell with a possible underlying mine.
	Flag

	// Unflag represents a kind of operation to unflag a flagged field cell.
	Unflag
)

// String returns stringified representation of GameState.
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

// MarshalJSON returns GameState value that can be part of JSON structure.
func (s GameState) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, s.String())), nil
}

func strToGameState(str string) (GameState, error) {
	switch str {
	case "InProgress":
		return InProgress, nil

	case "Cleared":
		return Cleared, nil

	case "Lost":
		return Lost, nil

	default:
		return 0, fmt.Errorf("unknown state is given: %s", str)

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
			panic(fmt.Errorf("invalid operation result is returned: %s", r.NewState))

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

func (g *Game) Save(w io.Writer) (int, error) {
	savable := struct {
		Field  *Field    `json:"field"`
		State  GameState `json:"state"`
		Quota  int       `json:"quota"`
		Opened int       `json:"opened"`
	}{
		Field:  g.field,
		State:  g.state,
		Quota:  g.quota,
		Opened: g.opened,
	}

	b, err := json.Marshal(savable)
	if err != nil {
		return 0, err
	}

	return w.Write(b)
}

func Restore(r io.Reader, options ...GameOption) (*Game, error) {
	// Construct game with given options
	game := &Game{}
	for _, opt := range options {
		err := opt(game)
		if err != nil {
			return nil, fmt.Errorf("failed to apply GameOption: %s", err.Error())
		}
	}

	// Setup ui if not set via GameOption
	if game.ui == nil {
		game.ui = &defaultUI{}
	}

	// Parse saved data
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	result := gjson.ParseBytes(b)

	// Set state
	stateValue := result.Get("state")
	if !stateValue.Exists() {
		return nil, errors.New(`"state" field is not given`)
	}
	state, err := strToGameState(stateValue.String())
	if err != nil {
		return nil, err
	}
	game.state = state

	// Set quota
	quotaValue := result.Get("quota")
	if !quotaValue.Exists() {
		return nil, errors.New(`"quota" field is not given`)
	}
	game.quota = int(quotaValue.Int())

	// Set opened
	openedValue := result.Get("opened")
	if !openedValue.Exists() {
		return nil, errors.New(`"opened" field is not given`)
	}
	game.opened = int(openedValue.Int())

	// Set field
	fieldValue := result.Get("field")
	if !fieldValue.Exists() {
		return nil, errors.New(`"field" field is not given`)
	}
	field := &Field{}
	err = json.Unmarshal([]byte(fieldValue.String()), field)
	if err != nil {
		return nil, fmt.Errorf("failed to construct Field: %s", err.Error())
	}
	game.field = field

	return game, nil
}
