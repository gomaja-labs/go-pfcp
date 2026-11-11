// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewMeasurementInformation creates a new MeasurementInformation IE.
func NewMeasurementInformation(flags uint8) *IE {
	return newUint8ValIE(MeasurementInformation, flags)
}

// MeasurementInformation returns MeasurementInformation in uint8 if the type of IE matches.
func (i *IE) MeasurementInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case MeasurementInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasISTM reports whether an IE has ISTM bit.
func (i *IE) HasISTM() bool {
	if i == nil {
		return false
	}

	if i.Type != MeasurementInformation {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has4thBit(i.Payload[0])
}

// HasRADI reports whether an IE has RADI bit.
func (i *IE) HasRADI() bool {
	if i == nil {
		return false
	}

	if i.Type != MeasurementInformation {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has3rdBit(i.Payload[0])
}

// HasINAM reports whether an IE has INAM bit.
func (i *IE) HasINAM() bool {
	if i == nil {
		return false
	}

	if i.Type != MeasurementInformation {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has2ndBit(i.Payload[0])
}

// HasMBQE reports whether an IE has MBQE bit.
func (i *IE) HasMBQE() bool {
	if i == nil {
		return false
	}

	if i.Type != MeasurementInformation {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has1stBit(i.Payload[0])
}

// HasSSPOC reports whether the SSPOC flag is set.
// TS 29.244 V20.1.0 clause 8.2.68, octet 5 bit 6.
func (i *IE) HasSSPOC() bool { return payloadFlag(i, MeasurementInformation, 0, 0x20) }

// HasASPOC reports whether the ASPOC flag is set.
// TS 29.244 V20.1.0 clause 8.2.68, octet 5 bit 7.
func (i *IE) HasASPOC() bool { return payloadFlag(i, MeasurementInformation, 0, 0x40) }

// HasCIAM reports whether the CIAM flag is set.
// TS 29.244 V20.1.0 clause 8.2.68, octet 5 bit 8.
func (i *IE) HasCIAM() bool { return payloadFlag(i, MeasurementInformation, 0, 0x80) }
