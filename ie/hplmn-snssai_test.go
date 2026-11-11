package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestHPLMNSNSSAICurrent(t *testing.T) {
	for _, sst := range []uint8{0, 1, 255} {
		for _, sd := range []uint32{0, 0x123456, 0xffffff} {
			v := ie.HPLMNSNSSAIFields{SST: sst, SD: sd}
			b := []byte{sst, byte(sd >> 16), byte(sd >> 8), byte(sd)}
			assertAddedFields(t, ie.NewHPLMNSNSSAI(v), &v, b, "HPLMNSNSSAI", true)
		}
	}

	var nilFields *ie.HPLMNSNSSAIFields
	assertAddedNilFields(t, nilFields)
}

func TestHPLMNSNSSAIInvalidFields(t *testing.T) {
	for _, v := range []ie.HPLMNSNSSAIFields{{SD: 0x1000000}} {
		assertAddedInvalid(t, ie.NewHPLMNSNSSAI(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}

var _ = registerConstructorFixture(ie.NewHPLMNSNSSAI, ie.HPLMNSNSSAIFields{})
