package fixers

import (
	"go/ast"
	"os"
	"strings"

	"golang.org/x/tools/go/ast/astutil"
)

func init() {
	AddFixer("use_path_join", func(options FixerOptions) (CsFixer, error) {
		return &UsePathJoinCsFixer{}, nil
	})
}

type UsePathJoinCsFixer struct {
}

func (l *UsePathJoinCsFixer) Lint(content string) (Problems, error) {
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

func (l *UsePathJoinCsFixer) Fix(content string) (string, error) {
	fset, file, err := ContentToAst(content)
	if err != nil {
		return content, err
	}

	wrongNodeCount := 0

	astutil.Apply(
		file,
		nil,
		func(cursor *astutil.Cursor) bool {
			if l.wrongNode(cursor.Node()) {
				wrongNodeCount++

				e, _ := cursor.Node().(*ast.CallExpr)

				e.Args[0] = l.processArg(e.Args[0])
			}

			return true
		},
	)

	if wrongNodeCount > 0 {
		astutil.AddImport(fset, file, "path")
	}

	return AstToContent(fset, file), nil
}

func (l *UsePathJoinCsFixer) wrongNode(n ast.Node) bool {
	selectors := map[string]bool{"os.Readlink": true}

	e, ok := n.(*ast.CallExpr)

	if !ok {
		return false // not a function call
	}

	selector, ok := e.Fun.(*ast.SelectorExpr)

	if !ok {
		// unusual selector, like interface{}(*to).(type)
		return false
	}

	ident, ok := selector.X.(*ast.Ident)

	if !ok {
		// unusual selector, like position := fset.Position(cursor.Node().Pos())
		return false
	}

	if _, ok := selectors[ident.Name+"."+selector.Sel.Name]; !ok {
		return false
	}

	return l.wrongArg(e.Args[0])
}

func (l *UsePathJoinCsFixer) wrongArg(n ast.Node) bool {
	binArg, ok := n.(*ast.BinaryExpr)

	if ok && binArg.Op.String() == "+" {
		// Left or right have path separator?
		return l.wrongArg(binArg.X) || l.wrongArg(binArg.Y)
	}

	arg, ok := n.(*ast.BasicLit)

	if !ok {
		return false
	}

	// Something like os.Readlink("foo/self") ?
	parts := strings.Split(arg.Value, string(os.PathSeparator))

	return len(parts) > 1
}

func (l *UsePathJoinCsFixer) processArg(n ast.Expr) ast.Expr {
	binArg, ok := n.(*ast.BinaryExpr)

	if ok {
		binArg.X = l.processArg(binArg.X)
		binArg.Y = l.processArg(binArg.Y)

		rightJoin := l.rightPathJoinCall(binArg.X)

		if rightJoin != nil {
			if l.emptyLeftPathJoinCall(binArg.Y) {
				rightJoin.(*ast.CallExpr).Args = append(rightJoin.(*ast.CallExpr).Args, binArg.Y.(*ast.CallExpr).Args[1:]...)
				return binArg.X
			}

			if l.emptyRightPathJoinCall(rightJoin) {
				rightArgs := rightJoin.(*ast.CallExpr).Args
				rightJoin.(*ast.CallExpr).Args = append(rightArgs[:len(rightArgs)-1], binArg.Y.(*ast.CallExpr).Args...)
				return binArg.X
			}
		} else if l.emptyLeftPathJoinCall(binArg.Y) {
			args := binArg.Y.(*ast.CallExpr).Args
			binArg.Y.(*ast.CallExpr).Args = append([]ast.Expr{binArg.X}, args[1:]...)
			return binArg.Y
		}

		return binArg
	}

	arg, ok := n.(*ast.BasicLit)

	if !ok {
		return n
	}

	parts := strings.Split(strings.Trim(arg.Value, "\""), string(os.PathSeparator))

	if len(parts) <= 1 {
		return n
	}

	var args []ast.Expr

	for _, part := range parts {
		args = append(args, &ast.BasicLit{Value: "\"" + part + "\""})
	}

	pathJoinCall := &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   &ast.Ident{Name: "path"},
			Sel: &ast.Ident{Name: "Join"},
		},
		Args: args,
	}

	return pathJoinCall
}

func (l *UsePathJoinCsFixer) pathJoinCall(n ast.Expr) bool {
	arg, ok := n.(*ast.CallExpr)

	if !ok {
		return false
	}

	selector, ok := arg.Fun.(*ast.SelectorExpr)

	if !ok {
		return false
	}

	ident, ok := selector.X.(*ast.Ident)

	if !ok {
		return false
	}

	if ident.Name != "path" {
		return false
	}

	return selector.Sel.Name == "Join"
}

func (l *UsePathJoinCsFixer) rightPathJoinCall(n ast.Expr) ast.Expr {
	binArg, ok := n.(*ast.BinaryExpr)

	if ok {
		return l.rightPathJoinCall(binArg.Y)
	}

	if !l.pathJoinCall(n) {
		return nil
	}

	return n
}

func (l *UsePathJoinCsFixer) emptyLeftPathJoinCall(n ast.Expr) bool {
	if !l.pathJoinCall(n) {
		return false
	}

	arg, ok := n.(*ast.CallExpr).Args[0].(*ast.BasicLit)

	return ok && arg.Value == `""`
}

func (l *UsePathJoinCsFixer) emptyRightPathJoinCall(n ast.Expr) bool {
	if !l.pathJoinCall(n) {
		return false
	}

	arg, ok := n.(*ast.CallExpr).Args[len(n.(*ast.CallExpr).Args)-1].(*ast.BasicLit)

	return ok && arg.Value == `""`
}

func (l *UsePathJoinCsFixer) String() string {
	return "Use path.Join"
}
