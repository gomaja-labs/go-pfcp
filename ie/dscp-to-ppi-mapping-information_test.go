package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestDSCPToPPIMappingInformationCurrent(t *testing.T) {
	for ppi := 0; ppi < 8; ppi++ {
		for _, n := range []int{1, 2, 63} {
			v := ie.DSCPToPPIMappingInformationFields{PPI: uint8(ppi), DSCP: make([]uint8, n)}
			for j := range v.DSCP {
				v.DSCP[j] = uint8(j + 1)
			}
			i := ie.NewDSCPToPPIMappingInformation(v)
			for j := 0; j < i.MarshalLen(); j++ {
				if e := i.MarshalTo(make([]byte, j)); !errors.Is(e, ie.ErrInvalidLength) {
					t.Fatal("short IE buffer", e)
				}
			}
			b := append([]byte{byte(ppi)}, v.DSCP...)
			wire, e := i.Marshal()
			if e != nil || !bytes.Equal(wire[4:], b) {
				t.Fatal("wire", e)
			}
			p, e := ie.Parse(wire)
			if e != nil {
				t.Fatal(e)
			}
			got, e := p.DSCPToPPIMappingInformation()
			if e != nil || !reflect.DeepEqual(got, &v) {
				t.Fatal("value", e)
			}
			again, e := p.Marshal()
			if e != nil || !bytes.Equal(wire, again) {
				t.Fatal("roundtrip", e)
			}
			spare := bytes.Clone(b)
			spare[0] |= 0xf8
			for j := 1; j < len(spare); j++ {
				spare[j] |= 0xc0
			}
			got, e = ie.New(ie.DSCPToPPIMappingInformation, spare).DSCPToPPIMappingInformation()
			if e != nil || !reflect.DeepEqual(got, &v) {
				t.Fatal("spares", e)
			}
			if v.MarshalLen() != len(b) {
				t.Fatal("length")
			}
			for j := 0; j < len(b); j++ {
				if v.MarshalTo(make([]byte, j)) == nil {
					t.Fatal("short buffer")
				}
			}
			dst := make([]byte, len(b))
			if e = v.MarshalTo(dst); e != nil || !bytes.Equal(dst, b) {
				t.Fatal("MarshalTo", e)
			}
			if e = v.UnmarshalBinary(b); e != nil {
				t.Fatal(e)
			}
		}
	}
	for _, n := range []int{0, 1, 65, 66} {
		_, e := ie.New(ie.DSCPToPPIMappingInformation, make([]byte, n)).DSCPToPPIMappingInformation()
		if !errors.Is(e, ie.ErrMalformed) {
			t.Fatalf("length %d: %v", n, e)
		}
	}
	for _, v := range []ie.DSCPToPPIMappingInformationFields{{}, {PPI: 8, DSCP: []byte{0}}, {DSCP: make([]byte, 64)}, {DSCP: []byte{64}}} {
		assertAddedInvalid(t, ie.NewDSCPToPPIMappingInformation(v))
	}
	var f *ie.DSCPToPPIMappingInformationFields
	assertAddedNilFields(t, f)
	var i *ie.IE
	if _, e := i.DSCPToPPIMappingInformation(); !errors.Is(e, ie.ErrIENotFound) {
		t.Fatal(e)
	}
	if _, e := ie.New(0, nil).DSCPToPPIMappingInformation(); e == nil {
		t.Fatal("wrong type")
	}
}

var _ = registerConstructorFixture(ie.NewDSCPToPPIMappingInformation, ie.DSCPToPPIMappingInformationFields{PPI: 7, DSCP: []byte{0, 63}})
