package ie_test

import (
	"errors"
	"net"
	"net/netip"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

type addressFields interface {
	Marshal() ([]byte, error)
	MarshalTo([]byte) error
	UnmarshalBinary([]byte) error
}
type addressSlot struct {
	name string
	flag uint64
	v6   bool
}

// Each selected address must match its family; absent fields cannot discard arguments.
// TS 29.244 V20.1.0 clauses 8.1.2 and the address clauses named in each case.
func TestAddressFieldFamilies(t *testing.T) {
	v4, v6 := net.ParseIP("192.0.2.1"), net.ParseIP("2001:db8::1")
	a4, a6 := netip.MustParseAddr("192.0.2.1"), netip.MustParseAddr("2001:db8::1")
	cases := []struct {
		name   string
		fields addressFields
		slots  []addressSlot
	}{
		{"8.2.37/FSEID", ie.NewFSEIDFields(1, v4, v6), []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.3/FTEID", ie.NewFTEIDFields(3, 1, v4, v6, 0), []addressSlot{{"IPv4Address", 1, false}, {"IPv6Address", 2, true}}},
		{"8.2.129/AlternativeSMF", ie.NewAlternativeSMFIPAddressFields(v4, v6), []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.135/CPPFCPEntity", ie.NewCPPFCPEntityIPAddressFields(v4, v6), []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.158/LinkSpecific", ie.NewLinkSpecificMultipathIPAddressFields(v4, v6, v4, v6), []addressSlot{{"UELinkSpecificIPv4AddressFor3GPPAccess", 1, false}, {"UELinkSpecificIPv6AddressFor3GPPAccess", 2, true}, {"UELinkSpecificIPv4AddressForNon3GPPAccess", 4, false}, {"UELinkSpecificIPv6AddressForNon3GPPAccess", 8, true}}},
		{"8.2.157/MPTCP", ie.NewMPTCPAddressInformationFields(1, 80, v4, v6), []addressSlot{{"MPTCPIPv4Address", 1, false}, {"MPTCPIPv6Address", 2, true}}},
		{"8.2.159/PMF", ie.NewPMFAddressInformationFields(v4, v6, 1, 2, nil, nil), []addressSlot{{"PMFIPv4Address", 1, false}, {"PMFIPv6Address", 2, true}}},
		{"8.2.138/Source", ie.NewSourceIPAddressFields(v4, v6, 0), []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.62/UE", ie.NewUEIPAddressFields(3, v4.String(), v6.String(), 0, 0), []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.56/OuterHeader", ie.NewOuterHeaderCreationFields(0x0300, 1, v4.String(), v6.String(), 0, 0, 0), []addressSlot{{"IPv4Address", 0x0100, false}, {"IPv6Address", 0x0200, true}}},
		{"8.2.70/RemoteGTPU", ie.NewRemoteGTPUPeerFields(3, v4.String(), v6.String(), 0, "", time.Time{}), []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.199/CPIP", &ie.CPIPAddressFields{Flags: 3, IPv4Address: a4, IPv6Address: a6}, []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.209/LocalIngress4", &ie.LocalIngressTunnelFields{Flags: 1, Address: a4}, []addressSlot{{"Address", 1, false}}},
		{"8.2.209/LocalIngress6", &ie.LocalIngressTunnelFields{Flags: 2, Address: a6}, []addressSlot{{"Address", 2, true}}},
		{"8.2.200/Replacement", &ie.IPAddressAndPortNumberReplacementFields{Flags: 27, DestinationIPv4: a4, DestinationIPv6: a6, SourceIPv4: a4, SourceIPv6: a6}, []addressSlot{{"DestinationIPv4", 1, false}, {"DestinationIPv6", 2, true}, {"SourceIPv4", 8, false}, {"SourceIPv6", 16, true}}},
		{"8.2.226/MPQUIC", &ie.MPQUICAddressInformationFields{Flags: 3, ProxyType: 2, IPv4Address: a4, IPv6Address: a6}, []addressSlot{{"IPv4Address", 1, false}, {"IPv6Address", 2, true}}},
		{"8.2.248/MeasurementEndpoint", ie.NewMeasurementEndpointAddressFields(3, a4, a6), []addressSlot{{"IPv4Address", 1, false}, {"IPv6Address", 2, true}}},
		{"8.2.243/N6Routing", ie.NewN6RoutingInformationFields(27, a4, a6, 0, a4, a6, 0), []addressSlot{{"SourceIPv4Address", 1, false}, {"SourceIPv6Address", 2, true}, {"DestinationIPv4Address", 8, false}, {"DestinationIPv6Address", 16, true}}},
		{"8.2.242/MappedN6", ie.NewMappedN6IPAddressFields(2, a4), []addressSlot{{"IPv4Address", 2, false}}},
		{"8.2.270/MoQRelay", ie.NewMoQRelayIPAddressFields(3, a4, a6), []addressSlot{{"IPv4Address", 2, false}, {"IPv6Address", 1, true}}},
		{"8.2.261/OnPath", ie.NewOnPathN6ConnectionInformationFields(24, a4, a6, 0), []addressSlot{{"IPv4Address", 8, false}, {"IPv6Address", 16, true}}},
		{"8.2.137/Multicast4", ie.NewIPMulticastAddressFields(net.ParseIP("239.1.1.1"), nil, net.ParseIP("239.1.1.2"), nil), []addressSlot{{"StartIPv4Address", 2, false}, {"EndIPv4Address", 2, false}}},
		{"8.2.137/Multicast6", ie.NewIPMulticastAddressFields(nil, net.ParseIP("ff3e::1"), nil, net.ParseIP("ff3e::2")), []addressSlot{{"StartIPv6Address", 1, true}, {"EndIPv6Address", 1, true}}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.fields.Marshal(); err != nil {
				t.Fatalf("valid fields: %v", err)
			}
			original := reflect.ValueOf(tc.fields).Elem()
			for _, slot := range tc.slots {
				t.Run(slot.name, func(t *testing.T) {
					for _, kind := range []string{"missing", "wrongFamily", "mapped", "unflagged", "malformed", "empty", "zone"} {
						if kind == "mapped" && !slot.v6 {
							continue
						}
						clone := reflect.New(original.Type())
						clone.Elem().Set(original)
						field := clone.Elem().FieldByName(slot.name)
						ipType := field.Type() == reflect.TypeFor[net.IP]()
						var value any
						switch kind {
						case "missing":
							value = netip.Addr{}
							if ipType {
								value = net.IP(nil)
							}
						case "wrongFamily":
							if !slot.v6 {
								value = a6
							} else {
								value = a4
							}
							if ipType {
								value = v6
								if slot.v6 {
									value = v4
								}
							}
						case "mapped":
							value = netip.MustParseAddr("::ffff:192.0.2.1")
							if ipType {
								value = v4.To16()
							}
						case "zone":
							if ipType || !slot.v6 {
								continue
							}
							value = a6.WithZone("eth0")
						case "malformed":
							if !ipType {
								continue
							}
							value = net.IP{1, 2, 3}
						case "empty":
							if !ipType {
								continue
							}
							value = net.IP{}
						case "unflagged":
							flags := clone.Elem().FieldByName("Flags")
							if !flags.IsValid() {
								flags = clone.Elem().FieldByName("OuterHeaderCreationDescription")
							}
							flags.SetUint(flags.Uint() &^ slot.flag)
							value = field.Interface()
						}
						field.Set(reflect.ValueOf(value))
						f := clone.Interface().(addressFields)
						if _, err := f.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
							t.Errorf("%s Marshal: %v", kind, err)
						}
						if err := f.MarshalTo(make([]byte, 256)); !errors.Is(err, ie.ErrInvalidValue) {
							t.Errorf("%s MarshalTo: %v", kind, err)
						}
					}
				})
			}
		})
	}
}

