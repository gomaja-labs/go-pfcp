package utils_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

func TestDecodePLMNRejectsInvalidDigitsAndLength(t *testing.T) {
	for n := 0; n <= 5; n++ {
		if n == 3 {
			continue
		}
		if _, _, e := utils.DecodePLMN(make([]byte, n)); e == nil {
			t.Fatalf("accepted %d octets", n)
		}
	}
	valid := []byte{0x21, 0xf3, 0x54}
	for pos := 0; pos < 6; pos++ {
		for digit := byte(10); digit < 16; digit++ {
			b := append([]byte(nil), valid...)
			if pos%2 == 0 {
				b[pos/2] = b[pos/2]&0xf0 | digit
			} else {
				b[pos/2] = b[pos/2]&0x0f | digit<<4
			}
			_, _, err := utils.DecodePLMN(b)
			if pos == 3 && digit == 15 {
				if err != nil {
					t.Fatal("valid MNC filler rejected")
				}
				continue
			}
			if err == nil {
				t.Fatalf("accepted nibble %d at %d", digit, pos)
			}
		}
	}
	for _, mcc := range []string{"000", "123", "999"} {
		for _, mnc := range []string{"00", "12", "99", "000", "456", "999"} {
			b, e := utils.EncodePLMN(mcc, mnc)
			if e != nil {
				t.Fatal(e)
			}
			c, n, e := utils.DecodePLMN(b)
			if e != nil || c != mcc || n != mnc {
				t.Fatal("PLMN digits not preserved", e)
			}
		}
	}
}

func TestFixedWidthConversionsRejectWrongLength(t *testing.T) {
	for _, tc := range []struct {
		fn    any
		width int
	}{{utils.Uint24To32, 3}, {utils.Uint40To64, 5}} {
		fn := reflect.ValueOf(tc.fn)
		for n := 0; n <= 7; n++ {
			b := make([]byte, n)
			for j := range b {
				b[j] = 255
			}
			got := fn.Call([]reflect.Value{reflect.ValueOf(b)})
			if len(got) != 2 {
				t.Fatalf("%d-bit conversion has no error result", tc.width*8)
			}
			if n == tc.width {
				if !got[1].IsNil() || got[0].Uint() != uint64(1)<<(tc.width*8)-1 {
					t.Fatal("valid conversion failed")
				}
			} else if got[1].IsNil() {
				t.Fatalf("accepted %d octets for width %d", n, tc.width)
			}
		}
	}
}
