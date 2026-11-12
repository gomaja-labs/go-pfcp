package ie_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMalformedFields(t *testing.T) {
	cases := []struct {
		name string
		call func() error
	}{
		{"SDFFilter", func() error { _, e := ie.ParseSDFFilterFields([]byte{1, 0, 255, 255, 0}); return e }},
		{"PFDContents", func() error { _, e := ie.ParsePFDContentsFields([]byte{1, 0, 255, 255, 0, 0}); return e }},
		{"FlowDescriptionShort", func() error { _, e := ie.New(ie.FlowInformation, nil).FlowDescription(); return e }},
		{"FlowDescriptionLength", func() error { _, e := ie.New(ie.FlowInformation, []byte{1, 0, 5, 1, 2}).FlowDescription(); return e }},
		{"RemoteGTPUPeerDI", func() error { _, e := ie.ParseRemoteGTPUPeerFields([]byte{4, 0, 0}); return e }},
		{"FQCSIDShortNode", func() error { _, e := ie.New(ie.FQCSID, []byte{0}).CSIDs(); return e }},
		{"TraceSingleOctet", func() error { _, e := ie.ParseTraceInformationFields(make([]byte, 6)); return e }},
		{"SourceMPL", func() error { _, e := ie.ParseSourceIPAddressFields([]byte{6, 127, 0, 0, 1}); return e }},
		{"SubsequentTimeQuota", func() error { _, e := ie.New(ie.SubsequentTimeQuota, nil).SubsequentTimeQuota(); return e }},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if p := panicOf(func() {
				err := tc.call()
				if err == nil {
					t.Error("accepted malformed fields")
				}
			}); p != nil {
				t.Fatalf("panic %v", p)
			}
		})
	}
}

func TestEmptyFlagHelpers(t *testing.T) {
	for _, tc := range []struct {
		name string
		call func()
	}{
		{"HasBID", func() {
			if ie.New(ie.FiveGSUserPlaneNodeID, nil).HasBID() {
				t.Error("unexpected MAC")
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if p := panicOf(tc.call); p != nil {
				t.Fatalf("panic %v", p)
			}
		})
	}
}

func TestPFDContentsFlagRoundTrip(t *testing.T) {
	for flags := 0; flags < 256; flags++ {
		t.Run(fmt.Sprintf("flags_%02x", flags), func(t *testing.T) {
			values := []string{"a", "bc", "def", "ghij", "klmno"}
			for k := range values {
				if flags&(1<<k) == 0 {
					values[k] = ""
				}
			}
			lists := [][]string{{"p", "qr", ""}, {"stu", "v"}, {"w", "xyz"}}
			for k := range lists {
				if flags&(1<<(k+5)) == 0 {
					lists[k] = nil
				}
			}
			want := ie.NewPFDContentsFields(values[0], values[1], values[2], []byte(values[3]), values[4], lists[0], lists[1], lists[2])
			if p := panicOf(func() {
				b, e := want.Marshal()
				if e != nil {
					t.Fatal(e)
				}
				got, e := ie.ParsePFDContentsFields(b)
				if e != nil {
					t.Fatal(e)
				}
				if !reflect.DeepEqual(got, want) {
					t.Error("PFD contents lost fields")
				}
			}); p != nil {
				t.Fatalf("panic %v", p)
			}
		})
	}
}

func TestFQCSIDCounts(t *testing.T) {
	for _, addr := range []string{"192.0.2.1", "2001:db8::1", "01020304"} {
		for n := 0; n <= 15; n++ {
			ids := make([]uint16, n)
			for k := range ids {
				ids[k] = uint16(k + 700)
			}
			t.Run(addr+fmt.Sprintf("count_%d", n), func(t *testing.T) {
				if p := panicOf(func() {
					i := ie.NewFQCSID(addr, 2, ids...)
					got, e := i.CSIDs()
					if e != nil {
						t.Fatal(e)
					}
					if !reflect.DeepEqual(got, ids) && len(ids) > 0 {
						t.Error("CSIDs lost")
					}
					a, e := i.NodeAddress()
					if e != nil {
						t.Error("node address", e)
					}
					if len(a) == 0 {
						t.Error("empty address")
					}
				}); p != nil {
					t.Fatalf("panic %v", p)
				}
			})
		}
	}
}

func TestOuterHeaderFlagRoundTrip(t *testing.T) {
	for flags := 0; flags < 256; flags++ {
		t.Run(fmt.Sprintf("flags_%02x", flags), func(t *testing.T) {
			want := ie.NewOuterHeaderCreationFields(uint16(flags)<<8, 42, "192.0.2.1", "2001:db8::1", 8805, 0x123456, 0xabcdef)
			if p := panicOf(func() {
				b, e := want.Marshal()
				if e != nil {
					t.Fatal(e)
				}
				got, e := ie.ParseOuterHeaderCreationFields(b)
				if e != nil {
					t.Fatal(e)
				}
				if !reflect.DeepEqual(got, want) {
					t.Error("outer header lost fields")
				}
			}); p != nil {
				t.Fatalf("panic %v", p)
			}
		})
	}
}

func TestMACListRoundTrip(t *testing.T) {
	for _, removed := range []bool{false, true} {
		for n := 0; n < 3; n++ {
			for tags := 0; tags < 4; tags++ {
				t.Run(fmt.Sprintf("count_%d", n)+fmt.Sprintf("tags_%d", tags)+map[bool]string{false: "detected", true: "removed"}[removed], func(t *testing.T) {
					ct, st := ie.New(ie.CTAG, nil), ie.New(ie.STAG, nil)
					if tags&1 != 0 {
						ct = ie.NewCTAG(7, 1, 1, 100)
					}
					if tags&2 != 0 {
						st = ie.NewSTAG(7, 3, 0, 200)
					}
					macs := make([]net.HardwareAddr, n)
					for k := range macs {
						macs[k] = net.HardwareAddr{2, 0, 0, 0, 0, byte(k)}
					}
					if p := panicOf(func() {
						if removed {
							want := ie.NewMACAddressesRemovedFields(ct, st, macs...)
							b, e := want.Marshal()
							if e != nil {
								t.Fatal(e)
							}
							got, e := ie.ParseMACAddressesRemovedFields(b)
							if e != nil {
								t.Fatal(e)
							}
							if !reflect.DeepEqual(got, want) {
								t.Error("removed MACs/tags lost")
							}
						} else {
							want := ie.NewMACAddressesDetectedFields(ct, st, macs...)
							b, e := want.Marshal()
							if e != nil {
								t.Fatal(e)
							}
							got, e := ie.ParseMACAddressesDetectedFields(b)
							if e != nil {
								t.Fatal(e)
							}
							if !reflect.DeepEqual(got, want) {
								t.Error("detected MACs/tags lost")
							}
						}
					}); p != nil {
						t.Fatalf("panic %v", p)
					}
				})
			}
		}
	}
}

func TestTraceInformationRoundTrip(t *testing.T) {
	for _, ip := range []net.IP{net.ParseIP("192.0.2.1"), net.ParseIP("2001:db8::1")} {
		want := ie.NewTraceInformationFields("001", "01", "abc", []byte{1, 2}, 3, []byte{4, 5, 6}, ip)
		b, e := want.Marshal()
		if e != nil {
			t.Fatal(e)
		}
		got, e := ie.ParseTraceInformationFields(b)
		if e != nil {
			t.Fatal(e)
		}
		if !reflect.DeepEqual(got, want) {
			t.Error("trace fields lost")
		}
		for n := 0; n < len(b); n++ {
			if p := panicOf(func() {
				_, err := ie.ParseTraceInformationFields(b[:n:n])
				if !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Errorf("prefix %d error %v", n, err)
				}
			}); p != nil {
				t.Errorf("prefix %d panic %v", n, p)
			}
		}
	}
}

