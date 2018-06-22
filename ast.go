package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
)

type ASTFile struct {
	src        []byte
	importPath string

	fset *token.FileSet
	f    *ast.File

	fileNode    *ast.File
	nameNode    *ast.Ident
	nameNodePos token.Position
	cmap        ast.CommentMap

	changed bool
}

func NewASTFile(src []byte, importPath string) *ASTFile {
	a := &ASTFile{
		src:        src,
		importPath: importPath,
	}

	return a

}

func (a *ASTFile) Init() error {
	fset := token.NewFileSet() // positions are relative to fset
	f, err := parser.ParseFile(fset, "src.go", a.src, parser.ParseComments)
	if err != nil {
		return err
	}

	a.fset = fset
	a.f = f

	ast.Inspect(f, func(n ast.Node) bool {
		if fn, ok := n.(*ast.File); ok {
			a.fileNode = fn
			return false
		}
		return true
	})

	a.nameNode = a.fileNode.Name
	a.nameNodePos = fset.Position(a.nameNode.End())

	cmap := ast.NewCommentMap(fset, f, f.Comments)
	if cmap == nil {
		cmap = make(ast.CommentMap)
	}
	a.cmap = cmap

	return nil
}

func (a *ASTFile) EnsureImportComment() error {
	if err := a.Init(); err != nil {
		return err
	}

	if cgs, ok := a.cmap[a.nameNode]; ok {
		for _, cg := range cgs {
			pos := a.fset.Position(cg.Pos())
			if pos.Line == a.nameNodePos.Line {
				for i, c := range cg.List {
					text := strings.TrimSpace(c.Text)
					//TODO(anarcher): need more checks
					if !strings.Contains(text, "import") {
						a.appendImportComment(a.nameNode, cg)
						a.changed = true
					} else if !a.isSameImportComment(text) {
						a.updateImportComment(a.nameNode, cg, i)
						a.changed = true

					}
				}
			}
		}
	} else {
		a.addImportComment(a.nameNode)
		a.changed = true
	}

	a.f.Comments = a.cmap.Comments()

	return nil
}

func (a *ASTFile) String() string {
	return string(a.Bytes())
}

func (a *ASTFile) Bytes() []byte {
	var buf bytes.Buffer
	if err := format.Node(&buf, a.fset, a.f); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func (a *ASTFile) updateImportComment(node *ast.Ident, cg *ast.CommentGroup, updateIdx int) {
	comment := &ast.Comment{
		Slash: token.Pos(node.End()),
		Text:  a.importComment(),
	}

	cg.List[updateIdx] = comment
}

func (a *ASTFile) appendImportComment(node *ast.Ident, cg *ast.CommentGroup) {
	comment := &ast.Comment{
		Slash: token.Pos(node.End()),
		Text:  a.importComment(),
	}

	var comments []*ast.Comment
	comments = append(comments, comment)

	var pos token.Pos = comment.End()
	for _, c := range cg.List {
		c.Slash = pos
		comments = append(comments, c)
	}

	cg.List = comments
}

func (a *ASTFile) addImportComment(node *ast.Ident) {
	a.cmap[node] = []*ast.CommentGroup{
		{
			List: []*ast.Comment{
				{
					Slash: token.Pos(node.End()),
					Text:  a.importComment(),
				},
			},
		},
	}
}

func (a *ASTFile) IsChanged() bool {
	return a.changed
}

func (a *ASTFile) importComment() string {
	return fmt.Sprintf("// %s", a.importStatement())
}

func (a *ASTFile) importStatement() string {
	return fmt.Sprintf("import \"%s\"", a.importPath)
}

func (a *ASTFile) isSameImportComment(comment string) bool {
	x := strings.TrimSpace(a.importComment())
	y := strings.TrimSpace(comment)

	if strings.Compare(x, y) == 0 {
		return true
	}

	return false
}
