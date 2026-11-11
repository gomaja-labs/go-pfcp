package ie

import (
	"io"
)

// RTP header extension flags and format values, TS 29.244 V20.1.0 clause 8.2.241.
const (
	// Second-octet presence values in Figure 8.2.241-1.
	RTPHeaderExtensionAdditionalInformationPSSA        uint8 = 4
	RTPHeaderExtensionAdditionalInformationNPDS        uint8 = 8
	RTPHeaderExtensionAdditionalInformationFI          uint8 = 1
	RTPHeaderExtensionAdditionalInformationPSSAI       uint8 = 2
	RTPHeaderExtensionAdditionalInformationNPDSI       uint8 = 4
	RTPHeaderExtensionAdditionalInformationShortFormat uint8 = 1
	RTPHeaderExtensionAdditionalInformationLongFormat  uint8 = 2
)

// NewRTPHeaderExtensionAdditionalInformation creates a RTP Header Extension Additional Information IE.
// TS 29.244 V20.1.0 clause 8.2.241.
func NewRTPHeaderExtensionAdditionalInformation(flags, format uint8, pssa, npds bool) *IE {
	return recentFieldsIE(RTPHeaderExtensionAdditionalInformation, NewRTPHeaderExtensionAdditionalInformationFields(flags, format, pssa, npds))
}

// RTPHeaderExtensionAdditionalInformation returns fields defined in TS 29.244 V20.1.0 clause 8.2.241.
func (i *IE) RTPHeaderExtensionAdditionalInformation() (*RTPHeaderExtensionAdditionalInformationFields, error) {
	b, err := recentIEPayload(i, RTPHeaderExtensionAdditionalInformation, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseRTPHeaderExtensionAdditionalInformationFields(b)
	if err != nil {
		return nil, &DecodeError{Type: RTPHeaderExtensionAdditionalInformation, Err: err}
	}
	return f, nil
}

// RTPHeaderExtensionAdditionalInformationFields models TS 29.244 V20.1.0 clause 8.2.241.
type RTPHeaderExtensionAdditionalInformationFields struct {
	Flags  uint8
	Format uint8
	PSSA   bool
	NPDS   bool
}

// NewRTPHeaderExtensionAdditionalInformationFields retains all supplied arguments; Marshal validates them.
func NewRTPHeaderExtensionAdditionalInformationFields(flags, format uint8, pssa, npds bool) *RTPHeaderExtensionAdditionalInformationFields {
	return &RTPHeaderExtensionAdditionalInformationFields{flags, format, pssa, npds}
}

// ParseRTPHeaderExtensionAdditionalInformationFields parses the known fields and ignores trailing extension octets.
func ParseRTPHeaderExtensionAdditionalInformationFields(b []byte) (*RTPHeaderExtensionAdditionalInformationFields, error) {
	f := &RTPHeaderExtensionAdditionalInformationFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *RTPHeaderExtensionAdditionalInformationFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(7) != 0 || f.Format > 2 || (f.Flags&1 == 0 && f.Format != 0) || (f.Flags&1 != 0 && f.Format == 0) || (f.Flags&2 == 0 && f.PSSA) || (f.Flags&4 == 0 && f.NPDS) {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *RTPHeaderExtensionAdditionalInformationFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2
}

// Marshal serializes the fields.
func (f *RTPHeaderExtensionAdditionalInformationFields) Marshal() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo serializes the fields into b, rejecting short buffers and invalid values.
func (f *RTPHeaderExtensionAdditionalInformationFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	b[1] = f.Format
	if f.PSSA {
		b[1] |= 4
	}
	if f.NPDS {
		b[1] |= 8
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded fields.
func (f *RTPHeaderExtensionAdditionalInformationFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = RTPHeaderExtensionAdditionalInformationFields{}
	if len(b) < 2 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 7
	if f.Flags&1 != 0 {
		f.Format = b[1] & 3
	}
	f.PSSA = f.Flags&2 != 0 && b[1]&4 != 0
	f.NPDS = f.Flags&4 != 0 && b[1]&8 != 0
	return f.validate()
}
