package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestEventNotificationURICurrent(t *testing.T) {
	for _, v := range []string{"https://example.com/events?x=%20#frag", "urn:example:events", "https://[2001:db8::1]:443/x", "x://user:pass@[v1.a:b]:/x", "x:", "x:/a", "x:a/b"} {
		assertAddedVector(t, ie.NewEventNotificationURI([]byte(v)), []byte(v), "EventNotificationURI", false)
	}
	for _, v := range []string{"relative", "https://bad host/x", "1x:abc", "x:%xx", "x:a\\b", "x://[invalid]/", "x://user@@example.com", "x://host:notaport", "x:a#b#c", "x:a[bad]"} {
		assertAddedInvalid(t, ie.NewEventNotificationURI([]byte(v)))
	}
	assertAddedNull(t, ie.NewEventNotificationURI(nil), "EventNotificationURI")
}

func TestEventNotificationURILengthBoundary(t *testing.T) {
	b := append([]byte("x:"), bytes.Repeat([]byte{65}, 65533)...)
	assertAddedVector(t, ie.NewEventNotificationURI(b), b, "EventNotificationURI", false)
	assertAddedInvalid(t, ie.NewEventNotificationURI(append(b, 65)))
}
