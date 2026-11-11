package ie_test

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestLeafValueContracts(t *testing.T) {
	// Fixed versus extendable classifications: TS 29.244 V20.1.0 Table 8.1.2-1.
	for _, tc := range []struct {
		i       *ie.IE
		minimum int
		fixed   bool
	}{
		{ie.NewCause(1), 1, true}, {ie.NewOffendingIE(1), 2, true},
		{ie.NewSequenceNumber(1), 4, true}, {ie.NewMetric(1), 1, true},
		{ie.NewURSEQN(1), 4, true}, {ie.NewMultiplier(1, 1), 12, true},
		{ie.NewAggregatedURRID(1), 4, true}, {ie.NewFramedRouting(1), 4, true},
		{ie.NewWeight(1), 1, true}, {ie.NewNumberOfReports(1), 2, true},
		{ie.NewPortNumber(1), 4, true}, {ie.NewNWTTPortNumber(1), 4, true},
		{ie.NewNFInstanceID(make([]byte, 16)), 16, true}, {ie.NewSNSSAI(1, 1), 4, true},
		{ie.NewSourceInterface(1), 1, false}, {ie.NewDestinationInterface(1), 1, false},
		{ie.NewPDNType(1), 1, false}, {ie.NewQFI(1), 1, false},
		{ie.NewTimer(2 * time.Second), 1, false}, {ie.NewDLBufferingDuration(2 * time.Second), 1, false},
		{ie.NewGracefulReleasePeriod(2 * time.Second), 1, false},
		{ie.NewCTAG(7, 3, 1, 0xabc), 3, false}, {ie.NewSTAG(7, 3, 1, 0xabc), 3, false},
		{ie.NewIPMulticastAddress(nil, nil, nil, nil), 1, false},
		{ie.NewFailedRuleID(5, 0x1234), 3, false},
		{ie.NewCreateBridgeRouterInfo(2), 1, false},
	} {
		t.Run(tc.i.Type.String(), func(t *testing.T) {
			name := tc.i.Type.String()
			invoke := func(i *ie.IE) (any, error) {
				t.Helper()
				v := reflect.ValueOf(i).MethodByName(name).Call(nil)
				var e error
				if !v[1].IsNil() {
					e = v[1].Interface().(error)
				}
				return v[0].Interface(), e
			}
			if _, e := invoke(nil); !errors.Is(e, ie.ErrIENotFound) {
				t.Fatalf("nil error=%v", e)
			}
			var wrong *ie.InvalidTypeError
			if _, e := invoke(ie.New(0xffff, nil)); !errors.As(e, &wrong) {
				t.Fatalf("wrong type error=%v", e)
			}
			for n := 0; n < tc.minimum; n++ {
				if _, e := invoke(ie.New(tc.i.Type, tc.i.Payload[:n])); e == nil {
					t.Fatalf("accepted truncated length %d", n)
				}
			}
			expected, e := invoke(tc.i)
			if e != nil {
				t.Fatal(e)
			}
			extra := ie.New(tc.i.Type, append(bytes.Clone(tc.i.Payload), 0xaa, 0x55))
			got, e := invoke(extra)
			if tc.fixed {
				if !errors.Is(e, ie.ErrInvalidLength) {
					t.Fatalf("fixed extension error=%v", e)
				}
				return
			}
			if e != nil || !reflect.DeepEqual(got, expected) {
				t.Fatalf("extension changed value, error=%v", e)
			}
			raw, e := extra.Marshal()
			if e != nil {
				t.Fatal(e)
			}
			parsed, e := ie.Parse(raw)
			if e != nil {
				t.Fatal(e)
			}
			again, e := parsed.Marshal()
			if e != nil || !bytes.Equal(raw, again) {
				t.Fatal("opaque IE did not retain extension")
			}
		})
	}
}

