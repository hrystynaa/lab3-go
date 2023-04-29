package lang

import (
	"io"

	"github.com/hrystynaa/lab3-go/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation

	// TODO: Реалізувати парсинг команд.
	res = append(res, painter.OperationFunc(painter.WhiteFill))
	res = append(res, painter.UpdateOp)

	return res, nil
}
