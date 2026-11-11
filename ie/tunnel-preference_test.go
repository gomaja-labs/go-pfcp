package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestTunnelPreferenceCurrent(t *testing.T) {
	for _, v := range []uint32{0, 1, 0x1234, 16777215} {
		assertAddedVector(t, ie.NewTunnelPreference(v), []byte{byte(v >> 16), byte(v >> 8), byte(v)}, "TunnelPreference", true)
		got, err := ie.NewTunnelPreference(v).TunnelPreference()
		if err != nil || got != v {
			t.Fatalf("value %d: %d %v", v, got, err)
		}
	}
	assertAddedInvalid(t, ie.NewTunnelPreference(16777216))
}