func TestAddressConstructorsRetainInvalidArguments(t *testing.T) {
	v4, v6 := net.ParseIP("192.0.2.1"), net.ParseIP("2001:db8::1")
	for name, item := range map[string]*ie.IE{
		"Alternative": ie.NewAlternativeSMFIPAddress(v6, nil), "CPPFCPEntity": ie.NewCPPFCPEntityIPAddress(nil, v4),
		"LinkSpecific": ie.NewLinkSpecificMultipathIPAddress(nil, nil, v6, nil), "MPTCP": ie.NewMPTCPAddressInformation(1, 80, nil, v4),
		"PMF": ie.NewPMFAddressInformation(nil, v4, 0, 0, nil, nil), "Source": ie.NewSourceIPAddress(v6, nil, 0),
		"MulticastEndOnly": ie.NewIPMulticastAddress(nil, nil, v4, nil), "MulticastWrongFamily": ie.NewIPMulticastAddress(v6, nil, nil, nil),
		"FTEIDUnflagged": ie.NewFTEID(1, 1, v4, v6, 0), "FTEIDChoose": ie.NewFTEID(5, 0, v4, nil, 0),
		"UEUnflagged": ie.NewUEIPAddress(2, v4.String(), v6.String(), 0, 0), "UEBadUnflagged": ie.NewUEIPAddress(2, v4.String(), "bad", 0, 0),
		"UEMapped": ie.NewUEIPAddress(1, "", v4.String(), 0, 0), "UEChoose": ie.NewUEIPAddress(0x10, v4.String(), "", 0, 0),
		"OuterUnflagged": ie.NewOuterHeaderCreation(0x0100, 1, v4.String(), v6.String(), 0, 0, 0), "OuterBadUnflagged": ie.NewOuterHeaderCreation(0x0100, 1, v4.String(), "bad", 0, 0, 0),
		"OuterMapped":  ie.NewOuterHeaderCreation(0x0200, 1, "", v4.String(), 0, 0, 0),
		"RemoteMapped": ie.NewRemoteGTPUPeer(1, "", v4.String(), 0, "", time.Time{}),
		"TraceMissing": ie.NewTraceInformation("001", "01", "abc", nil, 0, nil, nil),
		"TraceShort":   ie.NewTraceInformation("001", "01", "abc", nil, 0, nil, net.IP{1, 2, 3}),
	} {
		t.Run(name, func(t *testing.T) { assertInvalidConstructor(t, item) })
	}
}

