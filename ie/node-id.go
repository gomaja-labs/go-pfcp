// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
	"net"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

// NodeID definitions.
const (
	NodeIDIPv4Address uint8 = 0
	NodeIDIPv6Address uint8 = 1
	NodeIDFQDN        uint8 = 2
)

// NewNodeID creates a new NodeID IE.
//
// Only one of the parameters should have a non-empty value(!="").
// If multiple parameters are given, the first one is used(ipv4 > ipv6 > fqdn)
func NewNodeID(ipv4, ipv6, fqdn string) *IE {
	var p []byte

	switch {
	case ipv4 != "":
		v4 := net.ParseIP(ipv4).To4()
		if v4 == nil {
			return newInvalidIE(NodeID, fmt.Errorf("%w: %q is not an IPv4 address", ErrInvalidValue, ipv4))
		}
		p = append([]byte{NodeIDIPv4Address}, v4...)
	case ipv6 != "":
		v6 := net.ParseIP(ipv6)
		if v6 == nil {
			return newInvalidIE(NodeID, fmt.Errorf("%w: %q is not an IPv6 address", ErrInvalidValue, ipv6))
		}
		p = append([]byte{NodeIDIPv6Address}, v6.To16()...)
	case fqdn != "":
		encoded := utils.EncodeFQDN(fqdn)
		if encoded == nil {
			return newInvalidIE(NodeID, invalidFQDN(fqdn))
		}
		p = append([]byte{NodeIDFQDN}, encoded...)
	default:
		return newInvalidIE(NodeID, fmt.Errorf("%w: no Node ID value given", ErrInvalidValue))
	}

	return New(NodeID, p)
}

// NewNodeIDHeuristic creates a new NodeID IE.
func NewNodeIDHeuristic(nodeID string) *IE {
	if nodeID == "" {
		return newInvalidIE(NodeID, fmt.Errorf("%w: no Node ID value given", ErrInvalidValue))
	}

	var p []byte
	ip := net.ParseIP(nodeID)
	if ip != nil {
		if v4 := ip.To4(); v4 != nil { // IPv4
			p = make([]byte, 5)
			p[0] = NodeIDIPv4Address
			copy(p[1:], v4)
		} else { // IPv6
			p = make([]byte, 17)
			p[0] = NodeIDIPv6Address
			copy(p[1:], ip.To16())
		}
	} else { // FQDN
		encoded := utils.EncodeFQDN(nodeID)
		if encoded == nil {
			return newInvalidIE(NodeID, invalidFQDN(nodeID))
		}
		p = append([]byte{NodeIDFQDN}, encoded...)
	}

	return New(NodeID, p)
}

// NodeID returns NodeID in string if the type of IE matches.
func (i *IE) NodeID() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	if i.Type != NodeID {
		return "", &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return "", io.ErrUnexpectedEOF
	}

	// TS 29.244 V20.1.0 clauses 8.1.2 and 8.2.38: ignore the spare nibble.
	switch i.Payload[0] & 0x0f {
	case NodeIDIPv4Address:
		if len(i.Payload) < 5 {
			return "", io.ErrUnexpectedEOF
		}
		return net.IP(i.Payload[1:5:5]).String(), nil
	case NodeIDIPv6Address:
		if len(i.Payload) < 17 {
			return "", io.ErrUnexpectedEOF
		}
		return net.IP(i.Payload[1:17:17]).String(), nil
	case NodeIDFQDN:
		// TS 29.244 V20.1.0 clause 8.2.38; RFC 1035 section 3.1.
		return decodeFQDNValue(i.Payload[1:])

	default:
		return "", &InvalidNodeIDError{ID: i.Payload[0]}
	}
}
