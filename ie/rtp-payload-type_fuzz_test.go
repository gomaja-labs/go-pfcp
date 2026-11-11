package ie_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// FuzzNewNonGroupedIE exercises the typed value helpers, including every
// conditional-field parser introduced in TS 29.244 clauses 8.2.236–286.
func FuzzNewNonGroupedIE(f *testing.F) {
	types := []ie.IEType{344, 345, 346, 347, 348, 349, 350, 351, 352, 353, 354, 357, 358, 359, 360, 364, 365, 366, 367, 368, 369, 370, 371, 372, 373, 374, 375, 376, 377, 379, 380, 382, 384, 385, 387, 388, 389, 390, 391, 392, 393, 394, 395, 396, 398, 402, 403, 406, 407, 408}
	for n := range types {
		f.Add(uint16(n), []byte{})
		f.Add(uint16(n), bytes.Repeat([]byte{255}, 32))
	}
	// Valid conditional fields are also seeds, so their full parsers are reached.
	f.Add(uint16(4), []byte{15, 0x1d, 0x4c, 0x23, 0x28, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
	f.Add(uint16(5), []byte{7, 14})
	f.Add(uint16(9), []byte{2, 2, 1, 9, 0x12, 0x34, 0x56})
	f.Add(uint16(40), []byte{1, 'k', 1, 's'})
	f.Fuzz(func(t *testing.T, selector uint16, payload []byte) {
		if len(payload) > 65535 {
			return
		}
		typ := types[int(selector)%len(types)]
		i := ie.New(typ, append([]byte(nil), payload...))
		wire, err := i.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := ie.Parse(wire)
		if err != nil {
			t.Fatal(err)
		}
		again, err := parsed.Marshal()
		if err != nil || !bytes.Equal(wire, again) {
			t.Fatal("raw IE changed", err)
		}
		result := reflect.ValueOf(parsed).MethodByName(typ.String()).Call(nil)
		if !result[1].IsNil() {
			return
		}
		fields, ok := result[0].Interface().(newFields)
		if !ok {
			return
		}
		b, err := fields.Marshal()
		if err != nil {
			t.Fatal("decoded Fields cannot be encoded", err)
		}
		normalized := ie.New(typ, b)
		second := reflect.ValueOf(normalized).MethodByName(typ.String()).Call(nil)
		if !second[1].IsNil() || !reflect.DeepEqual(result[0].Interface(), second[0].Interface()) {
			t.Fatal("Fields are not stable")
		}
	})
}
