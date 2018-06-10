package fixers

import (
	"go/token"
	"go/ast"
	"bytes"
	"strings"
	"go/parser"
	"go/format"
	"os"
)

func init() {
	FixersMap["use_path_join"] = func () CsFixer {
		return &UsePathJoinCsFixer{}
	}
}

type UsePathJoinCsFixer struct {
	positions []token.Pos
	fset *token.FileSet
}

func (l *UsePathJoinCsFixer) Lint(content string) (Problems, error) {
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
		problems = append(problems, &Problem{Position: NewPosition(position.Line), Text: "Use path.Join", LineText: lines[position.Line-1]})
	}

	return problems, nil
}

func (l *UsePathJoinCsFixer) Fix(content string) (string, error) {
	/// FIXME XXX: implement me
	return content, nil
}

func (l *UsePathJoinCsFixer) check(n ast.Node) bool {
	selectors := map[string]bool {"os.Readlink": true}

	e, ok := n.(*ast.CallExpr)

	if !ok {
		return true // not a binary operation
	}

	selector := e.Fun.(*ast.SelectorExpr)
	ident := selector.X.(*ast.Ident)

	if _, ok := selectors[ident.Name + "." + selector.Sel.Name]; !ok {
		return true
	}

	binArg, ok := e.Args[0].(*ast.BinaryExpr)

	if ok && binArg.Op.String() == "+" {
		// Something like os.Readlink(gosigar.Procd + "self")
		l.positions = append(l.positions, e.Pos())
		return true
	}

	arg, ok := e.Args[0].(*ast.BasicLit)

	if ok {
		// Something like os.Readlink("foo/self")
		parts := strings.Split(arg.Value, string(os.PathSeparator))

		if len(parts) > 1 {
			l.positions = append(l.positions, e.Pos())
		}

		return true
	}

	return true
}

func (l *UsePathJoinCsFixer) String() string {
	return "Use path.Join"
}
