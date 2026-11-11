package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestTunnelPasswordCurrent(t *testing.T) {
	for _, v := range []string{"a", "12345", "é", "é\x00"} {
		assertAddedVector(t, ie.NewTunnelPassword(v), []byte(v), "TunnelPassword", false)
	}
	assertAddedNull(t, ie.NewTunnelPassword(""), "TunnelPassword")
	assertAddedVector(t, ie.NewTunnelPassword(string(bytes.Repeat([]byte{65}, 65535))), bytes.Repeat([]byte{65}, 65535), "TunnelPassword", false)
	assertAddedInvalid(t, ie.NewTunnelPassword(string(bytes.Repeat([]byte{65}, 65536))))
	assertAddedInvalid(t, ie.NewTunnelPassword(string([]byte{255})))
}

func TestTunnelPasswordInvalidUTF8Decode(t *testing.T) {
	_, err := ie.New(ie.TunnelPassword, []byte{255}).TunnelPassword()
	var decoded *ie.DecodeError
	if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &decoded) || decoded.Type != ie.TunnelPassword {
		t.Fatalf("invalid UTF8 did not return typed error: %v", err)
	}
}

var _ = registerConstructorFixture(ie.NewTunnelPassword, "example")
