package ie_test

import (
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestNodeIDBounds(t *testing.T) {
	for _, b := range [][]byte{{0, 1}, {1, 1}, {2, 1, 'a', 5, 'b'}, {2, 64, 'a'}} {
		if _, err := ie.New(ie.NodeID, b).NodeID(); err == nil {
			t.Error("accepted malformed node ID")
		}
	}
	long := strings.Repeat("x", 64)
	if ie.NewNodeID("", "", long) != nil {
		t.Error("accepted overlong node label")
	}
	if ie.NewNodeIDHeuristic(long) != nil {
		t.Error("heuristic accepted overlong node label")
	}
}
