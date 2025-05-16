package lang

import (
	"bufio"
	"fmt"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

type Parser struct {
}

var isReset = false

type CurState struct {
	Figures    []*painter.Figure
	BgRectFill []*painter.BgRect

	BgColorOp painter.OperationFunc
	UpdateOp  painter.Operation
	MoveOp    []painter.Operation
}

func UpdateState() *CurState { return &CurState{} }

func (p *Parser) Parse(in io.Reader, s *CurState) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var res []painter.Operation

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		operations, err := parseLine(line, s)
		if err != nil {
			return nil, fmt.Errorf("[Error]: parse error on line '%s': %w", line, err)
		}

		res = append(res, operations...)
	}

	return res, scanner.Err()
}

func parseLine(line string, s *CurState) ([]painter.Operation, error) {
	fields := strings.Fields(line)
	const size = 400

	if len(fields) == 0 {
		return nil, nil
	}

	switch fields[0] {
	case "update":
		s.UpdateOp = painter.UpdateOp

	case "white":
		s.BgColorOp = painter.WhiteBackgroundOp(color.White)

	case "green":
		s.BgColorOp = painter.GreenBackgroundOp(color.RGBA{0, 255, 0, 255})

	case "bgrect":
		vals, err := parseFloatNum(fields, 4)
		if err != nil {
			return nil, err
		}

		op := &painter.BgRect{
			X1: int(vals[0] * size), Y1: int(vals[1] * size),
			X2: int(vals[2] * size), Y2: int(vals[3] * size),
		}

		s.BgRectFill = append(s.BgRectFill, op)

	case "figure":
		vals, err := parseFloatNum(fields, 2)

		if err != nil {
			return nil, err
		}
		s.Figures = append(s.Figures, &painter.Figure{
			X: int(vals[0] * size), Y: int(vals[1] * size),
		})

	case "move":
		vals, err := parseFloatNum(fields, 2)

		if err != nil {
			return nil, err
		}
		dx, dy := int(vals[0]*400), int(vals[1]*400)

		moveOp := painter.Move(dx, dy, s.Figures)
		s.MoveOp = append(s.MoveOp, moveOp)

		return []painter.Operation{moveOp}, nil

	case "reset":
		*s = *UpdateState()
		isReset = true
		return []painter.Operation{painter.Reset()}, nil

	default:
		return nil, fmt.Errorf("[Error]: unknown command: %s", fields[0])
	}

	return buildOps(s), nil
}

func parseFloatNum(fields []string, count int) ([]float64, error) {
	if len(fields) != count+1 {
		return nil, fmt.Errorf("[Error]: expected %d args, got %d", count, len(fields)-1)
	}
	values := make([]float64, count)
	for i := 0; i < count; i++ {
		v, err := strconv.ParseFloat(fields[i+1], 64)
		if err != nil {
			return nil, fmt.Errorf("invalid numeric value: %s", fields[i+1])
		}
		values[i] = v
	}
	return values, nil
}

func buildOps(s *CurState) []painter.Operation {
	var ops []painter.Operation

	if s.BgColorOp != nil {
		ops = append(ops, s.BgColorOp)
	} else {
		// For move command without figure bg is green
		s.BgColorOp = painter.GreenBackgroundOp(color.RGBA{0, 255, 0, 255})
		ops = append(ops, s.BgColorOp)
	}

	if len(s.BgRectFill) > 0 {
		// draw last rectangle
		ops = append(ops, s.BgRectFill[len(s.BgRectFill)-1])
	}
	for _, fig := range s.Figures {
		ops = append(ops, fig)
	}
	if s.UpdateOp != nil {
		ops = append(ops, s.UpdateOp)
	}

	return ops
}
