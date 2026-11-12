package ie

import (
	"io"
)

// ReportingThresholds flags and values, TS 29.244 V20.1.0 clause 8.2.240.
const (
	// ReportingThresholdsDLCI selects DLCI.
	ReportingThresholdsDLCI uint8 = 1
	// ReportingThresholdsULCI selects ULCI.
	ReportingThresholdsULCI uint8 = 2
	// ReportingThresholdsDLDR selects DLDR.
	ReportingThresholdsDLDR uint8 = 4
	// ReportingThresholdsULDR selects ULDR.
	ReportingThresholdsULDR uint8 = 8
)

// NewReportingThresholds creates a Reporting Thresholds IE.
// TS 29.244 V20.1.0 clause 8.2.240.
func NewReportingThresholds(flags uint8, downlinkCongestion uint16, uplinkCongestion uint16, downlinkDataRate uint64, uplinkDataRate uint64) *IE {
	return recentFieldsIE(ReportingThresholds, NewReportingThresholdsFields(flags, downlinkCongestion, uplinkCongestion, downlinkDataRate, uplinkDataRate))
}

// ReportingThresholds returns the Reporting Thresholds fields.
// TS 29.244 V20.1.0 clause 8.2.240.
func (i *IE) ReportingThresholds() (*ReportingThresholdsFields, error) {
	b, err := recentIEPayload(i, ReportingThresholds, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseReportingThresholdsFields(b)
	if err != nil {
		return nil, &DecodeError{Type: ReportingThresholds, Err: err}
	}
	return f, nil
}

// ReportingThresholdsFields models TS 29.244 V20.1.0 clause 8.2.240.
// Flags select optional fields; nonzero values for absent fields are rejected.
// Rates are integer kilobits per second (1000 bits/s).
// Congestion values are hundredths of a percent, from 0 through 10000.
type ReportingThresholdsFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// DownlinkCongestion is the downlink congestion.
	DownlinkCongestion uint16
	// UplinkCongestion is the uplink congestion.
	UplinkCongestion uint16
	// DownlinkDataRate is the downlink data rate.
	DownlinkDataRate uint64
	// UplinkDataRate is the uplink data rate.
	UplinkDataRate uint64
}

// NewReportingThresholdsFields stores all arguments; Marshal validates their representation.
func NewReportingThresholdsFields(flags uint8, downlinkCongestion uint16, uplinkCongestion uint16, downlinkDataRate uint64, uplinkDataRate uint64) *ReportingThresholdsFields {
	return &ReportingThresholdsFields{flags, downlinkCongestion, uplinkCongestion, downlinkDataRate, uplinkDataRate}
}

// ParseReportingThresholdsFields decodes fields, ignoring spare bits and trailing extension octets.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseReportingThresholdsFields(b []byte) (*ReportingThresholdsFields, error) {
	f := &ReportingThresholdsFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *ReportingThresholdsFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if f.Flags & ^uint8(15) != 0 || f.Flags == 0 {
		return ErrInvalidValue
	}
	if f.Flags&1 == 0 && f.DownlinkCongestion != 0 {
		return ErrInvalidValue
	}
	if f.DownlinkCongestion > 10000 {
		return ErrInvalidValue
	}
	if f.Flags&2 == 0 && f.UplinkCongestion != 0 {
		return ErrInvalidValue
	}
	if f.UplinkCongestion > 10000 {
		return ErrInvalidValue
	}
	if f.Flags&4 == 0 && f.DownlinkDataRate != 0 {
		return ErrInvalidValue
	}
	if f.DownlinkDataRate > 1099511627775 {
		return ErrInvalidValue
	}
	if f.Flags&8 == 0 && f.UplinkDataRate != 0 {
		return ErrInvalidValue
	}
	if f.UplinkDataRate > 1099511627775 {
		return ErrInvalidValue
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *ReportingThresholdsFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	n := 1
	if f.Flags&1 != 0 {
		n += 2
	}
	if f.Flags&2 != 0 {
		n += 2
	}
	if f.Flags&4 != 0 {
		n += 5
	}
	if f.Flags&8 != 0 {
		n += 5
	}
	return n
}

// Marshal serializes fields as defined in clause 8.2.240.
func (f *ReportingThresholdsFields) Marshal() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	b := make([]byte, f.MarshalLen())
	if err := f.MarshalTo(b); err != nil {
		return nil, err
	}
	return b, nil
}

// MarshalTo writes fields to b, rejecting invalid values and short buffers.
func (f *ReportingThresholdsFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = f.Flags
	offset := 1
	if f.Flags&1 != 0 {
		putRecentUint(b[offset:offset+2], uint64(f.DownlinkCongestion))
		offset += 2
	}
	if f.Flags&2 != 0 {
		putRecentUint(b[offset:offset+2], uint64(f.UplinkCongestion))
		offset += 2
	}
	if f.Flags&4 != 0 {
		putRecentUint(b[offset:offset+5], f.DownlinkDataRate)
		offset += 5
	}
	if f.Flags&8 != 0 {
		putRecentUint(b[offset:offset+5], f.UplinkDataRate)
	}
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded known fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *ReportingThresholdsFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = ReportingThresholdsFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	f.Flags = b[0] & 15
	offset := 1
	if f.Flags&1 != 0 {
		if len(b)-offset < 2 {
			return io.ErrUnexpectedEOF
		}
		f.DownlinkCongestion = uint16(recentUint(b[offset : offset+2]))
		offset += 2
	}
	if f.Flags&2 != 0 {
		if len(b)-offset < 2 {
			return io.ErrUnexpectedEOF
		}
		f.UplinkCongestion = uint16(recentUint(b[offset : offset+2]))
		offset += 2
	}
	if f.Flags&4 != 0 {
		if len(b)-offset < 5 {
			return io.ErrUnexpectedEOF
		}
		f.DownlinkDataRate = recentUint(b[offset : offset+5])
		offset += 5
	}
	if f.Flags&8 != 0 {
		if len(b)-offset < 5 {
			return io.ErrUnexpectedEOF
		}
		f.UplinkDataRate = recentUint(b[offset : offset+5])
	}
	return f.validate()
}
