// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// SteeringFunctionality definitions.
const (
	// SteeringFunctionalityATSSSLL selects ATSSSLL.
	SteeringFunctionalityATSSSLL uint8 = 0
	// SteeringFunctionalityMPTCP selects MPTCP.
	SteeringFunctionalityMPTCP uint8 = 1
)

// NewSteeringFunctionality creates a new SteeringFunctionality IE.
func NewSteeringFunctionality(sfunc uint8) *IE {
	return newUint8ValIE(SteeringFunctionality, sfunc&0x0f)
}

// SteeringFunctionality returns SteeringFunctionality in uint8 if the type of IE matches.
func (i *IE) SteeringFunctionality() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SteeringFunctionality:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// TS 29.244 V20.1.0 Table 8.2.124-1.
const (
	// SteeringFunctionalityMPQUICUDP selects MPQUICUDP.
	SteeringFunctionalityMPQUICUDP uint8 = 2
	// SteeringFunctionalityMPQUICIP selects MPQUICIP.
	SteeringFunctionalityMPQUICIP uint8 = 3
	// SteeringFunctionalityMPQUICE selects MPQUICE.
	SteeringFunctionalityMPQUICE uint8 = 4
)
