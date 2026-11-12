package pfcp_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestExportedDocumentation keeps the public API discoverable in go doc.
func TestExportedDocumentation(t *testing.T) {
	for _, dir := range []string{"ie", "message"} {
		entries, err := os.ReadDir(dir)
		if err != nil {
			t.Fatal(err)
		}
		for _, entry := range entries {
			if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".go") || strings.HasSuffix(entry.Name(), "_test.go") {
				continue
			}
			name := filepath.Join(dir, entry.Name())
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, name, nil, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				if !ok || (fn.Name.Name != "UnmarshalBinary" && !strings.HasPrefix(fn.Name.Name, "Parse")) {
					continue
				}
				if !strings.Contains(fn.Doc.Text(), "borrow") {
					t.Errorf("%s: %s must document borrowed-input ownership", fset.Position(fn.Pos()), fn.Name.Name)
				}
			}
			for _, id := range undocumented(file) {
				t.Errorf("%s: %s needs a doc comment starting with its name", fset.Position(id.Pos()), id.Name)
			}
		}
	}
}

func undocumented(file *ast.File) []*ast.Ident {
	var missing []*ast.Ident
	check := func(id *ast.Ident, doc *ast.CommentGroup) {
		if id.IsExported() && !strings.HasPrefix(doc.Text(), id.Name+" ") && !strings.HasPrefix(doc.Text(), id.Name+"\n") {
			missing = append(missing, id)
		}
	}
	members := func(fields *ast.FieldList) {
		for _, field := range fields.List {
			doc := field.Doc
			if doc == nil {
				doc = field.Comment
			}
			if len(field.Names) == 0 {
				typ := field.Type
				if star, ok := typ.(*ast.StarExpr); ok {
					typ = star.X
				}
				switch typ := typ.(type) {
				case *ast.Ident:
					check(typ, doc)
				case *ast.SelectorExpr:
					check(typ.Sel, doc)
				}
			}
			for _, id := range field.Names {
				check(id, doc)
			}
		}
	}
	for _, decl := range file.Decls {
		switch decl := decl.(type) {
		case *ast.FuncDecl:
			check(decl.Name, decl.Doc)
		case *ast.GenDecl:
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					doc := spec.Doc
					if doc == nil {
						doc = decl.Doc
					}
					check(spec.Name, doc)
					if !spec.Name.IsExported() {
						continue
					}
					switch typ := spec.Type.(type) {
					case *ast.StructType:
						members(typ.Fields)
					case *ast.InterfaceType:
						members(typ.Methods)
					}
				case *ast.ValueSpec:
					doc := spec.Doc
					if doc == nil && len(decl.Specs) == 1 {
						doc = decl.Doc
					}
					for _, id := range spec.Names {
						check(id, doc)
					}
				}
			}
		}
	}
	return missing
}

func TestDocumentationChecker(t *testing.T) {
	const src = `package sample
// WrongName is an incorrect comment.
func Exported() {}
func hidden() {}
// Thing is public.
type Thing struct {
 Missing int
 // Good is documented.
 Good int
 *Embedded
}
// API describes an interface.
type API interface { Method() }
const (
 MissingConstant = 1
 // GoodConstant is documented.
 GoodConstant = 2
)
// Value is documented.
var Value int
`
	file, err := parser.ParseFile(token.NewFileSet(), "sample.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	var names []string
	for _, id := range undocumented(file) {
		names = append(names, id.Name)
	}
	if got := strings.Join(names, ","); got != "Exported,Missing,Embedded,Method,MissingConstant" {
		t.Fatalf("missing documentation: %s", got)
	}
}