var _ = registerConstructorFixture(ie.NewUEIPAddress, uint8(3), "192.0.2.1", "2001:db8::1", uint8(0), uint8(0))
var _ = registerConstructorFixture(ie.NewUEIPAddressFields, uint8(3), "192.0.2.1", "2001:db8::1", uint8(0), uint8(0))
var _ = registerConstructorFixture(ie.NewIPMulticastAddress, net.ParseIP("239.1.1.1"), net.IP(nil), net.ParseIP("239.1.1.2"), net.IP(nil))
var _ = registerConstructorFixture(ie.NewIPMulticastAddressFields, net.ParseIP("239.1.1.1"), net.IP(nil), net.ParseIP("239.1.1.2"), net.IP(nil))

func TestAddressSpecialCases(t *testing.T) {
	for name, item := range map[string]*ie.IE{
		"PMFUnflaggedPorts":     ie.NewPMFAddressInformation(nil, nil, 1, 2, nil, nil),
		"MulticastBothFamilies": ie.NewIPMulticastAddress(net.ParseIP("239.1.1.1"), net.ParseIP("ff3e::1"), nil, nil),
		"UEWrongFamily4":        ie.NewUEIPAddress(2, "2001:db8::1", "", 0, 0),
		"RemoteWrongFamily4":    ie.NewRemoteGTPUPeer(2, "2001:db8::1", "", 0, "", time.Time{}),
		"OuterWrongFamily4":     ie.NewOuterHeaderCreation(0x0100, 1, "2001:db8::1", "", 0, 0, 0),
	} {
		t.Run(name, func(t *testing.T) { assertInvalidConstructor(t, item) })
	}
	// Inferred FQ-CSID addresses choose an IPv4 encoding for an IPv4 input.
	// TS 29.244 V20.1.0 clause 8.2.43.
	for _, addr := range []string{"192.0.2.1", "2001:db8::1"} {
		f := ie.NewFQCSID(addr, 1, 2)
		if _, err := f.Marshal(); err != nil {
			t.Fatal(err)
		}
		got, err := f.NodeAddress()
		if err != nil || !net.IP(got).Equal(net.ParseIP(addr)) {
			t.Fatalf("FQ-CSID address: %v", err)
		}
	}
}

