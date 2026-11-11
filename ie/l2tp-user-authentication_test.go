package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestL2TPUserAuthenticationCurrent(t *testing.T) {
	for _, typ := range []uint16{0, 1, 2, 3, 4, 5, 65535} {
		for flags := 0; flags < 16; flags++ {
			required := 0
			switch typ {
			case 1:
				required = 5
			case 2, 5:
				required = 15
			case 3:
				required = 13
			}
			if flags&required != required {
				continue
			}
			v := ie.L2TPUserAuthenticationFields{Type: typ, Flags: ie.L2TPAuthenticationFlags(flags)}
			b := []byte{byte(typ >> 8), byte(typ), byte(flags)}
			if flags&1 != 0 {
				v.Name = []byte("name")
				b = append(b, 4, 'n', 'a', 'm', 'e')
			}
			if flags&2 != 0 {
				v.Challenge = []byte{0, 255}
				b = append(b, 2, 0, 255)
			}
			if flags&4 != 0 {
				v.Response = []byte{7}
				b = append(b, 1, 7)
			}
			if flags&8 != 0 {
				v.ID = 255
				b = append(b, 255)
			}
			assertAddedFields(t, ie.NewL2TPUserAuthentication(v), &v, b, "L2TPUserAuthentication", false)
			spare := bytes.Clone(b)
			spare[2] |= 0xf0
			got, e := ie.New(ie.L2TPUserAuthentication, spare).L2TPUserAuthentication()
			if e != nil || got.Flags != v.Flags {
				t.Fatal("spare flags", e)
			}
		}
	}
	var f *ie.L2TPUserAuthenticationFields
	assertAddedNilFields(t, f)
}
func TestL2TPUserAuthenticationValidation(t *testing.T) {
	for _, typ := range []uint16{1, 2, 3, 5} {
		for flags := 0; flags < 16; flags++ {
			required := 15
			if typ == 1 {
				required = 5
			}
			if typ == 3 {
				required = 13
			}
			if flags&required == required {
				continue
			}
			v := ie.L2TPUserAuthenticationFields{Type: typ, Flags: ie.L2TPAuthenticationFlags(flags)}
			assertAddedInvalid(t, ie.NewL2TPUserAuthentication(v))
		}
	}
	for _, v := range []ie.L2TPUserAuthenticationFields{{Flags: 16}, {Name: []byte{1}}, {Challenge: []byte{1}}, {Response: []byte{1}}, {ID: 1}, {Flags: 1, Name: make([]byte, 256)}, {Flags: 2, Challenge: make([]byte, 256)}, {Flags: 4, Response: make([]byte, 256)}} {
		assertAddedInvalid(t, ie.NewL2TPUserAuthentication(v))
	}
	for _, flag := range []byte{1, 2, 4} {
		v := ie.L2TPUserAuthenticationFields{Flags: ie.L2TPAuthenticationFlags(flag)}
		if flag == 1 {
			v.Name = make([]byte, 255)
		}
		if flag == 2 {
			v.Challenge = make([]byte, 255)
		}
		if flag == 4 {
			v.Response = make([]byte, 255)
		}
		if _, e := ie.NewL2TPUserAuthentication(v).Marshal(); e != nil {
			t.Fatal(e)
		}
	}
	// Receive flags are decoded as sent, independently of sender type rules.
	for _, typ := range []uint16{1, 2, 3, 5} {
		v, e := ie.New(ie.L2TPUserAuthentication, []byte{0, byte(typ), 0}).L2TPUserAuthentication()
		if e != nil || v.Type != typ || v.Flags != 0 {
			t.Fatal("receive flags inferred", e)
		}
	}
	_, e := ie.NewL2TPUserAuthentication(ie.L2TPUserAuthenticationFields{Type: 2}).Marshal()
	if !errors.Is(e, ie.ErrInvalidValue) {
		t.Fatal(e)
	}
}

var _ = registerConstructorFixture(ie.NewL2TPUserAuthentication, ie.L2TPUserAuthenticationFields{Type: 4})
