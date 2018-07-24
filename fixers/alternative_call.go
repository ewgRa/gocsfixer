package fixers

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"strings"
)

func init() {
	AddFixer("alternative_call", func(options FixerOptions) (CsFixer, error) {
		selector, err := options.extractRequiredString("selector")

		if err != nil {
			return nil, err
		}

		alternative, err := options.extractRequiredString("alternative")

		if err != nil {
			return nil, err
		}

		return &AlternativeCallCsFixer{selector: selector, alternative: alternative}, nil
	})
}

type AlternativeCallCsFixer struct {
	selector    string
	alternative string
	positions   []token.Pos
	fset        *token.FileSet
}

func (l *AlternativeCallCsFixer) Lint(content string) (Problems, error) {
	l.positions = []token.Pos{}

	l.fset = token.NewFileSet()

	file, err := parser.ParseFile(l.fset, "", content, parser.ParseComments)
	if err != nil {
		return Problems{}, err
	}

	var buf bytes.Buffer
	format.Node(&buf, l.fset, file)

	ast.Inspect(file, l.inspect)

	lines := strings.Split(content, "\n")

	var problems []*Problem

	for _, tokenPos := range l.positions {
		position := l.fset.Position(tokenPos)
		problems = append(problems, &Problem{Position: NewPosition(position.Line), Text: l.String(), LineText: lines[position.Line-1]})
	}

	return problems, nil
}

func (l *AlternativeCallCsFixer) Fix(content string) (string, error) {
	l.fset = token.NewFileSet()

	file, err := parser.ParseFile(l.fset, "", content, parser.ParseComments)
	if err != nil {
		return content, err
	}

	astutil.Apply(
		file,
		nil,
		func(cursor *astutil.Cursor) bool {
			if l.wrongNode(cursor.Node()) {
				e, _ := cursor.Node().(*ast.CallExpr)
				l.processArg(e)
			}

			return true
		},
	)

	var buf bytes.Buffer
	format.Node(&buf, l.fset, file)

	return buf.String(), nil
}

func (l *AlternativeCallCsFixer) inspect(n ast.Node) bool {
	if l.wrongNode(n) {
		l.positions = append(l.positions, n.Pos())
	}

	return true
}

func (l *AlternativeCallCsFixer) wrongNode(n ast.Node) bool {
	e, ok := n.(*ast.CallExpr)

	if !ok {
		return false // not a function call
	}

	selector, ok := e.Fun.(*ast.SelectorExpr)

	if !ok {
		return false
	}

	ident := selector.X.(*ast.Ident)

	return l.selector == ident.Name+"."+selector.Sel.Name
}

func (l *AlternativeCallCsFixer) processArg(n *ast.CallExpr) {
	alternativePart := strings.Split(l.alternative, ".")

	selector := n.Fun.(*ast.SelectorExpr)
	selector.Sel.Name = alternativePart[1]
	selector.X.(*ast.Ident).Name = alternativePart[0]
}

func (l *AlternativeCallCsFixer) String() string {
	return "Instead " + l.selector + " use alternative call " + l.alternative
}
