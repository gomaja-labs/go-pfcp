package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestTLContainerCurrent(t *testing.T) {
	for _, v := range [][]byte{{1}, {0, 255, 128}} {
		assertAddedVector(t, ie.NewTLContainer(v), v, "TLContainer", false)
	}
	assertAddedNull(t, ie.NewTLContainer(nil), "TLContainer")
	assertAddedVector(t, ie.NewTLContainer((bytes.Repeat([]byte{65}, 65535))), bytes.Repeat([]byte{65}, 65535), "TLContainer", false)
	assertAddedInvalid(t, ie.NewTLContainer((bytes.Repeat([]byte{65}, 65536))))

}
