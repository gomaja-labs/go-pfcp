package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// TS 29.244 V20.1.0 clause 8.2.218: every defined value and every received spare-bit pattern.
func TestConfiguredTimeDomainCurrent(t *testing.T) {
	for v := 0; v <= 1; v++ {
		assertAddedVector(t, ie.NewConfiguredTimeDomain(ie.ConfiguredTimeDomainFlags(v)), []byte{byte(v)}, "ConfiguredTimeDomain", false)
	}
	for v := 0; v < 256; v++ {
		got, err := ie.New(ie.ConfiguredTimeDomain, []byte{byte(v)}).ConfiguredTimeDomain()
		if err != nil || uint8(got) != uint8(v&1) {
			t.Fatalf("received %d: %v, %v", v, got, err)
		}
	}
	for _, v := range []int{2, 255} {
		if v >= 0 && v <= 1 {
			continue
		}
		assertAddedInvalid(t, ie.NewConfiguredTimeDomain(ie.ConfiguredTimeDomainFlags(v)))
	}
}
