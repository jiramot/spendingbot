package rest

import "github.com/jiramot/spendingbot/line"

type rest struct {
	line              *line.Line
	lineChannelSecret string
}

func NewRestHandler(line *line.Line, lineChannelSecret string) *rest {
	return &rest{line: line, lineChannelSecret: lineChannelSecret}
}
