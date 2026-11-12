package ie_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"net"
	"net/netip"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// The registry is checked against the package declarations, so new APIs cannot
// silently escape the constructor and fields-parser harnesses.
func TestAPIRegistry(t *testing.T) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	functions := map[string]bool{}
	for _, path := range files {
		if strings.HasSuffix(path, "_test.go") {
			continue
		}
		f, e := parser.ParseFile(token.NewFileSet(), path, nil, 0)
		if e != nil {
			t.Fatal(e)
		}
		for _, d := range f.Decls {
			fn, ok := d.(*ast.FuncDecl)
			if !ok || fn.Recv != nil {
				continue
			}
			name := fn.Name.Name
			if strings.HasPrefix(name, "New") {
				if _, ok := constructors[name]; !ok {
					t.Errorf("constructor missing from registry: %s", name)
				}
				functions[name] = true
			}
			if strings.HasPrefix(name, "Parse") && strings.HasSuffix(name, "Fields") {
				if _, ok := fieldParsers[name]; !ok {
					t.Errorf("fields parser missing from registry: %s", name)
				}
				functions[name] = true
			}
		}
	}
	if len(functions) != len(constructors)+len(fieldParsers) {
		t.Error("API registry contains stale entries")
	}
}

func constructorArgument(typ reflect.Type, name string) reflect.Value {
	switch typ {
	case reflect.TypeFor[time.Time]():
		return reflect.ValueOf(time.Unix(1700000000, 0).UTC())
	case reflect.TypeFor[time.Duration]():
		return reflect.ValueOf(10 * time.Second)
	case reflect.TypeFor[ie.IEType]():
		return reflect.ValueOf(ie.Cause)
	case reflect.TypeFor[net.IP]():
		addr := "192.0.2.1"
		if strings.Contains(name, "6") {
			addr = "2001:db8::1"
		}
		return reflect.ValueOf(net.ParseIP(addr))
	case reflect.TypeFor[net.HardwareAddr]():
		return reflect.ValueOf(net.HardwareAddr{2, 0, 0, 0, 0, 1})
	case reflect.TypeFor[*ie.IE]():
		if name == "ctag" {
			return reflect.ValueOf(ie.NewCTAG(7, 1, 1, 100))
		}
		if name == "stag" {
			return reflect.ValueOf(ie.NewSTAG(7, 3, 0, 200))
		}
		return reflect.ValueOf(ie.NewCause(1))
	}
	v := reflect.New(typ).Elem()
	switch typ.Kind() {
	case reflect.String:
		value := "abc"
		switch name {
		case "v4", "ipv4", "nodeAddr", "nodeID", "addr", "addrs":
			value = "192.0.2.1"
		case "v6", "ipv6":
			value = "2001:db8::1"
		case "mcc":
			value = "001"
		case "mnc":
			value = "01"
		case "imsi", "imei", "msisdn":
			value = "123456"
		case "ttc":
			value = "ab"
		case "spi":
			value = "abcd"
		case "fl":
			value = "abc"
		}
		v.SetString(value)
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		value := uint64(1)
		if name == "flags" {
			value = 255
		}
		if name == "desc" {
			value = 0xff00
		}
		if name == "addrType" {
			value = 0
		}
		v.SetUint(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(1)
	case reflect.Slice:
		if typ == reflect.TypeFor[[]byte]() {
			v = reflect.MakeSlice(typ, 16, 16)
			return v
		}
		v = reflect.MakeSlice(typ, 1, 1)
		v.Index(0).Set(constructorArgument(typ.Elem(), name))
	}
	return v
}

// Constructor fixtures are registered beside the tests that own each API.
// Registration happens during package initialization, before tests run.
var constructorFixtures = map[uintptr][]any{}

// A flag fixture supplies valid arguments for each supported flag combination.
// Nil arguments exclude combinations forbidden by that constructor's clause.
var constructorFlagFixtures = map[uintptr]func(uint8) []any{}

func registerConstructorFlagFixture(fn any, args func(uint8) []any) bool {
	key := reflect.ValueOf(fn).Pointer()
	if _, exists := constructorFlagFixtures[key]; exists {
		panic("duplicate constructor flag fixture")
	}
	constructorFlagFixtures[key] = args
	return true
}

var constructorComparers = cmp.Options{
	cmp.Comparer(func(a, b net.IP) bool { return a.Equal(b) }),
	cmp.Comparer(func(a, b netip.Addr) bool { return a == b }),
}

func registerConstructorFixture(fn any, args ...any) bool {
	key := reflect.ValueOf(fn).Pointer()
	if _, exists := constructorFixtures[key]; exists {
		panic("duplicate constructor fixture")
	}
	constructorFixtures[key] = args
	return true
}

func callConstructor(entry constructorEntry) []reflect.Value {
	fn := reflect.ValueOf(entry.fn)
	if fixture, ok := constructorFixtures[fn.Pointer()]; ok {
		args := make([]reflect.Value, len(fixture))
		for k, arg := range fixture {
			args[k] = reflect.ValueOf(arg)
		}
		if fn.Type().IsVariadic() {
			return fn.CallSlice(args)
		}
		return fn.Call(args)
	}
	args := make([]reflect.Value, fn.Type().NumIn())
	for k := range args {
		args[k] = constructorArgument(fn.Type().In(k), entry.params[k])
		name := entry.params[k]
		if name == "chid" {
			args[k].SetUint(0)
		}
		if fn.Type().In(k) == reflect.TypeFor[ie.IEType]() {
			switch reflect.ValueOf(entry.fn).Pointer() {
			case reflect.ValueOf(ie.NewGroupedIE).Pointer():
				args[k] = reflect.ValueOf(ie.CreatePDR)
			case reflect.ValueOf(ie.NewVendorSpecificIE).Pointer(), reflect.ValueOf(ie.NewVendorSpecificGroupedIE).Pointer():
				args[k] = reflect.ValueOf(ie.IEType(0x8001))
			}
		}
		if name == "flags" && (fn.Pointer() == reflect.ValueOf(ie.NewPacketRateStatus).Pointer() || fn.Pointer() == reflect.ValueOf(ie.NewPacketRateStatusFields).Pointer()) {
			args[k].SetUint(7)
		}
		if name == "flags" && (fn.Pointer() == reflect.ValueOf(ie.NewFTEID).Pointer() || fn.Pointer() == reflect.ValueOf(ie.NewFTEIDFields).Pointer()) {
			args[k].SetUint(3)
		}
	}
	if fn.Type().IsVariadic() {
		return fn.CallSlice(args)
	}
	return fn.Call(args)
}

func TestConstructorRoundTrips(t *testing.T) {
	for name, entry := range constructors {
		t.Run(name, func(t *testing.T) {
			var result []reflect.Value
			if p := panicOf(func() { result = callConstructor(entry) }); p != nil {
				t.Fatalf("constructor panic: %v", p)
			}
			got := result[0]
			if got.IsNil() {
				t.Fatal("constructor rejected synthetic inputs")
			}
			value := got.Interface().(interface{ Marshal() ([]byte, error) })
			wire, err := value.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if strings.HasSuffix(name, "Fields") {
				parse := reflect.ValueOf(fieldParsers["Parse"+strings.TrimPrefix(name, "New")])
				result := parse.Call([]reflect.Value{reflect.ValueOf(wire)})
				if !result[1].IsNil() {
					t.Fatal(result[1].Interface())
				}
				if !cmp.Equal(result[0].Interface(), value, constructorComparers) {
					t.Error("fields constructor inputs lost on decode:", cmp.Diff(value, result[0].Interface(), constructorComparers))
				}
			} else {
				decoded, err := ie.Parse(wire)
				if err != nil {
					t.Fatal(err)
				}
				again, err := decoded.Marshal()
				if err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(wire, again) {
					t.Error("constructor changed wire bytes")
				}
				// Every available helper is compared across the wire round trip. The
				// fields cases above additionally compare against the constructor fields.
				before, after := reflect.ValueOf(value), reflect.ValueOf(decoded)
				for m := 0; m < before.NumMethod(); m++ {
					if name == "NewVendorSpecificGroupedIE" {
						break
					}
					method := before.Method(m)
					if method.Type().NumIn() != 0 {
						continue
					}
					a, b := method.Call(nil), after.Method(m).Call(nil)
					if len(a) != 2 || a[1].Type() != reflect.TypeFor[error]() {
						continue
					}
					if !a[1].IsNil() || !b[1].IsNil() {
						continue
					}
					if !reflect.DeepEqual(a[0].Interface(), b[0].Interface()) {
						t.Errorf("helper %s changed across encoding", before.Type().Method(m).Name)
					}
				}
			}
		})
	}
	t.Logf("constructors exercised: %d", len(constructors))
}

func TestFieldsInconsistentLengths(t *testing.T) {
	for name, entry := range constructors {
		if !strings.HasSuffix(name, "Fields") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			base := callConstructor(entry)[0]
			if base.IsNil() {
				return
			}
			for n := 0; n < base.Elem().NumField(); n++ {
				field := base.Elem().Type().Field(n)
				if !strings.Contains(field.Name, "Length") && field.Name != "NumberOfMACAddresses" {
					continue
				}
				v := reflect.New(base.Elem().Type())
				v.Elem().Set(base.Elem())
				value := v.Elem().Field(n)
				if value.Kind() < reflect.Uint || value.Kind() > reflect.Uint64 {
					continue
				}
				value.SetUint(0)
				if p := panicOf(func() { _, _ = v.Interface().(interface{ Marshal() ([]byte, error) }).Marshal() }); p != nil {
					t.Errorf("field %s set to zero: panic %v", field.Name, p)
				}
			}
		})
	}
}

