package fixers

import (
	"go/token"
	"go/parser"
	"bytes"
	"go/format"
	"go/ast"
	"strings"
)

func init() {
	AddFixer("no_new_line_before_error_check", func (options FixerOptions) (CsFixer, error) {
		return &NoNewLineBeforeErrorCsFixer{}, nil
	})
}

type NoNewLineBeforeErrorCsFixer struct {
	positions []token.Pos
	fset *token.FileSet
}

func (l *NoNewLineBeforeErrorCsFixer) Lint(content string) (Problems, error) {
	l.positions = []token.Pos{}

	l.fset = token.NewFileSet()

	file, err := parser.ParseFile(l.fset, "", content, parser.ParseComments)
	if err != nil {
		return Problems{}, err
	}

	var buf bytes.Buffer
	format.Node(&buf, l.fset, file)

	ast.Inspect(file, l.check)

	lines := strings.Split(content, "\n")

	var problems []*Problem

	for _, tokenPos := range l.positions {
		position := l.fset.Position(tokenPos)

		if position.Line < 2 {
			continue
		}

		if lines[position.Line-2] == "" {
			problems = append(problems, &Problem{Position: NewPosition(position.Line-1), Text: "No newline before check error", LineText: lines[position.Line-2]})
		}
	}

	return problems, nil
}

func (l *NoNewLineBeforeErrorCsFixer) Fix(content string) (string, error) {
	problems, err := l.Lint(content)

	if nil != err {
		return "", err
	}

	if len(problems) == 0 {
		return content, nil
	}

	lines := strings.Split(content, "\n")

	for i := len(problems)-1; i >=0; i-- {
		lines = append(lines[:problems[i].Position.Line-1], lines[problems[i].Position.Line:]...)
	}

	return strings.Join(lines, "\n"), nil
}

func (l *NoNewLineBeforeErrorCsFixer) check(n ast.Node) bool {
	e, ok := n.(*ast.BinaryExpr)
	if !ok {
		return true // not a binary operation
	}
	if e.Op != token.EQL && e.Op != token.NEQ {
		return true // not a comparison
	}

	var buf bytes.Buffer
	format.Node(&buf, l.fset, e.X)

	var bufY bytes.Buffer
	format.Node(&bufY, l.fset, e.Y)

	if (buf.String() == "err" || bufY.String() == "err") {
		l.positions = append(l.positions, e.Pos())
	}

	return true
}

func (l *NoNewLineBeforeErrorCsFixer) String() string {
	return "No new line before error check"
}
