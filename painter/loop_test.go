package painter

import (
	"image"
	"image/color"
	"image/draw"
	"reflect"
	"testing"

	"golang.org/x/exp/shiny/screen"
)

func makeTestOp(name string, testOps *[]string) OperationFunc {
	return OperationFunc(func(t screen.Texture) {
		if testOps != nil {
			*testOps = append(*testOps, name)
		}
	})
}

func TestLoop_Post(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	var testOps []string

	l.Start(mockScreen{})

	l.Post(makeTestOp("op 1", &testOps))
	l.Post(makeTestOp("op 2", &testOps))
	l.Post(makeTestOp("op 3", &testOps))

	l.Post(WhiteBackgroundOp(color.White))
	l.Post(UpdateOp)
	l.StopAndWait()

	if tr.lastTexture == nil {
		t.Fatal("Texture was not updated")
	}

	mt, ok := tr.lastTexture.(*mockTexture)
	if !ok {
		t.Fatal("Unexpected texture type", reflect.TypeOf(tr.lastTexture))
	}

	if len(mt.Colors) == 0 || mt.Colors[0] != color.White {
		t.Error("First color is not white or Colors is empty:", mt.Colors)
	}

	if !reflect.DeepEqual(testOps, []string{"op 1", "op 2", "op 3"}) {
		t.Error("Bad order of operations:", testOps)
	}
}

func TestLoop_StopAndWait(t *testing.T) {
	var (
		l  Loop
		tr testReceiver
	)
	l.Receiver = &tr

	l.Start(mockScreen{})
	l.Post(WhiteBackgroundOp(color.White))

	l.StopAndWait()

	if !l.stopReq {
		t.Error("Loop.stopReq should be true after StopAndWait")
	}

	select {
	case <-l.stopped:
	default:
		t.Error("Loop.stopped channel should be closed after StopAndWait")
	}
}

func TestBgRect_Do(t *testing.T) {
	mt := &mockTexture{}
	op := BgRect{10, 20, 30, 40}

	op.Do(mt)

	if len(mt.Colors) != 1 {
		t.Fatalf("expected 1 fill, got %d", len(mt.Colors))
	}

	expectedRect := image.Rect(10, 20, 30, 40)
	if mt.Rects[0] != expectedRect {
		t.Errorf("expected rect %v, got %v", expectedRect, mt.Rects[0])
	}
	if mt.Colors[0] != color.Black {
		t.Errorf("expected color Black, got %v", mt.Colors[0])
	}
}

func TestFigure_Do(t *testing.T) {
	mt := &mockTexture{}
	op := Figure{X: 100, Y: 100}

	op.Do(mt)

	if len(mt.Colors) != 2 {
		t.Fatalf("expected 2 fills (vertical and horizontal), got %d", len(mt.Colors))
	}

	expectedColor := color.RGBA{R: 255, G: 230, B: 69, A: 255}
	for i, c := range mt.Colors {
		rgba, ok := c.(color.RGBA)
		if !ok {
			t.Errorf("fill %d color is not RGBA: %v", i, c)
			continue
		}
		if rgba != expectedColor {
			t.Errorf("fill %d color = %v, want %v", i, rgba, expectedColor)
		}
	}
}

func TestMove(t *testing.T) {
	figs := []*Figure{
		{X: 10, Y: 20},
		{X: 30, Y: 40},
	}

	mt := &mockTexture{}
	moveOp := Move(5, -5, figs)
	moveOp(mt)

	if figs[0].X != 15 || figs[0].Y != 15 {
		t.Errorf("figure 0 position expected (15,15), got (%d,%d)", figs[0].X, figs[0].Y)
	}
	if figs[1].X != 35 || figs[1].Y != 35 {
		t.Errorf("figure 1 position expected (35,35), got (%d,%d)", figs[1].X, figs[1].Y)
	}
}

func TestResetOp_Do(t *testing.T) {
	mt := &mockTexture{}

	op := ResetOp{}
	ok := op.Do(mt)

	if !ok {
		t.Error("ResetOp.Do should return true")
	}

	if len(mt.Colors) == 0 {
		t.Fatal("Fill was not called on texture")
	}

	if mt.Colors[0] != color.Black {
		t.Errorf("expected fill color Black, got %v", mt.Colors[0])
	}

	expectedRect := mt.Bounds()
	if len(mt.Rects) == 0 {
		t.Fatal("Fill was not called with any rectangle")
	} else if mt.Rects[0] != expectedRect {
		t.Errorf("expected fill rect %v, got %v", expectedRect, mt.Rects[0])
	}
}

type testReceiver struct {
	lastTexture screen.Texture
}

func (tr *testReceiver) Update(t screen.Texture) {
	tr.lastTexture = t
}

type mockScreen struct{}

func (m mockScreen) NewBuffer(size image.Point) (screen.Buffer, error) {
	panic("implement me")
}

func (m mockScreen) NewTexture(size image.Point) (screen.Texture, error) {
	return new(mockTexture), nil
}

func (m mockScreen) NewWindow(opts *screen.NewWindowOptions) (screen.Window, error) {
	panic("implement me")
}

type mockTexture struct {
	Colors []color.Color
	Rects  []image.Rectangle
}

func (m *mockTexture) Release() {}

func (m *mockTexture) Size() image.Point { return size }

func (m *mockTexture) Bounds() image.Rectangle {
	return image.Rectangle{Max: m.Size()}
}

func (m *mockTexture) Upload(dp image.Point, src screen.Buffer, sr image.Rectangle) {}

func (m *mockTexture) Fill(dr image.Rectangle, src color.Color, op draw.Op) {
	m.Colors = append(m.Colors, src)
	m.Rects = append(m.Rects, dr)
}
