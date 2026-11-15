package ie_test

import (
	"bytes"
	"errors"
	"net"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.37: each presence bit selects its address family.
func TestFSEIDAddressValidation(t *testing.T) {
	v4, v6 := net.ParseIP("192.0.2.1"), net.ParseIP("2001:db8::1")
	for name, args := range map[string][2]net.IP{
		"IPv6AsIPv4": {v6, nil}, "IPv4AsIPv6": {nil, v4},
		"MalformedIPv4": {net.IP{192, 0, 2}, nil}, "MalformedIPv6": {nil, net.IP{1}},
		"EmptyIPv4": {net.IP{}, nil}, "MissingBoth": {nil, nil},
	} {
		t.Run(name, func(t *testing.T) { assertInvalidConstructor(t, ie.NewFSEID(1, args[0], args[1])) })
	}
	for name, f := range map[string]*ie.FSEIDFields{
		"MissingIPv4": {Flags: 2}, "MissingIPv6": {Flags: 1},
		"UnflaggedIPv4": {Flags: 1, IPv4Address: v4, IPv6Address: v6},
		"UnflaggedIPv6": {Flags: 2, IPv4Address: v4, IPv6Address: v6},
		"SpareFlags":    {Flags: 0x82, IPv4Address: v4},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := f.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("Marshal: %v", err)
			}
			if err := f.MarshalTo(make([]byte, 64)); !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("MarshalTo: %v", err)
			}
		})
	}
	for _, args := range [][2]net.IP{{v4, nil}, {nil, v6}, {v4, v6}} {
		f := ie.NewFSEIDFields(1, args[0], args[1])
		wire, err := f.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if len(wire) != f.MarshalLen() {
			t.Fatal("length disagrees")
		}
		wire[0] |= 0xfc
		parsed, err := ie.ParseFSEIDFields(wire)
		if err != nil {
			t.Fatal(err)
		}
		if parsed.Flags != f.Flags {
			t.Errorf("spare flags evaluated: %d", parsed.Flags)
		}
		wire[0] &= 3
		got, err := parsed.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, wire) {
			t.Error("valid fields changed")
		}
	}
}
