package fixers

import "fmt"

type FixerOptions map[interface{}]interface{}
type FixerCreateFunc func (options FixerOptions) (CsFixer, error)

var FixersMap map[string]FixerCreateFunc

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
	return fmt.Sprintf("at line %v: %s", p.Position.Line, p.Text)
}

type Problems []*Problem

func NewPosition(line int) *Position {
	return &Position{Line: line}
}

type Position struct {
	Line int
}

func AddFixer(name string, createFunc FixerCreateFunc)  {
	if FixersMap == nil {
		FixersMap = make(map[string]FixerCreateFunc, 0)
	}

	FixersMap[name] = createFunc
}

