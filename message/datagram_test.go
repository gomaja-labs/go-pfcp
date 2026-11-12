// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

func allMessages() []message.Message {
	return []message.Message{
		message.NewHeartbeatRequest(1, nil, nil, nil), message.NewHeartbeatResponse(1, nil, nil),
		message.NewPFDManagementRequest(1, nil), message.NewPFDManagementResponse(1, nil, nil, nil),
		message.NewAssociationSetupRequest(1, nil), message.NewAssociationSetupResponse(1, nil),
		message.NewAssociationUpdateRequest(1, nil), message.NewAssociationUpdateResponse(1, nil),
		message.NewAssociationReleaseRequest(1, nil), message.NewAssociationReleaseResponse(1, nil, nil, nil),
		message.NewVersionNotSupportedResponse(1, nil), message.NewNodeReportRequest(1, nil), message.NewNodeReportResponse(1, nil, nil, nil, nil),
		message.NewSessionSetDeletionRequest(1, nil, nil, nil), message.NewSessionSetDeletionResponse(1, nil, nil, nil, nil),
		message.NewSessionEstablishmentRequest(42, 1, nil), message.NewSessionEstablishmentResponse(42, 1, nil),
		message.NewSessionModificationRequest(42, 1, nil), message.NewSessionModificationResponse(42, 1, nil),
		message.NewSessionDeletionRequest(42, 1, nil), message.NewSessionDeletionResponse(42, 1, nil),
		message.NewSessionReportRequest(42, 1, nil), message.NewSessionReportResponse(42, 1, nil),
		message.NewGeneric(200, 42, 1, nil),
		message.NewSessionSetModificationRequest(1), message.NewSessionSetModificationResponse(1),
	}
}

