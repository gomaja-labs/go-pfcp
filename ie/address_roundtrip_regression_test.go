package ie_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"reflect"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.20; RFC 8506 section 8.38;
// RFC 5952 section 5 explicitly permits IPv4-mapped IPv6 mixed notation.
func TestRedirectMappedIPv6(t *testing.T) {
	const addr = "::ffff:192.0.2.1"
	item := ie.NewRedirectInformation(ie.RedirectAddrIPv6, addr)
	wire, err := item.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	want := append([]byte{0, 38, 0, byte(3 + len(addr)), 1, 0, byte(len(addr))}, addr...)
	if !bytes.Equal(wire, want) {
		t.Fatal("mapped text encoding changed")
	}
	fields, err := ie.ParseRedirectInformationFields(want[4:])
	if err != nil || fields.RedirectServerAddress != addr {
		t.Fatalf("mapped text receive: %v", err)
	}
	for _, item := range []*ie.IE{
		ie.NewRedirectInformation(ie.RedirectAddrIPv4, addr),
		ie.NewRedirectInformation(ie.RedirectAddrIPv6, "192.0.2.1"),
	} {
		assertInvalidConstructor(t, item)
	}
}

// Each receive payload is independent of the encoder under test. Binary IPv6
// fields reject mapped addresses on send and preserve them on receive.
// The case names identify the applicable TS 29.244 V20.1.0 clauses.
func TestMappedIPv6AddressRoundTrip(t *testing.T) {
	mapped := netip.MustParseAddr("::ffff:192.0.2.1")
	absent := netip.Addr{}
	addr := mapped.AsSlice()
	payload := func(prefix ...byte) []byte { return append(prefix, addr...) }
	ssm := append([]byte{0x50}, netip.MustParseAddr("ff3e::1").AsSlice()...)
	ssm = append(ssm, payload(0x50)...)
	cases := []struct {
		name   string
		typ    ie.IEType
		fields addressFields
		item   *ie.IE
		value  []byte
	}{
		{"8.2.199/CPIP", ie.CPIPAddress, &ie.CPIPAddressFields{Flags: 1, IPv6Address: mapped}, ie.NewCPIPAddress(ie.CPIPAddressFields{Flags: 1, IPv6Address: mapped}), payload(1)},
		{"8.2.209/Local", ie.LocalIngressTunnel, &ie.LocalIngressTunnelFields{Flags: 2, Address: mapped}, ie.NewLocalIngressTunnel(ie.LocalIngressTunnelFields{Flags: 2, Address: mapped}), payload(2, 0, 0)},
		{"8.2.226/MPQUIC", ie.MPQUICAddressInformation, &ie.MPQUICAddressInformationFields{Flags: 2, ProxyType: 2, Port: 80, IPv6Address: mapped}, ie.NewMPQUICAddressInformation(ie.MPQUICAddressInformationFields{Flags: 2, ProxyType: 2, Port: 80, IPv6Address: mapped}), payload(2, 2, 0, 80)},
		{"8.2.206/MBS", ie.MBSSessionIdentifier, &ie.MBSSessionIdentifierFields{Flags: ie.MBSSessionSSMI, MulticastAddress: netip.MustParseAddr("ff3e::1"), SourceAddress: mapped}, ie.NewMBSSessionIdentifier(ie.MBSSessionIdentifierFields{Flags: ie.MBSSessionSSMI, MulticastAddress: netip.MustParseAddr("ff3e::1"), SourceAddress: mapped}), append([]byte{2}, ssm...)},
		{"8.2.207/Transport", ie.MulticastTransportInformation, &ie.MulticastTransportInformationFields{CommonTEID: 1, MulticastAddress: netip.MustParseAddr("ff3e::1"), SourceAddress: mapped}, ie.NewMulticastTransportInformation(ie.MulticastTransportInformationFields{CommonTEID: 1, MulticastAddress: netip.MustParseAddr("ff3e::1"), SourceAddress: mapped}), append([]byte{0, 0, 0, 0, 1}, ssm...)},
		{"8.2.248/Endpoint", ie.MeasurementEndpointAddress, ie.NewMeasurementEndpointAddressFields(2, absent, mapped), ie.NewMeasurementEndpointAddress(2, absent, mapped), payload(2)},
		{"8.2.243/N6RoutingSource", ie.N6RoutingInformation, ie.NewN6RoutingInformationFields(2, absent, mapped, 0, absent, absent, 0), ie.NewN6RoutingInformation(2, absent, mapped, 0, absent, absent, 0), payload(2)},
		{"8.2.243/N6RoutingDestination", ie.N6RoutingInformation, ie.NewN6RoutingInformationFields(16, absent, absent, 0, absent, mapped, 0), ie.NewN6RoutingInformation(16, absent, absent, 0, absent, mapped, 0), payload(16)},
		{"8.2.270/MoQ", ie.MoQRelayIPAddress, ie.NewMoQRelayIPAddressFields(1, absent, mapped), ie.NewMoQRelayIPAddress(1, absent, mapped), payload(1)},
		{"8.2.261/OnPath", ie.OnPathN6ConnectionInformation, ie.NewOnPathN6ConnectionInformationFields(16, absent, mapped, 0), ie.NewOnPathN6ConnectionInformation(16, absent, mapped, 0), payload(16)},
		{"8.2.200/ReplacementDestination", ie.IPAddressAndPortNumberReplacement, &ie.IPAddressAndPortNumberReplacementFields{Flags: 2, DestinationIPv6: mapped}, ie.NewIPAddressAndPortNumberReplacement(ie.IPAddressAndPortNumberReplacementFields{Flags: 2, DestinationIPv6: mapped}), payload(2)},
		{"8.2.200/ReplacementSource", ie.IPAddressAndPortNumberReplacement, &ie.IPAddressAndPortNumberReplacementFields{Flags: 16, SourceIPv6: mapped}, ie.NewIPAddressAndPortNumberReplacement(ie.IPAddressAndPortNumberReplacementFields{Flags: 16, SourceIPv6: mapped}), payload(16)},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("send", func(t *testing.T) {
				if _, err := tc.fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
					t.Errorf("Marshal mapped address: %v", err)
				}
				if err := tc.fields.MarshalTo(make([]byte, len(tc.value))); !errors.Is(err, ie.ErrInvalidValue) {
					t.Errorf("MarshalTo mapped address: %v", err)
				}
				assertInvalidConstructor(t, tc.item)
			})
			t.Run("receive", func(t *testing.T) {
				got := reflect.New(reflect.TypeOf(tc.fields).Elem()).Interface().(addressFields)
				if err := got.UnmarshalBinary(tc.value); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, tc.fields) {
					t.Fatal("mapped address fields changed")
				}
				// Forward the original IE; typed encoders apply send validation.
				item := ie.New(tc.typ, tc.value)
				wire, err := item.Marshal()
				if err != nil {
					t.Fatal(err)
				}
				parsed, err := ie.Parse(wire)
				if err != nil {
					t.Fatal(err)
				}
				method := strings.TrimSuffix(reflect.TypeOf(tc.fields).Elem().Name(), "Fields")
				values := reflect.ValueOf(parsed).MethodByName(method).Call(nil)
				if !values[1].IsNil() || !reflect.DeepEqual(values[0].Interface(), tc.fields) {
					t.Fatalf("typed receive %s changed mapped fields: %v", method, values[1].Interface())
				}
				wire, err = parsed.Marshal()
				if err != nil || !bytes.Equal(wire[4:], tc.value) {
					t.Fatalf("mapped receive re-encode: %v", err)
				}
			})
		})
	}
}

