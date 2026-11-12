package ie_test

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"net"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

type bufferEncoder interface {
	MarshalLen() int
	MarshalTo([]byte) error
}

// TS 29.244 V20.1.0 clause 8.1.2: send every spare octet as zero.
func TestMarshalToInitializesCoveredBytes(t *testing.T) {
	check := func(t *testing.T, value bufferEncoder) {
		t.Helper()
		n := value.MarshalLen()
		zero, dirty := make([]byte, n+8), bytes.Repeat([]byte{255}, n+8)
		if err := value.MarshalTo(zero[:n]); err != nil {
			t.Fatal(err)
		}
		if err := value.MarshalTo(dirty); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(zero[:n], dirty[:n]) {
			t.Error("encoding depends on prior buffer contents")
		}
		if !bytes.Equal(dirty[n:], bytes.Repeat([]byte{255}, 8)) {
			t.Error("wrote beyond MarshalLen")
		}
	}
	for name, entry := range constructors {
		t.Run(name, func(t *testing.T) {
			value := callConstructor(entry)[0].Interface().(bufferEncoder)
			check(t, value)
		})
	}
	// DI lengths may include future extension octets (clauses 8.2.70/8.2.24).
	t.Run("RemoteInterfaceExtension", func(t *testing.T) {
		check(t, &ie.RemoteGTPUPeerFields{Flags: 6, IPv4Address: net.IP{192, 0, 2, 1}, DILength: 3, DestinationInterface: 1})
	})
}

func assertInvalidConstructor(t *testing.T, item *ie.IE) {
	t.Helper()
	if item == nil {
		t.Fatal("constructor returned nil")
	}
	group := ie.NewCreatePDR(item)
	msg := message.NewSessionEstablishmentRequest(1, 2, group)
	for name, marshal := range map[string]func() ([]byte, error){
		"IE": item.Marshal, "group": group.Marshal, "message": msg.Marshal,
		"IETo":      func() ([]byte, error) { return nil, item.MarshalTo(make([]byte, item.MarshalLen())) },
		"groupTo":   func() ([]byte, error) { return nil, group.MarshalTo(make([]byte, group.MarshalLen())) },
		"messageTo": func() ([]byte, error) { return nil, msg.MarshalTo(make([]byte, msg.MarshalLen())) },
	} {
		_, err := marshal()
		var enc *ie.EncodeError
		if !errors.As(err, &enc) || !errors.Is(err, ie.ErrInvalidValue) || enc.Type != item.Type {
			t.Errorf("%s: invalid argument lost: %v", name, err)
		}
	}
}

// A rejected argument must survive grouping and message construction.
func TestConstructorsNeverReturnNil(t *testing.T) {
	for name, entry := range constructors {
		t.Run(name, func(t *testing.T) {
			fn := reflect.ValueOf(entry.fn)
			typ := fn.Type()
			for variant := 0; variant < 4*(typ.NumIn()+1); variant++ {
				args := make([]reflect.Value, typ.NumIn())
				for k := range args {
					argType := typ.In(k)
					v := reflect.New(argType).Elem()
					if variant%4 != 0 {
						switch argType.Kind() {
						case reflect.String:
							v.SetString(strings.Repeat("!", 256))
						case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
							v.SetUint(^uint64(0) >> (64 - argType.Bits()))
						case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
							v.SetInt(-1)
						case reflect.Slice:
							size := 1
							if variant%4 == 2 {
								size = 17
							}
							if variant%4 == 3 {
								size = 256
							}
							v.Set(reflect.MakeSlice(argType, size, size))
						}
					}
					args[k] = v
					if variant >= 4 && k != variant/4-1 {
						args[k] = constructorArguments(entry)[k]
					}
				}
				var results []reflect.Value
				if p := panicOf(func() {
					if typ.IsVariadic() {
						results = fn.CallSlice(args)
					} else {
						results = fn.Call(args)
					}
				}); p != nil {
					t.Fatalf("variant %d: panic %v", variant, p)
				}
				if results[0].IsNil() {
					t.Errorf("variant %d: constructor returned nil", variant)
					continue
				}
				if item, ok := results[0].Interface().(*ie.IE); ok {
					if _, err := item.Marshal(); err != nil {
						assertInvalidConstructor(t, item)
					}
				}
			}
		})
	}
}

