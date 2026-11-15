package message_test

import (
	"errors"
	"io"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

func TestDatagramRetainsChainFailureFrame(t *testing.T) {
	for _, tc := range []struct {
		name   string
		follow bool
		tail   []byte
		want   error
	}{
		{"missing_follow_on", true, nil, io.ErrUnexpectedEOF},
		{"stray_octet", false, []byte{0}, message.ErrInvalidLength},
		{"unannounced_frame", false, wireHeader(0x21, 54, nil), message.ErrInvalidLength},
	} {
		t.Run(tc.name, func(t *testing.T) {
			first := wireHeader(0x25, 54, nil)
			last := wireHeader(0x21, 54, nil)
			if tc.follow {
				last[0] |= 4
			}
			b := append(append(first, last...), tc.tail...)
			ms, err := message.ParseDatagram(b)
			var de *message.DatagramError
			wantMessages, wantIndex, wantOffset := 1, 1, len(first)
			if tc.follow {
				wantMessages, wantIndex, wantOffset = 2, 2, len(b)
			}
			if len(ms) != wantMessages || !errors.Is(err, tc.want) || !errors.As(err, &de) || de.Index != wantIndex || de.Offset != wantOffset {
				t.Fatalf("messages=%d error=%v", len(ms), err)
			}
			if (de.Header == nil) != tc.follow {
				t.Fatal("only the missing follow-on lacks a header")
			}
		})
	}
	// A bad body and a missing follow-on are failures of separate frames.
	b := wireHeader(0x25, 54, []byte{0, 19, 0, 2, 1})
	ms, err := message.ParseDatagram(b)
	var de *message.DatagramError
	var ieErr *ie.DecodeError
	if len(ms) != 0 || !errors.Is(err, io.ErrUnexpectedEOF) || !errors.Is(err, ie.ErrInvalidLength) || !errors.As(err, &de) || !errors.As(err, &ieErr) || ieErr.Type != ie.Cause {
		t.Fatalf("messages=%d error=%v", len(ms), err)
	}
	joined := err.(interface{ Unwrap() []error })
	if len(joined.Unwrap()) != 2 {
		t.Fatal("body and missing follow-on must have separate details")
	}
	for n, detail := range joined.Unwrap() {
		var frame *message.DatagramError
		if !errors.As(detail, &frame) || frame.Index != n || frame.Offset != n*len(b) {
			t.Fatalf("detail %d: %v", n, detail)
		}
		if n == 0 {
			if frame.Header == nil || !errors.Is(frame.Err, ie.ErrInvalidLength) {
				t.Fatal("body failure lacks header or cause")
			}
		} else if frame.Header != nil || !errors.Is(frame.Err, io.ErrUnexpectedEOF) {
			t.Fatal("missing frame lacks incomplete-header context")
		}
	}
}
