package message_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/message"
)

func FuzzParse(f *testing.F) {
	f.Fuzz(func(t *testing.T, b []byte) {
		if _, err := message.Parse(b); err != nil {
			t.Skip()
		}
	})
}
