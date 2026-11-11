package utils_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

func TestFQDNBounds(t *testing.T) {
	for _, name := range []string{strings.Repeat("x", 64), "a..b", strings.Repeat("x.", 128)} {
		if got := utils.EncodeFQDN(name); got != nil {
			t.Errorf("accepted invalid name of length %d", len(name))
		}
	}
	for _, name := range []string{"a.b", strings.Repeat("x", 63) + ".b", strings.Repeat("x", 63) + "." + strings.Repeat("y", 63) + "." + strings.Repeat("z", 63) + "." + strings.Repeat("w", 61)} {
		wire := utils.EncodeFQDN(name)
		if wire == nil || utils.DecodeFQDN(wire) != name {
			t.Errorf("valid name of length %d lost", len(name))
		}
	}
	for _, wire := range [][]byte{{3, 'a'}, {1, 'a', 3, 'b'}, append([]byte{64}, bytes.Repeat([]byte{'a'}, 64)...), {0xc0, 0}} {
		if got := utils.DecodeFQDN(wire); got != "" {
			t.Error("accepted malformed name")
		}
	}
}

func TestPLMNBounds(t *testing.T) {
	for _, v := range [][2]string{{"", ""}, {"1", "123"}, {"123", ""}, {"1234", "12"}} {
		if _, err := utils.EncodePLMN(v[0], v[1]); err == nil {
			t.Error("accepted invalid PLMN")
		}
	}
	for n := 0; n < 3; n++ {
		if _, _, err := utils.DecodePLMN(make([]byte, n)); err == nil {
			t.Error("accepted short PLMN")
		}
	}
}
