package minesweeper

import (
	"testing"
)

func TestNewConfig(t *testing.T) {
	config := NewConfig()

	if config.FieldWidth == 0 {
		t.Errorf("Config.FieldWidth is not set.")
	}

	if config.FieldHeight == 0 {
		t.Errorf("Config.FieldHeight is not set.")
	}

	if config.MineCnt == 0 {
		t.Errorf("Config.MineCnt is not set.")
	}
}

func TestNewFiled(t *testing.T) {
	var configs = []*Config{
		{
			FieldWidth:  12,
			FieldHeight: 0,
			MineCnt:     9,
		},
		{
			FieldWidth:  0,
			FieldHeight: 12,
			MineCnt:     9,
		},
		{
			FieldWidth:  12,
			FieldHeight: 12,
			MineCnt:     0,
		},
		{
			FieldWidth:  12,
			FieldHeight: 12,
			MineCnt:     9,
		},
	}

	for i, config := range configs {
		field, err := NewField(config)

		if config.FieldWidth == 0 || config.FieldHeight == 0 || config.MineCnt == 0 {
			if err == nil {
				t.Errorf("Error is not returned on invalid *Config. Test #%d.", i+1)
			}

			continue
		}

		if field == nil {
			t.Fatal("Field is nil.")
		}

		mineCnt := 0
		for _, row := range field.cells {
			for _, c := range row {
				if c.hasMine {
					mineCnt++
				}
			}
		}
		if config.MineCnt != mineCnt {
			t.Errorf("Expected mine count of %d, but was %d.", config.MineCnt, mineCnt)
		}
	}
}
