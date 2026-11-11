package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestGroupIDCurrent(t *testing.T) {
	for _, v := range [][]byte{{1}, {0, 255, 128}} {
		assertAddedVector(t, ie.NewGroupID(v), v, "GroupID", false)
	}
	assertAddedNull(t, ie.NewGroupID(nil), "GroupID")
	assertAddedVector(t, ie.NewGroupID((bytes.Repeat([]byte{65}, 65535))), bytes.Repeat([]byte{65}, 65535), "GroupID", false)
	assertAddedInvalid(t, ie.NewGroupID((bytes.Repeat([]byte{65}, 65536))))

}
