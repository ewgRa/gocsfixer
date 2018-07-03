package fixers

import (
	"go/token"
	"go/ast"
	"bytes"
	"strings"
	"go/parser"
	"go/format"
	"os"
	"golang.org/x/tools/go/ast/astutil"
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
	l.fset = token.NewFileSet()

	file, err := parser.ParseFile(l.fset, "", content, parser.ParseComments)
	if err != nil {
		return content, err
	}

	var wrongNodeCount int

	astutil.Apply(
		file,
		nil,
		func(cursor *astutil.Cursor) bool {
			if l.isWrongNode(cursor.Node()) {
				wrongNodeCount++

				e, _ := cursor.Node().(*ast.CallExpr)

				pathJoinCall := &ast.CallExpr{
					Fun: &ast.SelectorExpr{
						X: &ast.Ident{Name: "path"},
						Sel: &ast.Ident{Name: "Join"},
					},
					Args: l.processArg(e.Args[0]),
				}

				e.Args = []ast.Expr{pathJoinCall}
			}

			return true
		},
	)

	if (wrongNodeCount > 0) {
		astutil.AddImport(l.fset, file, "path")
	}

	var buf bytes.Buffer
	format.Node(&buf, l.fset, file)

	return buf.String(), nil
}

func (l *UsePathJoinCsFixer) check(n ast.Node) bool {
	if l.isWrongNode(n) {
		l.positions = append(l.positions, n.Pos())
	}

	return true
}

func (l *UsePathJoinCsFixer) isWrongNode(n ast.Node) bool {
	selectors := map[string]bool {"os.Readlink": true}

	e, ok := n.(*ast.CallExpr)

	if !ok {
		return false // not a function call
	}

	selector := e.Fun.(*ast.SelectorExpr)
	ident := selector.X.(*ast.Ident)

	if _, ok := selectors[ident.Name + "." + selector.Sel.Name]; !ok {
		return false
	}

	binArg, ok := e.Args[0].(*ast.BinaryExpr)

	if ok && binArg.Op.String() == "+" {
		return true
	}

	arg, ok := e.Args[0].(*ast.BasicLit)

	if ok {
		// Something like os.Readlink("foo/self")
		parts := strings.Split(arg.Value, string(os.PathSeparator))

		if len(parts) > 1 {
			return true
		}
	}

	return false
}

func (l *UsePathJoinCsFixer) processArg(n ast.Node) []ast.Expr {
	var result []ast.Expr

	binArg, ok := n.(*ast.BinaryExpr)

	if ok {
		result = append(result, l.processArg(binArg.X)...)
		result = append(result, l.processArg(binArg.Y)...)
		return result
	}

	arg, ok := n.(*ast.BasicLit)

	if ok {
		parts := strings.Split(strings.Trim(arg.Value, "\""), string(os.PathSeparator))

		for _, part := range parts {
			result = append(result, &ast.BasicLit{Value: "\"" + part + "\""})
		}

		return result
	}

	result = append(result, n.(ast.Expr))

	return result
}

func (l *UsePathJoinCsFixer) String() string {
	return "Use path.Join"
}
