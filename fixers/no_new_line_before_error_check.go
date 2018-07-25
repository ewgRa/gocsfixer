package fixers

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"strings"
)

func init() {
	AddFixer("no_new_line_before_error_check", func(options FixerOptions) (CsFixer, error) {
		return &NoNewLineBeforeErrorCsFixer{}, nil
	})
}

type NoNewLineBeforeErrorCsFixer struct {
}

func (l *NoNewLineBeforeErrorCsFixer) Lint(content string) (Problems, error) {
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
			e, ok := cursor.Node().(*ast.BinaryExpr)

			if !ok {
				return true // not a binary operation
			}

			if e.Op != token.EQL && e.Op != token.NEQ {
				return true // not a comparison
			}

			xContent := AstToContent(fset, e.X)
			yContent := AstToContent(fset, e.Y)

			if xContent == "err" || yContent == "err" {
				position := fset.Position(e.Pos())

				if position.Line >= 2 && lines[position.Line-2] == "" {
					problems = append(problems, &Problem{Position: NewPosition(position.Line - 1), Text: l.String(), LineText: lines[position.Line-2]})
				}
			}

			return true
		},
	)

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

	for i := len(problems) - 1; i >= 0; i-- {
		lines = append(lines[:problems[i].Position.Line-1], lines[problems[i].Position.Line:]...)
	}

	return strings.Join(lines, "\n"), nil
}

func (l *NoNewLineBeforeErrorCsFixer) String() string {
	return "No new line before error check"
}
