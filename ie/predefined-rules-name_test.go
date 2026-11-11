package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestPredefinedRulesNameCurrent(t *testing.T) {
	for _, v := range [][]byte{{1}, {0, 255, 128}} {
		assertAddedVector(t, ie.NewPredefinedRulesName(v), v, "PredefinedRulesName", false)
	}
	assertAddedNull(t, ie.NewPredefinedRulesName(nil), "PredefinedRulesName")
	assertAddedVector(t, ie.NewPredefinedRulesName((bytes.Repeat([]byte{65}, 65535))), bytes.Repeat([]byte{65}, 65535), "PredefinedRulesName", false)
	assertAddedInvalid(t, ie.NewPredefinedRulesName((bytes.Repeat([]byte{65}, 65536))))

}

var _ = registerConstructorFixture(ie.NewPredefinedRulesName, []byte{1})
