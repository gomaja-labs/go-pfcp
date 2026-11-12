// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewRequestedQoSMonitoring creates a new RequestedQoSMonitoring IE.
func NewRequestedQoSMonitoring(rp, ul, dl int) *IE {
	if rp < 0 || rp > 1 || ul < 0 || ul > 1 || dl < 0 || dl > 1 {
		return newInvalidIE(RequestedQoSMonitoring, ErrInvalidValue)
	}
	return NewRequestedQoSMonitoringFlags(uint16((rp << 2) | (ul << 1) | dl))
}

// RequestedQoSMonitoring returns RequestedQoSMonitoring in uint16 if the type of IE matches.
func (i *IE) RequestedQoSMonitoring() (uint16, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case RequestedQoSMonitoring:
		if len(i.Payload) < 2 {
			return 0, io.ErrUnexpectedEOF
		}
		return uint16(i.Payload[0]) | uint16(i.Payload[1]&3)<<8, nil
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}

}

// HasRP reports whether an IE has RP bit.
func (i *IE) HasRP() bool {
	if i == nil {
		return false
	}

	v, err := i.RequestedQoSMonitoring()
	if err != nil {
		return false
	}

	return v&4 != 0
}

// HasUL reports whether an IE has UL bit.
func (i *IE) HasUL() bool {
	if i == nil {
		return false
	}

	v, err := i.RequestedQoSMonitoring()
	if err != nil {
		return false
	}

	return v&2 != 0
}

// HasDL reports whether an IE has DL bit.
func (i *IE) HasDL() bool {
	if i == nil {
		return false
	}

	v, err := i.RequestedQoSMonitoring()
	if err != nil {
		return false
	}

	return v&1 != 0
}

// NewRequestedQoSMonitoringFlags encodes both request octets. Bits 0..7 are
// octet 5 and bits 8..9 are DLAB/ULAB in octet 6.
// TS 29.244 V20.1.0 clause 8.2.167.
func NewRequestedQoSMonitoringFlags(flags uint16) *IE {
	if flags == 0 || flags&^uint16(0x3ff) != 0 {
		return newInvalidIE(RequestedQoSMonitoring, ErrInvalidValue)
	}
	return New(RequestedQoSMonitoring, []byte{uint8(flags), uint8(flags >> 8)})
}

// HasGTPUPM reports whether GTPUPM is requested.
func (i *IE) HasGTPUPM() bool { v, e := i.RequestedQoSMonitoring(); return e == nil && v&(1<<3) != 0 }

// HasDLCI reports whether DLCI is requested.
func (i *IE) HasDLCI() bool { v, e := i.RequestedQoSMonitoring(); return e == nil && v&(1<<4) != 0 }

// HasULCI reports whether ULCI is requested.
func (i *IE) HasULCI() bool { v, e := i.RequestedQoSMonitoring(); return e == nil && v&(1<<5) != 0 }

// HasULDR reports whether ULDR is requested.
func (i *IE) HasULDR() bool { v, e := i.RequestedQoSMonitoring(); return e == nil && v&(1<<7) != 0 }

// HasDLAB reports whether DLAB is requested.
func (i *IE) HasDLAB() bool { v, e := i.RequestedQoSMonitoring(); return e == nil && v&(1<<8) != 0 }

// HasULAB reports whether ULAB is requested.
func (i *IE) HasULAB() bool { v, e := i.RequestedQoSMonitoring(); return e == nil && v&(1<<9) != 0 }

// RequestsDLDR reports whether monitoring of downlink data rate is requested.
func (i *IE) RequestsDLDR() bool { v, e := i.RequestedQoSMonitoring(); return e == nil && v&0x40 != 0 }

// Requested QoS monitoring flags, TS 29.244 V20.1.0 clause 8.2.167.
const (
	RequestedQoSMonitoringDLPD   uint16 = 1 << 0
	RequestedQoSMonitoringULPD   uint16 = 1 << 1
	RequestedQoSMonitoringRPPD   uint16 = 1 << 2
	RequestedQoSMonitoringGTPUPM uint16 = 1 << 3
	RequestedQoSMonitoringDLCI   uint16 = 1 << 4
	RequestedQoSMonitoringULCI   uint16 = 1 << 5
	RequestedQoSMonitoringDLDR   uint16 = 1 << 6
	RequestedQoSMonitoringULDR   uint16 = 1 << 7
	RequestedQoSMonitoringDLAB   uint16 = 1 << 8
	RequestedQoSMonitoringULAB   uint16 = 1 << 9
)