func TestLeafSemanticBoundaries(t *testing.T) {
	// TS 29.244 V20.1.0 clauses 8.2.140 and 8.2.89.
	if v, e := ie.NewCreateBridgeRouterInfo(2).CreateBridgeRouterInfo(); e != nil || v != 2 {
		t.Fatalf("RII constructor value=%d error=%v", v, e)
	}
	if v, e := ie.New(ie.QFI, []byte{0xff}).QFI(); e != nil || v != 63 {
		t.Fatalf("direct QFI spare bits value=%d error=%v", v, e)
	}

	// TS 29.244 V20.1.0 clause 8.2.26: MBS and ordinary forwarding are exclusive.
	for _, f := range []byte{8, 16, 24} {
		if e := ie.NewApplyAction(2, f).ValidateApplyAction(); e == nil {
			t.Fatal("accepted FORW with MBS")
		}
	}
	if e := ie.NewApplyAction().ValidateApplyAction(); e == nil {
		t.Fatal("accepted null action")
	}
	// Clauses 8.2.27 and 8.2.89: absence and spare bits are independent.
	if _, e := ie.NewDownlinkDataServiceInformation(false, true, 0, 9).PPI(); !errors.Is(e, ie.ErrIENotFound) {
		t.Fatalf("PPI absence error=%v", e)
	}
	if _, e := ie.New(ie.DownlinkDataServiceInformation, []byte{1}).PPI(); !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatalf("PPI truncation error=%v", e)
	}
	if v, e := ie.New(ie.DownlinkDataServiceInformation, []byte{3, 0xff, 0xff}).PPI(); e != nil || v != 7 {
		t.Fatalf("PPI spare bits value=%d error=%v", v, e)
	}
	if v, e := ie.New(ie.DownlinkDataServiceInformation, []byte{3, 0, 0xff}).QFI(); e != nil || v != 63 {
		t.Fatalf("QFI spare bits value=%d error=%v", v, e)
	}
	// Clause 8.2.62: choosing IPv4 must not suppress a supplied IPv6 address.
	i := ie.NewUEIPAddress(0x11, "", "2001:db8::1", 0, 0)
	if !i.HasIPv6() {
		t.Fatal("CHV4 suppressed IPv6")
	}
	if ie.NewUEIPAddress(0x21, "", "", 0, 0).HasIPv6() {
		t.Fatal("CHV6 did not suppress IPv6")
	}
	// Clause 8.2.80: reject values that cannot fit their chosen rule type.
	for _, c := range []struct {
		typ byte
		id  uint32
	}{{0, 65536}, {5, 65536}, {4, 256}, {6, 256}, {7, 1}, {255, 1}} {
		requireInvalidValue(t, "rule identifier", ie.FailedRuleID, ie.NewFailedRuleID(c.typ, c.id))
	}
	for _, typ := range []byte{7, 31, 0xe7} {
		i := ie.New(ie.FailedRuleID, []byte{typ, 1, 2, 3, 4, 0xaa})
		v, e := i.FailedRuleID()
		kind, ke := i.RuleIDType()
		if e != nil || ke != nil || v != 0x01020304 || kind != ie.RuleIDTypeFAR {
			t.Fatalf("future rule type %d value=%d kind=%d errors=%v/%v", typ, v, kind, e, ke)
		}
	}
	// Clauses 8.2.131 and 8.2.25: required framing and null value behavior.
	if _, e := ie.New(ie.SMFSetID, []byte{0}).SMFSetID(); e == nil {
		t.Fatal("accepted missing SMF name")
	}
	if _, e := ie.New(ie.UPFunctionFeatures, nil).UPFunctionFeatures(); e == nil {
		t.Fatal("accepted null feature value")
	}
	if i := ie.NewSMFSetID("smf.example"); i.Payload[0] != 0 {
		t.Fatal("missing SMF spare")
	}
	for n := 0; n <= 17; n++ {
		input := bytes.Repeat([]byte{0x55}, n)
		got := ie.NewUPFunctionFeatures(input...)
		if len(got.Payload) < 2 || len(got.Payload)%2 != 0 || !bytes.Equal(got.Payload[:n], input) {
			t.Fatalf("lost features at length %d", n)
		}
	}
	// Clause 8.2.129: Alternative SMF IP Address has no F-TEID fields.
	fields := reflect.TypeFor[ie.AlternativeSMFIPAddressFields]()
	for _, name := range []string{"TEID", "ChooseID"} {
		if _, ok := fields.FieldByName(name); ok {
			t.Errorf("unrelated field %s", name)
		}
	}
}

func TestAlternativeSMFPreference(t *testing.T) {
	typ := reflect.TypeFor[*ie.AlternativeSMFIPAddressFields]()
	for _, name := range []string{"HasPPE", "SetPPEFlag"} {
		if _, ok := typ.MethodByName(name); !ok {
			t.Fatalf("missing fields helper %s", name)
		}
	}
	call := func(f *ie.AlternativeSMFIPAddressFields, name string) []reflect.Value {
		return reflect.ValueOf(f).MethodByName(name).Call(nil)
	}
	if call(nil, "HasPPE")[0].Bool() {
		t.Fatal("nil preferred")
	}
	call(nil, "SetPPEFlag")
	f := &ie.AlternativeSMFIPAddressFields{Flags: 2, IPv4Address: []byte{192, 0, 2, 1}}
	if call(f, "HasPPE")[0].Bool() {
		t.Fatal("unexpected preference")
	}
	call(f, "SetPPEFlag")
	if !call(f, "HasPPE")[0].Bool() {
		t.Fatal("preference setter lost flag")
	}
	b, e := f.Marshal()
	if e != nil {
		t.Fatal(e)
	}
	g, e := ie.ParseAlternativeSMFIPAddressFields(b)
	if e != nil || !call(g, "HasPPE")[0].Bool() {
		t.Fatalf("preference did not round trip: %v", e)
	}
}

// SRR IDs use the one-octet format of clause 8.2.151, including when carried
// in Failed Rule ID (TS 29.244 V20.1.0 clause 8.2.80).
func TestFailedSRRIDWidth(t *testing.T) {
	i := ie.NewFailedRuleID(6, 0x12)
	if len(i.Payload) != 2 {
		t.Fatalf("SRR ID payload length %d; want 2", len(i.Payload))
	}
	for _, p := range [][]byte{{6, 0x12}, {6, 0x12, 0x34, 0x56, 0x78}} {
		v, e := ie.New(ie.FailedRuleID, p).FailedRuleID()
		if e != nil || v != 0x12 {
			t.Fatalf("SRR ID=%d error=%v", v, e)
		}
	}
	if _, e := ie.New(ie.FailedRuleID, []byte{6}).FailedRuleID(); !errors.Is(e, io.ErrUnexpectedEOF) {
		t.Fatalf("truncated SRR ID error=%v", e)
	}
}
