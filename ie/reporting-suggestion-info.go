package ie

import (
	"time"
)

// ReportingUrgency is defined in Table 8.2.229-1.
type ReportingUrgency uint8

const (
	// ReportingDelayTolerant selects Delay Tolerant.
	ReportingDelayTolerant ReportingUrgency = 0
	// ReportingNonDelayTolerant selects Non Delay Tolerant.
	ReportingNonDelayTolerant ReportingUrgency = 1
)

// ReportingSuggestionInfoFields is the structured value of Reporting Suggestion Info.
// TS 29.244 V20.1.0 clause 8.2.229.
type ReportingSuggestionInfoFields struct {
	// Urgency selects delay-tolerant or non-delay-tolerant reporting.
	Urgency ReportingUrgency
	// ReportingTime is the delay-tolerant reporting interval.
	ReportingTime time.Duration
}

// NewReportingSuggestionInfo creates a Reporting Suggestion Info IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.229.
func NewReportingSuggestionInfo(fields ReportingSuggestionInfoFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(ReportingSuggestionInfo, b, err)
}

// ReportingSuggestionInfo returns the structured Reporting Suggestion Info value.
// TS 29.244 V20.1.0 clauses 8.2.229 and 8.1.2.
func (i *IE) ReportingSuggestionInfo() (*ReportingSuggestionInfoFields, error) {
	b, err := addedValuePayload(i, ReportingSuggestionInfo, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseReportingSuggestionInfoFields(b)
	if err != nil {
		return nil, &DecodeError{Type: ReportingSuggestionInfo, Err: err}
	}
	return f, nil
}

// ParseReportingSuggestionInfoFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.229 and 8.1.2.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseReportingSuggestionInfoFields(b []byte) (*ReportingSuggestionInfoFields, error) {
	f := &ReportingSuggestionInfoFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes Reporting Suggestion Info and resets previously decoded fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *ReportingSuggestionInfoFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = ReportingSuggestionInfoFields{}
	r := &addedReader{b: b}
	f.Urgency = ReportingUrgency(r.u8() & 15)
	if f.Urgency == ReportingDelayTolerant {
		f.ReportingTime = time.Duration(r.u32()) * time.Second
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *ReportingSuggestionInfoFields) valid() bool {
	return f.Urgency <= ReportingNonDelayTolerant && addedDurationFits(f.ReportingTime, time.Second, 0xffffffff) && (f.Urgency == ReportingDelayTolerant || f.ReportingTime == 0)
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *ReportingSuggestionInfoFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := []byte{byte(f.Urgency)}
	if f.Urgency == ReportingDelayTolerant {
		b = addedPut32(b, uint32(f.ReportingTime/time.Second))
	}
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *ReportingSuggestionInfoFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *ReportingSuggestionInfoFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}
