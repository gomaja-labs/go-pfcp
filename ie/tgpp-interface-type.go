// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// TGPPInterfaceType definitons.
const (
	// TGPPInterfaceTypeS1U selects S1 U.
	TGPPInterfaceTypeS1U uint8 = 0
	// TGPPInterfaceTypeS5S8U selects S5 S8 U.
	TGPPInterfaceTypeS5S8U uint8 = 1
	// TGPPInterfaceTypeS4U selects S4 U.
	TGPPInterfaceTypeS4U uint8 = 2
	// TGPPInterfaceTypeS11U selects S11 U.
	TGPPInterfaceTypeS11U uint8 = 3
	// TGPPInterfaceTypeS12U selects S12 U.
	TGPPInterfaceTypeS12U uint8 = 4
	// TGPPInterfaceTypeGnGpU selects Gn Gp U.
	TGPPInterfaceTypeGnGpU uint8 = 5
	// TGPPInterfaceTypeS2aU selects S2a U.
	TGPPInterfaceTypeS2aU uint8 = 6
	// TGPPInterfaceTypeS2bU selects S2b U.
	TGPPInterfaceTypeS2bU uint8 = 7
	// TGPPInterfaceTypeENBDL selects ENBDL.
	TGPPInterfaceTypeENBDL uint8 = 8
	// TGPPInterfaceTypeENBUL selects ENBUL.
	TGPPInterfaceTypeENBUL uint8 = 9
	// TGPPInterfaceTypeSGWUPFDL selects SGWUPFDL.
	TGPPInterfaceTypeSGWUPFDL uint8 = 10
	// TGPPInterfaceTypeN33GPPAccess selects N33 GPP Access.
	TGPPInterfaceTypeN33GPPAccess uint8 = 11
	// TGPPInterfaceTypeN3TrustedNon3GPPAccess selects N3 Trusted Non3 GPP Access.
	TGPPInterfaceTypeN3TrustedNon3GPPAccess uint8 = 12
	// TGPPInterfaceTypeN3UnTrustedNon3GPPAccess selects N3 Un Trusted Non3 GPP Access.
	TGPPInterfaceTypeN3UnTrustedNon3GPPAccess uint8 = 13
	// TGPPInterfaceTypeN3ForDataForwarding selects N3 For Data Forwarding.
	TGPPInterfaceTypeN3ForDataForwarding uint8 = 14
	// TGPPInterfaceTypeN9 selects N9.
	TGPPInterfaceTypeN9 uint8 = 15
	// TGPPInterfaceTypeSGi selects S Gi.
	TGPPInterfaceTypeSGi uint8 = 16
	// TGPPInterfaceTypeN6 selects N6.
	TGPPInterfaceTypeN6 uint8 = 17
	// TGPPInterfaceTypeN19 selects N19.
	TGPPInterfaceTypeN19 uint8 = 18
	// TGPPInterfaceTypeS8U selects S8 U.
	TGPPInterfaceTypeS8U uint8 = 19
	// TGPPInterfaceTypeGpU selects Gp U.
	TGPPInterfaceTypeGpU uint8 = 20
)

// NewTGPPInterfaceType creates a new TGPPInterfaceType IE.
func NewTGPPInterfaceType(intf uint8) *IE {
	return newUint8ValIE(TGPPInterfaceType, intf)
}

// TGPPInterfaceType returns TGPPInterfaceType in uint8 if the type of IE matches.
func (i *IE) TGPPInterfaceType() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case TGPPInterfaceType:
		if len(i.Payload) < 1 {
			return 0, io.ErrUnexpectedEOF
		}

		return i.Payload[0] & 0x3f, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// TS 29.244 V20.1.0 Table 8.2.118-1.
const (
	// TGPPInterfaceTypeN9Roaming selects N9 Roaming.
	TGPPInterfaceTypeN9Roaming uint8 = 21
	// TGPPInterfaceTypeIuU selects Iu U.
	TGPPInterfaceTypeIuU uint8 = 22
	// TGPPInterfaceTypeN9DataForwarding selects N9 Data Forwarding.
	TGPPInterfaceTypeN9DataForwarding uint8 = 23
	// TGPPInterfaceTypeSxaU selects Sxa U.
	TGPPInterfaceTypeSxaU uint8 = 24
	// TGPPInterfaceTypeSxbU selects Sxb U.
	TGPPInterfaceTypeSxbU uint8 = 25
	// TGPPInterfaceTypeSxcU selects Sxc U.
	TGPPInterfaceTypeSxcU uint8 = 26
	// TGPPInterfaceTypeN4U selects N4 U.
	TGPPInterfaceTypeN4U uint8 = 27
	// TGPPInterfaceTypeSGWUPFUL selects SGWUPFUL.
	TGPPInterfaceTypeSGWUPFUL uint8 = 28
	// TGPPInterfaceTypeN6mbNmb9 selects N6mb Nmb9.
	TGPPInterfaceTypeN6mbNmb9 uint8 = 29
	// TGPPInterfaceTypeN3mb selects N3mb.
	TGPPInterfaceTypeN3mb uint8 = 30
	// TGPPInterfaceTypeN19mb selects N19mb.
	TGPPInterfaceTypeN19mb uint8 = 31
)
