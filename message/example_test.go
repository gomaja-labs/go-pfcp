package message_test

import (
	"errors"
	"fmt"
	"net"
	"net/netip"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

func ExampleNewSessionEstablishmentRequest() {
	// A synthetic uplink PDR and a FAR that drops matching packets.
	pdr := (&ie.CreatePDRFields{
		PDRID:      ie.NewPDRID(7),
		Precedence: ie.NewPrecedence(100),
		PDI: (&ie.PDIFields{
			SourceInterface: ie.NewSourceInterface(ie.SrcInterfaceAccess),
			// TS 29.244 V20.1.0 clause 8.2.3: V4 flag, TEID 42, IPv4 192.0.2.1.
			LocalFTEID: ie.New(ie.FTEID, []byte{1, 0, 0, 0, 42, 192, 0, 2, 1}),
		}).IE(),
		FARID: ie.NewFARID(9),
	}).IE()
	// TS 29.244 V20.1.0 Tables 7.5.2.1-1 and 7.5.2.3-1.
	request := message.NewSessionEstablishmentRequest(0, 1,
		ie.NewNodeIDAddress(netip.MustParseAddr("192.0.2.10")),
		ie.NewFSEID(100, net.ParseIP("192.0.2.10"), nil),
		pdr,
		ie.NewCreateFAR(ie.NewFARID(9), ie.NewApplyAction(0x01)), // DROP, clause 8.2.26.
	)
	wire, err := request.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := message.ParseSessionEstablishmentRequest(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	fields, err := decoded.CreatePDR[0].CreatePDR()
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := fields.PDRID.PDRID()
	if err != nil {
		fmt.Println(err)
		return
	}
	pdi, err := fields.PDI.PDI()
	if err != nil {
		fmt.Println(err)
		return
	}
	tunnel, err := pdi.LocalFTEID.FTEID()
	if err != nil {
		fmt.Println(err)
		return
	}
	far, err := fields.FARID.FARID()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("PDR %d: TEID %d at %s, FAR %d\n", id, tunnel.TEID, tunnel.IPv4Address, far)
	// Output: PDR 7: TEID 42 at 192.0.2.1, FAR 9
}

func ExampleNewHeartbeatRequest() {
	started := time.Date(2026, time.November, 1, 0, 0, 0, 0, time.UTC)
	request := message.NewHeartbeatRequest(12, ie.NewRecoveryTimeStamp(started))
	wire, err := request.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	received, err := message.ParseHeartbeatRequest(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	// TS 29.244 V20.1.0 clause 6.4: match the request's sequence number.
	response := message.NewHeartbeatResponse(received.Sequence(), ie.NewRecoveryTimeStamp(started))
	wire, err = response.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	reply, err := message.ParseHeartbeatResponse(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	recovery, err := reply.RecoveryTimeStamp.RecoveryTimeStamp()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(reply.Sequence(), recovery.UTC().Format(time.RFC3339))
	// Output: 12 2026-11-01T00:00:00Z
}

func ExampleMarshalDatagram() {
	// The peer has already advertised bundling support; both requests target it.
	first := message.NewSessionDeletionRequest(100, 1)
	second := message.NewSessionDeletionRequest(200, 2)
	wire, err := message.MarshalDatagram(first, second)
	if err != nil {
		fmt.Println(err)
		return
	}
	received, err := message.ParseDatagram(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, msg := range received {
		fmt.Println(msg.MessageTypeName(), msg.SEID(), msg.Sequence())
	}
	fmt.Println("input FO:", first.HasFO(), second.HasFO())
	// Output:
	// Session Deletion Request 100 1
	// Session Deletion Request 200 2
	// input FO: false false
}

func ExampleUnsupportedVersionError() {
	// Synthetic version-2 node header, with a four-octet body length.
	wire := []byte{0x40, message.MsgTypeHeartbeatRequest, 0, 4, 0, 0, 1, 0}
	_, err := message.Parse(wire)
	var unsupported *message.UnsupportedVersionError
	if !errors.As(err, &unsupported) {
		fmt.Println("expected unsupported version:", err)
		return
	}
	fmt.Println("received version:", unsupported.Header.Version())
	response := message.NewVersionNotSupportedResponse(unsupported.Header.SequenceNumber)
	encoded, err := response.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := message.ParseVersionNotSupportedResponse(encoded)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("supported version:", decoded.Version())
	// Output:
	// received version: 2
	// supported version: 1
}

func ExampleNewSessionModificationRequest_nullLength() {
	// TS 29.244 V20.1.0 clause 5.6.3: stop applying every URR for PDR 7.
	update := ie.NewUpdatePDR(ie.NewPDRID(7), ie.New(ie.URRID, nil))
	request := message.NewSessionModificationRequest(100, 2, update)
	wire, err := request.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := message.ParseSessionModificationRequest(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	fields, err := decoded.UpdatePDR[0].UpdatePDR()
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := fields.PDRID.PDRID()
	if err != nil {
		fmt.Println(err)
		return
	}
	urr := fields.URRID[0]
	fmt.Println("PDR:", id, "remove URRs:", len(urr.Payload) == 0 && len(urr.ChildIEs) == 0)
	_, err = urr.URRID()
	fmt.Println("has URR value:", err == nil)
	// Output:
	// PDR: 7 remove URRs: true
	// has URR value: false
}

func ExampleInvalidLengthError() {
	request := message.NewSessionDeletionRequest(100, 83)
	wire, err := request.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	wire[3]++ // Synthetic Length is one larger than the received datagram.
	_, err = message.ParseDatagram(wire)
	var invalid *message.InvalidLengthError
	if !errors.As(err, &invalid) {
		fmt.Println("expected invalid length:", err)
		return
	}
	h := invalid.Header
	fmt.Println(invalid)
	fmt.Println(message.MessageTypeName(h.Type), h.SequenceNumber, h.IsRequest())
	// TS 29.244 V20.1.0 clause 7.2.2.4.2: the response carries the peer's
	// SEID, looked up from the session that the received SEID identifies, or 0.
	// Here the peer's SEID is the synthetic constant 200.
	response := message.NewSessionDeletionResponse(200, h.SequenceNumber, ie.NewCause(ie.CauseInvalidLength))
	wire, err = response.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := message.ParseSessionDeletionResponse(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	cause, err := decoded.Cause.Cause()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(decoded.MessageTypeName(), decoded.Sequence(), cause)
	// Output:
	// invalid PFCP message length: Length field 13 but 12 octets follow it
	// Session Deletion Request 83 true
	// Session Deletion Response 83 68
}
