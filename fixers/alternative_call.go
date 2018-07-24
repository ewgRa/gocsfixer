package fixers

import (
	"go/ast"
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
}

func (l *AlternativeCallCsFixer) Lint(content string) (Problems, error) {
	problems := Problems{}

	fset, file, err := ContentToAst(content)
	if err != nil {
		return problems, err
	}

	lines := strings.Split(content, "\n")

	astutil.Apply(
		file,
		nil,
		func(cursor *astutil.Cursor) bool {
			if l.wrongNode(cursor.Node()) {
				position := fset.Position(cursor.Node().Pos())
				problems = append(problems, &Problem{Position: NewPosition(position.Line), Text: l.String(), LineText: lines[position.Line-1]})
			}

			return true
		},
	)

	return problems, nil
}

func (l *AlternativeCallCsFixer) Fix(content string) (string, error) {
	fset, file, err := ContentToAst(content)
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

	return AstToContent(fset, file), nil
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
