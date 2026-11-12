// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// PDNType definitions.
const (
	_ uint8 = 0
	// PDNTypeIPv4 selects IPv4.
	PDNTypeIPv4 uint8 = 1
	// PDNTypeIPv6 selects IPv6.
	PDNTypeIPv6 uint8 = 2
	// PDNTypeIPv4v6 selects IPv4v6.
	PDNTypeIPv4v6 uint8 = 3
	// PDNTypeNonIP selects Non IP.
	PDNTypeNonIP uint8 = 4
	// PDNTypeEthernet selects Ethernet.
	PDNTypeEthernet uint8 = 5
)

// NewPDNType creates a new PDNType IE.
func NewPDNType(typ uint8) *IE {
	return newUint8ValIE(PDNType, typ)
}

// PDNType returns PDNType in uint8 if the type of IE matches.
func (i *IE) PDNType() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != PDNType {
		return 0, &InvalidTypeError{Type: i.Type}
	}

	// TS 29.244 V20.1.0 clauses 8.1.2 and 8.2.79: ignore spare bits.
	return i.valueAsUint8Masked(0x07)
}
