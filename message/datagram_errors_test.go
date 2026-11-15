// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

// TS 29.244 V20.1.0 clauses 6.5 and 7.6.7: each malformed Request
// needs its own error response, including the seventeenth and later frames.
func TestDatagramReportsEveryFailedRequest(t *testing.T) {
	const count = 17
	var wire []byte
	for n := 0; n < count; n++ {
		frame, err := message.NewSessionDeletionRequest(42, uint32(n+1)).Marshal()
		if err != nil {
			t.Fatal(err)
		}
		frame = append(frame, 0, 56, 0, 1)
		binary.BigEndian.PutUint16(frame[2:4], uint16(len(frame)-4))
		if n+1 < count {
			frame[0] |= 4
		}
		wire = append(wire, frame...)
	}
	messages, err := message.ParseDatagram(wire)
	if len(messages) != 0 || !errors.Is(err, ie.ErrInvalidLength) {
		t.Fatalf("messages=%d error=%v", len(messages), err)
	}
	failures, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("missing frame list: %v", err)
	}
	details := failures.Unwrap()
	if len(details) != count {
		t.Fatalf("details=%d want=%d", len(details), count)
	}
	for n, detail := range details {
		var frame *message.DatagramError
		var invalid *ie.DecodeError
		if !errors.As(detail, &frame) || frame.Header == nil {
			t.Fatalf("frame %d has no received header: %v", n, detail)
		}
		if frame.Index != n || frame.Offset != n*20 || frame.Header.SequenceNumber != uint32(n+1) || frame.Header.SEID != 42 || frame.Header.Type != message.MsgTypeSessionDeletionRequest {
			t.Fatalf("frame %d: index=%d offset=%d header=%v", n, frame.Index, frame.Offset, frame.Header)
		}
		if !errors.Is(detail, ie.ErrInvalidLength) || !errors.As(detail, &invalid) || invalid.Type != ie.PDRID || invalid.TypeUnknown {
			t.Fatalf("frame %d lost its IE cause: %v", n, detail)
		}
	}
}

func TestDatagramNearMaximumFailureDetails(t *testing.T) {
	const count = 3250
	for _, validTail := range []bool{false, true} {
		t.Run(fmt.Sprintf("valid_tail_%t", validTail), func(t *testing.T) {
			var wire []byte
			for n := 0; n < count; n++ {
				frame := wireHeader(0x25, message.MsgTypeSessionDeletionRequest, []byte{0, 56, 0, 1})
				binary.BigEndian.PutUint64(frame[4:12], 42)
				frame[12], frame[13], frame[14] = byte((n+1)>>16), byte((n+1)>>8), byte(n+1)
				if n+1 == count && !validTail {
					frame[0] &^= 4
				}
				wire = append(wire, frame...)
			}
			wantMessages := 0
			if validTail {
				wire = append(wire, wireHeader(0x21, message.MsgTypeSessionDeletionRequest, nil)...)
				wantMessages = 1
			}
			if len(wire) != 65000+wantMessages*16 {
				t.Fatal(len(wire))
			}
			messages, err := message.ParseDatagram(wire)
			var failures message.DatagramErrors
			if len(messages) != wantMessages || !errors.As(err, &failures) || len(failures) != count {
				t.Fatalf("messages=%d details=%d error=%v", len(messages), len(failures), err)
			}
			for n, frame := range failures {
				var invalid *ie.DecodeError
				if frame == nil || frame.Index != n || frame.Offset != n*20 || frame.Header == nil || frame.Header.SEID != 42 || frame.Header.SequenceNumber != uint32(n+1) || frame.Header.Type != message.MsgTypeSessionDeletionRequest || !errors.Is(frame, ie.ErrInvalidLength) || !errors.As(frame, &invalid) || invalid.Type != ie.PDRID {
					t.Fatalf("frame %d: %v", n, frame)
				}
			}
			want := failures[0].Error() + " (and 3249 more failed frames)"
			if err.Error() != want {
				t.Fatalf("error text has %d octets, want %d", len(err.Error()), len(want))
			}
			// Retaining every failure keeps decoding linear in the input; about ten
			// allocations per failed frame, more under the race detector.
			allocs := testing.AllocsPerRun(3, func() {
				_, err := message.ParseDatagram(wire)
				_ = err.Error()
			})
			if allocs > 15*count {
				t.Fatalf("allocations=%.0f for %d failed frames", allocs, count)
			}
			t.Logf("input_octets=%d details=%d error_octets=%d allocations=%.0f", len(wire), len(failures), len(err.Error()), allocs)
		})
	}
}

func TestDatagramFailureCountBound(t *testing.T) {
	for _, count := range []int{0, 1, 2, 17, 4095} {
		// Each 16-octet frame fails in its body parser and announces another.
		// The missing final frame reaches the documented len(input)/16 + 1 bound.
		wire := bytes.Repeat(wireHeader(0x25, message.MsgTypeHeartbeatRequest, nil), count)
		_, err := message.ParseDatagram(wire)
		var failures message.DatagramErrors
		if !errors.As(err, &failures) || len(failures) != len(wire)/16+1 {
			t.Fatalf("input=%d details=%d error=%v", len(wire), len(failures), err)
		}
		last := failures[len(failures)-1]
		if last.Index != count || last.Offset != len(wire) || last.Header != nil || !errors.Is(last, io.ErrUnexpectedEOF) {
			t.Fatalf("missing frame: %v", last)
		}
	}
	// A standalone node failure and an unsupported prefix each fit the extra detail.
	for _, wire := range [][]byte{wireHeader(0x20, 1, []byte{0}), {0x40}} {
		_, err := message.ParseDatagram(wire)
		var failures message.DatagramErrors
		if !errors.As(err, &failures) || len(failures) != 1 || len(failures) > len(wire)/16+1 {
			t.Fatalf("input=%d details=%d error=%v", len(wire), len(failures), err)
		}
	}
}

