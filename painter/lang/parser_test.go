package lang

import (
	"strings"
	"testing"

	"github.com/hrystynaa/lab3-go/painter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parse_func(t *testing.T) {
	tests := []struct {
		name    string
		command string
		op      painter.Operation
	}{
		{
			name:    "white command",
			command: "white\nupdate",
			op:      painter.OperationFunc(painter.WhiteFill),
		},
		{
			name:    "green command",
			command: "green\nupdate",
			op:      painter.OperationFunc(painter.GreenFill),
		},
		{
			name:    "reset command",
			command: "reset\nupdate",
			op:      painter.OperationFunc(painter.ResetScreen),
		},
	}
	parser := &Parser{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ops, err := parser.Parse(strings.NewReader(tt.command))
			require.NoError(t, err)
			assert.IsType(t, tt.op, ops[0])
		})
	}
}

func Test_parse_struct(t *testing.T) {
	tests := []struct {
		name    string
		command string
		op      painter.Operation
	}{
		{
			name:    "bgrect command",
			command: "bgrect 0.25 0.25 0.75 0.75\nupdate ",
			op:      &painter.BgRectangle{X1: 200, Y1: 200, X2: 600, Y2: 600},
		},
		{
			name:    "figure command",
			command: "figure 0.5 0.5\nupdate",
			op:      &painter.Figure{X: 400, Y: 400},
		},
		{
			name:    "move command",
			command: "move 0.3 0.3\nupdate",
			op:      &painter.Move{X: 240, Y: 240, Figures: []*painter.Figure(nil)},
		},
		{
			name:    "update command",
			command: "update",
			op:      painter.UpdateOp,
		},
		{
			name:    "invalid command",
			command: "invalidcommand\nupdate",
			op:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser := &Parser{}
			ops, err := parser.Parse(strings.NewReader(tt.command))
			if tt.op == nil {
				assert.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.IsType(t, tt.op, ops[1])
				assert.Equal(t, tt.op, ops[1])
			}
		})
	}
}
