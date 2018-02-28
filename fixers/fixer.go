package fixers

import "fmt"

type CsFixer interface {
}

type Linter interface {
	Lint(content string) (Problems, error)
}

type Fixer interface {
	Fix(content string) (Problems, string, error)
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