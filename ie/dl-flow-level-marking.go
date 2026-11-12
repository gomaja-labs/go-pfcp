// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"fmt"
	"io"
)

// NewDLFlowLevelMarking creates a new DLFlowLevelMarking IE.
func NewDLFlowLevelMarking(flags uint8, ttc, sci []byte) *IE {
	f := NewDLFlowLevelMarkingFields(flags, ttc, sci)
	b, err := f.Marshal()
	if err != nil {
		return newInvalidIE(DLFlowLevelMarking, fmt.Errorf("%w: %w", ErrInvalidValue, err))
	}

	return New(DLFlowLevelMarking, b)
}

// DLFlowLevelMarking returns DLFlowLevelMarking in *DLFlowLevelMarkingFields if the type of IE matches.
func (i *IE) DLFlowLevelMarking() (*DLFlowLevelMarkingFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case DLFlowLevelMarking:
		f, err := ParseDLFlowLevelMarkingFields(i.Payload)
		if err != nil {
			return nil, err
		}

		return f, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// HasTTC reports whether an IE has TTC bit.
func (i *IE) HasTTC() bool {
	if i == nil {
		return false
	}

	if i.Type != DLFlowLevelMarking {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	// TS 29.244 V20.1.0 clause 8.2.66: TTC is bit 1.
	return has1stBit(i.Payload[0])
}

// HasSCI reports whether an IE has SCI bit.
func (i *IE) HasSCI() bool {
	if i == nil {
		return false
	}

	if i.Type != DLFlowLevelMarking {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	// TS 29.244 V20.1.0 clause 8.2.66: SCI is bit 2.
	return has2ndBit(i.Payload[0])
}

// DLFlowLevelMarkingFields represents a f contained in DLFlowLevelMarking IE.
type DLFlowLevelMarkingFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// ToSTrafficClass contains the ToS/traffic class and mask octets.
	ToSTrafficClass []byte
	// ServiceClassIdentifier contains the service class identifier octets.
	ServiceClassIdentifier []byte
}

// NewDLFlowLevelMarkingFields creates a new DLFlowLevelMarkingFields.
func NewDLFlowLevelMarkingFields(flags uint8, ttc, sci []byte) *DLFlowLevelMarkingFields {
	f := &DLFlowLevelMarkingFields{Flags: flags}

	f.ToSTrafficClass = ttc

	f.ServiceClassIdentifier = sci

	return f
}

// ParseDLFlowLevelMarkingFields parses b into DLFlowLevelMarkingFields.
func ParseDLFlowLevelMarkingFields(b []byte) (*DLFlowLevelMarkingFields, error) {
	f := &DLFlowLevelMarkingFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *DLFlowLevelMarkingFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = DLFlowLevelMarkingFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	// TS 29.244 V20.1.0 clauses 8.1.2 and 8.2.66: ignore spare bits.
	f.Flags = b[0] & 3
	offset := 1

	if has1stBit(f.Flags) {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		f.ToSTrafficClass = b[offset : offset+2]
		offset += 2
	}

	if has2ndBit(f.Flags) {
		if l < offset+2 {
			return io.ErrUnexpectedEOF
		}
		f.ServiceClassIdentifier = b[offset : offset+2]
	}

	return nil
}

// Marshal returns the serialized bytes of DLFlowLevelMarkingFields.
func (f *DLFlowLevelMarkingFields) Marshal() ([]byte, error) {
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
func (f *DLFlowLevelMarkingFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	// TS 29.244 V20.1.0 clause 8.2.66: each binary marking is exactly two octets.
	if f.Flags&^uint8(3) != 0 || f.Flags&1 != 0 && len(f.ToSTrafficClass) != 2 || f.Flags&2 != 0 && len(f.ServiceClassIdentifier) != 2 || f.Flags&1 == 0 && len(f.ToSTrafficClass) != 0 || f.Flags&2 == 0 && len(f.ServiceClassIdentifier) != 0 {
		return ErrInvalidValue
	}
	if f.MarshalLen() > 65535 {
		return ErrInvalidLength
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}

	l := len(b)
	// TS 29.244 V20.1.0 clause 8.2.66: the presence flags are one octet.
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	b[0] = f.Flags
	offset := 1

	if has1stBit(f.Flags) {
		copy(b[offset:offset+2], f.ToSTrafficClass)
		offset += 2
	}

	if has2ndBit(f.Flags) {
		copy(b[offset:offset+2], f.ServiceClassIdentifier)
	}

	return nil
}

// MarshalLen returns field length in integer.
func (f *DLFlowLevelMarkingFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if has1stBit(f.Flags) {
		l += 2
	}
	if has2ndBit(f.Flags) {
		l += 2
	}

	return l
}
