// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"net"
)

// NewMACAddressesDetected creates a new MACAddressesDetected IE.
func NewMACAddressesDetected(ctag, stag *IE, macs ...net.HardwareAddr) *IE {
	fields := NewMACAddressesDetectedFields(ctag, stag, macs...)

	b, err := fields.Marshal()
	if err != nil {
		return nil
	}

	return New(MACAddressesDetected, b)
}

// MACAddressesDetected returns MACAddressesDetected in structured format if the type of IE matches.
func (i *IE) MACAddressesDetected() (*MACAddressesDetectedFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case MACAddressesDetected:
		fields, err := ParseMACAddressesDetectedFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	case EthernetTrafficInformation:
		ies, err := i.EthernetTrafficInformation()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == MACAddressesDetected {
				return x.MACAddressesDetected()
			}
		}
		return nil, ErrIENotFound
	case EthernetContextInformation:
		ies, err := i.EthernetContextInformation()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == MACAddressesDetected {
				return x.MACAddressesDetected()
			}
		}
		return nil, ErrIENotFound
	case UsageReportWithinSessionModificationResponse,
		UsageReportWithinSessionDeletionResponse,
		UsageReportWithinSessionReportRequest:
		ies, err := i.UsageReport()
		if err != nil {
			return nil, err
		}
		for _, x := range ies {
			if x.Type == EthernetTrafficInformation {
				return x.MACAddressesDetected()
			}
		}
		return nil, ErrIENotFound
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// MACAddressesDetectedFields represents a fields contained in MACAddressesDetected IE.
type MACAddressesDetectedFields struct {
	NumberOfMACAddresses uint8
	MACAddresses         []net.HardwareAddr
	CTAGLength           uint8
	CTAG                 []byte
	STAGLength           uint8
	STAG                 []byte
}

// NewMACAddressesDetectedFields creates a new NewMACAddressesDetectedFields.
func NewMACAddressesDetectedFields(ctag, stag *IE, macs ...net.HardwareAddr) *MACAddressesDetectedFields {
	var ct, st []byte
	if ctag != nil {
		ct = ctag.Payload
	}
	if stag != nil {
		st = stag.Payload
	}
	return &MACAddressesDetectedFields{
		NumberOfMACAddresses: uint8(len(macs)),
		MACAddresses:         macs,
		CTAGLength:           uint8(len(ct)),
		CTAG:                 ct,
		STAGLength:           uint8(len(st)),
		STAG:                 st,
	}
}

// ParseMACAddressesDetectedFields parses b into MACAddressesDetectedFields.
func ParseMACAddressesDetectedFields(b []byte) (*MACAddressesDetectedFields, error) {
	f := &MACAddressesDetectedFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
// TS 29.244 V20.1.0 clause 8.2.103: exactly count MAC addresses,
// followed by independently length-prefixed C-TAG and S-TAG values.
func (f *MACAddressesDetectedFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = MACAddressesDetectedFields{}

	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	*f = MACAddressesDetectedFields{NumberOfMACAddresses: b[0]}
	f.MACAddresses = make([]net.HardwareAddr, 0, int(b[0]))
	offset := 1
	for n := 0; n < int(f.NumberOfMACAddresses); n++ {
		mac, err := takePayload(b, &offset, 6)
		if err != nil {
			return err
		}
		f.MACAddresses = append(f.MACAddresses, net.HardwareAddr(mac))
	}
	for _, field := range []struct {
		length *uint8
		value  *[]byte
	}{{&f.CTAGLength, &f.CTAG}, {&f.STAGLength, &f.STAG}} {
		length, err := takePayload(b, &offset, 1)
		if err != nil {
			return err
		}
		*field.length = length[0]
		value, err := takePayload(b, &offset, int(length[0]))
		if err != nil {
			return err
		}
		if len(value) > 0 {
			*field.value = value
		}
	}
	return nil
}

// Marshal returns the serialized bytes of MACAddressesDetectedFields.
func (f *MACAddressesDetectedFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}

	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
func (f *MACAddressesDetectedFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if int(f.NumberOfMACAddresses) != len(f.MACAddresses) || int(f.CTAGLength) != len(f.CTAG) || int(f.STAGLength) != len(f.STAG) {
		return ErrInvalidLength
	}
	for _, mac := range f.MACAddresses {
		if len(mac) != 6 {
			return ErrInvalidLength
		}
	}

	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.NumberOfMACAddresses
	offset := 1

	for _, mac := range f.MACAddresses {
		copy(b[offset:offset+6], mac)
		offset += 6
	}

	b[offset] = f.CTAGLength
	offset++
	copy(b[offset:offset+int(f.CTAGLength)], f.CTAG)
	offset += int(f.CTAGLength)

	b[offset] = f.STAGLength
	offset++
	copy(b[offset:offset+int(f.STAGLength)], f.STAG)

	return nil
}

// MarshalLen returns field length in integer.
func (f *MACAddressesDetectedFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 3
	l += int(f.NumberOfMACAddresses) * 6
	l += int(f.CTAGLength)
	l += int(f.STAGLength)

	return l
}
