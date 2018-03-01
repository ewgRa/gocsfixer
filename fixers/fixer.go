package fixers

import "fmt"

var FixersMap map[string]func () CsFixer

func init()  {
	FixersMap = make(map[string]func () CsFixer, 0)
}

type CsFixer interface {
}

type Linter interface {
	Lint(content string) (Problems, error)
}

type Fixer interface {
	Fix(content string) (string, error)
}

type Problem struct {
	Position *Position
	Text string
	LineText string

}

func (p *Problem) String() string {
	return fmt.Sprintf("at line %v: %s", p.Position.line, p.Text)
}

type Problems []*Problem

func NewPosition(line int) *Position {
	return &Position{line: line}
}

type Position struct {
	line int
}