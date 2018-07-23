package fixers

import (
	"bytes"
	"errors"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"strconv"
	"strings"
)

func init() {
	AddFixer("group_import", func(options FixerOptions) (CsFixer, error) {
		lintText, err := options.extractString("lintText", "Group import")

		if err != nil {
			return nil, err
		}

		groupStdLib, err := options.extractBool("stdLib", false)

		if err != nil {
			return nil, err
		}

		return &GroupImportFixer{lintText: lintText, groupStdLib: groupStdLib}, nil
	})
}

type GroupImportFixer struct {
	lintText    string
	groupStdLib bool
}

func (l *GroupImportFixer) Lint(content string) (Problems, error) {
	var problems Problems

	if !l.groupStdLib {
		return problems, nil
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return problems, err
	}

	lines := strings.Split(content, "\n")

	importList := astutil.Imports(fset, file)

	// astutil.Imports think that imports that have docs start new group, we consider it as one group, so, we need merge them
	importList = l.mergeDocGroups(fset, importList)

	firstGroupMixed := false

	if len(importList) > 0 {
		firstGroupMixed = l.mixedGroup(importList[0])

		// Check that it is not "alone" 'import "stdlib"' that for us is invalid in any case
		if !firstGroupMixed && len(importList[0]) == 1 && l.isStdLibImport(file.Imports[0].Path.Value) {
			for i := 0; i < len(file.Decls); i++ {
				decl := file.Decls[i]
				gen, ok := decl.(*ast.GenDecl)
				if !ok || gen.Tok != token.IMPORT {
					continue
				}

				if !gen.Lparen.IsValid() {
					firstGroupMixed = true
				}

				break
			}
		}
	}

	for k, importSpecs := range importList {
		for _, importSpec := range importSpecs {
			if l.isStdLibImport(importSpec.Path.Value) && (firstGroupMixed || k != 0) {
				line := fset.Position(importSpec.Pos()).Line

				problems = append(problems, &Problem{
					Position: NewPosition(line),
					Text:     l.lintText,
					LineText: lines[line],
				})
			}
		}
	}

	return problems, nil
}

