package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestCallingNumberCurrent(t *testing.T) {
	for _, v := range []string{"a", "12345", "é", "é\x00"} {
		assertAddedVector(t, ie.NewCallingNumber(v), []byte(v), "CallingNumber", false)
	}
	assertAddedNull(t, ie.NewCallingNumber(""), "CallingNumber")
	assertAddedVector(t, ie.NewCallingNumber(string(bytes.Repeat([]byte{65}, 65535))), bytes.Repeat([]byte{65}, 65535), "CallingNumber", false)
	assertAddedInvalid(t, ie.NewCallingNumber(string(bytes.Repeat([]byte{65}, 65536))))
	assertAddedInvalid(t, ie.NewCallingNumber(string([]byte{255})))
}

func TestCallingNumberInvalidUTF8Decode(t *testing.T) {
	_, err := ie.New(ie.CallingNumber, []byte{255}).CallingNumber()
	var decoded *ie.DecodeError
	if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &decoded) || decoded.Type != ie.CallingNumber {
		t.Fatalf("invalid UTF8 did not return typed error: %v", err)
	}
}

var _ = registerConstructorFixture(ie.NewCallingNumber, "example")
