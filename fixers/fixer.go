package fixers

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/parser"
	"bytes"
	"go/format"
)

type FixerOptions map[interface{}]interface{}

func (o FixerOptions) extractRequiredString(name string) (string, error) {
	_, ok := o[name]

	if !ok {
		return "", errors.New(name + " option is required")
	}

	return o.extractString(name, "")
}

func (o FixerOptions) extractString(name string, defaultValue string) (string, error) {
	v, ok := o[name]

	if !ok {
		return defaultValue, nil
	}

	value, ok := v.(string)

	if !ok {
		return "", errors.New("Wrong " + name + " option")
	}

	return value, nil
}

func (o FixerOptions) extractBool(name string, defaultValue bool) (bool, error) {
	v, ok := o[name]

	if !ok {
		return defaultValue, nil
	}

	value, ok := v.(bool)

	if !ok {
		return false, errors.New("Wrong " + name + " option")
	}

	return value, nil
}

type FixerCreateFunc func(options FixerOptions) (CsFixer, error)

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
	Text     string
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

func AddFixer(name string, createFunc FixerCreateFunc) {
	if FixersMap == nil {
		FixersMap = make(map[string]FixerCreateFunc, 0)
	}

	FixersMap[name] = createFunc
}

func CreateFixer(name string, options FixerOptions) (CsFixer, error) {
	if createFunc, ok := FixersMap[name]; ok {
		fixer, err := createFunc(options)

		if err != nil {
			return nil, err
		}

		return fixer, nil
	}

	return nil, nil
}

func ContentToAst(content string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	return fset, file, err
}

func AstToContent(fset *token.FileSet, node interface{}) string {
	var buf bytes.Buffer
	format.Node(&buf, fset, node)
	return buf.String()
}