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

// TS 29.244 V20.1.0 clauses 8.2.144 and 8.2.182 specify opaque
// TS 24.539 management messages, including arbitrary non-text octets.
func TestBinaryContainerAPIs(t *testing.T) {
	for _, tc := range []struct {
		name        string
		typ         ie.IEType
		constructor any
		getter      any
	}{
		{"PortManagementInformationContainer", ie.PortManagementInformationContainer, ie.NewPortManagementInformationContainer, (*ie.IE).PortManagementInformationContainer},
		{"UserPlaneNodeManagementInformationContainer", ie.UserPlaneNodeManagementInformationContainer, ie.NewUserPlaneNodeManagementInformationContainer, (*ie.IE).UserPlaneNodeManagementInformationContainer},
	} {
		t.Run(tc.name, func(t *testing.T) {
			newIE, ok := tc.constructor.(func([]byte) *ie.IE)
			if !ok {
				t.Fatalf("constructor must accept []byte, got %T", tc.constructor)
			}
			get, ok := tc.getter.(func(*ie.IE) ([]byte, error))
			if !ok {
				t.Fatalf("getter must return []byte, got %T", tc.getter)
			}
			for _, payload := range [][]byte{{0}, {0xff, 0, 0x80, 0xc0, 0xaf}, {1, 2, 3}} {
				item := newIE(payload)
				got, err := get(item)
				if err != nil || !bytes.Equal(got, payload) {
					t.Fatalf("value %v, error %v", got, err)
				}
				wire, err := item.Marshal()
				if err != nil {
					t.Fatal(err)
				}
				parsed, err := ie.Parse(wire)
				if err != nil {
					t.Fatal(err)
				}
				got, err = get(parsed)
				if err != nil || !bytes.Equal(got, payload) {
					t.Fatalf("parsed value %v, error %v", got, err)
				}
				replay, err := parsed.Marshal()
				if err != nil || !bytes.Equal(replay, wire) {
					t.Fatalf("round trip error %v", err)
				}
				for n := 0; n < len(wire); n++ {
					if _, err := ie.Parse(wire[:n]); err == nil {
						t.Fatalf("accepted truncated IE length %d", n)
					}
				}
			}
			for _, payload := range [][]byte{nil, {}} {
				item := newIE(payload)
				wire, err := item.Marshal()
				if err != nil || len(wire) != 4 {
					t.Fatalf("null encoding length %d, error %v", len(wire), err)
				}
				if _, err := get(item); !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Fatalf("null value error %v", err)
				}
			}
			if _, err := get(nil); !errors.Is(err, ie.ErrIENotFound) {
				t.Fatalf("nil error %v", err)
			}
			var wrong *ie.InvalidTypeError
			if _, err := get(ie.New(ie.Cause, []byte{1})); !errors.As(err, &wrong) {
				t.Fatalf("wrong type error %v", err)
			}
		})
	}
}