func TestAddressSpareBits(t *testing.T) {
	v4 := net.ParseIP("192.0.2.1")
	cases := []struct {
		name  string
		f     addressFields
		spare byte
	}{
		{"AlternativeSMF", ie.NewAlternativeSMFIPAddressFields(v4, nil), 0xf8},
		{"CPPFCPEntity", ie.NewCPPFCPEntityIPAddressFields(v4, nil), 0xfc},
		{"MPTCP", ie.NewMPTCPAddressInformationFields(1, 80, v4, nil), 0xfc},
		{"PMF", ie.NewPMFAddressInformationFields(v4, nil, 1, 2, nil, nil), 0xf8},
		{"Source", ie.NewSourceIPAddressFields(v4, nil, 0), 0xf8},
		{"LinkSpecific", ie.NewLinkSpecificMultipathIPAddressFields(v4, nil, nil, nil), 0xf0},
		{"Multicast", ie.NewIPMulticastAddressFields(net.ParseIP("239.1.1.1"), nil, nil, nil), 0xf0},
		{"UE", ie.NewUEIPAddressFields(2, "192.0.2.1", "", 0, 0), 0x80},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.f.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			want := b[0]
			b[0] |= tc.spare
			if err := tc.f.UnmarshalBinary(b); err != nil {
				t.Fatal(err)
			}
			if got := reflect.ValueOf(tc.f).Elem().FieldByName("Flags").Uint(); got != uint64(want) {
				t.Fatalf("spare bits changed flags: %d", got)
			}
		})
	}
	f := ie.NewOuterHeaderCreationFields(0x0100, 1, "192.0.2.1", "", 0, 0, 0)
	b, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	b[1] |= 0xf8
	if err := f.UnmarshalBinary(b); err != nil || f.OuterHeaderCreationDescription != 0x0100 {
		t.Fatalf("outer header spares: %v", err)
	}
}

func TestAddressConditionalArguments(t *testing.T) {
	v4 := net.ParseIP("192.0.2.1")
	for name, item := range map[string]*ie.IE{
		"UEUnflaggedDelegation":  ie.NewUEIPAddress(2, "192.0.2.1", "", 1, 0),
		"UEUnflaggedPrefix":      ie.NewUEIPAddress(2, "192.0.2.1", "", 0, 1),
		"UEConflictingChoose4":   ie.NewUEIPAddress(0x12, "", "", 0, 0),
		"UEConflictingChoose6":   ie.NewUEIPAddress(0x21, "", "", 0, 0),
		"OuterUnflaggedTEID":     ie.NewOuterHeaderCreation(0x1000, 1, "192.0.2.1", "", 0, 0, 0),
		"OuterUnflaggedPort":     ie.NewOuterHeaderCreation(0x1000, 0, "192.0.2.1", "", 1, 0, 0),
		"OuterUnflaggedCTag":     ie.NewOuterHeaderCreation(0x1000, 0, "192.0.2.1", "", 0, 1, 0),
		"OuterUnflaggedSTag":     ie.NewOuterHeaderCreation(0x1000, 0, "192.0.2.1", "", 0, 0, 1),
		"OuterCTagOverflow":      ie.NewOuterHeaderCreation(0x4000, 0, "", "", 0, 0x1000000, 0),
		"OuterSTagOverflow":      ie.NewOuterHeaderCreation(0x8000, 0, "", "", 0, 0, 0x1000000),
		"OuterLowLayerConflict":  ie.NewOuterHeaderCreation(0x0104, 1, "192.0.2.1", "", 0, 0, 0),
		"FTEIDChooseTEID":        ie.NewFTEID(5, 1, nil, nil, 0),
		"FTEIDUnflaggedChooseID": ie.NewFTEID(1, 1, v4, nil, 1),
	} {
		t.Run(name, func(t *testing.T) { assertInvalidConstructor(t, item) })
	}
	source := ie.NewSourceIPAddressFields(v4, nil, 0)
	source.MaskPrefixLength = 1
	if _, err := source.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("Source unflagged prefix: %v", err)
	}
}

