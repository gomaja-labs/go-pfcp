package ie_test

import (
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

// TestMitigationActions checks TS 29.244 V20.1.0 clause 8.2.283 and its figure.
func TestMitigationActions(t *testing.T) {
	if ie.MitigationActionsDROP != 1 || ie.MitigationActionsTHRO != 2 || ie.MitigationActionsIDDN != 4 {
		t.Fatal("action constants disagree with Figure 8.2.283-1")
	}
	for n := 0; n < 256; n++ {
		v := uint8(n)
		i := ie.NewMitigationActions(v)
		if v == 0 || v > 7 {
			assertNewIEInvalid(t, i, ie.MitigationActions)
			continue
		}
		checkNewIEWire(t, i, []byte{v})
		got, err := i.MitigationActions()
		if err != nil || got != v {
			t.Fatalf("actions %d: got %d, error %v", v, got, err)
		}
	}
	checkNewIEFailures(t, ie.MitigationActions, func(i *ie.IE) (any, error) { return i.MitigationActions() }, 1)
}

func TestMitigationActionsReceive(t *testing.T) {
	for n := 0; n < 256; n++ {
		// Clause 8.1.2 requires ignoring spares and unknown extension octets.
		i := ie.New(ie.MitigationActions, []byte{byte(n), 0xa5, 0x5a})
		got, err := i.MitigationActions()
		want := uint8(n) & 7
		if want == 0 {
			var dec *ie.DecodeError
			if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &dec) || dec.Type != ie.MitigationActions {
				t.Fatalf("empty actions %d: %v", n, err)
			}
			continue
		}
		if err != nil || got != want {
			t.Fatalf("actions %d with spares/extensions: got %d, want %d, error %v", n, got, want, err)
		}
		checkNewIEWire(t, i, i.Payload)
	}
}

func TestMitigationActionsDeferredError(t *testing.T) {
	i := ie.NewMitigationActions(0)
	values := []interface{ Marshal() ([]byte, error) }{
		i,
		ie.NewCreatePDR(i),
		message.NewSessionEstablishmentRequest(1, 2, i),
	}
	for _, value := range values {
		_, err := value.Marshal()
		var enc *ie.EncodeError
		if !errors.Is(err, ie.ErrInvalidValue) || !errors.As(err, &enc) || enc.Type != ie.MitigationActions {
			t.Fatalf("%T lost invalid action context: %v", value, err)
		}
	}
}
