// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

// Command hb-server receives Heartbeat Requests and sends matching responses.
//
// Heartbeat exchanging feature is planned be included in the go-pfcp package's
// built-in functions in the future.
package main

import (
	"flag"
	"log"
	"net"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/gomaja-labs/go-pfcp/message"
)

func main() {
	var (
		listen = flag.String("s", "127.0.0.2:8805", "addr/port to listen on")
	)
	flag.Parse()

	laddr, err := net.ResolveUDPAddr("udp", *listen)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", laddr)
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = conn.Close() }()

	started := time.Now() // Recovery time is when this entity started.
	buf := make([]byte, 65535)
	for {
		log.Printf("waiting for messages to come on: %s", laddr)
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}

		msg, err := message.Parse(buf[:n])
		if err != nil {
			log.Printf("ignored undecodable message: %s", err)
			continue
		}

		hbreq, ok := msg.(*message.HeartbeatRequest)
		if !ok {
			log.Printf("got unexpected message: %s, from: %s", msg.MessageTypeName(), addr)
			continue
		}

		if hbreq.RecoveryTimeStamp == nil {
			log.Printf("ignored Heartbeat Request without Recovery Time Stamp, from: %s", addr)
			continue
		}
		ts, err := hbreq.RecoveryTimeStamp.RecoveryTimeStamp()
		if err != nil {
			log.Printf("got Heartbeat Request with invalid TS: %s, from: %s", err, addr)
			continue
		} else {
			log.Printf("got Heartbeat Request with TS: %s, from: %s", ts, addr)
		}

		// TS 29.244 V20.1.0 clause 6.4: a response uses the request's sequence number.
		seq := hbreq.Sequence()
		hbres, err := message.NewHeartbeatResponse(seq, ie.NewRecoveryTimeStamp(started)).Marshal()
		if err != nil {
			log.Fatal(err)
		}

		if _, err := conn.WriteTo(hbres, addr); err != nil {
			log.Fatal(err)
		}
		log.Printf("sent Heartbeat Response to: %s", addr)
	}
}
