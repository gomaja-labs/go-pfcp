package message_test

import (
	"bytes"
	"net/netip"
	"reflect"
	"sync"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

func TestMarshalDoesNotMutate(t *testing.T) {
	m := message.NewSessionDeletionRequest(1, 2, ie.NewNodeIDAddress(netip.MustParseAddr("192.0.2.1")))
	before := *m.Header
	out, err := m.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(before, *m.Header) {
		t.Fatal("Marshal changed the header")
	}
	var wg sync.WaitGroup
	for range 8 {
		wg.Go(func() {
			for range 50 {
				b, e := m.Marshal()
				if e != nil || !bytes.Equal(b, out) {
					t.Error("concurrent marshal changed output", e)
				}
			}
		})
	}
	wg.Wait()
}