// TS 29.244 V20.1.0 clause 8.2.62: IP6PL requires V6 or CHV6
// and excludes IPv6D; IPv6D also requires V6 or CHV6. Receive stays tolerant.
func TestUEIPAddressPrefixConditions(t *testing.T) {
	for _, family := range []uint8{0, 1, 0x20} {
		for _, prefix := range []uint8{0, 8, 0x40, 0x48} {
			flags := family | prefix
			t.Run(fmt.Sprintf("flags_%02x", flags), func(t *testing.T) {
				address := ""
				raw := []byte{flags}
				if family == 1 {
					address = "2001:db8::1"
					raw = append(raw, netip.MustParseAddr(address).AsSlice()...)
				}
				var delegation, length uint8
				if prefix&8 != 0 {
					delegation = 4
					raw = append(raw, delegation)
				}
				if prefix&0x40 != 0 {
					length = 72
					raw = append(raw, length)
				}
				f := ie.NewUEIPAddressFields(flags, "", address, delegation, length)
				item := ie.NewUEIPAddress(flags, "", address, delegation, length)
				invalid := prefix == 0x48 || family == 0 && prefix != 0
				if invalid {
					assertInvalidConstructor(t, item)
					if _, err := f.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
						t.Errorf("Marshal: %v", err)
					}
					if err := f.MarshalTo(make([]byte, 64)); !errors.Is(err, ie.ErrInvalidValue) {
						t.Errorf("MarshalTo: %v", err)
					}
				} else {
					wire, err := item.Marshal()
					if err != nil || !bytes.Equal(wire[4:], raw) {
						t.Fatalf("valid prefix encoding: %v", err)
					}
				}
				parsed, err := ie.ParseUEIPAddressFields(raw)
				if err != nil || parsed.Flags != flags {
					t.Fatalf("receive tolerance: %v", err)
				}
			})
		}
	}
}

