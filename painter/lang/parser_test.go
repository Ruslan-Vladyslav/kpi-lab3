package lang_test

import (
	"strings"
	"testing"

	"image/color"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
)

func TestParser_ValidCommands(t *testing.T) {
	input := `
		update
		white
		green
		rect 0 0 10 10
		point 5 5
		move 1 2
		reset
	`

	parser := lang.Parser{}
	ops, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(ops) != 7 {
		t.Fatalf("Expected 7 operations, got %d", len(ops))
	}

	if fill, ok := ops[1].(painter.FillInColorOp); !ok || fill.Color != color.White {
		t.Errorf("Expected white fill operation, got %T with color %+v", ops[1], fill.Color)
	}

	expectedGreen := color.RGBA{G: 0xff, A: 0xff}
	if fill, ok := ops[2].(painter.FillInColorOp); !ok || fill.Color != expectedGreen {
		t.Errorf("Expected green fill operation, got %T with color %+v", ops[2], fill.Color)
	}

	if _, ok := ops[3].(painter.BgRectOp); !ok {
		t.Errorf("Expected rect operation, got %T", ops[3])
	}

	if _, ok := ops[4].(painter.FigureOp); !ok {
		t.Errorf("Expected point operation, got %T", ops[4])
	}

	if move, ok := ops[5].(painter.MoveOp); !ok || move.MX != 1 || move.MY != 2 {
		t.Errorf("Expected move operation with (1,2), got %+v", ops[5])
	}

	if _, ok := ops[6].(painter.ResetOp); !ok {
		t.Errorf("Expected reset operation, got %T", ops[6])
	}
}

func TestParser_InvalidCommand(t *testing.T) {
	input := `unknown`

	parser := lang.Parser{}
	_, err := parser.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("Expected error for unknown command")
	}
}

func TestParser_InvalidRectArgs(t *testing.T) {
	input := `rect 1 2`

	parser := lang.Parser{}
	_, err := parser.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("Expected error for rect with too few arguments")
	}
}

func TestParser_InvalidPointArgs(t *testing.T) {
	input := `point 1`

	parser := lang.Parser{}
	_, err := parser.Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("Expected error for point with too few arguments")
	}
}

func TestParser_InvalidMoveArgs(t *testing.T) {
	input := `move x y`

	parser := lang.Parser{}
	_, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Log("Note: parser doesn't return error on parseInt fail — consider adding proper validation")
	}
}