func TestConstructorsHaveNoNilReturn(t *testing.T) {
	paths, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	for _, path := range paths {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv != nil || !strings.HasPrefix(fn.Name.Name, "New") {
				continue
			}
			ast.Inspect(fn.Body, func(n ast.Node) bool {
				ret, ok := n.(*ast.ReturnStmt)
				if !ok || len(ret.Results) != 1 {
					return true
				}
				if id, ok := ret.Results[0].(*ast.Ident); ok && id.Name == "nil" {
					t.Errorf("%s: %s returns nil", path, fn.Name.Name)
				}
				return true
			})
		}
	}
}

func TestRemotePeerTimestampConstructor(t *testing.T) {
	// Require the timestamp in both convenience APIs, rather than a hidden FromFields-only path.
	ts := time.Unix(1700000000, 0).UTC()
	for name, fn := range map[string]any{"IE": ie.NewRemoteGTPUPeer, "Fields": ie.NewRemoteGTPUPeerFields} {
		t.Run(name, func(t *testing.T) {
			ctor := reflect.ValueOf(fn)
			if ctor.Type().NumIn() != 6 {
				t.Fatal("constructor cannot accept Recovery Timestamp")
			}
			got := ctor.Call([]reflect.Value{reflect.ValueOf(uint8(0x1e)), reflect.ValueOf("192.0.2.1"), reflect.ValueOf(""), reflect.ValueOf(uint8(1)), reflect.ValueOf("net"), reflect.ValueOf(ts)})[0].Interface().(interface{ Marshal() ([]byte, error) })
			wire, err := got.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if name == "IE" {
				item, e := ie.Parse(wire)
				if e != nil {
					t.Fatal(e)
				}
				wire = item.Payload
			}
			want := []byte{0x1e, 192, 0, 2, 1, 0, 1, 1, 0, 3, 'n', 'e', 't', 0xe8, 0xfe, 0x6f, 0x80}
			if !bytes.Equal(wire, want) {
				t.Error("timestamp not encoded at the specified offset")
			}
			f, e := ie.ParseRemoteGTPUPeerFields(wire)
			if e != nil || !f.RecoveryTimestamp.Equal(ts) {
				t.Fatalf("timestamp decode: %v", e)
			}
			wire[7] |= 0xf0
			f, e = ie.ParseRemoteGTPUPeerFields(wire)
			if e != nil || f.DestinationInterface != 1 {
				t.Fatalf("DI spare bits changed value: %v", e)
			}
		})
	}
}

