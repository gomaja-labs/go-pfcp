// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/gomaja-labs/go-pfcp/message"
)

func wireHeader(flags, typ byte, payload []byte) []byte {
	n := 8
	if flags&1 != 0 {
		n = 16
	}
	b := make([]byte, n+len(payload))
	b[0], b[1] = flags, typ
	binary.BigEndian.PutUint16(b[2:4], uint16(len(b)-4))
	copy(b[n:], payload)
	return b
}

func TestHeaderLengthBoundaries(t *testing.T) {
	for _, flags := range []byte{0x20, 0x21} {
		valid := wireHeader(flags, 200, []byte{0, 24, 0, 1, 42})
		t.Run(string(rune(flags)), func(t *testing.T) {
			for n := 0; n < len(valid); n++ {
				if _, err := message.ParseHeader(valid[:n]); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Errorf("truncation at %d of %d returned %v, want unexpected EOF", n, len(valid), err)
				}
			}
			minimum := 4
			if flags&1 != 0 {
				minimum = 12
			}
			for length := 0; length < minimum; length++ {
				b := make([]byte, 4+length)
				b[0] = flags
				binary.BigEndian.PutUint16(b[2:4], uint16(length))
				if _, err := message.ParseHeader(b); !errors.Is(err, message.ErrInvalidLength) {
					t.Fatalf("length below minimum: %v", err)
				}
			}
			short := append([]byte(nil), valid...)
			binary.BigEndian.PutUint16(short[2:4], uint16(len(valid)-10))
			if _, err := message.ParseHeader(short); err == nil {
				t.Error("accepted length below header minimum")
			}
			if _, err := message.ParseHeader(append(valid, 0)); err == nil {
				t.Error("accepted trailing bytes")
			}
		})
	}
}

func TestHeaderFlagIsolation(t *testing.T) {
	h := &message.Header{Flags: 0x21}
	if err := h.SetMP(7); err != nil {
		t.Fatal(err)
	}
	if !h.HasMP() || h.HasFO() || !h.HasSEID() || h.MP() != 7 {
		t.Fatal("priority corrupted flag bits")
	}
}

func TestHeaderDecodedVersion(t *testing.T) {
	for version := 0; version < 8; version++ {
		h, err := message.ParseHeader(wireHeader(byte(version<<5), 200, nil))
		if err != nil {
			t.Fatal(err)
		}
		if h.Version() != version {
			t.Errorf("version = %d, want %d", h.Version(), version)
		}
	}
}

func TestHeaderShortBuffers(t *testing.T) {
	for _, flags := range []byte{0x20, 0x21} {
		h := &message.Header{Flags: flags, Type: 200, Payload: []byte{1, 2, 3}}
		for n := 0; n < h.MarshalLen(); n++ {
			t.Run(string(rune(flags))+"/"+string(rune(n+'a')), func(t *testing.T) {
				defer func() {
					if recover() != nil {
						t.Error("MarshalTo panicked on short buffer")
					}
				}()
				if err := h.MarshalTo(make([]byte, n)); err == nil {
					t.Error("accepted short buffer")
				}
			})
		}
	}
}

func TestMessageShortBuffers(t *testing.T) {
	for typ := byte(1); typ <= 57; typ++ {
		if typ > 15 && typ < 50 {
			continue
		}
		flags := byte(0x20)
		if typ >= 50 {
			flags |= 1
		}
		m, err := message.Parse(wireHeader(flags, typ, nil))
		if err != nil {
			t.Fatal(err)
		}
		for n := 0; n < m.MarshalLen(); n++ {
			t.Run(m.MessageTypeName()+"/"+string(rune(n+'a')), func(t *testing.T) {
				defer func() {
					if recover() != nil {
						t.Error("message MarshalTo panicked on short buffer")
					}
				}()
				if err := m.MarshalTo(make([]byte, n)); err == nil {
					t.Error("accepted short buffer")
				}
			})
		}
	}
}

func TestHeaderComputedLength(t *testing.T) {
	h := &message.Header{Flags: 0x21, Type: 200, Payload: []byte{1, 2, 3}}
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if int(binary.BigEndian.Uint16(b[2:4])) != len(b)-4 {
		t.Error("encoded stale Message Length")
	}
}

func TestHeaderSpareBits(t *testing.T) {
	h := &message.Header{Flags: 0x39, Type: 200, MessagePriority: 0x0f}
	b, err := h.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if b[0]&0x18 != 0 || b[15] != 0 {
		t.Error("sent nonzero spare bits or priority without MP")
	}
}

func TestHeaderSequenceOverflow(t *testing.T) {
	h := &message.Header{Flags: 0x20, Type: 200, SequenceNumber: 1 << 24}
	if _, err := h.Marshal(); err == nil {
		t.Error("silently truncated sequence number")
	}
}

func TestUnknownMessageSilent(t *testing.T) {
	if os.Getenv("PFCP_SILENCE_CHILD") == "1" {
		if _, err := message.Parse(wireHeader(0x21, 200, nil)); err != nil {
			t.Fatal(err)
		}
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=^TestUnknownMessageSilent$")
	cmd.Env = append(os.Environ(), "PFCP_SILENCE_CHILD=1")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(out, []byte("unknown type")) {
		t.Error("Parse logged an unknown message")
	}
}

func TestMessageTruncatedIE(t *testing.T) {
	for _, payload := range [][]byte{{0}, {0, 24}, {0, 24, 0}, {0, 24, 0, 1}} {
		if _, err := message.Parse(wireHeader(0x20, 1, payload)); err == nil {
			t.Errorf("accepted truncated IE of %d bytes", len(payload))
		}
	}
}

func TestHeaderConstructorMasksFlags(t *testing.T) {
	h := message.NewHeader(200, 0, 0, nil)
	if h.Flags != 0x21 {
		t.Errorf("unmasked flags: %d", h.Flags)
	}
}

func TestMessageWireOrder(t *testing.T) {
	// Optional, unknown, then recognized singleton repeated: preserve every occurrence.
	payload := []byte{0, 24, 0, 1, 42, 0, 96, 0, 4, 0, 0, 0, 1, 0, 96, 0, 4, 0, 0, 0, 2}
	b := wireHeader(0x20, 1, payload)
	m, err := message.Parse(b)
	if err != nil {
		t.Fatal(err)
	}
	out := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(out); err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, b) {
		t.Fatal("decode/encode changed IE order or dropped repeated singleton")
	}
}