func TestDatagramErrorsMethods(t *testing.T) {
	var nilFrame *message.DatagramError
	if nilFrame.Error() != "PFCP frame error" || nilFrame.Unwrap() != nil {
		t.Fatal("nil frame methods")
	}
	first := &message.DatagramError{Index: 1, Offset: 16, Err: io.ErrUnexpectedEOF}
	last := &message.DatagramError{Index: 2, Offset: 32, Err: message.ErrInvalidHeader}
	for _, tc := range []struct {
		name      string
		failures  message.DatagramErrors
		text      string
		unwrapped []*message.DatagramError
	}{
		{"nil", nil, "no failed PFCP frames", nil},
		{"empty", message.DatagramErrors{}, "no failed PFCP frames", nil},
		{"one_nil", message.DatagramErrors{nil}, "PFCP frame error", nil},
		{"one", message.DatagramErrors{first}, first.Error(), []*message.DatagramError{first}},
		{"two", message.DatagramErrors{first, last}, first.Error() + " (and 1 more failed frame)", []*message.DatagramError{first, last}},
		{"nil_first", message.DatagramErrors{nil, first}, "PFCP frame error (and 1 more failed frame)", []*message.DatagramError{first}},
		{"nil_middle", message.DatagramErrors{first, nil, last}, first.Error() + " (and 2 more failed frames)", []*message.DatagramError{first, last}},
		{"nil_last", message.DatagramErrors{first, nil}, first.Error() + " (and 1 more failed frame)", []*message.DatagramError{first}},
		{"all_nil", message.DatagramErrors{nil, nil, nil}, "PFCP frame error (and 2 more failed frames)", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.failures.Error(); got != tc.text {
				t.Fatalf("text=%q want=%q", got, tc.text)
			}
			got := tc.failures.Unwrap()
			if len(got) != len(tc.unwrapped) {
				t.Fatalf("unwrapped=%d want=%d", len(got), len(tc.unwrapped))
			}
			for n, want := range tc.unwrapped {
				if got[n] != want { //nolint:errorlint // Unwrap must preserve the original error pointers.
					t.Fatalf("unwrapped element %d lost identity/order", n)
				}
				if !errors.Is(tc.failures, want.Err) {
					t.Fatalf("cause %d not reachable", n)
				}
			}
		})
	}
}

func TestDatagramErrorsReachEveryCause(t *testing.T) {
	// Exercise real parser failures after an earlier IE body failure. Traversal
	// must reach the later frame, rather than only the first error's causes.
	for _, tc := range []struct {
		name string
		tail []byte
		want error
	}{
		{"message_length", append(wireHeader(0x21, 54, nil), 0), message.ErrInvalidLength},
		{"incomplete_header", nil, io.ErrUnexpectedEOF},
		{"unsupported_version", []byte{0x40}, nil},
		{"invalid_header", wireHeader(0x20, 1, nil), message.ErrInvalidHeader},
		{"IE_length", wireHeader(0x21, 54, []byte{0, 56, 0, 1}), ie.ErrInvalidLength},
	} {
		t.Run(tc.name, func(t *testing.T) {
			// A distinct cause in the first frame prevents it from masking a lost tail.
			prefix := wireHeader(0x25, message.MsgTypeHeartbeatRequest, nil)
			if errors.Is(tc.want, message.ErrInvalidHeader) {
				prefix = wireHeader(0x25, 54, []byte{0, 19, 0, 2, 1})
			}
			_, err := message.ParseDatagram(append(prefix, tc.tail...))
			var failures message.DatagramErrors
			if !errors.As(err, &failures) || len(failures) != 2 {
				t.Fatalf("details=%d error=%v", len(failures), err)
			}
			if tc.want != nil && !errors.Is(err, tc.want) {
				t.Fatalf("lost cause: %v", err)
			}
			var frame *message.DatagramError
			if !errors.As(err, &frame) || frame != failures[0] {
				t.Fatal("first frame not reachable")
			}
			switch tc.name {
			case "message_length":
				var invalid *message.InvalidLengthError
				if !errors.As(err, &invalid) || invalid.Header != failures[1].Header {
					t.Fatal("lost invalid length/header identity")
				}
			case "unsupported_version":
				var unsupported *message.UnsupportedVersionError
				if !errors.As(err, &unsupported) || unsupported.Header != failures[1].Header || unsupported.Header.Version() != 2 {
					t.Fatal("lost received version")
				}
			case "IE_length":
				var invalid *ie.DecodeError
				if !errors.As(err, &invalid) || invalid.Type != ie.PDRID {
					t.Fatal("lost offending IE")
				}
			}
		})
	}
}

func TestDatagramSuccessHasNilError(t *testing.T) {
	for _, messages := range [][]message.Message{
		{message.NewHeartbeatRequest(1)},
		{message.NewSessionDeletionRequest(42, 1)},
		{message.NewSessionDeletionRequest(42, 1), message.NewSessionDeletionRequest(43, 2)},
	} {
		wire, err := message.MarshalDatagram(messages...)
		if err != nil {
			t.Fatal(err)
		}
		decoded, err := message.ParseDatagram(wire)
		if err != nil || len(decoded) != len(messages) {
			t.Fatalf("messages=%d error=%v", len(decoded), err)
		}
	}
}
