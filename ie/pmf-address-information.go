// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

// NewPMFAddressInformation creates a new PMFAddressInformation IE.
func NewPMFAddressInformation(v4, v6 net.IP, port1, port2 uint16, mac1, mac2 net.HardwareAddr) *IE {
	fields := NewPMFAddressInformationFields(v4, v6, port1, port2, mac1, mac2)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(PMFAddressInformation, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(PMFAddressInformation, b)
}

// PMFAddressInformation returns PMFAddressInformation in structured format if the type of IE matches.
func (i *IE) PMFAddressInformation() (*PMFAddressInformationFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case PMFAddressInformation:
		fields, err := ParsePMFAddressInformationFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return fields, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// PMFAddressInformationFields represents a fields contained in PMFAddressInformation IE.
type PMFAddressInformationFields struct {
	Flags                         uint8
	PMFIPv4Address                net.IP
	PMFIPv6Address                net.IP
	PMFPortFor3GPPAccess          uint16
	PMFPortForNon3GPPAccess       uint16
	PMFMACAddressFor3GPPAccess    net.HardwareAddr
	PMFMACAddressForNon3GPPAccess net.HardwareAddr
}

// NewPMFAddressInformationFields creates a new NewPMFAddressInformationFields.
func NewPMFAddressInformationFields(v4, v6 net.IP, port1, port2 uint16, mac1, mac2 net.HardwareAddr) *PMFAddressInformationFields {
	f := &PMFAddressInformationFields{
		Flags:                         0x00,
		PMFPortFor3GPPAccess:          port1,
		PMFPortForNon3GPPAccess:       port2,
		PMFMACAddressFor3GPPAccess:    mac1,
		PMFMACAddressForNon3GPPAccess: mac2,
	}

	if mac1 != nil || mac2 != nil {
		f.SetMACFlag()
	}
	if v4 != nil {
		f.SetIPv4Flag()
		f.PMFIPv4Address = v4
	}
	if v6 != nil {
		f.SetIPv6Flag()
		f.PMFIPv6Address = v6
	}

	if v4 == nil && v6 == nil {
		f.PMFPortFor3GPPAccess = 0
		f.PMFPortForNon3GPPAccess = 0
	}
	return f
}

// HasMAC reports whether MAC flag is set.
func (f *PMFAddressInformationFields) HasMAC() bool {
	if f == nil {
		return false
	}

	return has3rdBit(f.Flags)
}

// SetMACFlag sets MAC flag in PMFAddressInformation.
func (f *PMFAddressInformationFields) SetMACFlag() {
	if f == nil {
		return
	}

	f.Flags |= 0x04
}

// HasIPv6 reports whether IPv6 flag is set.
func (f *PMFAddressInformationFields) HasIPv6() bool {
	if f == nil {
		return false
	}

	return has2ndBit(f.Flags)
}

// SetIPv6Flag sets IPv6 flag in PMFAddressInformation.
func (f *PMFAddressInformationFields) SetIPv6Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x02
}

// HasIPv4 reports whether IPv4 flag is set.
func (f *PMFAddressInformationFields) HasIPv4() bool {
	if f == nil {
		return false
	}

	return has1stBit(f.Flags)
}

// SetIPv4Flag sets IPv4 flag in PMFAddressInformation.
func (f *PMFAddressInformationFields) SetIPv4Flag() {
	if f == nil {
		return
	}

	f.Flags |= 0x01
}

// ParsePMFAddressInformationFields parses b into PMFAddressInformationFields.
func ParsePMFAddressInformationFields(b []byte) (*PMFAddressInformationFields, error) {
	f := &PMFAddressInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *PMFAddressInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = PMFAddressInformationFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0]
	offset := 1

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		f.PMFIPv4Address = net.IP(b[offset : offset+4])
		offset += 4
	}

	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		f.PMFIPv6Address = net.IP(b[offset : offset+16])
		offset += 16
	}

	// TS 29.244 V20.1.0 clause 8.2.159: ports accompany IP addresses;
	// the two MAC addresses are present only when the MAC flag is set.
	if f.HasIPv4() || f.HasIPv6() {
		ports, err := takePayload(b, &offset, 4)
		if err != nil {
			return err
		}
		f.PMFPortFor3GPPAccess = binary.BigEndian.Uint16(ports[:2])
		f.PMFPortForNon3GPPAccess = binary.BigEndian.Uint16(ports[2:])
	}
	if f.HasMAC() {
		macs, err := takePayload(b, &offset, 12)
		if err != nil {
			return err
		}
		f.PMFMACAddressFor3GPPAccess = net.HardwareAddr(macs[:6:6])
		f.PMFMACAddressForNon3GPPAccess = net.HardwareAddr(macs[6:12:12])
	}

	return nil
}

// Marshal returns the serialized bytes of PMFAddressInformationFields.
func (f *PMFAddressInformationFields) Marshal() ([]byte, error) {
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
func (f *PMFAddressInformationFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags
	offset := 1

	if f.HasIPv4() {
		if l < offset+4 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+4], f.PMFIPv4Address.To4())
		offset += 4
	}

	if f.HasIPv6() {
		if l < offset+16 {
			return io.ErrUnexpectedEOF
		}
		copy(b[offset:offset+16], f.PMFIPv6Address.To16())
		offset += 16
	}

	if f.HasIPv4() || f.HasIPv6() {
		binary.BigEndian.PutUint16(b[offset:offset+2], f.PMFPortFor3GPPAccess)
		binary.BigEndian.PutUint16(b[offset+2:offset+4], f.PMFPortForNon3GPPAccess)
		offset += 4
	}
	if f.HasMAC() {
		if len(f.PMFMACAddressFor3GPPAccess) != 6 || len(f.PMFMACAddressForNon3GPPAccess) != 6 {
			return ErrInvalidLength
		}
		copy(b[offset:offset+6], f.PMFMACAddressFor3GPPAccess)
		copy(b[offset+6:offset+12], f.PMFMACAddressForNon3GPPAccess)
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *PMFAddressInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if f.HasIPv4() || f.HasIPv6() {
		l += 4
	}
	if f.HasMAC() {
		l += 12
	}
	if f.HasIPv4() {
		l += 4
	}
	if f.HasIPv6() {
		l += 16
	}
	return l
}
