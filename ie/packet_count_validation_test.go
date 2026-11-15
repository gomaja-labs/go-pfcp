package ie_test

import (
	"errors"
	"io"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.30 and Table 8.1.2-1: Variable, exactly 1 or 2 octets.
func TestPacketCountExactLength(t *testing.T) {
	for n := 0; n <= 8; n++ {
		b := make([]byte, n)
		for k := range b {
			b[k] = byte(k + 1)
		}
		item := ie.New(ie.DLBufferingSuggestedPacketCount, b)
		wire, err := item.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		parsed, err := ie.Parse(wire)
		if err != nil {
			t.Fatal(err)
		} // Null length remains valid IE framing.
		count, err := parsed.DLBufferingSuggestedPacketCount()
		switch n {
		case 0:
			if !errors.Is(err, io.ErrUnexpectedEOF) {
				t.Errorf("null: %v", err)
			}
		case 1:
			if err != nil || count != 1 {
				t.Errorf("one: %d, %v", count, err)
			}
		case 2:
			if err != nil || count != 258 {
				t.Errorf("two: %d, %v", count, err)
			}
		default:
			var dec *ie.DecodeError
			if count != 0 || !errors.As(err, &dec) || !errors.Is(err, ie.ErrInvalidLength) || dec.Type != ie.DLBufferingSuggestedPacketCount {
				t.Errorf("length %d: count %d, error %v", n, count, err)
			}
		}
	}
	for _, count := range []uint16{0, 1, 255, 256, 65535} {
		item := ie.NewDLBufferingSuggestedPacketCount(count)
		wantLen := 1
		if count > 255 {
			wantLen = 2
		}
		if len(item.Payload) != wantLen {
			t.Errorf("count %d: length", count)
		}
		got, err := item.DLBufferingSuggestedPacketCount()
		if got != count || err != nil {
			t.Errorf("count %d: %d, %v", count, got, err)
		}
	}
	var nilIE *ie.IE
	if _, err := nilIE.DLBufferingSuggestedPacketCount(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	var typ *ie.InvalidTypeError
	if _, err := ie.New(ie.Cause, []byte{1}).DLBufferingSuggestedPacketCount(); !errors.As(err, &typ) {
		t.Fatal(err)
	}
}
