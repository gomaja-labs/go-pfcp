package ie

import (
	"encoding/json"
	"unicode/utf8"
)

// DNSQueryResponseFilterFields is the structured value of DNS Query/Response Filter.
// TS 29.244 V20.1.0 clause 8.2.201.
type DNSQueryResponseFilterFields struct {
	// DomainNamePattern contains the exact JSON octet string defined in
	// TS 29.244 V20.1.0 clause 8.2.201. Validation preserves its spelling.
	DomainNamePattern string
}

// NewDNSQueryResponseFilter creates a DNS Query/Response Filter IE from its fields.
// TS 29.244 V20.1.0 clause 8.2.201.
func NewDNSQueryResponseFilter(fields DNSQueryResponseFilterFields) *IE {
	b, err := fields.Marshal()
	return addedFieldsIE(DNSQueryResponseFilter, b, err)
}

// DNSQueryResponseFilter returns the structured DNS Query/Response Filter value.
// TS 29.244 V20.1.0 clauses 8.2.201 and 8.1.2.
func (i *IE) DNSQueryResponseFilter() (*DNSQueryResponseFilterFields, error) {
	b, err := addedValuePayload(i, DNSQueryResponseFilter, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseDNSQueryResponseFilterFields(b)
	if err != nil {
		return nil, &DecodeError{Type: DNSQueryResponseFilter, Err: err}
	}
	return f, nil
}

// ParseDNSQueryResponseFilterFields decodes the known fields; unknown extension octets are ignored.
// TS 29.244 V20.1.0 clauses 8.2.201 and 8.1.2.
// Decoding may borrow b; keep b unchanged while the result is in use.
func ParseDNSQueryResponseFilterFields(b []byte) (*DNSQueryResponseFilterFields, error) {
	f := &DNSQueryResponseFilterFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary decodes DNS Query/Response Filter and resets previously decoded fields.
// Decoding may borrow b; keep b unchanged while the result is in use.
func (f *DNSQueryResponseFilterFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = DNSQueryResponseFilterFields{}
	r := &addedReader{b: b}
	n := int(r.u16())
	v := r.take(n)
	if r.err == nil {
		f.DomainNamePattern = string(v)
	}
	if r.err != nil {
		return r.err
	}

	if !f.valid() {
		return ErrInvalidValue
	}
	return nil
}
func (f *DNSQueryResponseFilterFields) valid() bool {
	return len(f.DomainNamePattern) <= 65533 && addedDNSPatternValid([]byte(f.DomainNamePattern))
}

// Marshal encodes the known fields with spare bits set to zero.
func (f *DNSQueryResponseFilterFields) Marshal() ([]byte, error) {
	if f == nil {
		return nil, ErrIENotFound
	}
	if !f.valid() {
		return nil, ErrInvalidValue
	}
	b := addedPut16(nil, uint16(len(f.DomainNamePattern)))
	b = append(b, f.DomainNamePattern...)
	if len(b) > 65535 {
		return nil, ErrInvalidValue
	}
	return b, nil
}

// MarshalLen returns the encoded value length, or zero for nil or invalid fields.
func (f *DNSQueryResponseFilterFields) MarshalLen() int { b, _ := f.Marshal(); return len(b) }

// MarshalTo writes the encoded fields to b, rejecting an undersized buffer.
func (f *DNSQueryResponseFilterFields) MarshalTo(b []byte) error {
	wire, err := f.Marshal()
	return addedCopyTo(b, wire, err)
}

// addedDNSPatternValid checks the JSON datatype without rewriting or dropping
// unknown object members. TS 29.571 clause 5.2.4.23.
func addedDNSPatternValid(b []byte) bool {
	if !utf8.Valid(b) {
		return false
	}
	var obj map[string]json.RawMessage
	if json.Unmarshal(b, &obj) != nil || obj == nil {
		return false
	}
	regex, hasRegex := obj["regex"]
	rule, hasRule := obj["stringMatchingRule"]
	if hasRegex == hasRule {
		return false
	}
	if hasRegex {
		var value string
		return json.Unmarshal(regex, &value) == nil && string(regex) != "null"
	}
	var match map[string]json.RawMessage
	if json.Unmarshal(rule, &match) != nil || match == nil {
		return false
	}
	var conditions []map[string]json.RawMessage
	if json.Unmarshal(match["stringMatchingConditions"], &conditions) != nil || len(conditions) == 0 {
		return false
	}
	for _, condition := range conditions {
		var operator string
		raw, ok := condition["matchingOperator"]
		if !ok || string(raw) == "null" || json.Unmarshal(raw, &operator) != nil {
			return false
		}
		if raw, ok = condition["matchingString"]; ok {
			var value string
			if string(raw) == "null" || json.Unmarshal(raw, &value) != nil {
				return false
			}
		}
	}
	return true
}
