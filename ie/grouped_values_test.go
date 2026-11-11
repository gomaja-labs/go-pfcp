package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestGroupedNetworkInstanceFQDN(t *testing.T) {
	parent := ie.NewPDI(ie.NewFQDNIE(ie.NetworkInstance, "a.b"))
	for _, i := range []*ie.IE{parent, parent.ChildIEs[0]} {
		v, err := i.NetworkInstanceFQDN()
		if err != nil || v != "a.b" {
			t.Errorf("FQDN value %q error %v", v, err)
		}
		v, err = i.NetworkInstanceHeuristic()
		if err != nil || v != "a.b" {
			t.Errorf("heuristic value %q error %v", v, err)
		}
	}
}
