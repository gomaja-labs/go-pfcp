package ie_test

import (
	"bytes"
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestCalledNumberCurrent(t *testing.T) {
	for _, v := range []string{"a", "12345", "é", "é\x00"} {
		assertAddedVector(t, ie.NewCalledNumber(v), []byte(v), "CalledNumber", false)
	}
	assertAddedNull(t, ie.NewCalledNumber(""), "CalledNumber")
	assertAddedVector(t, ie.NewCalledNumber(string(bytes.Repeat([]byte{65}, 65535))), bytes.Repeat([]byte{65}, 65535), "CalledNumber", false)
	assertAddedInvalid(t, ie.NewCalledNumber(string(bytes.Repeat([]byte{65}, 65536))))
	assertAddedInvalid(t, ie.NewCalledNumber(string([]byte{255})))
}

func TestCalledNumberInvalidUTF8Decode(t *testing.T) {
	_, err := ie.New(ie.CalledNumber, []byte{255}).CalledNumber()
	var decoded *ie.DecodeError
	if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &decoded) || decoded.Type != ie.CalledNumber {
		t.Fatalf("invalid UTF8 did not return typed error: %v", err)
	}
}

var _ = registerConstructorFixture(ie.NewCalledNumber, "example")
