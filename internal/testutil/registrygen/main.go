package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func main() {
	fs := token.NewFileSet()
	dir := "."
	if _, err := os.Stat("ie.go"); err != nil {
		dir = "ie"
	}
	files, _ := filepath.Glob(filepath.Join(dir, "*.go"))
	ctors := map[string][]string{}
	var parsers []string
	for _, p := range files {
		if strings.HasSuffix(p, "_test.go") {
			continue
		}
		f, e := parser.ParseFile(fs, p, nil, 0)
		if e != nil {
			panic(e)
		}
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}
			name := fn.Name.Name
			if strings.HasPrefix(name, "Parse") && strings.HasSuffix(name, "Fields") {
				parsers = append(parsers, name)
			}
			if strings.HasPrefix(name, "New") {
				names := make([]string, 0, len(ctors))
				for _, arg := range fn.Type.Params.List {
					for _, n := range arg.Names {
						names = append(names, n.Name)
					}
				}
				ctors[name] = names
			}
		}
	}
	var out bytes.Buffer
	out.WriteString("//go:generate go run ../internal/testutil/registrygen\n\npackage ie_test\n\nimport \"github.com/gomaja-labs/go-pfcp/ie\"\n\nvar fieldParsers=map[string]any{\n")
	sort.Strings(parsers)
	for _, name := range parsers {
		fmt.Fprintf(&out, "%q:ie.%s,\n", name, name)
	}
	out.WriteString("}\n\ntype constructorEntry struct{fn any;params []string}\nvar constructors=map[string]constructorEntry{\n")
	names := make([]string, 0, len(ctors))
	for n := range ctors {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Fprintf(&out, "%q:{ie.%s,[]string{", n, n)
		for _, param := range ctors[n] {
			fmt.Fprintf(&out, "%q,", param)
		}
		out.WriteString("}},\n")
	}
	out.WriteString("}\n")
	b, e := format.Source(out.Bytes())
	if e != nil {
		panic(e)
	}
	if e = os.WriteFile(filepath.Join(dir, "fields_registry_test.go"), b, 0600); e != nil {
		panic(e)
	}
	fmt.Println("constructors", len(ctors), "fields parsers", len(parsers))
}
