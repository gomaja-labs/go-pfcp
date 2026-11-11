package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestMetadataCurrent(t *testing.T) {
	for _, v := range []string{"AA==", "AAE=", "AQID", "/+79"} {
		assertAddedVector(t, ie.NewMetadata([]byte(v)), []byte(v), "Metadata", false)
	}
	for _, v := range []string{"%%%", "A", "AB==", "AA=", "AA==\n"} {
		assertAddedInvalid(t, ie.NewMetadata([]byte(v)))
	}
	assertAddedNull(t, ie.NewMetadata(nil), "Metadata")
}

func TestMetadataLengthBoundary(t *testing.T) {
	b := bytes.Repeat([]byte("AAAA"), 16383)
	assertAddedVector(t, ie.NewMetadata(b), b, "Metadata", false)
	assertAddedInvalid(t, ie.NewMetadata(append(b, []byte("AAAA")...)))
}

var _ = registerConstructorFixture(ie.NewMetadata, []byte("YQ=="))
