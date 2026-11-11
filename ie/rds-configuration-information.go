// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewRDSConfigurationInformation creates a new RDSConfigurationInformation IE.
func NewRDSConfigurationInformation(rds uint8) *IE {
	return newUint8ValIE(RDSConfigurationInformation, rds&0x01)
}

// RDSConfigurationInformation returns RDSConfigurationInformation in uint8 if the type of IE matches.
func (i *IE) RDSConfigurationInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case RDSConfigurationInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasRDS reports whether an IE has RDS bit.
func (i *IE) HasRDS() bool {
	// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 8.
	if i != nil && i.Type == UPFunctionFeatures {
		return payloadFlag(i, UPFunctionFeatures, 4, 0x80)
	}

	if i == nil {
		return false
	}

	switch i.Type {
	case RDSConfigurationInformation:
		if len(i.Payload) < 1 {
			return false
		}

		return has1stBit(i.Payload[0])
	default:
		return false
	}
}