func TestOuterHeaderCreationRequiresDescription(t *testing.T) {
	// TS 29.244 V20.1.0 clause 8.2.56 requires a nonzero description.
	assertInvalidConstructor(t, ie.NewOuterHeaderCreation(0, 0, "", "", 0, 0, 0))
	fields := ie.NewOuterHeaderCreationFields(0, 0, "", "", 0, 0, 0)
	if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("Marshal: %v", err)
	}
	if err := fields.MarshalTo(make([]byte, 2)); !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("MarshalTo: %v", err)
	}
	// Clause 8.1.2 still permits a null IE, and receive does not enforce send rules.
	wire, err := ie.New(ie.OuterHeaderCreation, nil).Marshal()
	if err != nil || len(wire) != 4 {
		t.Fatalf("null IE: %v", err)
	}
	if _, err := ie.ParseOuterHeaderCreationFields([]byte{0, 0}); err != nil {
		t.Fatalf("receive tolerance: %v", err)
	}
}

func TestSourceIPAddressFlagsOnly(t *testing.T) {
	// TS 29.244 V20.1.0 clause 8.2.138: all fields after flags are optional.
	item := ie.NewSourceIPAddress(nil, nil, 0)
	wire, err := item.Marshal()
	if err != nil || !bytes.Equal(wire, []byte{0, 192, 0, 1, 0}) {
		t.Fatalf("flags-only encoding: %v", err)
	}
	fields, err := item.SourceIPAddress()
	if err != nil {
		t.Fatal(err)
	}
	if fields.Flags != 0 || fields.IPv4Address != nil || fields.IPv6Address != nil || fields.MaskPrefixLength != 0 {
		t.Fatal("flags-only fields changed")
	}
	value, err := fields.Marshal()
	if err != nil || !bytes.Equal(value, []byte{0}) {
		t.Fatalf("re-encode: %v", err)
	}
	for _, raw := range [][]byte{nil, {1}, {2}, {4}} {
		if _, err := ie.ParseSourceIPAddressFields(raw); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Errorf("missing selected fields: %v", err)
		}
	}
}

func TestMappedIPv6ExplicitConstructors(t *testing.T) {
	// TS 29.244 V20.1.0 clauses 8.2.38 and 8.2.188. Reject mapped
	// values on send so incidental IPv4 representations cannot change identity.
	mapped := netip.MustParseAddr("::ffff:192.0.2.1")
	t.Run("NodeID", func(t *testing.T) {
		assertInvalidConstructor(t, ie.NewNodeIDAddress(mapped))
		item := ie.New(ie.NodeID, append([]byte{ie.NodeIDIPv6Address}, mapped.AsSlice()...))
		value, err := item.NodeID()
		if err != nil || value.Address != mapped || value.Type != ie.NodeIDIPv6Address {
			t.Fatalf("Node ID receive: %v", err)
		}
	})
	t.Run("LNS", func(t *testing.T) {
		assertInvalidConstructor(t, ie.NewLNSAddress(mapped))
		item := ie.New(ie.LNSAddress, mapped.AsSlice())
		value, err := item.LNSAddress()
		if err != nil || value != mapped {
			t.Fatalf("LNS receive: %v", err)
		}
	})
}

func TestMappedNetIPAddressReceive(t *testing.T) {
	// TS 29.244 V20.1.0 clauses 8.2.3, 8.2.37, 8.2.56, 8.2.62 and
	// 8.2.138: receive preserves the sixteen IPv6 address octets even though
	// net.IP cannot retain their explicit family for subsequent typed encoding.
	mapped := netip.MustParseAddr("::ffff:192.0.2.1").AsSlice()
	cases := []struct {
		name   string
		fields addressFields
		prefix []byte
	}{
		{"FSEID", &ie.FSEIDFields{}, []byte{1, 0, 0, 0, 0, 0, 0, 0, 1}},
		{"FTEID", &ie.FTEIDFields{}, []byte{2, 0, 0, 0, 1}},
		{"Source", &ie.SourceIPAddressFields{}, []byte{1}},
		{"UE", &ie.UEIPAddressFields{}, []byte{1}},
		{"Outer", &ie.OuterHeaderCreationFields{}, []byte{2, 0, 0, 0, 0, 1}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			value := append(tc.prefix, mapped...)
			if err := tc.fields.UnmarshalBinary(value); err != nil {
				t.Fatal(err)
			}
			got := reflect.ValueOf(tc.fields).Elem().FieldByName("IPv6Address").Bytes()
			if !bytes.Equal(got, mapped) {
				t.Fatal("received IPv6 octets changed")
			}
			if _, err := tc.fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
				t.Fatalf("ambiguous net.IP send: %v", err)
			}
		})
	}
}