func TestAddressFieldsRejectUnencodedValues(t *testing.T) {
	v4 := net.ParseIP("192.0.2.1")
	entity := ie.NewCPPFCPEntityIPAddressFields(v4, nil)
	entity.TEID = 1
	if _, err := entity.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("unencoded CP entity TEID: %v", err)
	}
	entity.TEID = 0
	entity.ChooseID = []byte{1}
	if _, err := entity.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("unencoded CP entity ChooseID: %v", err)
	}
	mac := net.HardwareAddr{2, 0, 0, 0, 0, 1}
	pmf := ie.NewPMFAddressInformationFields(v4, nil, 0, 0, mac, mac)
	pmf.Flags &^= 4
	if _, err := pmf.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Errorf("unflagged PMF MACs: %v", err)
	}
}

func TestAddressConditionalFieldsValidation(t *testing.T) {
	v4 := net.ParseIP("192.0.2.1")
	for _, fields := range []addressFields{
		ie.NewUEIPAddressFields(2, "192.0.2.1", "", 1, 0),
		ie.NewUEIPAddressFields(0x12, "", "", 0, 0),
		ie.NewOuterHeaderCreationFields(0x1000, 1, "192.0.2.1", "", 0, 0, 0),
		ie.NewOuterHeaderCreationFields(0x4000, 0, "", "", 0, 0x1000000, 0),
		ie.NewFTEIDFields(5, 1, nil, nil, 0),
		ie.NewFTEIDFields(1, 1, v4, nil, 1),
	} {
		if _, err := fields.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
			t.Errorf("%T Marshal: %v", fields, err)
		}
		if err := fields.MarshalTo(make([]byte, 64)); !errors.Is(err, ie.ErrInvalidValue) {
			t.Errorf("%T MarshalTo: %v", fields, err)
		}
	}
	// Selected scalar arguments survive at their width boundaries.
	f := ie.NewOuterHeaderCreationFields(0xcf00, 0xffffffff, "192.0.2.1", "2001:db8::1", 65535, 0xffffff, 0xffffff)
	b, err := f.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	got, err := ie.ParseOuterHeaderCreationFields(b)
	if err != nil || got.TEID != f.TEID || got.PortNumber != f.PortNumber || got.CTag != f.CTag || got.STag != f.STag {
		t.Fatal("outer header scalar changed")
	}
	for _, ue := range []*ie.UEIPAddressFields{
		ie.NewUEIPAddressFields(0x09, "", "2001:db8::1", 4, 0),
		ie.NewUEIPAddressFields(0x41, "", "2001:db8::1", 0, 72),
	} {
		b, err = ue.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := ie.ParseUEIPAddressFields(b)
		if err != nil || parsed.IPv6PrefixDelegationBits != ue.IPv6PrefixDelegationBits || parsed.IPv6PrefixLength != ue.IPv6PrefixLength {
			t.Fatal("UE prefix argument changed")
		}
	}
}
