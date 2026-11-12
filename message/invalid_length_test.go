package message_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/message"
)

// TS 29.244 V20.1.0 clause 7.6.3: a complete header must remain
// available for rejecting a request or ending an outstanding transaction.
func lengthFixtures() []*message.Header {
	return []*message.Header{
		{Flags: 0x20, Type: message.MsgTypeAssociationSetupRequest, SequenceNumber: 0x123456},
		{Flags: 0x20, Type: message.MsgTypeAssociationSetupResponse, SequenceNumber: 0xabcdef},
		{Flags: 0x23, Type: message.MsgTypeSessionDeletionRequest, SEID: 0x0123456789abcdef, SequenceNumber: 0x2468ac, MessagePriority: 9},
		{Flags: 0x21, Type: message.MsgTypeSessionDeletionResponse, SEID: 0xfedcba9876543210, SequenceNumber: 0x13579b},
	}
}

func lengthParsers(typ uint8) map[string]func([]byte) error {
	var typed func([]byte) error
	var receiver message.Message
	switch typ {
	case message.MsgTypeAssociationSetupRequest:
		typed = func(b []byte) error { _, err := message.ParseAssociationSetupRequest(b); return err }
		receiver = message.NewAssociationSetupRequest(1)
	case message.MsgTypeAssociationSetupResponse:
		typed = func(b []byte) error { _, err := message.ParseAssociationSetupResponse(b); return err }
		receiver = message.NewAssociationSetupResponse(1)
	case message.MsgTypeSessionDeletionRequest:
		typed = func(b []byte) error { _, err := message.ParseSessionDeletionRequest(b); return err }
		receiver = message.NewSessionDeletionRequest(1, 1)
	case message.MsgTypeSessionDeletionResponse:
		typed = func(b []byte) error { _, err := message.ParseSessionDeletionResponse(b); return err }
		receiver = message.NewSessionDeletionResponse(1, 1)
	}
	return map[string]func([]byte) error{
		"Parse":                   func(b []byte) error { _, err := message.Parse(b); return err },
		"ParseHeader":             func(b []byte) error { _, err := message.ParseHeader(b); return err },
		"Header.UnmarshalBinary":  (&message.Header{}).UnmarshalBinary,
		"ParseGeneric":            func(b []byte) error { _, err := message.ParseGeneric(b); return err },
		"Generic.UnmarshalBinary": (&message.Generic{}).UnmarshalBinary,
		"ParseTyped":              typed,
		"Typed.UnmarshalBinary":   receiver.UnmarshalBinary,
	}
}

func requireInvalidLength(t *testing.T, err error, wire []byte, want *message.Header) *message.Header {
	t.Helper()
	var invalid *message.InvalidLengthError
	if !errors.Is(err, message.ErrInvalidLength) || !errors.As(err, &invalid) {
		t.Fatalf("expected InvalidLengthError wrapping ErrInvalidLength, got %v", err)
	}
	if !reflect.DeepEqual(invalid.Header, want) {
		t.Fatalf("error header = %+v, want %+v (nil payload required)", invalid.Header, want)
	}
	length, received := int(binary.BigEndian.Uint16(wire[2:4])), len(wire)-4
	if invalid.Length != length || invalid.Received != received {
		t.Fatalf("Length/Received = %d/%d, want %d/%d", invalid.Length, invalid.Received, length, received)
	}
	minimum := 4
	if want.Flags&1 != 0 {
		minimum = 12
	}
	text := fmt.Sprintf("invalid PFCP message length: Length field %d but %d octets follow it", length, received)
	if length < minimum {
		text = fmt.Sprintf("invalid PFCP message length: Length field %d does not cover the header", length)
	}
	if invalid.Error() != text {
		t.Fatalf("error text = %q, want %q", invalid.Error(), text)
	}
	return invalid.Header
}

