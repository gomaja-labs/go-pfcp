// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/message"
)

func seedMessages(f *testing.F) {
	f.Helper()
	for typ := byte(1); typ <= 17; typ++ {
		f.Add(wireHeader(0x20, typ, nil))
	}
	for typ := byte(50); typ <= 57; typ++ {
		f.Add(wireHeader(0x21, typ, nil))
	}
	for _, flags := range []byte{0x20, 0x21, 0x23, 0x25, 0x27} {
		f.Add(wireHeader(flags, 200, nil))
	}
}

// zeroSpareBits tests only header spares. IE values are opaque here; their
// semantic helpers are exercised separately by the IE package.
func zeroSpareBits(b []byte) bool {
	for len(b) > 0 {
		if len(b) < 8 {
			return false
		}
		n := 4 + int(binary.BigEndian.Uint16(b[2:4]))
		offset := 7
		if b[0]&1 != 0 {
			offset = 15
		}
		if n > len(b) || n <= offset {
			return false
		}
		if b[0]&0x18 != 0 {
			return false
		}
		mask := byte(0xff)
		if b[0]&3 == 3 {
			mask = 0x0f
		}
		if b[offset]&mask != 0 {
			return false
		}
		b = b[n:]
	}
	return true
}

func FuzzParse(f *testing.F) {
	seedMessages(f)
	f.Fuzz(func(t *testing.T, b []byte) {
		m, err := message.Parse(b)
		if err != nil {
			return
		}
		out := make([]byte, m.MarshalLen())
		if err := m.MarshalTo(out); err != nil {
			t.Fatal("successful parse could not encode", err)
		}
		d, err := message.Parse(out)
		if err != nil {
			t.Fatal("encoding could not decode", err)
		}
		if !reflect.DeepEqual(m, d) {
			t.Fatal("decode/encode changed message value")
		}
		if zeroSpareBits(b) && !bytes.Equal(b, out) {
			t.Fatal("decode/encode changed message bytes with zero header spares")
		}
	})
}

func FuzzParseDatagram(f *testing.F) {
	seedMessages(f)
	bundle, err := message.MarshalDatagram(message.NewSessionDeletionRequest(1, 1), message.NewSessionDeletionRequest(2, 2))
	if err != nil {
		f.Fatal(err)
	}
	f.Add(bundle)
	for bad := 0; bad < 3; bad++ {
		var b []byte
		for n := 0; n < 3; n++ {
			body := []byte(nil)
			if n == bad {
				body = []byte{0, 19, 0, 2, 1}
			}
			frame := wireHeader(0x21, 54, body)
			if n < 2 {
				frame[0] |= 4
			}
			b = append(b, frame...)
		}
		f.Add(b)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		ms, err := message.ParseDatagram(b)
		if err != nil {
			var frame *message.DatagramError
			if !errors.As(err, &frame) {
				t.Fatal("missing frame error context", err)
			}
			for _, m := range ms {
				out := make([]byte, m.MarshalLen())
				if e := m.MarshalTo(out); e != nil {
					t.Fatal("retained frame could not encode", e)
				}
				decoded, e := message.Parse(out)
				if e != nil || decoded.MessageType() != m.MessageType() || decoded.Sequence() != m.Sequence() || decoded.SEID() != m.SEID() {
					t.Fatal("retained frame lost values", e)
				}
			}
			return
		}
		out, err := message.MarshalDatagram(ms...)
		if err != nil {
			t.Fatal("successful datagram parse could not encode", err)
		}
		ds, err := message.ParseDatagram(out)
		if err != nil {
			t.Fatal("encoded datagram could not decode", err)
		}
		if !reflect.DeepEqual(ms, ds) {
			t.Fatal("decode/encode changed datagram values")
		}
		if zeroSpareBits(b) && !bytes.Equal(b, out) {
			t.Fatal("decode/encode changed datagram bytes with zero header spares")
		}
	})
}
