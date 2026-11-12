package pfcp_test

import (
	"fmt"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

func Example() {
	started := time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)
	request := message.NewHeartbeatRequest(1, ie.NewRecoveryTimeStamp(started))
	wire, err := request.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := message.ParseHeartbeatRequest(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(decoded.MessageTypeName(), decoded.Sequence())
	// Output: Heartbeat Request 1
}
