package ie_test

import (
	"errors"
	"io"
	"net/netip"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestExplicitVendorGroupOperations(t *testing.T) {
	i := ie.NewVendorSpecificGroupedIE(0x8001, 1234, ie.NewCause(1))
	if i.IsGrouped() {
		t.Fatal("vendor type classified as automatically grouped")
	}
	if _, err := i.FindByType(ie.Cause); err != nil {
		t.Fatal(err)
	}
	if err := i.Add(ie.NewQFI(7)); err != nil {
		t.Fatal(err)
	}
	if err := i.Remove(ie.Cause); err != nil {
		t.Fatal(err)
	}
	if _, err := i.FindByType(ie.Cause); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal("removed child still present", err)
	}
	if child, err := i.FindByType(ie.QFI); err != nil || child == nil {
		t.Fatal("added child absent", err)
	}
}

func TestRedirectAddressLengthOverflow(t *testing.T) {
	for _, addrs := range [][]string{{strings.Repeat("x", 65536)}, {"192.0.2.1", strings.Repeat("x", 65536)}} {
		typ := ie.RedirectAddrURL
		if len(addrs) == 2 {
			typ = ie.RedirectAddrIPv4AndIPv6
		}
		if i := ie.NewRedirectInformation(typ, addrs...); i != nil {
			t.Error("constructor accepted overflowing address length")
		}
		f := ie.NewRedirectInformationFields(typ, addrs...)
		if _, err := f.Marshal(); !errors.Is(err, ie.ErrInvalidLength) {
			t.Error("fields accepted overflowing address length", err)
		}
	}
}

func TestForwardingPolicyLengthBoundary(t *testing.T) {
	for _, n := range []int{0, 1, 255, 256, 65536} {
		input := strings.Repeat("x", n)
		i := ie.NewForwardingPolicy(input)
		if n > 255 {
			if i != nil {
				t.Errorf("accepted identifier length %d", n)
			}
			continue
		}
		b, err := i.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := ie.Parse(b)
		if err != nil {
			t.Fatal(err)
		}
		got, err := parsed.ForwardingPolicyIdentifier()
		if err != nil || got != input {
			t.Errorf("identifier length %d changed, error %v", n, err)
		}
	}
}

func TestNodeIDSpareBits(t *testing.T) {
	for _, i := range []*ie.IE{ie.NewNodeIDAddress(netip.MustParseAddr("192.0.2.1")), ie.NewNodeIDAddress(netip.MustParseAddr("2001:db8::1")), ie.NewNodeIDFQDN("node.example")} {
		want, err := i.NodeID()
		if err != nil {
			t.Fatal(err)
		}
		i.Payload[0] |= 0xf0
		if got, err := i.NodeID(); err != nil || got != want {
			t.Error("spare bits changed Node ID", err)
		}
	}
}

func TestNullStringValue(t *testing.T) {
	for _, typ := range []ie.IEType{ie.ApplicationID, ie.DataNetworkAccessIdentifier, ie.ActivatePredefinedRules} {
		b, err := ie.New(typ, nil).Marshal()
		if err != nil {
			t.Fatal(err)
		}
		i, err := ie.Parse(b)
		if err != nil {
			t.Fatal("null IE must parse", err)
		}
		if got, err := i.ValueAsString(); got != "" || !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Error("null IE exposed as a valid string", err)
		}
	}
}
