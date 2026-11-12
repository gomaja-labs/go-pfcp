package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestPredefinedRulesNameCurrent(t *testing.T) {
	for _, v := range [][]byte{{1}, {0, 255, 128}} {
		assertAddedVector(t, ie.NewPredefinedRulesName(string(v)), v, "PredefinedRulesName", false)
	}
	assertAddedNull(t, ie.NewPredefinedRulesName(""), "PredefinedRulesName")
	assertAddedVector(t, ie.NewPredefinedRulesName(string((bytes.Repeat([]byte{65}, 65535)))), bytes.Repeat([]byte{65}, 65535), "PredefinedRulesName", false)
	assertAddedInvalid(t, ie.NewPredefinedRulesName(string((bytes.Repeat([]byte{65}, 65536)))))

}

var _ = registerConstructorFixture(ie.NewPredefinedRulesName, string([]byte{1}))
