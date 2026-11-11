// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewNodeReportType creates a new NodeReportType IE.
func NewNodeReportType(flags uint8) *IE {
	return newUint8ValIE(NodeReportType, flags)
}

// NodeReportType returns NodeReportType in uint8 if the type of IE matches.
func (i *IE) NodeReportType() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != NodeReportType {
		return 0, &InvalidTypeError{Type: i.Type}
	}

	return i.ValueAsUint8()
}

// HasUPFR reports whether an IE has UPFR bit.
func (i *IE) HasUPFR() bool {
	if i == nil {
		return false
	}

	if i.Type != NodeReportType {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has1stBit(i.Payload[0])
}

// HasUPRR reports whether the UPRR flag is set.
// TS 29.244 V20.1.0 clause 8.2.69, octet 5 bit 2.
func (i *IE) HasUPRR() bool { return payloadFlag(i, NodeReportType, 0, 0x02) }

// HasCKDR reports whether the CKDR flag is set.
// TS 29.244 V20.1.0 clause 8.2.69, octet 5 bit 3.
func (i *IE) HasCKDR() bool { return payloadFlag(i, NodeReportType, 0, 0x04) }

// HasGPQR reports whether the GPQR flag is set.
// TS 29.244 V20.1.0 clause 8.2.69, octet 5 bit 4.
func (i *IE) HasGPQR() bool { return payloadFlag(i, NodeReportType, 0, 0x08) }

// HasPURR reports whether the PURR flag is set.
// TS 29.244 V20.1.0 clause 8.2.69, octet 5 bit 5.
func (i *IE) HasPURR() bool { return payloadFlag(i, NodeReportType, 0, 0x10) }

// HasVSR reports whether the VSR flag is set.
// TS 29.244 V20.1.0 clause 8.2.69, octet 5 bit 6.
func (i *IE) HasVSR() bool { return payloadFlag(i, NodeReportType, 0, 0x20) }

// HasN6DR reports whether the N6DR flag is set.
// TS 29.244 V20.1.0 clause 8.2.69, octet 5 bit 7.
func (i *IE) HasN6DR() bool { return payloadFlag(i, NodeReportType, 0, 0x40) }
