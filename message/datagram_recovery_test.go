package message_test

import (
	"bytes"
	"errors"
	"io"
	"strings"
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
			if len(ms) != 2 || !errors.Is(err, tc.want) || !errors.As(err, &de) || de.Index != 1 || de.Offset != len(first) {
				t.Fatalf("messages=%d error=%v", len(ms), err)
			}
		})
	}
	// A bad body and a chain failure on the same frame are one failed frame.
	b := wireHeader(0x25, 54, []byte{0, 19, 0, 2, 1})
	ms, err := message.ParseDatagram(b)
	var de *message.DatagramError
	var ieErr *ie.DecodeError
	if len(ms) != 0 || !errors.Is(err, io.ErrUnexpectedEOF) || !errors.Is(err, ie.ErrInvalidLength) || !errors.As(err, &de) || !errors.As(err, &ieErr) || ieErr.Type != ie.Cause {
		t.Fatalf("messages=%d error=%v", len(ms), err)
	}
	joined := err.(interface{ Unwrap() []error })
	if len(joined.Unwrap()) != 1 {
		t.Fatal("one failed frame must have one detail")
	}
}

func TestDatagramBoundsFailureDetails(t *testing.T) {
	for _, tc := range []struct {
		name         string
		tailType     byte
		wantMessages int
		summary      string
	}{
		{"valid_tail", 54, 1, "3104 further frames failed"},
		{"all_invalid", 1, 0, "3105 further frames failed"}, // S=1 contradicts this node message type.
	} {
		t.Run(tc.name, func(t *testing.T) {
			bad := wireHeader(0x25, 54, []byte{0, 19, 0, 2, 1})
			// 3,120 failed frames and a final frame fill exactly 64 KiB.
			b := bytes.Repeat(bad, 3120)
			b = append(b, wireHeader(0x21, tc.tailType, nil)...)
			if len(b) != 65536 {
				t.Fatal(len(b))
			}
			ms, err := message.ParseDatagram(b)
			if len(ms) != tc.wantMessages {
				t.Fatalf("messages=%d want=%d", len(ms), tc.wantMessages)
			}
			joined, ok := err.(interface{ Unwrap() []error })
			if !ok {
				t.Fatalf("missing joined error: %v", err)
			}
			failures := joined.Unwrap()
			if len(failures) != 17 {
				t.Fatalf("details=%d want 16 plus summary", len(failures))
			}
			for n, failure := range failures[:16] {
				var de *message.DatagramError
				if !errors.As(failure, &de) || de.Index != n || de.Offset != n*len(bad) {
					t.Fatalf("detail %d: %v", n, failure)
				}
			}
			if !strings.Contains(failures[16].Error(), tc.summary) {
				t.Fatal(failures[16])
			}
			if size := len(err.Error()); size >= 4096 {
				t.Fatalf("error text=%d octets", size)
			}
			allocs := testing.AllocsPerRun(5, func() {
				_, err := message.ParseDatagram(b)
				_ = err.Error()
			})
			t.Logf("64 KiB: allocations=%.0f error_octets=%d", allocs, len(err.Error()))
			// Allow race instrumentation and fmt pool variation within this fixed-input budget.
			if allocs > 35000 {
				t.Fatalf("excess allocations: %.0f", allocs)
			}
		})
	}
}

func TestDatagramFailureDetailLimit(t *testing.T) {
	for _, count := range []int{1, 15, 16, 17} {
		bad := wireHeader(0x25, 54, []byte{0, 19, 0, 2, 1})
		b := bytes.Repeat(bad, count)
		b[len(b)-len(bad)] &^= 4
		ms, err := message.ParseDatagram(b)
		if len(ms) != 0 || err == nil {
			t.Fatalf("count=%d messages=%d error=%v", count, len(ms), err)
		}
		failures := err.(interface{ Unwrap() []error }).Unwrap()
		want := count
		if count > 16 {
			want = 17
		}
		if len(failures) != want {
			t.Fatalf("count=%d details=%d", count, len(failures))
		}
		var de *message.DatagramError
		if count > 16 && errors.As(failures[16], &de) {
			t.Fatal("summary invented frame context")
		}
	}
}
