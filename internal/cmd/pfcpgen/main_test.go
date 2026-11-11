package main

import (
	"bytes"
	"go/parser"
	"go/token"
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
