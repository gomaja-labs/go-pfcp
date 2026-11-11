package ie

import (
	"io"
	"unicode/utf8"
)

// The key ID is UTF-8, at most 80 octets without padding. The shared secret
// is a UTF-8 encoded ASCII passphrase without CR or LF: TS 29.244 V20.1.0
// clause 8.2.276, RFC 4656 section 3.1 and RFC 5357 section 3.1.
// NewKeyIDAndSharedSecretOWAMPOrTWAMP creates a Key ID and Shared Secret (OWAMP or TWAMP) IE.
// TS 29.244 V20.1.0 clause 8.2.276.
func NewKeyIDAndSharedSecretOWAMPOrTWAMP(keyID, sharedSecret string) *IE {
	return recentFieldsIE(KeyIDAndSharedSecretOWAMPOrTWAMP, NewKeyIDAndSharedSecretOWAMPOrTWAMPFields(keyID, sharedSecret))
}

// KeyIDAndSharedSecretOWAMPOrTWAMP returns fields defined in TS 29.244 V20.1.0 clause 8.2.276.
func (i *IE) KeyIDAndSharedSecretOWAMPOrTWAMP() (*KeyIDAndSharedSecretOWAMPOrTWAMPFields, error) {
	b, err := recentIEPayload(i, KeyIDAndSharedSecretOWAMPOrTWAMP, 1, false)
	if err != nil {
		return nil, err
	}
	f, err := ParseKeyIDAndSharedSecretOWAMPOrTWAMPFields(b)
	if err != nil {
		return nil, &DecodeError{Type: KeyIDAndSharedSecretOWAMPOrTWAMP, Err: err}
	}
	return f, nil
}

// KeyIDAndSharedSecretOWAMPOrTWAMPFields models TS 29.244 V20.1.0 clause 8.2.276.
type KeyIDAndSharedSecretOWAMPOrTWAMPFields struct {
	KeyID        string
	SharedSecret string
}

// NewKeyIDAndSharedSecretOWAMPOrTWAMPFields retains all supplied arguments; Marshal validates them.
func NewKeyIDAndSharedSecretOWAMPOrTWAMPFields(keyID, sharedSecret string) *KeyIDAndSharedSecretOWAMPOrTWAMPFields {
	return &KeyIDAndSharedSecretOWAMPOrTWAMPFields{keyID, sharedSecret}
}

// ParseKeyIDAndSharedSecretOWAMPOrTWAMPFields parses the known fields and ignores trailing extension octets.
func ParseKeyIDAndSharedSecretOWAMPOrTWAMPFields(b []byte) (*KeyIDAndSharedSecretOWAMPOrTWAMPFields, error) {
	f := &KeyIDAndSharedSecretOWAMPOrTWAMPFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}
func (f *KeyIDAndSharedSecretOWAMPOrTWAMPFields) validate() error {
	if f == nil {
		return ErrIENotFound
	}
	if len(f.KeyID) > 80 || len(f.SharedSecret) > 255 || !utf8.ValidString(f.KeyID) || !utf8.ValidString(f.SharedSecret) {
		return ErrInvalidValue
	}
	for _, v := range f.SharedSecret {
		if v > 127 || v == '\r' || v == '\n' {
			return ErrInvalidValue
		}
	}
	return nil
}

// MarshalLen returns the known field size in octets.
func (f *KeyIDAndSharedSecretOWAMPOrTWAMPFields) MarshalLen() int {
	if f == nil {
		return 0
	}
	return 2 + len(f.KeyID) + len(f.SharedSecret)
}

// Marshal serializes the fields.
func (f *KeyIDAndSharedSecretOWAMPOrTWAMPFields) Marshal() ([]byte, error) {
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
func (f *KeyIDAndSharedSecretOWAMPOrTWAMPFields) MarshalTo(b []byte) error {
	if err := f.validate(); err != nil {
		return err
	}
	if len(b) < f.MarshalLen() {
		return io.ErrUnexpectedEOF
	}
	b[0] = byte(len(f.KeyID))
	copy(b[1:], f.KeyID)
	offset := 1 + len(f.KeyID)
	b[offset] = byte(len(f.SharedSecret))
	copy(b[offset+1:], f.SharedSecret)
	return nil
}

// UnmarshalBinary replaces the receiver with the decoded fields.
func (f *KeyIDAndSharedSecretOWAMPOrTWAMPFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = KeyIDAndSharedSecretOWAMPOrTWAMPFields{}
	if len(b) < 1 {
		return io.ErrUnexpectedEOF
	}
	n := int(b[0])
	if len(b) < 2+n {
		return io.ErrUnexpectedEOF
	}
	m := int(b[1+n])
	if len(b) < 2+n+m {
		return io.ErrUnexpectedEOF
	}
	f.KeyID = string(b[1 : 1+n])
	f.SharedSecret = string(b[2+n : 2+n+m])
	return f.validate()
}