func TestFlagDrivenConstructorRoundTrips(t *testing.T) {
	exercised, excluded := 0, 0
	defer func() { t.Logf("flag combinations exercised=%d excluded_by_presence_rule=%d", exercised, excluded) }()
	for name, entry := range constructors {
		if !strings.HasSuffix(name, "Fields") {
			continue
		}
		flag := -1
		for n, param := range entry.params {
			if param == "flags" {
				flag = n
			}
		}
		if flag < 0 {
			continue
		}
		t.Run(name, func(t *testing.T) {
			fn := reflect.ValueOf(entry.fn)
			for flags := 0; flags < 256; flags++ {
				// TS 29.244 V20.1.0 clauses 8.2.13, 8.2.16, 8.2.44,
				// 8.2.50, 8.2.63 and 8.2.86 require a nonzero presence mask.
				switch name {
				case "NewVolumeThresholdFields", "NewSubsequentVolumeThresholdFields", "NewVolumeQuotaFields", "NewSubsequentVolumeQuotaFields":
					if flags&7 == 0 {
						excluded++
						continue
					}
				case "NewVolumeMeasurementFields":
					if flags&63 == 0 {
						excluded++
						continue
					}
				case "NewPacketRateStatusFields":
					if flags&^7 != 0 {
						excluded++
						continue
					}
				case "NewPacketRateFields":
					if flags&7 == 0 {
						excluded++
						continue
					}
				}

				args := make([]reflect.Value, fn.Type().NumIn())
				for n := range args {
					args[n] = constructorArgument(fn.Type().In(n), entry.params[n])
				}
				args[flag].SetUint(uint64(flags))
				if fixture, ok := constructorFlagFixtures[fn.Pointer()]; ok {
					values := fixture(uint8(flags))
					if values == nil {
						excluded++
						continue
					}
					for n, value := range values {
						args[n] = reflect.ValueOf(value)
					}
				}
				exercised++
				// Values whose presence flag is clear are not constructor inputs
				// to the wire encoding. Keep them zero when comparing fields.
				var masks map[string]int
				switch name {
				case "NewCTAGFields", "NewSTAGFields":
					masks = map[string]int{"pcp": 1, "deiFlag": 2, "cvid": 4}
				case "NewNumberOfUEIPAddressesFields":
					masks = map[string]int{"v4": 1, "v6": 2}
				case "NewQoSMonitoringMeasurementFields", "NewPacketDelayThresholdsFields":
					masks = map[string]int{"dl": 1, "ul": 2, "rp": 4}
				}
				for n, param := range entry.params {
					if mask, ok := masks[param]; ok && flags&mask == 0 {
						args[n].SetUint(0)
					}
				}

				if name == "NewPacketRateStatusFields" {
					for n, param := range entry.params {
						if (param == "ul" && flags&1 == 0) || (param == "dl" && flags&2 == 0) || (param == "aul" && (flags&5 != 5)) || (param == "adl" && (flags&6 != 6)) {
							args[n].SetUint(0)
						}
						if param == "vtime" && flags&3 == 0 {
							args[n] = reflect.Zero(args[n].Type())
						}
					}
				}
				if name == "NewFTEIDFields" {
					if flags&4 != 0 {
						args[1].SetUint(0)
						args[2] = reflect.Zero(args[2].Type())
						args[3] = reflect.Zero(args[3].Type())
					} else {
						if flags&1 == 0 {
							args[2] = reflect.Zero(args[2].Type())
						}
						if flags&2 == 0 {
							args[3] = reflect.Zero(args[3].Type())
						}
					}
					if flags&12 != 12 {
						args[4].SetUint(0)
					}
				}
				v := fn.Call(args)[0]
				if v.IsNil() {
					t.Fatalf("flag combination %d rejected", flags)
				}
				fields := v.Interface().(interface{ Marshal() ([]byte, error) })
				wire, err := fields.Marshal()
				if err != nil {
					t.Fatalf("flag combination %d: %v", flags, err)
				}
				parse := reflect.ValueOf(fieldParsers["Parse"+strings.TrimPrefix(name, "New")])
				parsed := parse.Call([]reflect.Value{reflect.ValueOf(wire)})
				if !parsed[1].IsNil() {
					t.Fatalf("flag combination %d: %v", flags, parsed[1].Interface())
				}
				if !cmp.Equal(fields, parsed[0].Interface(), constructorComparers) {
					t.Fatalf("flag combination %d lost fields: %s", flags, cmp.Diff(fields, parsed[0].Interface(), constructorComparers))
				}
			}
		})
	}
}

func TestOversizedFieldInputs(t *testing.T) {
	for name, entry := range constructors {
		if !strings.HasSuffix(name, "Fields") {
			continue
		}
		t.Run(name, func(t *testing.T) {
			fn := reflect.ValueOf(entry.fn)
			for n := 0; n < fn.Type().NumIn(); n++ {
				typ := fn.Type().In(n)
				if typ.Kind() != reflect.String {
					continue
				}
				args := make([]reflect.Value, fn.Type().NumIn())
				for k := range args {
					args[k] = constructorArgument(fn.Type().In(k), entry.params[k])
				}
				if name == "NewUEIPAddressFields" {
					args[0].SetUint(3)
				}
				args[n].SetString(strings.Repeat("1", 65536))
				if p := panicOf(func() {
					v := fn.Call(args)[0]
					if v.IsNil() {
						return
					}
					if _, err := v.Interface().(interface{ Marshal() ([]byte, error) }).Marshal(); err == nil {
						t.Errorf("oversized %s silently encoded", entry.params[n])
					}
				}); p != nil {
					t.Errorf("oversized %s panic %v", entry.params[n], p)
				}
			}
		})
	}
}
