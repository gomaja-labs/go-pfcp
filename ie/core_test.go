package ie_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestIEFramingLengthAndEOF(t *testing.T) {
	for _, b := range [][]byte{{}, {0}, {0, 19}, {0, 19, 0}, {0, 19, 0, 2, 1}, {0, 1, 0, 7, 0, 19, 0, 4, 1, 2, 3}, {0x80, 1, 0, 3, 0, 1}, {0x80, 1, 0, 1}} {
		for _, parse := range []func([]byte) error{
			func(b []byte) error { _, e := ie.Parse(b); return e },
			func(b []byte) error { _, e := ie.ParseMultiIEs(b); return e },
			func(b []byte) error { _, e := ie.ParseMessageIEs(b); return e },
		} {
			if len(b) == 0 {
				continue
			} // Empty sequences are valid.
			err := parse(b)
			if len(b) >= 4 && b[1] == 1 && b[0] == 0 && err == nil {
				continue
			} // Deferred malformed group.
			if !errors.Is(err, ie.ErrInvalidLength) || !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Fatalf("%d-octet input: %v", len(b), err)
			}
			var de *ie.DecodeError
			if !errors.As(err, &de) {
				t.Fatalf("missing DecodeError: %v", err)
			}
			want := ie.IEType(0)
			if len(b) >= 2 {
				want = ie.IEType(b[0])<<8 | ie.IEType(b[1])
				if want == ie.CreatePDR {
					want = ie.Cause
				}
			}
			if de.Type != want {
				t.Fatalf("type=%d want=%d", de.Type, want)
			}
		}
	}
	_, err := ie.New(ie.Cause, nil).Cause()
	if !errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, ie.ErrInvalidLength) {
		t.Fatal("value error misclassified", err)
	}
}

func TestVendorConstructorRejectsStandardType(t *testing.T) {
	for _, typ := range []ie.IEType{0, 1, 19, 0x7fff} {
		for _, eid := range []uint16{0, 1, 10415, 65535} {
			for _, v := range []*ie.IE{ie.NewVendorSpecificIE(typ, eid, []byte{1}), ie.NewVendorSpecificGroupedIE(typ, eid, ie.NewCause(1))} {
				if v == nil {
					t.Fatal("constructor returned nil")
				}
				_, err := v.Marshal()
				var ee *ie.EncodeError
				if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &ee) || ee.Type != typ {
					t.Fatalf("type %d EID %d accepted: %v", typ, eid, err)
				}
			}
		}
	}
	for _, typ := range []ie.IEType{0x8000, 0xffff} {
		v := ie.New(typ, []byte{1})
		v.EnterpriseID = 123
		b, e := v.Marshal()
		if e != nil {
			t.Fatal(e)
		}
		d, e := ie.Parse(b)
		if e != nil || d.EnterpriseID != 123 {
			t.Fatal("New vendor EID not preserved", e)
		}
	}
}

func TestEqualInvalidStructure(t *testing.T) {
	bad := func() *ie.IE { return ie.NewRecoveryTimeStamp(time.Unix(-1, 0)) }
	if bad().Equal(ie.NewRecoveryTimeStamp(time.Unix(-2, 0))) {
		t.Fatal("different constructor error text ignored")
	}
	nestedA, nestedB := ie.NewCreatePDR(bad()), ie.NewCreatePDR(bad())
	nestedB.ChildIEs[0].Payload = []byte{42}
	if nestedA.Equal(nestedB) {
		t.Fatal("invalid child structure ignored")
	}
	g := ie.NewCreatePDR(bad())
	h := ie.NewCreatePDR(bad())
	if !g.Equal(g) || !g.Equal(h) || !h.Equal(g) {
		t.Fatal("equivalent invalid groups differ")
	}
	a, b := bad(), bad()
	a.Payload = []byte{1}
	b.Payload = []byte{2}
	if a.Equal(b) {
		t.Fatal("invalid payload ignored")
	}
	b.Payload = bytes.Clone(a.Payload)
	if !a.Equal(b) {
		t.Fatal("equivalent invalid leaf differs")
	}
	b.EnterpriseID = 1
	if a.Equal(b) {
		t.Fatal("invalid enterprise ID ignored")
	}
	h.ChildIEs = append(h.ChildIEs, ie.NewCause(1))
	if g.Equal(h) {
		t.Fatal("extra child ignored")
	}
	if g.Equal(ie.NewCreatePDR(ie.NewCause(1))) {
		t.Fatal("encodable and invalid groups equal")
	}
}

func TestTimestampFractionalLastSecond(t *testing.T) {
	for _, nsec := range []int64{0, 1, 500000000, 999999999} {
		want := time.Unix(1<<32-1, 0)
		v := ie.NewRecoveryTimeStamp(want.Add(time.Duration(nsec)))
		if _, e := v.Marshal(); e != nil {
			t.Fatal("last second rejected", e)
		}
		got, e := v.RecoveryTimeStamp()
		if e != nil || !got.Equal(want) {
			t.Fatal("not truncated to seconds", e)
		}
	}
	for _, ts := range []time.Time{time.Unix(1<<32, 0), time.Unix(-1, 999999999)} {
		if _, e := ie.NewRecoveryTimeStamp(ts).Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
			t.Fatal("outside window accepted", e)
		}
	}
}

func TestGroupedAddDoesNotOverwriteCopies(t *testing.T) {
	a, b, c := ie.NewCause(1), ie.NewCause(2), ie.NewCause(3)
	g := ie.NewCreatePDR(a)
	g.ChildIEs = make([]*ie.IE, 1, 8)
	g.ChildIEs[0] = a
	cp := *g
	if err := g.Add(b, nil); err != nil {
		t.Fatal(err)
	}
	if err := cp.Add(c); err != nil {
		t.Fatal(err)
	}
	if g.ChildIEs[1] != b || cp.ChildIEs[1] != c {
		t.Fatal("Add overwrote shared children")
	}
}

func TestVendorEnterpriseIDReserved(t *testing.T) {
	for _, v := range []*ie.IE{
		ie.NewVendorSpecificIE(0x8001, 10415, []byte{1}),
		ie.NewVendorSpecificGroupedIE(0x8001, 10415, ie.NewCause(1)),
	} {
		if _, e := v.Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
			t.Fatal("3GPP enterprise ID accepted for vendor IE", e)
		}
	}
}
