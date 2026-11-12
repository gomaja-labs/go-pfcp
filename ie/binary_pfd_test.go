// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 5.2.1A.4 note 2 leaves Custom PFD
// Content interpretation implementation specific; clause 8.2.39 frames it.
func TestBinaryCustomPFDContent(t *testing.T) {
	newIE, ok := any(ie.NewPFDContents).(func(string, string, string, []byte, string, []string, []string, []string) *ie.IE)
	if !ok {
		t.Fatalf("custom PFD constructor must accept []byte, got %T", ie.NewPFDContents)
	}
	newFields, ok := any(ie.NewPFDContentsFields).(func(string, string, string, []byte, string, []string, []string, []string) *ie.PFDContentsFields)
	if !ok {
		t.Fatalf("custom PFD fields constructor must accept []byte, got %T", ie.NewPFDContentsFields)
	}
	payload := []byte{0xff, 0, 0x80, 0xc0, 0xaf}
	fields := newFields("", "", "", payload, "", nil, nil, nil)
	value, ok := any(fields.CustomPFDContent).([]byte)
	if !ok || !bytes.Equal(value, payload) {
		t.Fatalf("custom PFD field must be []byte, got %T", fields.CustomPFDContent)
	}
	item := newIE("", "", "", payload, "", nil, nil, nil)
	got, err := item.PFDContents()
	if err != nil || !bytes.Equal(got.CustomPFDContent, payload) {
		t.Fatalf("custom PFD decode error %v", err)
	}
	wire, err := item.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ie.Parse(wire)
	if err != nil {
		t.Fatal(err)
	}
	replay, err := parsed.Marshal()
	if err != nil || !bytes.Equal(wire, replay) {
		t.Fatalf("round trip error %v", err)
	}
	for n := 0; n < len(item.Payload); n++ {
		if _, err := ie.ParsePFDContentsFields(item.Payload[:n]); !errors.Is(err, io.ErrUnexpectedEOF) {
			t.Fatalf("truncated payload length %d error %v", n, err)
		}
	}
	empty := newIE("", "", "", nil, "", nil, nil, nil)
	emptyFields, err := empty.PFDContents()
	if err != nil || emptyFields.HasCP() {
		t.Fatalf("absent custom content error %v", err)
	}
	invalid := newIE("", "", "", make([]byte, 65536), "", nil, nil, nil)
	if invalid == nil {
		t.Fatal("oversized custom content must return an invalid IE")
	}
	if _, err := invalid.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("oversized custom content error %v", err)
	}
}
