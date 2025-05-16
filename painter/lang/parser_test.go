package lang_test

import (
	"strings"
	"testing"

	"github.com/roman-mazur/architecture-lab-3/painter"
	"github.com/roman-mazur/architecture-lab-3/painter/lang"
)

func TestParseFigureAndUpdate(t *testing.T) {
	input := `figure 0.1 0.2
update`

	state := lang.UpdateState()
	parser := &lang.Parser{}

	ops, err := parser.Parse(strings.NewReader(input), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	foundFigure := false
	foundUpdate := false

	for _, op := range ops {
		if _, ok := op.(*painter.Figure); ok {
			foundFigure = true
		}
		if op == painter.UpdateOp {
			foundUpdate = true
		}
	}

	if !foundFigure {
		t.Error("expected figure operation, not found")
	}
	if !foundUpdate {
		t.Error("expected update operation, not found")
	}
}

func TestParseBgColorAndRect(t *testing.T) {
	input := `white
bgrect 0.1 0.1 0.3 0.3
update`

	state := lang.UpdateState()
	parser := &lang.Parser{}

	ops, err := parser.Parse(strings.NewReader(input), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	hasBgColor := false
	hasBgRect := false

	for _, op := range ops {
		switch op.(type) {
		case painter.OperationFunc:
			hasBgColor = true
		case *painter.BgRect:
			hasBgRect = true
		}
	}

	if !hasBgColor {
		t.Error("expected background color operation")
	}
	if !hasBgRect {
		t.Error("expected background rectangle operation")
	}
}

func TestParseMove(t *testing.T) {
	state := lang.UpdateState()

	state.Figures = []*painter.Figure{
		{X: 100, Y: 100},
	}

	input := `move 0.1 0.2`

	parser := &lang.Parser{}
	ops, err := parser.Parse(strings.NewReader(input), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ops) != 1 {
		t.Fatalf("expected 1 move operation, got %d", len(ops))
	}

	ok := ops[0].Do(nil)
	if !ok {
		t.Errorf("expected Do to return true")
	}

	fig := state.Figures[0]
	expectedX := 100 + int(0.1*400)
	expectedY := 100 + int(0.2*400)

	if fig.X != expectedX || fig.Y != expectedY {
		t.Errorf("expected figure at (%d,%d), got (%d,%d)", expectedX, expectedY, fig.X, fig.Y)
	}
}

func TestParseReset(t *testing.T) {
	input := `figure 0.1 0.1
reset`

	state := lang.UpdateState()
	parser := &lang.Parser{}

	ops, err := parser.Parse(strings.NewReader(input), state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(ops) == 0 {
		t.Fatal("expected reset operation")
	}

	if _, ok := ops[0].(painter.ResetOp); ok {
		t.Error("expected first operation to be ResetOp")
	}

	if len(state.Figures) != 0 {
		t.Error("expected state to be reset (no figures), but figures still present")
	}
}

func TestErrorCommand(t *testing.T) {
	input := `invalidcmd 1 2`

	state := lang.UpdateState()
	parser := &lang.Parser{}

	_, err := parser.Parse(strings.NewReader(input), state)
	if err == nil {
		t.Fatal("expected error for unknown command")
	}
}
