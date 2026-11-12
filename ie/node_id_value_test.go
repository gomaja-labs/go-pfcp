package ie_test

import (
	"bytes"
	"errors"
	"io"
	"net/netip"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestNodeIDTypedValue(t *testing.T) {
	for _, tc := range []struct {
		name    string
		value   *ie.IE
		want    ie.NodeIDValue
		payload []byte
	}{
		{"IPv4", ie.NewNodeIDAddress(netip.MustParseAddr("192.0.2.1")), ie.NodeIDValue{Type: ie.NodeIDIPv4Address, Address: netip.MustParseAddr("192.0.2.1")}, []byte{0, 192, 0, 2, 1}},
		{"IPv6", ie.NewNodeIDAddress(netip.MustParseAddr("2001:db8::1")), ie.NodeIDValue{Type: ie.NodeIDIPv6Address, Address: netip.MustParseAddr("2001:db8::1")}, []byte{1, 0x20, 1, 0x0d, 0xb8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}},
		{"mapped IPv6", ie.NewNodeIDAddress(netip.MustParseAddr("::ffff:192.0.2.1")), ie.NodeIDValue{Type: ie.NodeIDIPv6Address, Address: netip.MustParseAddr("::ffff:192.0.2.1")}, []byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff, 192, 0, 2, 1}},
		{"numeric domain", ie.NewNodeIDFQDN("192.0.2.1"), ie.NodeIDValue{Type: ie.NodeIDFQDN, FQDN: "192.0.2.1"}, []byte{2, 3, '1', '9', '2', 1, '0', 1, '2', 1, '1'}},
		{"domain", ie.NewNodeIDFQDN("node.example"), ie.NodeIDValue{Type: ie.NodeIDFQDN, FQDN: "node.example"}, []byte{2, 4, 'n', 'o', 'd', 'e', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e'}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if !bytes.Equal(tc.value.Payload, tc.payload) {
				t.Fatal("incorrect wire value")
			}
			wire, err := tc.value.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			parsed, err := ie.Parse(wire)
			if err != nil {
				t.Fatal(err)
			}
			got, err := parsed.NodeID()
			if err != nil || got != tc.want {
				t.Fatalf("got %+v, %v; want %+v", got, err, tc.want)
			}
			again, err := parsed.Marshal()
			if err != nil || !bytes.Equal(wire, again) {
				t.Fatal("wire did not round trip")
			}
			parsed.Payload[0] |= 0xf0
			got, err = parsed.NodeID()
			if err != nil || got != tc.want {
				t.Fatal("spare bits changed value")
			}
			// Address alternatives are extendable (TS 29.244 V20.1.0 clauses 8.1.2/8.2.38).
			if tc.want.Type != ie.NodeIDFQDN {
				parsed.Payload = append(parsed.Payload, 0xaa)
				got, err = parsed.NodeID()
				if err != nil || got != tc.want {
					t.Fatal("extension changed address")
				}
			}
			limit := len(tc.payload)
			if tc.want.Type == ie.NodeIDFQDN {
				limit = 3
			}
			for n := 0; n < limit; n++ {
				if _, err := ie.New(ie.NodeID, tc.payload[:n]).NodeID(); err == nil {
					t.Errorf("accepted prefix length %d", n)
				}
			}
		})
	}
	for _, tc := range []struct {
		name  string
		value *ie.IE
		want  error
	}{
		{"nil", nil, ie.ErrIENotFound}, {"wrong type", ie.New(ie.Cause, []byte{0}), nil}, {"null", ie.New(ie.NodeID, nil), io.ErrUnexpectedEOF}, {"unknown kind", ie.New(ie.NodeID, []byte{3, 0}), nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.value.NodeID()
			if err == nil || got != (ie.NodeIDValue{}) {
				t.Fatalf("got %+v, %v", got, err)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("got %v", err)
			}
		})
	}
	for name, value := range map[string]*ie.IE{"invalid address": ie.NewNodeIDAddress(netip.Addr{}), "zone": ie.NewNodeIDAddress(netip.MustParseAddr("fe80::1%eth0")), "empty name": ie.NewNodeIDFQDN(""), "root only": ie.NewNodeIDFQDN("."), "root suffix": ie.NewNodeIDFQDN("node.example."), "overlong label": ie.NewNodeIDFQDN(strings.Repeat("x", 64))} {
		requireInvalidValue(t, name, ie.NodeID, value)
	}
}

func TestNodeIDDomainOmitsRootOctet(t *testing.T) {
	// TS 29.244 V20.1.0 clause 8.2.38 explicitly excludes the trailing zero.
	for _, payload := range [][]byte{{2, 0}, {2, 1, 'a', 0}} {
		if _, err := ie.New(ie.NodeID, payload).NodeID(); err == nil {
			t.Error("accepted DNS root octet")
		}
	}
	// RFC 2181 section 11 permits arbitrary label octets: zero inside a
	// length-delimited label is not a root terminator.
	value := ie.NewNodeIDFQDN("a\x00")
	got, err := value.NodeID()
	if err != nil || got.FQDN != "a\x00" {
		t.Fatalf("binary label: %+v, %v", got, err)
	}
}
