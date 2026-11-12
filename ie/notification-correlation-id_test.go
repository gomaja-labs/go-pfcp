package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestNotificationCorrelationIDCurrent(t *testing.T) {
	for _, v := range [][]byte{{1}, {0, 255, 128}} {
		assertAddedVector(t, ie.NewNotificationCorrelationID(string(v)), v, "NotificationCorrelationID", false)
	}
	assertAddedNull(t, ie.NewNotificationCorrelationID(""), "NotificationCorrelationID")
	assertAddedVector(t, ie.NewNotificationCorrelationID(string((bytes.Repeat([]byte{65}, 65535)))), bytes.Repeat([]byte{65}, 65535), "NotificationCorrelationID", false)
	assertAddedInvalid(t, ie.NewNotificationCorrelationID(string((bytes.Repeat([]byte{65}, 65536)))))

}

var _ = registerConstructorFixture(ie.NewNotificationCorrelationID, string([]byte{1}))
