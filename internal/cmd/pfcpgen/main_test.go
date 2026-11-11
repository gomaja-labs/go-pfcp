package main

import (
	"bytes"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/internal/spec"
)

func TestGenerate(t *testing.T) {
	c, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	a, err := generate(c)
	if err != nil {
		t.Fatal(err)
	}
	b, err := generate(c)
	if err != nil {
		t.Fatal(err)
	}
	if len(a) != 30 {
		t.Fatalf("got %d generated files", len(a))
	}
	for name, src := range a {
		if !bytes.Equal(src, b[name]) {
			t.Fatal("nondeterministic", name)
		}
		if !bytes.HasPrefix(src, []byte(banner)) {
			t.Fatal("missing provenance", name)
		}
		if _, err := parser.ParseFile(token.NewFileSet(), name, src, parser.AllErrors); err != nil {
			t.Fatal(err)
		}
	}
}

// TS 29.244 V20.1.0 Table 8.1.2-1 allocates no IE semantics to these kinds.
func TestGenerateAssignedTypes(t *testing.T) {
	c, err := spec.Load()
	if err != nil {
		t.Fatal(err)
	}
	files, err := generate(c)
	if err != nil {
		t.Fatal(err)
	}
	src := string(files["ie/types_gen.go"])
	for _, e := range c.IETypes {
		if e.Kind == "reserved" || e.Kind == "spare" || e.Kind == "can-be-used" {
			if strings.Contains(src, e.GoName+" IEType =") {
				t.Errorf("unassigned constant %s", e.GoName)
			}
			if e.ExistingGoName != nil && strings.Contains(src, *e.ExistingGoName+" IEType =") {
				t.Errorf("historical unassigned constant for type %d", e.Type)
			}
		}
	}
	// The current catalogue name takes precedence over the historical surface.
	for j := range c.IETypes {
		if c.IETypes[j].Type == 19 {
			c.IETypes[j].GoName = "CurrentName"
		}
	}
	files, err = generate(c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(files["ie/types_gen.go"]), "CurrentName IEType = 19") {
		t.Fatal("historical name overrides current name")
	}
}
