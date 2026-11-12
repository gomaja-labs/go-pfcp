// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
	"net/netip"
	"strings"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

// Node ID kinds from TS 29.244 V20.1.0 Table 8.2.38-2.
const (
	NodeIDIPv4Address uint8 = 0
	NodeIDIPv6Address uint8 = 1
	NodeIDFQDN        uint8 = 2
)

// NodeIDValue identifies the received kind without interpreting a domain name
// as an address. Address is set for IPv4/IPv6; FQDN is set for a domain name.
type NodeIDValue struct {
	Type    uint8
	Address netip.Addr
	FQDN    string
}

// NewNodeIDAddress creates a Node ID from an IPv4 or IPv6 address.
// IPv4-mapped IPv6 addresses retain their IPv6 kind. Zoned addresses cannot be
// represented on the wire and produce an IE that reports ErrInvalidValue.
func NewNodeIDAddress(addr netip.Addr) *IE {
	// TS 29.244 V20.1.0 clause 8.2.38: four or sixteen address octets.
	if !addr.IsValid() || addr.Zone() != "" {
		return newInvalidIE(NodeID, fmt.Errorf("%w: invalid or zoned Node ID address", ErrInvalidValue))
	}
	kind := NodeIDIPv6Address
	if addr.Is4() {
		kind = NodeIDIPv4Address
	}
	return New(NodeID, append([]byte{kind}, addr.AsSlice()...))
}

// NewNodeIDFQDN creates a Node ID from a domain name, including names that look
// like addresses. The name must omit the trailing dot; its wire encoding uses
// DNS labels without a trailing root octet.
func NewNodeIDFQDN(fqdn string) *IE {
	// TS 29.244 V20.1.0 clause 8.2.38; RFC 1035 section 3.1.
	if fqdn == "" || strings.HasSuffix(fqdn, ".") {
		return newInvalidIE(NodeID, fmt.Errorf("%w: Node ID name must be nonempty and omit the trailing dot", ErrInvalidValue))
	}
	encoded := utils.EncodeFQDN(fqdn)
	if encoded == nil {
		return newInvalidIE(NodeID, invalidFQDN(fqdn))
	}
	return New(NodeID, append([]byte{NodeIDFQDN}, encoded...))
}

// NodeID returns the Node ID kind and its address or domain name.
func (i *IE) NodeID() (NodeIDValue, error) {
	if i == nil {
		return NodeIDValue{}, ErrIENotFound
	}
	if i.Type != NodeID {
		return NodeIDValue{}, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 2 {
		return NodeIDValue{}, io.ErrUnexpectedEOF
	}
	// TS 29.244 V20.1.0 clauses 8.1.2/8.2.38: ignore spare bits and
	// unknown extension octets after a known address.
	kind := i.Payload[0] & 0x0f
	switch kind {
	case NodeIDIPv4Address:
		if len(i.Payload) < 5 {
			return NodeIDValue{}, io.ErrUnexpectedEOF
		}
		return NodeIDValue{Type: kind, Address: netip.AddrFrom4([4]byte(i.Payload[1:5]))}, nil
	case NodeIDIPv6Address:
		if len(i.Payload) < 17 {
			return NodeIDValue{}, io.ErrUnexpectedEOF
		}
		return NodeIDValue{Type: kind, Address: netip.AddrFrom16([16]byte(i.Payload[1:17]))}, nil
	case NodeIDFQDN:
		name, err := decodeFQDNValue(i.Payload[1:])
		if err != nil {
			return NodeIDValue{}, err
		}
		return NodeIDValue{Type: kind, FQDN: name}, nil
	default:
		return NodeIDValue{}, &InvalidNodeIDError{ID: kind}
	}
}
