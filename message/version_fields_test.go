package message_test

import (
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/message"
)

func TestUnsupportedVersionAvailableFields(t *testing.T) {
	for _, flags := range []byte{0x40, 0x41, 0x43} {
		b := wireHeader(flags, 54, []byte{0xff})
		// Nonzero field values distinguish a missing field from one actually read.
		for n := 4; n < len(b); n++ {
			b[n] = byte(n)
		}
		for size := 1; size <= len(b); size++ {
			_, err := message.ParseHeader(b[:size])
			var v *message.UnsupportedVersionError
			if !errors.As(err, &v) || v.Header == nil {
				t.Fatal(err)
			}
			h := v.Header
			wantType := byte(0)
			if size >= 2 {
				wantType = 54
			}
			var wantSEID uint64
			if flags&1 != 0 && size >= 12 {
				wantSEID = 0x0405060708090a0b
			}
			var wantSeq uint32
			if flags&1 == 0 && size >= 7 {
				wantSeq = 0x040506
			}
			if flags&1 != 0 && size >= 15 {
				wantSeq = 0x0c0d0e
			}
			if h.Flags != flags || h.Type != wantType || h.SEID != wantSEID || h.SequenceNumber != wantSeq || h.MessagePriority != 0 || h.Payload != nil {
				t.Fatalf("flags=%d size=%d: unexpected available fields: %+v", flags, size, h)
			}
		}
	}
}
