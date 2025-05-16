package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

func WhiteBackgroundOp(color color.Color) OperationFunc {
	return OperationFunc(func(t screen.Texture) {
		t.Fill(t.Bounds(), color, screen.Src)
	})
}

func GreenBackgroundOp(color color.Color) OperationFunc {
	return OperationFunc(func(t screen.Texture) {
		t.Fill(t.Bounds(), color, screen.Src)
	})
}

type BgRect struct {
	X1, Y1, X2, Y2 int
}

func (op BgRect) Do(t screen.Texture) bool {
	op.BackgroundRect()(t)
	return false
}

func (op *BgRect) BackgroundRect() OperationFunc {
	return func(t screen.Texture) {
		bounds := image.Rect(op.X1, op.Y1, op.X2, op.Y2)
		t.Fill(bounds, color.Black, screen.Src)
	}
}

type Figure struct {
	X, Y int
}

func (op Figure) Do(t screen.Texture) bool {
	op.Figure()(t)
	return false
}

func (op *Figure) Figure() OperationFunc {
	return func(t screen.Texture) {
		size := 100
		thickness := 20

		shapeColor := color.RGBA{R: 255, G: 230, B: 69, A: 255}

		x, y := op.X, op.Y

		vertical := image.Rect(
			x-thickness/2,
			y-size/2,
			x+thickness/2,
			y+size/2,
		)

		horizontal := image.Rect(
			x-thickness/2-size/2,
			y-thickness/2,
			x-thickness/2,
			y+thickness/2,
		)

		t.Fill(vertical, shapeColor, screen.Src)
		t.Fill(horizontal, shapeColor, screen.Src)
	}
}

type MoveOp struct {
	Mx, My  int
	Figures []*Figure
}

func (op MoveOp) Do(t screen.Texture) bool {
	for _, f := range op.Figures {
		f.X += op.Mx
		f.Y += op.My
	}
	return true
}

type ResetOp struct{}

func (ResetOp) Do(t screen.Texture) bool {
	t.Fill(t.Bounds(), color.Black, screen.Src)
	return true
}

func Reset() Operation {
	return ResetOp{}
}
