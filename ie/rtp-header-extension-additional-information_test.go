package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestRTPHeaderExtensionAdditionalInformation(t *testing.T) {
	for flags := uint8(0); flags < 8; flags++ {
		for format := uint8(0); format < 4; format++ {
			for values := 0; values < 4; values++ {
				pssa, npds := values&1 != 0, values&2 != 0
				f := ie.NewRTPHeaderExtensionAdditionalInformationFields(flags, format, pssa, npds)
				i := ie.NewRTPHeaderExtensionAdditionalInformation(flags, format, pssa, npds)
				valid := (flags&1 == 0 && format == 0 || flags&1 != 0 && (format == 1 || format == 2)) && (flags&2 != 0 || !pssa) && (flags&4 != 0 || !npds)
				if !valid {
					assertNewIEInvalid(t, i, ie.RTPHeaderExtensionAdditionalInformation)
					continue
				}
				second := format
				if pssa {
					second |= 4
				}
				if npds {
					second |= 8
				}
				checkNewIEWire(t, i, []byte{flags, second})
				checkNewFields(t, f, func(b []byte) (newFields, error) { return ie.ParseRTPHeaderExtensionAdditionalInformationFields(b) })
				got, err := i.RTPHeaderExtensionAdditionalInformation()
				if err != nil || !reflect.DeepEqual(f, got) {
					t.Fatal("fields", err)
				}
				// Inactive bit fields, as well as spare bits, are not evaluated on receipt.
				extra := []byte{flags | 0xf8, second | 0xf0, 0xa5}
				if flags&1 == 0 {
					extra[1] |= 3
				}
				if flags&2 == 0 {
					extra[1] |= 4
				}
				if flags&4 == 0 {
					extra[1] |= 8
				}
				got, err = ie.ParseRTPHeaderExtensionAdditionalInformationFields(extra)
				if err != nil || !reflect.DeepEqual(f, got) {
					t.Fatal("ignored bits", err)
				}
			}
		}
	}
	for flags := 8; flags <= 255; flags++ {
		assertNewIEInvalid(t, ie.NewRTPHeaderExtensionAdditionalInformation(uint8(flags), 0, false, false), ie.RTPHeaderExtensionAdditionalInformation)
	}
	checkNewFieldsNil(t, (*ie.RTPHeaderExtensionAdditionalInformationFields)(nil))
	checkNewIEFailures(t, ie.RTPHeaderExtensionAdditionalInformation, func(i *ie.IE) (any, error) { return i.RTPHeaderExtensionAdditionalInformation() }, 2)
}