func TestEveryMessageEncodeSafety(t *testing.T) {
	for _, m := range allMessages() {
		t.Run(m.MessageTypeName(), func(t *testing.T) {
			// Every optional field and list accepts nil and contributes only nonnil IEs.
			v := reflect.ValueOf(m).Elem()
			for n := 1; n < v.NumField(); n++ {
				f := v.Field(n)
				if f.Type() == reflect.TypeOf([]*ie.IE{}) {
					f.Set(reflect.ValueOf([]*ie.IE{nil, ie.NewApplicationID([]byte("x")), nil}))
				}
			}
			size := m.MarshalLen()
			for n := 0; n < size; n++ {
				b := bytes.Repeat([]byte{0xa5}, n)
				if err := m.MarshalTo(b); !errors.Is(err, io.ErrShortBuffer) {
					t.Fatalf("size %d: expected short buffer, got %v", n, err)
				}
				if !bytes.Equal(b, bytes.Repeat([]byte{0xa5}, n)) {
					t.Fatal("short buffer was modified")
				}
			}
			b := make([]byte, size+4)
			for n := size; n < len(b); n++ {
				b[n] = 0xa5
			}
			if err := m.MarshalTo(b); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(b[size:], bytes.Repeat([]byte{0xa5}, 4)) {
				t.Fatal("wrote past message")
			}
			if int(binary.BigEndian.Uint16(b[2:4])) != size-4 {
				t.Fatal("wrong encoded length")
			}
			if _, err := message.Parse(b[:size]); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestConstructorFraming(t *testing.T) {
	for _, m := range allMessages() {
		b := make([]byte, m.MarshalLen())
		if err := m.MarshalTo(b); err != nil {
			t.Fatal(err)
		}
		wantS := m.MessageType() > 17
		if (b[0]&1 != 0) != wantS || b[0]&6 != 0 || m.Version() != 1 {
			t.Fatalf("invalid constructor header for %s", m.MessageTypeName())
		}
	}
	// Session Set Modification is node-related even though its typed codec is pending.
	for _, typ := range []byte{16, 17} {
		h := message.NewHeader(typ, 42, 1, nil)
		if h.HasSEID() {
			t.Fatal("node type has SEID")
		}
	}
}

func TestPriorityValues(t *testing.T) {
	for pri := uint8(0); pri < 16; pri++ {
		h := message.NewHeader(54, 42, 0xffffff, nil)
		h.SetFO(true)
		if err := h.SetMP(pri); err != nil {
			t.Fatal(err)
		}
		b, err := h.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if b[0] != 0x27 || b[15] != pri<<4 {
			t.Fatalf("bad priority %d encoding", pri)
		}
		decoded, err := message.ParseHeader(b)
		if err != nil {
			t.Fatal(err)
		}
		if decoded.MP() != pri || !decoded.HasMP() {
			t.Fatal("priority decode mismatch")
		}
		decoded.ClearMP()
		if decoded.HasMP() || decoded.MP() != 0 || !decoded.HasFO() || !decoded.HasSEID() {
			t.Fatal("clear priority corrupted flags")
		}
	}
	for _, pri := range []byte{16, 255} {
		h := message.NewHeader(54, 42, 1, nil)
		old := *h
		if err := h.SetMP(pri); !errors.Is(err, message.ErrInvalidPriority) {
			t.Fatal("accepted invalid priority")
		}
		if !reflect.DeepEqual(*h, old) {
			t.Fatal("failed setter modified header")
		}
	}
	node := message.NewHeader(1, 0, 1, nil)
	if err := node.SetMP(0); !errors.Is(err, message.ErrInvalidPriority) {
		t.Fatal("node priority accepted")
	}
	if err := node.SetSEID(1); !errors.Is(err, message.ErrInvalidHeader) {
		t.Fatal("node SEID accepted")
	}
}

func TestSequenceAndLengthLimits(t *testing.T) {
	for _, seq := range []uint32{0, 1, 0xffffff, 0x1000000, 0xffffffff} {
		h := message.NewHeader(54, 0, seq, nil)
		b, err := h.Marshal()
		if seq > 0xffffff {
			if !errors.Is(err, message.ErrSequenceOverflow) {
				t.Fatal("accepted sequence overflow")
			}
			continue
		}
		if err != nil {
			t.Fatal(err)
		}
		d, err := message.ParseHeader(b)
		if err != nil || d.Sequence() != seq {
			t.Fatal("sequence not retained")
		}
	}
	for _, typ := range []byte{1, 54} {
		h := message.NewHeader(typ, 0, 1, nil)
		maxPayload := 0xffff - (h.MarshalLen() - 4)
		h.Payload = make([]byte, maxPayload)
		b, err := h.Marshal()
		if err != nil || len(b) != 0xffff+4 {
			t.Fatal("maximum length rejected")
		}
		h.Payload = append(h.Payload, 0)
		if _, err := h.Marshal(); !errors.Is(err, message.ErrInvalidLength) {
			t.Fatal("accepted overflowing length")
		}
	}
	m := message.NewGeneric(200, 0, 1, ie.NewApplicationID([]byte(string(make([]byte, 65535)))))
	if err := m.MarshalTo(nil); !errors.Is(err, message.ErrInvalidLength) {
		t.Fatal("oversized message accepted")
	}
}

func TestDatagramFraming(t *testing.T) {
	for count := 1; count <= 32; count++ {
		ms := make([]message.Message, count)
		for n := range ms {
			ms[n] = message.NewSessionDeletionRequest(uint64(n), uint32(n), ie.NewApplicationID([]byte("x")))
		}
		b, err := message.MarshalDatagram(ms...)
		if err != nil {
			t.Fatal(err)
		}
		ds, err := message.ParseDatagram(b)
		if err != nil || len(ds) != count {
			t.Fatalf("bundle count %d: %v", count, err)
		}
		for n, d := range ds {
			if d.Sequence() != uint32(n) || d.SEID() != uint64(n) {
				t.Fatal("bundle boundary or value mismatch")
			}
		}
		again, err := message.MarshalDatagram(ds...)
		if err != nil || !bytes.Equal(again, b) {
			t.Fatal("bundle bytes changed")
		}
		if count > 1 {
			if _, err := message.Parse(b); !errors.Is(err, message.ErrBundledMessage) {
				t.Fatal("Parse did not identify a bundle")
			}
		}
		for _, m := range ms {
			if reflect.ValueOf(m).Elem().FieldByName("Header").Interface().(*message.Header).HasFO() {
				t.Fatal("encoder changed input FO")
			}
		}
	}
	one := wireHeader(0x21, 54, nil)
	two := append(append([]byte(nil), one...), one...)
	two[0] |= 4
	cases := [][]byte{nil, {}, one[:7], append(append([]byte(nil), one...), 0), append(append([]byte(nil), one...), one...), two[:len(two)-1]}
	missing := append([]byte(nil), one...)
	missing[0] |= 4
	cases = append(cases, missing)
	node := wireHeader(0x20, 1, nil)
	cases = append(cases, append(append([]byte(nil), node...), one...))
	for n, b := range cases {
		if _, err := message.ParseDatagram(b); err == nil {
			t.Errorf("accepted invalid datagram %d", n)
		}
	}
	if _, err := message.MarshalDatagram(); err == nil {
		t.Fatal("empty datagram accepted")
	}
	if _, err := message.MarshalDatagram(nil); err == nil {
		t.Fatal("nil message accepted")
	}
	if _, err := message.MarshalDatagram(message.NewHeartbeatRequest(1, nil, nil), message.NewSessionDeletionRequest(1, 1)); err == nil {
		t.Fatal("bundled node message")
	}
	var nilMessage *message.SessionDeletionRequest
	if _, err := message.MarshalDatagram(nilMessage); err == nil {
		t.Fatal("typed nil message accepted")
	}
	// A standalone node message remains valid through the datagram API.
	if b, err := message.MarshalDatagram(message.NewHeartbeatRequest(1, nil, nil)); err != nil {
		t.Fatal(err)
	} else if _, err := message.ParseDatagram(b); err != nil {
		t.Fatal(err)
	}
}

func TestUnsupportedVersion(t *testing.T) {
	for version := 0; version < 8; version++ {
		if version == 1 {
			continue
		}
		b := wireHeader(byte(version<<5)|1, 54, []byte{0xff}) // Never parse future-version IEs.
		for _, parse := range []func([]byte) error{
			func(b []byte) error { _, err := message.Parse(b); return err },
			func(b []byte) error { _, err := message.ParseDatagram(b); return err },
			func(b []byte) error { _, err := message.ParseSessionDeletionRequest(b); return err },
		} {
			var v *message.UnsupportedVersionError
			if err := parse(b); !errors.As(err, &v) || v.Header.Version() != version || v.Header.Type != 54 {
				t.Fatalf("version %d: missing answerable header: %v", version, err)
			}
			response := message.NewVersionNotSupportedResponse(v.Header.Sequence())
			if _, err := response.Marshal(); err != nil {
				t.Fatal(err)
			}
		}
	}
}

func TestSFlagAndTypeValidation(t *testing.T) {
	for _, typ := range []byte{1, 15, 16, 17, 50, 57} {
		wrong := byte(0x20)
		if typ <= 17 {
			wrong |= 1
		}
		if _, err := message.Parse(wireHeader(wrong, typ, nil)); !errors.Is(err, message.ErrInvalidHeader) {
			t.Fatalf("type %d accepted contradictory S: %v", typ, err)
		}
	}
	for _, flags := range []byte{0x22, 0x24} {
		if _, err := message.ParseDatagram(wireHeader(flags, 1, nil)); err == nil {
			t.Fatal("node header accepted MP or FO")
		}
	}
	if _, err := message.ParseHeartbeatRequest(wireHeader(0x20, 2, nil)); !errors.Is(err, message.ErrInvalidHeader) {
		t.Fatal("typed parser accepted wrong message type")
	}
	for _, flags := range []byte{0x20, 0x21} {
		if _, err := message.Parse(wireHeader(flags, 200, nil)); err != nil {
			t.Fatal("unknown type rejected valid framing", err)
		}
	}
}

func TestSpareBitsIgnored(t *testing.T) {
	for _, flags := range []byte{0x20, 0x21, 0x23} {
		b := wireHeader(flags, 200, nil)
		b[0] |= 0x18
		b[len(b)-1] = 0xff
		h, err := message.ParseHeader(b)
		if err != nil {
			t.Fatal(err)
		}
		out, err := h.Marshal()
		if err != nil {
			t.Fatal(err)
		}
		if out[0]&0x18 != 0 || out[len(out)-1]&0x0f != 0 {
			t.Fatal("spare bits were sent")
		}
		if flags&2 == 0 && out[len(out)-1] != 0 {
			t.Fatal("priority sent with MP clear")
		}
	}
}

func TestMessageReuseAndFieldEdits(t *testing.T) {
	b := wireHeader(0x20, 1, []byte{0, 24, 0, 1, 42, 0, 96, 0, 4, 0, 0, 0, 1})
	m, err := message.ParseHeartbeatRequest(b)
	if err != nil {
		t.Fatal(err)
	}
	m.RecoveryTimeStamp = nil
	out, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 13 {
		t.Fatal("removed field persisted in payload")
	}
	m.RecoveryTimeStamp = ie.New(ie.RecoveryTimeStamp, []byte{0, 0, 0, 1})
	if err := m.UnmarshalBinary(wireHeader(0x20, 1, nil)); err != nil {
		t.Fatal(err)
	}
	if len(m.OtherIEs) != 0 || m.RecoveryTimeStamp != nil {
		t.Fatal("reuse retained old IEs")
	}
}

func TestHeaderPayloadCapacity(t *testing.T) {
	first := wireHeader(0x21, 54, []byte{0, 24, 0, 1, 42})
	second := wireHeader(0x21, 54, nil)
	b := append(append([]byte(nil), first...), second...)
	h, err := message.ParseHeader(b[:len(first)])
	if err != nil {
		t.Fatal(err)
	}
	before := append([]byte(nil), b[len(first):]...)
	h.Payload = append(h.Payload, 0xff)
	if !bytes.Equal(before, b[len(first):]) {
		t.Fatal("appending payload overwrote the following frame")
	}
}

func TestHeaderAllFlagAndTypeCombinations(t *testing.T) {
	for flags := 0; flags < 256; flags++ {
		for typ := 0; typ < 256; typ++ {
			b := wireHeader(byte(flags), byte(typ), nil)
			b[len(b)-1] = 0xff
			h, err := message.ParseHeader(b)
			if err != nil {
				t.Fatal(err)
			}
			if h.Version() != flags>>5 || h.Flags&0x18 != 0 {
				t.Fatal("incorrect version or spare decoding")
			}
			out, err := h.Marshal()
			if err != nil {
				continue
			} // Invalid versions and contradictory flags cannot be sent.
			if out[0] != byte(flags)&^0x18 {
				t.Fatal("encoder changed defined flags")
			}
			if out[len(out)-1]&15 != 0 {
				t.Fatal("encoder retained spare bits")
			}
			if flags&3 != 3 && out[len(out)-1] != 0 {
				t.Fatal("encoder retained absent priority")
			}
			d, err := message.ParseHeader(out)
			if err != nil || !reflect.DeepEqual(h, d) {
				t.Fatal("header round trip changed value")
			}
		}
	}
}

func TestNilAndZeroMessageEncoding(t *testing.T) {
	for _, m := range allMessages() {
		typ := reflect.TypeOf(m)
		for _, v := range []reflect.Value{reflect.Zero(typ), reflect.New(typ.Elem())} {
			x := v.Interface().(message.Message)
			if err := x.MarshalTo(nil); !errors.Is(err, message.ErrNilMessage) {
				t.Fatalf("%s: nil receiver/header: %v", typ, err)
			}
			if x.MarshalLen() != 0 {
				t.Fatal("nil receiver/header had encoded size")
			}
			if !v.IsNil() {
				if err := x.UnmarshalBinary(nil); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatal(err)
				}
			} else if err := x.UnmarshalBinary(nil); !errors.Is(err, message.ErrNilMessage) {
				t.Fatal(err)
			}
		}
	}
	var h *message.Header
	if _, err := h.Marshal(); !errors.Is(err, message.ErrNilMessage) {
		t.Fatal(err)
	}
	if err := h.MarshalTo(nil); !errors.Is(err, message.ErrNilMessage) {
		t.Fatal(err)
	}
	if err := h.UnmarshalBinary(nil); !errors.Is(err, message.ErrNilMessage) {
		t.Fatal(err)
	}
	if err := h.SetMP(0); !errors.Is(err, message.ErrNilMessage) {
		t.Fatal(err)
	}
	if err := h.SetSEID(0); !errors.Is(err, message.ErrNilMessage) {
		t.Fatal(err)
	}
	h.SetFO(true)
	h.ClearMP()
	h.SetSequenceNumber(0)
	if h.HasFO() || h.HasMP() || h.HasSEID() || h.Sequence() != 0 || h.MessageType() != 0 || h.Version() != 0 || h.MP() != 0 || h.MarshalLen() != 0 || h.String() != "<nil>" {
		t.Fatal("incorrect nil header accessors")
	}
}