func (l *GroupImportFixer) Fix(content string) (string, error) {
	if !l.groupStdLib {
		return content, nil
	}

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var stdLibImports []*ast.ImportSpec

	for _, importSpec := range file.Imports {
		if l.isStdLibImport(importSpec.Path.Value) {
			stdLibImports = append(stdLibImports, importSpec)
		}
	}

	for _, importSpec := range stdLibImports {
		path, err := strconv.Unquote(importSpec.Path.Value)

		if err != nil {
			return "", err
		}

		if importSpec.Name != nil {
			if !astutil.DeleteNamedImport(fset, file, importSpec.Name.Name, path) {
				return "", errors.New("Can't delete named " + importSpec.Name.Name + " import " + importSpec.Path.Value)
			}
		} else {
			if !astutil.DeleteImport(fset, file, path) {
				return "", errors.New("Can't delete import " + importSpec.Path.Value)
			}
		}

		// Delete comments and docs as well, we will restore it later with same imports
		if importSpec.Comment != nil {
			for i, cg := range file.Comments {
				if cg == importSpec.Comment {
					copy(file.Comments[i:], file.Comments[i+1:])
					file.Comments = file.Comments[:len(file.Comments)-1]
					break
				}
			}
		}

		if importSpec.Doc != nil {
			for i, cg := range file.Comments {
				if cg == importSpec.Doc {
					copy(file.Comments[i:], file.Comments[i+1:])
					file.Comments = file.Comments[:len(file.Comments)-1]
					break
				}
			}
		}
	}

	if len(stdLibImports) > 0 {
		f := file
		lastImport := -1

		// Mostly taken from astutil/imports.go::AddNamedImport
		impDecl := &ast.GenDecl{
			Tok: token.IMPORT,
		}

		impDecl.TokPos = f.Package

		ffile := fset.File(f.Package)
		pkgLine := ffile.Line(f.Package)
		for _, c := range f.Comments {
			if ffile.Line(c.Pos()) > pkgLine {
				break
			}
			impDecl.TokPos = c.End()
		}
		f.Decls = append(f.Decls, nil)
		copy(f.Decls[lastImport+2:], f.Decls[lastImport+1:])
		f.Decls[lastImport+1] = impDecl

		for _, newImport := range stdLibImports {
			// Mostly taken from astutil/imports.go::AddNamedImport

			// Insert new import at insertAt.
			insertAt := 0
			impDecl.Specs = append(impDecl.Specs, nil)
			copy(impDecl.Specs[insertAt+1:], impDecl.Specs[insertAt:])
			impDecl.Specs[insertAt] = newImport
			pos := impDecl.Pos()

			if newImport.Name != nil {
				newImport.Name.NamePos = pos
			}
			newImport.Path.ValuePos = pos
			newImport.EndPos = pos

			impDecl.Lparen = impDecl.Specs[0].Pos()

			file.Imports = append(file.Imports, newImport)
		}
	}

	var buf bytes.Buffer
	format.Node(&buf, fset, file)

	res := buf.String()

	/*
		rexp := regexp.MustCompile("(?s)import \\((.*?)\\)")
		importS := rexp.Find(buf.Bytes())

		fmt.Println(string(importS))

	*/
	res = strings.Replace(res, ")\n\nimport (", "", 1)

	if len(stdLibImports) > 0 {
		// Restore comment and doc blocks

		fset = token.NewFileSet()

		file, err = parser.ParseFile(fset, "", res, parser.ParseComments)
		if err != nil {
			return "", err
		}

		for _, impSpec := range stdLibImports {
			if impSpec.Comment != nil || impSpec.Doc != nil {
				for _, newImportSpec := range file.Imports {
					if newImportSpec.Path.Value != impSpec.Path.Value {
						continue
					}

					if (newImportSpec.Name != nil && impSpec.Name == nil) || (newImportSpec.Name == nil && impSpec.Name != nil) || (newImportSpec.Name != nil && impSpec.Name != nil && newImportSpec.Name.Name != impSpec.Name.Name) {
						continue
					}

					if impSpec.Comment != nil {
						newImportSpec.Comment = impSpec.Comment
						for _, cmt := range newImportSpec.Comment.List {
							cmt.Slash = newImportSpec.Path.End()
						}

						file.Comments = append(file.Comments, impSpec.Comment)
					}

					if impSpec.Doc != nil {
						newImportSpec.Doc = impSpec.Doc

						for _, cmt := range newImportSpec.Doc.List {
							cmt.Slash = newImportSpec.Pos() - 1
						}

						file.Comments = append(file.Comments, impSpec.Doc)
					}
				}
			}
		}

		// Sort comments, looks like it is important for printer
		cmts := file.Comments
		file.Comments = []*ast.CommentGroup{}

		for len(cmts) > 0 {
			var minCmt *ast.CommentGroup
			var minCmtI int

			for i, cmt := range cmts {
				if minCmt == nil || minCmt.Pos() > cmt.Pos() {
					minCmt = cmt
					minCmtI = i
				}
			}

			cmts = append(cmts[:minCmtI], cmts[minCmtI+1:]...)
			file.Comments = append(file.Comments, minCmt)
		}
		/*
			for i := 0; i < len(file.Decls); i++ {
				decl := file.Decls[i]
				gen, ok := decl.(*ast.GenDecl)
				if !ok || gen.Tok != token.IMPORT {
					continue
				}

				for _, impSpec := range stdLibImports {
					if impSpec.Comment != nil {
						file.Comments = append(file.Comments, stdLibImports[1].Comment)
						file.Imports[0].Comment = stdLibImports[1].Comment
						file.Imports[0].Comment.List[0].Slash = file.Imports[0].End()
						break
					}
				}

				break
			}
		*/
		buf.Reset()
		format.Node(&buf, fset, file)
		res = buf.String()
		//fmt.Println(res)
	}

	return res, nil
}

func (l *GroupImportFixer) isStdLibImport(importPath string) bool {
	// Taken from imports.go::isThirdParty
	return !strings.Contains(importPath, ".")
}

func (l *GroupImportFixer) mixedGroup(specs []*ast.ImportSpec) bool {
	haveStdImport := false
	haveNotStdImport := false

	for _, importSpec := range specs {
		if l.isStdLibImport(importSpec.Path.Value) {
			haveStdImport = true
		} else {
			haveNotStdImport = true
		}

		if haveStdImport && haveNotStdImport {
			return true
		}
	}

	return false
}

func (l *GroupImportFixer) mergeDocGroups(fset *token.FileSet, importGroups [][]*ast.ImportSpec) [][]*ast.ImportSpec {
	var res [][]*ast.ImportSpec

	for i, importGroup := range importGroups {
		if i == 0 || importGroup[0].Doc == nil {
			res = append(res, importGroup)
			continue
		}

		prevImportGroup := importGroups[i-1]

		if fset.Position(prevImportGroup[len(prevImportGroup)-1].End()).Line == fset.Position(importGroup[0].Doc.Pos()).Line-1 {
			res[len(res)-1] = append(res[len(res)-1], importGroup...)
		} else {
			res = append(res, importGroup)
		}
	}

	return res
}

func (l *GroupImportFixer) String() string {
	return "Group import"
}
