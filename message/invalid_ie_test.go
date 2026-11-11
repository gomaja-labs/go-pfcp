// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message_test

import (
	"errors"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

// TestMessageReportsUnencodableIE checks that an IE whose constructor could not
// encode its argument makes the message fail to encode, instead of the IE
// being left out of the message.
func TestMessageReportsUnencodableIE(t *testing.T) {
	m := message.NewHeartbeatRequest(1, ie.NewRecoveryTimeStamp(time.Unix(-1, 0)))
	if _, err := m.Marshal(); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("Marshal returned %v, want ErrInvalidValue", err)
	}
	b := make([]byte, m.MarshalLen())
	if err := m.MarshalTo(b); !errors.Is(err, ie.ErrInvalidValue) {
		t.Fatalf("MarshalTo returned %v, want ErrInvalidValue", err)
	}
}
