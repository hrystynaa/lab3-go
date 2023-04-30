package lang

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/hrystynaa/lab3-go/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
	bgRect  *painter.BgRectangle
	backOp  painter.Operation
	move    painter.Operation
	figures []*painter.Figure
	res     []painter.Operation
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	p.bgRect, p.backOp, p.figures, p.res = nil, nil, nil, nil

	for scanner.Scan() {
		commandLine := scanner.Text()
		if len(commandLine) == 0 {
			continue
		}
		err := p.parse(commandLine) // parse the line to get Operation
		if err != nil {
			return nil, err
		}

		if len(p.res) == 1 {
			if p.backOp != nil {
				p.res = append(p.res, p.backOp)
			}

			if p.bgRect != nil {
				p.res = append(p.res, p.bgRect)
			}

			if p.move != nil {
				p.res = append(p.res, p.move)
			}

			for _, figure := range p.figures {
				p.res = append(p.res, figure)
			}
			p.res = append(p.res[1:], p.res[0])
		}
	}
	return p.res, scanner.Err()
}

func (p *Parser) parse(commandLine string) error {
	fields := strings.Fields(commandLine)
	operation := fields[0]
	var args []int

	for i := 1; i < len(fields); i++ {
		arg, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return err
		}
		arg = arg * 800.0
		args = append(args, int(arg))
	}
	switch operation {
	case "white":
		p.backOp = painter.OperationFunc(painter.WhiteFill)
	case "green":
		p.backOp = painter.OperationFunc(painter.GreenFill)
	case "update":
		p.res = p.res[:0]
		p.res = append(p.res, painter.UpdateOp)
	case "bgrect":
		p.bgRect = &painter.BgRectangle{X1: args[0], Y1: args[1], X2: args[2], Y2: args[3]}
	case "figure":
		figure := &painter.Figure{X: args[0], Y: args[1]}
		p.figures = append(p.figures, figure)
	case "move":
		p.move = &painter.Move{X: args[0], Y: args[1], Figures: p.figures}
	case "reset":
		p.figures = p.figures[:0]
		p.bgRect = nil
		p.backOp = painter.OperationFunc(painter.ResetScreen)
	default:
		return errors.New("Failed")
	}
	return nil
}
