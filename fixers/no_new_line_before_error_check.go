package fixers

import (
	"go/token"
	"go/parser"
	"bytes"
	"go/format"
	"go/ast"
	"strings"
	"sort"
)

type NoNewLineBeforeErrorCheck struct {
	positions []token.Pos
	fset *token.FileSet
}

func (l *NoNewLineBeforeErrorCheck) Lint(content string) (Problems, error) {
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

	var checkLines []int
	checkLinesMap := make(map[int]token.Position, 0)

	for _, tokenPos := range l.positions {
		position := l.fset.Position(tokenPos)
		checkLines = append(checkLines, position.Line)
		checkLinesMap[position.Line] = position
	}

	sort.Sort(sort.Reverse(sort.IntSlice(checkLines)))

	var problems []*Problem

	for _, line := range checkLines {
		if line < 2 {
			continue
		}

		if lines[line-2] == "" {
			problems = append(problems, &Problem{Position: NewPosition(checkLinesMap[line].Line), Text: "No newline before check error", LineText: lines[line-1]})
		}
	}

	return problems, nil
}

func (l *NoNewLineBeforeErrorCheck) Fix(content string) (Problems, string) {
	return Problems{}, content
}

func (l *NoNewLineBeforeErrorCheck) check(n ast.Node) bool {
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

func (l *NoNewLineBeforeErrorCheck) String() string {
	return "No new line before error check"
}
