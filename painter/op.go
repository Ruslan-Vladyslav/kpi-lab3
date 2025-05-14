package painter

import (
	"image"
	"image/color"

	"golang.org/x/exp/shiny/screen"
)

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture, s *State) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture, s *State) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t, s) || ready
	}
	return
}

var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture, s *State) bool { return true }

type OperationFunc func(t screen.Texture, s *State)

func (f OperationFunc) Do(t screen.Texture, s *State) bool {
	f(t, s)
	return false
}

type State struct {
	BackgroundColor color.Color
	BgRect          *image.Rectangle
	Figures         []image.Point
}

type FillInColorOp struct {
	Color color.Color
}

func (op FillInColorOp) Do(t screen.Texture, s *State) bool {
	s.BackgroundColor = op.Color
	return false
}

type BgRectOp struct {
	X1, Y1, X2, Y2 float64
}

func (op BgRectOp) Do(t screen.Texture, s *State) bool {

	r := image.Rect(
		int(op.X1*float64(t.Size().X)),
		int(op.Y1*float64(t.Size().Y)),
		int(op.X2*float64(t.Size().X)),
		int(op.Y2*float64(t.Size().Y)))

	s.BgRect = &r
	return false
}

type FigureOp struct {
	X, Y float64
}

func (op FigureOp) Do(t screen.Texture, s *State) bool {
	x := int(op.X * float64(t.Size().X))
	y := int(op.Y * float64(t.Size().Y))

	s.Figures = append(s.Figures, image.Pt(x, y))
	return false
}

type MoveOp struct {
	MX, MY int
}

func (op MoveOp) Do(t screen.Texture, s *State) bool {
	for i := range s.Figures {
		s.Figures[i].X += op.MX
		s.Figures[i].Y += op.MY
	}
	return false
}

type ResetOp struct{}

func (op ResetOp) Do(t screen.Texture, s *State) bool {
	s.BackgroundColor = color.Black
	s.BgRect = nil
	s.Figures = nil
	return false
}