func TestFTEIDRequiresFlaggedAddressFamily(t *testing.T) {
	for _, tc := range []struct {
		name   string
		flags  uint8
		v4, v6 net.IP
	}{
		{"missing4", 1, nil, nil}, {"missing6", 2, nil, nil}, {"missingBoth", 3, nil, nil},
		{"wrong4", 1, net.ParseIP("2001:db8::1"), nil}, {"wrong6", 2, nil, net.ParseIP("192.0.2.1")},
		{"short4", 1, net.IP{1, 2, 3}, nil}, {"short6", 2, nil, net.IP{1, 2, 3}},
		{"missing4with6", 3, nil, net.ParseIP("2001:db8::1")}, {"missing6with4", 3, net.ParseIP("192.0.2.1"), nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			f := ie.NewFTEIDFields(tc.flags, 1, tc.v4, tc.v6, 0)
			if err := f.MarshalTo(make([]byte, f.MarshalLen())); !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("Fields.MarshalTo error %v", err)
			}
			assertInvalidConstructor(t, ie.NewFTEID(tc.flags, 1, tc.v4, tc.v6, 0))
		})
	}
	for flags := uint8(5); flags <= 7; flags++ {
		t.Run(fmt.Sprint("choose", flags), func(t *testing.T) {
			if _, err := ie.NewFTEID(flags, 0, nil, nil, 0).Marshal(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestPMFEmptyReceivedList(t *testing.T) {
	assertInvalidConstructor(t, ie.NewPMFControlInformation(4))
	f, err := ie.ParsePMFControlInformationFields([]byte{0xfc, 0, 0xff})
	if err != nil {
		t.Fatal(err)
	}
	if f.Flags != 4 || len(f.QFIs) != 0 {
		t.Fatal("received empty list or spare bits changed")
	}
}

func TestObsoleteFlagHelpersRemoved(t *testing.T) {
	typ := reflect.TypeFor[*ie.IE]()
	for _, name := range []string{"HasMODE", "HasNORD", "HasSESRL"} {
		if _, ok := typ.MethodByName(name); ok {
			t.Errorf("obsolete helper %s remains", name)
		}
	}
	if reflect.TypeOf(ie.NewQERControlIndications).NumIn() != 1 || reflect.TypeOf(ie.NewReportingFrequency).NumIn() != 2 {
		t.Error("constructors still accept obsolete flag arguments")
	}

}

// TS 29.244 V20.1.0 clause 8.2.43: the count is four bits and
// Node Type is a mandatory separate octet after the advertised CSIDs.
func TestFQCSIDCountBoundary(t *testing.T) {
	for _, addr := range []string{"192.0.2.1", "2001:db8::1", "00000001"} {
		for count := 0; count <= 16; count++ {
			t.Run(fmt.Sprintf("%s/%d", addr, count), func(t *testing.T) {
				ids := make([]uint16, count)
				for n := range ids {
					ids[n] = uint16(0x102 + n)
				}
				item := ie.NewFQCSID(addr, 5, ids...)
				if count > 15 {
					assertInvalidConstructor(t, item)
					return
				}
				got, err := item.CSIDs()
				if err != nil || !reflect.DeepEqual(got, ids) {
					t.Fatalf("CSIDs changed: %v", err)
				}
				item.Payload[len(item.Payload)-1] |= 0xf0
				node, err := item.FQCSIDNodeType()
				if err != nil || node != 5 {
					t.Fatalf("Node Type/spare bits: %v", err)
				}
				item.Payload = item.Payload[:len(item.Payload)-1]
				if _, err := item.CSIDs(); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("missing Node Type: %v", err)
				}
			})
		}
	}
}

// TS 29.244 V20.1.0 clause 8.2.3 requires V4 or V6 and permits CHID only with CH.
func TestFTEIDFlagRules(t *testing.T) {
	for flags := 0; flags < 256; flags++ {
		var v4, v6 net.IP
		if flags&4 == 0 {
			if flags&1 != 0 {
				v4 = net.ParseIP("192.0.2.1")
			}
			if flags&2 != 0 {
				v6 = net.ParseIP("2001:db8::1")
			}
		}
		valid := flags&^15 == 0 && flags&3 != 0 && (flags&8 == 0 || flags&4 != 0)
		item := ie.NewFTEID(uint8(flags), 0, v4, v6, 0)
		_, err := item.Marshal()
		if !valid {
			if !errors.Is(err, ie.ErrInvalidValue) {
				t.Errorf("invalid flags %d: %v", flags, err)
			}
			continue
		}
		if err != nil {
			t.Fatalf("valid flags %d: %v", flags, err)
		}
		payload := bytes.Clone(item.Payload)
		payload[0] |= 0xf0
		f, err := ie.ParseFTEIDFields(payload)
		if err != nil || int(f.Flags) != flags {
			t.Fatalf("received spare bits changed flags %d: %v", flags, err)
		}
	}
}

func TestCurrentFlagValuesIgnoreSpare(t *testing.T) {
	for _, tc := range []struct {
		typ  ie.IEType
		get  func(*ie.IE) (uint8, error)
		want uint8
	}{
		{ie.QERControlIndications, (*ie.IE).QERControlIndications, 1},
		{ie.ReportingFrequency, (*ie.IE).ReportingFrequency, 3},
	} {
		got, err := tc.get(ie.New(tc.typ, []byte{255}))
		if err != nil || got != tc.want {
			t.Errorf("received spare bits affect %s: %v", tc.typ, err)
		}
	}
	f, err := ie.ParseDLFlowLevelMarkingFields([]byte{0xff, 0x12, 0x34, 0x56, 0x78})
	if err != nil || f.Flags != 3 {
		t.Errorf("received spare bits affect DL flags: %v", err)
	}
}

func TestCurrentFlagConstructors(t *testing.T) {
	qer, ok := any(ie.NewQERControlIndications).(func(bool) *ie.IE)
	if !ok {
		t.Fatal("QER constructor must expose only RCSR")
	}
	freq, ok := any(ie.NewReportingFrequency).(func(bool, bool) *ie.IE)
	if !ok {
		t.Fatal("reporting constructor must expose only PERIO/EVETT")
	}
	for _, on := range []bool{false, true} {
		want := byte(0)
		if on {
			want = 1
		}
		if !bytes.Equal(qer(on).Payload, []byte{want}) {
			t.Error("RCSR or spare bits changed")
		}
	}
	for flags := 0; flags < 4; flags++ {
		item := freq(flags&2 != 0, flags&1 != 0)
		if flags == 0 {
			assertInvalidConstructor(t, item)
			continue
		}
		if !bytes.Equal(item.Payload, []byte{byte(flags)}) {
			t.Error("PERIO/EVETT or reserved bits changed")
		}
	}
}