func TestInvalidLengthCompleteHeader(t *testing.T) {
	for _, want := range lengthFixtures() {
		t.Run(message.MessageTypeName(want.Type), func(t *testing.T) {
			withBody := *want
			withBody.Payload = []byte{0, 24, 0, 1, 42}
			valid, err := withBody.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			// Receiving spare bits must not change any retained field.
			valid[0] |= 0x18
			headerLen := 8
			if want.HasSEID() {
				headerLen = 16
			}
			for _, kind := range []string{"too_large", "zero", "below_header", "trailing"} {
				t.Run(kind, func(t *testing.T) {
					wire := bytes.Clone(valid)
					switch kind {
					case "too_large":
						binary.BigEndian.PutUint16(wire[2:4], uint16(len(wire)-3))
					case "zero":
						binary.BigEndian.PutUint16(wire[2:4], 0)
					case "below_header":
						binary.BigEndian.PutUint16(wire[2:4], uint16(headerLen-5))
					case "trailing":
						wire = append(wire, 0)
					}
					for name, parse := range lengthParsers(want.Type) {
						t.Run(name, func(t *testing.T) { requireInvalidLength(t, parse(wire), wire, want) })
					}
					ms, err := message.ParseDatagram(wire)
					errorHeader := requireInvalidLength(t, err, wire, want)
					var frame *message.DatagramError
					if len(ms) != 0 || !errors.As(err, &frame) || frame.Index != 0 || frame.Offset != 0 || frame.Header != errorHeader {
						t.Fatalf("messages=%d; missing shared frame header: %v", len(ms), err)
					}
					// Every failed UnmarshalBinary preserves its populated receiver.
					h := withBody
					before := h
					requireInvalidLength(t, h.UnmarshalBinary(wire), wire, want)
					if !reflect.DeepEqual(h, before) {
						t.Fatal("failed decode changed header")
					}
					for _, m := range []message.Message{message.NewGeneric(want.Type, 123, 456)} {
						before := reflect.ValueOf(m).Elem().Interface()
						requireInvalidLength(t, m.UnmarshalBinary(wire), wire, want)
						if !reflect.DeepEqual(reflect.ValueOf(m).Elem().Interface(), before) {
							t.Fatal("failed decode changed message")
						}
					}
				})
			}
			for n := 0; n < headerLen; n++ {
				for name, parse := range lengthParsers(want.Type) {
					err := parse(valid[:n])
					var invalid *message.InvalidLengthError
					if !errors.Is(err, io.ErrUnexpectedEOF) || errors.As(err, &invalid) {
						t.Fatalf("%s prefix %d: %v", name, n, err)
					}
				}
				ms, err := message.ParseDatagram(valid[:n])
				var frame *message.DatagramError
				if len(ms) != 0 || !errors.As(err, &frame) || frame.Header != nil || !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("prefix %d: messages=%d error=%v", n, len(ms), err)
				}
			}
		})
	}
}

func TestInvalidLengthBundle(t *testing.T) {
	first := wireHeader(0x25, message.MsgTypeSessionDeletionRequest, nil)
	want := lengthFixtures()[2]
	for _, kind := range []string{"too_large", "zero", "below_header", "trailing", "missing_follow_on"} {
		t.Run(kind, func(t *testing.T) {
			withBody := *want
			withBody.Payload = []byte{0, 24, 0, 1, 42}
			last, err := withBody.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			switch kind {
			case "too_large":
				binary.BigEndian.PutUint16(last[2:4], uint16(len(last)-3))
			case "zero":
				binary.BigEndian.PutUint16(last[2:4], 0)
			case "below_header":
				binary.BigEndian.PutUint16(last[2:4], 11)
			case "trailing":
				last = append(last, 0)
			case "missing_follow_on":
				last[0] |= 4
			}
			b := append(bytes.Clone(first), last...)
			ms, err := message.ParseDatagram(b)
			var frame *message.DatagramError
			if !errors.As(err, &frame) {
				t.Fatal("missing frame error", err)
			}
			details := err.(interface{ Unwrap() []error }).Unwrap()
			if len(details) != 1 {
				t.Fatalf("details=%d", len(details))
			}
			if kind == "missing_follow_on" {
				if len(ms) != 2 || frame.Index != 2 || frame.Offset != len(b) || frame.Header != nil || !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("missing follow-on: messages=%d error=%v", len(ms), err)
				}
				return
			}
			errorHeader := requireInvalidLength(t, frame.Err, last, want)
			if len(ms) != 1 || ms[0].Sequence() != 0 || frame.Index != 1 || frame.Offset != len(first) || frame.Header != errorHeader {
				t.Fatalf("invalid last frame: messages=%d error=%v", len(ms), err)
			}
		})
	}
}

