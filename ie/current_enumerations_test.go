package ie_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"
)

func TestCurrentEnumerations(t *testing.T) {
	for _, tc := range []struct {
		file, name string
		value      uint64
	}{
		{"cause.go", "CauseMoreUsageReportsToSend", 2},
		{"cause.go", "CauseRequestPartiallyAccepted", 3},
		{"cause.go", "CauseAllDynamicAddressesOccupied", 79},
		{"cause.go", "CauseUnknownPredefinedRule", 80},
		{"cause.go", "CauseUnknownApplicationID", 81},
		{"cause.go", "CauseL2TPTunnelEstablishmentFailure", 82},
		{"cause.go", "CauseL2TPSessionEstablishmentFailure", 83},
		{"cause.go", "CauseL2TPTunnelRelease", 84},
		{"cause.go", "CauseL2TPSessionRelease", 85},
		{"cause.go", "CauseSessionRestorationResourceUnavailable", 86},
		{"cause.go", "CauseL2TPTunnelAuthenticationFailure", 87},
		{"cause.go", "CauseL2TPSessionAuthenticationFailure", 88},
		{"cause.go", "CauseL2TPLNSUnreachable", 89},
		{"cause.go", "CausePFDContentsSyntaxError", 90},
		{"cause.go", "CausePFDContentsSemanticsError", 91},
		{"cause.go", "CausePFDApplicationIDUnknown", 92},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeN9Roaming", 21},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeIuU", 22},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeN9DataForwarding", 23},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeSxaU", 24},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeSxbU", 25},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeSxcU", 26},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeN4U", 27},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeSGWUPFUL", 28},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeN6mbNmb9", 29},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeN3mb", 30},
		{"tgpp-interface-type.go", "TGPPInterfaceTypeN19mb", 31},
		{"steering-functionality.go", "SteeringFunctionalityMPQUICUDP", 2},
		{"steering-functionality.go", "SteeringFunctionalityMPQUICIP", 3},
		{"steering-functionality.go", "SteeringFunctionalityMPQUICE", 4},
		{"steering-mode.go", "SteeringModeRedundant", 4},
		{"priority.go", "PriorityPrimary", 5},
		{"priority.go", "PrioritySecondary", 6},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f, e := parser.ParseFile(token.NewFileSet(), tc.file, nil, 0)
			if e != nil {
				t.Fatal(e)
			}
			found := false
			ast.Inspect(f, func(n ast.Node) bool {
				v, ok := n.(*ast.ValueSpec)
				if !ok {
					return true
				}
				for j, name := range v.Names {
					if name.Name != tc.name {
						continue
					}
					found = true
					if j >= len(v.Values) {
						t.Fatal("constant needs an explicit value")
					}
					lit, ok := v.Values[j].(*ast.BasicLit)
					if !ok {
						t.Fatal("constant value is not a literal")
					}
					value, e := strconv.ParseUint(lit.Value, 0, 64)
					if e != nil || value != tc.value {
						t.Fatalf("constant = %s; want %d", lit.Value, tc.value)
					}
				}
				return true
			})
			if !found {
				t.Fatal("missing current enumeration constant")
			}
		})
	}
}
