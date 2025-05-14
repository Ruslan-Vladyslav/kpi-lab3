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
	State *painter.State
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	var res []painter.Operation

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		op, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("[Error]: parse error on line '%s': %w", line, err)
		}

		res = append(res, op)
	}

	return res, scanner.Err()
}

func parseLine(line string) (painter.Operation, error) {
	fields := strings.Fields(line)
	if len(fields) == 0 {
		return nil, nil
	}

	switch fields[0] {
	case "update":
		return painter.UpdateOp, nil
	case "white":
		return painter.FillInColorOp{Color: color.White}, nil

	case "green":
		return painter.FillInColorOp{Color: color.RGBA{G: 0xff, A: 0xff}}, nil

	case "rect":
		if len(fields) != 5 {
			return nil, fmt.Errorf("[Error]: invalid rect args, expected 4 got %d", len(fields)-1)
		}
		x0 := parseFloat(fields[1])
		y0 := parseFloat(fields[2])
		x1 := parseFloat(fields[3])
		y1 := parseFloat(fields[4])

		return painter.BgRectOp{
			X1: x0,
			Y1: y0,
			X2: x1,
			Y2: y1,
		}, nil

	case "point":
		if len(fields) != 3 {
			return nil, fmt.Errorf("[Error]: invalid point args")
		}
		x := parseFloat(fields[1])
		y := parseFloat(fields[2])
		return painter.FigureOp{
			X: x,
			Y: y,
		}, nil

	case "move":
		if len(fields) != 3 {
			return nil, fmt.Errorf("[Error]: invalid move args")
		}
		dx := parseInt(fields[1])
		dy := parseInt(fields[2])
		return painter.MoveOp{
			MX: dx,
			MY: dy,
		}, nil

	case "reset":
		return painter.ResetOp{}, nil

	default:
		return nil, fmt.Errorf("unknown command: %s", fields[0])
	}
}

func parseFloat(s string) float64 {
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
