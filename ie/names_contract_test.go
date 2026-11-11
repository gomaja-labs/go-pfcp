package ie_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestCurrentNames(t *testing.T) {
	for typ, name := range map[ie.IEType]string{
		194: "CreateBridgeRouterInfo", 196: "PortNumber", 198: "FiveGSUserPlaneNodeID",
		206: "TimeDomainNumber", 229: "LinkSpecificMultipathIPAddress",
		265: "MultipathApplicableIndication", 266: "UserPlaneNodeManagementInformationContainer",
	} {
		if got := typ.String(); got != name {
			t.Errorf("type %d name=%s; want %s", typ, got, name)
		}
	}
}

func TestReservedAPIRemoved(t *testing.T) {
	paths, e := filepath.Glob("*.go")
	if e != nil {
		t.Fatal(e)
	}
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, e := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if e != nil {
			t.Fatal(e)
		}
		ast.Inspect(f, func(n ast.Node) bool {
			id, ok := n.(*ast.Ident)
			if ok && strings.Contains(id.Name, "UserPlaneIPResourceInformation") {
				t.Errorf("reserved type 116 API remains in %s", path)
			}
			return true
		})
	}
}
