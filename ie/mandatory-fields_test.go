package ie_test

import (
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestExplicitMandatoryFields(t *testing.T) {
	for _, b := range [][]byte{nil, {0}, {0, 2}, {0, 2, 0}, {1, 2, 0, 0}, {2, 2, 0, 0}} {
		_, e := ie.ParseMPQUICAddressInformationFields(b)
		if !errors.Is(e, ie.ErrMalformed) {
			t.Fatalf("MPQUIC missing field len %d: %v", len(b), e)
		}
	}
	for _, b := range [][]byte{nil, {1}, {2}, {4}} {
		_, e := ie.ParseMBSSessionIdentifierFields(b)
		if !errors.Is(e, ie.ErrMalformed) {
			t.Fatalf("MBS missing field: %v", e)
		}
	}
	if _, e := ie.NewMBSSessionIdentifier(ie.MBSSessionIdentifierFields{}).Marshal(); e != nil {
		t.Fatal("inferred MBS field", e)
	}
	if _, e := ie.NewMPQUICAddressInformation(ie.MPQUICAddressInformationFields{ProxyType: ie.MPQUICProxyUDP}).Marshal(); e != nil {
		t.Fatal("inferred address", e)
	}
}