func TestInvalidLengthErrorNil(t *testing.T) {
	var err *message.InvalidLengthError
	if err.Error() != message.ErrInvalidLength.Error() || !errors.Is(err, message.ErrInvalidLength) {
		t.Fatal("nil error methods")
	}
}

func parserError[T any](parse func([]byte) (T, error)) func([]byte) error {
	return func(b []byte) error { _, err := parse(b); return err }
}

func TestEveryMessageInvalidLengthPreservesReceiver(t *testing.T) {
	parsers := map[uint8]func([]byte) error{
		1:   parserError(message.ParseHeartbeatRequest),
		2:   parserError(message.ParseHeartbeatResponse),
		3:   parserError(message.ParsePFDManagementRequest),
		4:   parserError(message.ParsePFDManagementResponse),
		5:   parserError(message.ParseAssociationSetupRequest),
		6:   parserError(message.ParseAssociationSetupResponse),
		7:   parserError(message.ParseAssociationUpdateRequest),
		8:   parserError(message.ParseAssociationUpdateResponse),
		9:   parserError(message.ParseAssociationReleaseRequest),
		10:  parserError(message.ParseAssociationReleaseResponse),
		11:  parserError(message.ParseVersionNotSupportedResponse),
		12:  parserError(message.ParseNodeReportRequest),
		13:  parserError(message.ParseNodeReportResponse),
		14:  parserError(message.ParseSessionSetDeletionRequest),
		15:  parserError(message.ParseSessionSetDeletionResponse),
		16:  parserError(message.ParseSessionSetModificationRequest),
		17:  parserError(message.ParseSessionSetModificationResponse),
		50:  parserError(message.ParseSessionEstablishmentRequest),
		51:  parserError(message.ParseSessionEstablishmentResponse),
		52:  parserError(message.ParseSessionModificationRequest),
		53:  parserError(message.ParseSessionModificationResponse),
		54:  parserError(message.ParseSessionDeletionRequest),
		55:  parserError(message.ParseSessionDeletionResponse),
		56:  parserError(message.ParseSessionReportRequest),
		57:  parserError(message.ParseSessionReportResponse),
		200: parserError(message.ParseGeneric),
	}
	for _, m := range allMessages() {
		t.Run(m.MessageTypeName(), func(t *testing.T) {
			h := reflect.ValueOf(m).Elem().FieldByName("Header").Interface().(*message.Header)
			h.SequenceNumber = 0x123456
			beforeHeader := *h
			beforeHeader.Payload = bytes.Clone(h.Payload)
			beforeWire := encodeLengthMessage(t, m)
			want := *h
			want.Payload = nil
			parse, ok := parsers[h.Type]
			if !ok {
				t.Fatal("missing typed parser")
			}
			for _, kind := range []string{"too_large", "zero", "below_header", "trailing"} {
				bad := bytes.Clone(beforeWire)
				switch kind {
				case "too_large":
					binary.BigEndian.PutUint16(bad[2:4], uint16(len(bad)-3))
				case "zero":
					binary.BigEndian.PutUint16(bad[2:4], 0)
				case "below_header":
					binary.BigEndian.PutUint16(bad[2:4], uint16(h.MarshalLen()-5))
				case "trailing":
					bad = append(bad, 0)
				}
				requireInvalidLength(t, parse(bad), bad, &want)
				requireInvalidLength(t, m.UnmarshalBinary(bad), bad, &want)
				if !reflect.DeepEqual(*h, beforeHeader) || !bytes.Equal(encodeLengthMessage(t, m), beforeWire) {
					t.Fatal("failed decode changed populated receiver")
				}
			}
			for n := 0; n < h.MarshalLen(); n++ {
				var invalid *message.InvalidLengthError
				err := parse(beforeWire[:n])
				if !errors.Is(err, io.ErrUnexpectedEOF) || errors.As(err, &invalid) {
					t.Fatalf("prefix %d: %v", n, err)
				}
			}
		})
	}
}

func encodeLengthMessage(t *testing.T, m message.Message) []byte {
	t.Helper()
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); err != nil {
		t.Fatal(err)
	}
	return b
}