func TestDownlinkDataServiceFlags(t *testing.T) {
	for flags := 0; flags < 4; flags++ {
		ppi, qfi := uint8(0), uint8(0)
		if flags&1 != 0 {
			ppi = 5
		}
		if flags&2 != 0 {
			qfi = 9
		}
		i := ie.NewDownlinkDataServiceInformation(flags&1 != 0, flags&2 != 0, ppi, qfi)
		if flags&2 != 0 {
			if p := panicOf(func() {
				v, e := i.QFI()
				if e != nil || v != 9 {
					t.Error("QFI", v, e)
				}
			}); p != nil {
				t.Error("QFI panic", p)
			}
		}
	}
}

func TestSDFFilterEmptyFields(t *testing.T) {
	_, e := ie.ParseSDFFilterFields([]byte{0, 0})
	if e != nil {
		t.Fatal("two flag/spare octets are valid", e)
	}
	f, e := ie.ParseSDFFilterFields([]byte{1, 0, 0, 0})
	if e != nil || f.FlowDescription != "" {
		t.Fatal("empty FD", e)
	}
}

func TestFieldsShortMarshal(t *testing.T) {
	for name, parse := range fieldParsers {
		t.Run(name, func(t *testing.T) {
			for _, payload := range payloadFamily() {
				var values []reflect.Value
				if p := panicOf(func() { values = reflect.ValueOf(parse).Call([]reflect.Value{reflect.ValueOf(payload)}) }); p != nil {
					continue
				}
				if !values[1].IsNil() {
					continue
				}
				f := values[0].Interface().(interface {
					MarshalLen() int
					MarshalTo([]byte) error
				})
				for n := 0; n < f.MarshalLen(); n++ {
					if p := panicOf(func() {
						if err := f.MarshalTo(bytes.Repeat([]byte{0}, n)); err == nil {
							t.Errorf("accepted buffer %d of %d", n, f.MarshalLen())
						}
					}); p != nil {
						t.Fatalf("short MarshalTo length %d: %v", n, p)
					}
				}
			}
		})
	}
}

// TS 29.244 V20.1.0 clause 8.1.2 requires ignoring spare bits and unknown
// extension octets. These payloads triggered panics, but are not invalid lengths.
func TestSpareBitsAndExtensions(t *testing.T) {
	if _, err := ie.ParseOuterHeaderCreationFields([]byte{0, 0x40, 0, 0, 0}); err != nil {
		t.Fatal(err)
	}
	if _, err := ie.ParseMACAddressesDetectedFields([]byte{0, 0, 0, 0, 0, 0, 0}); err != nil {
		t.Fatal(err)
	}
	if _, err := ie.ParseMACAddressesRemovedFields([]byte{0, 0, 0, 0, 0, 0, 0}); err != nil {
		t.Fatal(err)
	}
}
