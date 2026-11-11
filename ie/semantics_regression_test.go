package ie_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestConformanceTagValue(t *testing.T) {
	// Clauses 8.2.94 and 8.2.95: the 12 VID bits straddle the flag nibble.
	t.Run("CTAG", func(t *testing.T) {
		f, e := ie.NewCTAG(7, 3, 1, 0xabc).CTAG()
		if e != nil || f.CVID != 0xabc {
			t.Fatalf("VID got %d, error %v; want 2748", f.CVID, e)
		}
	})
	t.Run("STAG", func(t *testing.T) {
		f, e := ie.NewSTAG(7, 3, 1, 0xabc).STAG()
		if e != nil || f.CVID != 0xabc {
			t.Fatalf("VID got %d, error %v; want 2748", f.CVID, e)
		}
	})
}

func TestConformanceAnyMulticast(t *testing.T) {
	// Clause 8.2.137: A=1 forbids all address fields.
	f, e := ie.NewIPMulticastAddress(nil, nil, nil, nil).IPMulticastAddress()
	if e != nil {
		t.Fatal(e)
	}
	if !f.HasAny() {
		t.Fatal("Any flag lost")
	}
}

func TestConformanceTimerReservedUnit(t *testing.T) {
	// Clauses 8.2.29, 8.2.35 and 8.2.78: other units mean one minute.
	for _, typ := range []ie.IEType{ie.Timer, ie.DLBufferingDuration, ie.GracefulReleasePeriod} {
		t.Run(typ.String(), func(t *testing.T) {
			for _, unit := range []byte{5, 6} {
				i := ie.New(typ, []byte{unit<<5 | 3})
				var got time.Duration
				var e error
				switch typ {
				case ie.Timer:
					got, e = i.Timer()
				case ie.DLBufferingDuration:
					got, e = i.DLBufferingDuration()
				case ie.GracefulReleasePeriod:
					got, e = i.GracefulReleasePeriod()
				}
				if e != nil || got != 3*time.Minute {
					t.Errorf("unit %d got %v, error %v; want 3m", unit, got, e)
				}
			}
		})
	}
}

func TestConformanceUPFeatureLength(t *testing.T) {
	// Clause 8.2.25: feature groups extend through octet 16.
	want := bytes.Repeat([]byte{1}, 12)
	got := ie.NewUPFunctionFeatures(want...)
	if !bytes.Equal(got.Payload, want) {
		t.Fatalf("payload length %d; want 12", len(got.Payload))
	}
}

func TestConformanceApplyActionMBS(t *testing.T) {
	// Clause 8.2.26 includes FSSM and/or MBSU as a forwarding choice.
	for _, flags := range []byte{8, 16, 24} {
		if e := ie.NewApplyAction(0, flags).ValidateApplyAction(); e != nil {
			t.Errorf("MBS flags %d: %v", flags, e)
		}
	}
}

func TestConformanceSMFSetID(t *testing.T) {
	// Clause 8.2.131: a spare octet precedes the label-encoded FQDN.
	i := ie.NewSMFSetID("smf.example")
	if i.Payload[0] != 0 {
		t.Errorf("first octet %d; want spare 0", i.Payload[0])
	}
	p := append([]byte{0}, []byte{3, 's', 'm', 'f', 7, 'e', 'x', 'a', 'm', 'p', 'l', 'e'}...)
	got, e := ie.New(ie.SMFSetID, p).SMFSetID()
	if e != nil || got != "smf.example" {
		t.Errorf("got %q, error %v; want smf.example", got, e)
	}
}

func TestConformanceSourceDestinationSpare(t *testing.T) {
	// Clauses 8.2.2, 8.2.24 and 8.2.79: spare bits are not enum bits.
	s, e := ie.New(ie.SourceInterface, []byte{0xf1}).SourceInterface()
	if e != nil || s != 1 {
		t.Errorf("source %d error %v; want 1", s, e)
	}
	d, e := ie.New(ie.DestinationInterface, []byte{0xf1}).DestinationInterface()
	if e != nil || d != 1 {
		t.Errorf("destination %d error %v; want 1", d, e)
	}
	p, e := ie.New(ie.PDNType, []byte{0xf9}).PDNType()
	if e != nil || p != 1 {
		t.Errorf("PDN type %d error %v; want 1", p, e)
	}
}

func TestConformanceFailedRuleID(t *testing.T) {
	// Clause 8.2.80: MAR is uint16; SRR is uint8 (8.2.151); unknown types mean FAR.
	for _, tc := range []struct {
		typ   byte
		id    uint32
		width int
	}{{5, 0x1234, 2}, {6, 0x12, 1}, {7, 0x12345678, 4}} {
		t.Run(fmt.Sprint(tc.typ), func(t *testing.T) {
			p := make([]byte, 1+tc.width)
			p[0] = tc.typ
			switch tc.width {
			case 1:
				p[1] = byte(tc.id)
			case 2:
				binary.BigEndian.PutUint16(p[1:], uint16(tc.id))
			default:
				binary.BigEndian.PutUint32(p[1:], tc.id)
			}
			got, e := ie.New(ie.FailedRuleID, p).FailedRuleID()
			if e != nil || got != tc.id {
				t.Errorf("got %d error %v; want %d", got, e, tc.id)
			}
			if tc.typ < 7 && len(ie.NewFailedRuleID(tc.typ, tc.id).Payload) != len(p) {
				t.Error("constructor omitted rule ID")
			}
		})
	}
}

func TestConformanceFixedLength(t *testing.T) {
	// TS 29.244 V20.1.0 Table 8.1.2-1 marks these 14 IE types Fixed.
	fixed := map[ie.IEType]int{ie.Cause: 1, ie.OffendingIE: 2, ie.SequenceNumber: 4, ie.Metric: 1, ie.URSEQN: 4, ie.Multiplier: 12, ie.AggregatedURRID: 4, ie.FramedRouting: 4, ie.Weight: 1, ie.NumberOfReports: 2, ie.PortNumber: 4, ie.NWTTPortNumber: 4, ie.NFInstanceID: 16, ie.SNSSAI: 4}
	for typ, size := range fixed {
		t.Run(typ.String(), func(t *testing.T) {
			i := ie.New(typ, make([]byte, size+1))
			m := reflect.ValueOf(i).MethodByName(typ.String())
			v := m.Call(nil)
			if v[1].IsNil() {
				t.Errorf("accepted %d-octet payload; fixed length is %d", size+1, size)
			}
		})
	}
}
