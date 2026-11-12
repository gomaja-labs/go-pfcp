// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"
	"strings"
	"unicode/utf8"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

// NewUserID creates a new UserID IE.
func NewUserID(flags uint8, imsi, imei, msisdn, nai string) *IE {
	fields := NewUserIDFields(flags, imsi, imei, msisdn, nai)

	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(UserID, ErrInvalidValue)
	}

	return New(UserID, b)
}

// UserID returns UserID in structured format if the type of IE matches.
func (i *IE) UserID() (*UserIDFields, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	if i.Type != UserID {
		return nil, &InvalidTypeError{Type: i.Type}
	}

	s, err := ParseUserIDFields(i.Payload)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// UserIDFields represents a fields contained in UserID IE.
type UserIDFields struct {
	// Flags contains the presence and option bits.
	Flags uint8
	// IMSILength is the encoded IMSI length in octets.
	IMSILength uint8
	// IMSI is the international mobile subscriber identity.
	IMSI string
	// IMEILength is the encoded IMEI length in octets.
	IMEILength uint8
	// IMEI is the international mobile equipment identity.
	IMEI string
	// MSISDNLength is the encoded MSISDN length in octets.
	MSISDNLength uint8
	// MSISDN is the mobile subscriber number.
	MSISDN string
	// NAILength is the encoded NAI length in octets.
	NAILength uint8
	// NAI is the network access identifier.
	NAI string
	// SUPI is the subscription permanent identifier.
	SUPI string
	// GPSI is the generic public subscription identifier.
	GPSI string
	// PEI is the permanent equipment identifier.
	PEI string
}

// NewUserIDFields creates a new NewUserIDFields.
func NewUserIDFields(flags uint8, imsi, imei, msisdn, nai string) *UserIDFields {
	f := &UserIDFields{
		Flags: flags,
	}

	f.IMSI = imsi
	f.IMSILength = uint8((len(imsi) + 1) / 2)
	f.IMEI = imei
	f.IMEILength = uint8((len(imei) + 1) / 2)
	f.MSISDN = msisdn
	f.MSISDNLength = uint8((len(msisdn) + 1) / 2)
	f.NAI = nai
	f.NAILength = uint8(len(nai))

	return f
}

// ParseUserIDFields parses b into UserIDFields.
func ParseUserIDFields(b []byte) (*UserIDFields, error) {
	f := &UserIDFields{}
	if err := f.UnmarshalBinary(b); err != nil {
		return nil, err
	}
	return f, nil
}

// UnmarshalBinary parses b into IE.
func (f *UserIDFields) UnmarshalBinary(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	*f = UserIDFields{}

	l := len(b)
	if l < 1 {
		return io.ErrUnexpectedEOF
	}

	f.Flags = b[0] & 0x7f
	offset := 1

	if has1stBit(f.Flags) {
		if l < offset+1 {
			return io.ErrUnexpectedEOF
		}
		f.IMSILength = b[offset]
		offset++

		if l < offset+int(f.IMSILength) {
			return io.ErrUnexpectedEOF
		}
		f.IMSI = utils.SwappedBytesToStr(b[offset:offset+int(f.IMSILength)], true)
		offset += int(f.IMSILength)
	}

	if has2ndBit(f.Flags) {
		if l < offset+1 {
			return io.ErrUnexpectedEOF
		}
		f.IMEILength = b[offset]
		offset++

		if l < offset+int(f.IMEILength) {
			return io.ErrUnexpectedEOF
		}
		f.IMEI = utils.SwappedBytesToStr(b[offset:offset+int(f.IMEILength)], true)
		offset += int(f.IMEILength)
	}

	if has3rdBit(f.Flags) {
		if l < offset+1 {
			return io.ErrUnexpectedEOF
		}
		f.MSISDNLength = b[offset]
		offset++

		if l < offset+int(f.MSISDNLength) {
			return io.ErrUnexpectedEOF
		}
		f.MSISDN = utils.SwappedBytesToStr(b[offset:offset+int(f.MSISDNLength)], true)
		offset += int(f.MSISDNLength)
	}

	if has4thBit(f.Flags) {
		if l < offset+1 {
			return io.ErrUnexpectedEOF
		}
		f.NAILength = b[offset]
		offset++

		if l < offset+int(f.NAILength) {
			return io.ErrUnexpectedEOF
		}
		f.NAI = string(b[offset : offset+int(f.NAILength)])
		offset += int(f.NAILength)
	}

	// TS 29.244 V20.1.0 clause 8.2.101: SUPI/GPSI/PEI use TS 29.571
	// V20.1.0 clause 5.3.2 string data types, rather than the legacy TBCD fields.
	for _, v := range []struct {
		bit uint8
		dst *string
	}{{0x10, &f.SUPI}, {0x20, &f.GPSI}, {0x40, &f.PEI}} {
		if f.Flags&v.bit != 0 {
			if l < offset+1 {
				return io.ErrUnexpectedEOF
			}
			n := int(b[offset])
			offset++
			if l < offset+n {
				return io.ErrUnexpectedEOF
			}
			*v.dst = string(b[offset : offset+n])
			offset += n
		}
	}
	return f.validateIdentities()
}

// Marshal returns the serialized bytes of UserIDFields.
func (f *UserIDFields) Marshal() ([]byte, error) {
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
func (f *UserIDFields) MarshalTo(b []byte) error {
	if f == nil {
		return ErrIENotFound
	}
	if err := f.validateIdentities(); err != nil {
		return err
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

	b[0] = f.Flags & 0x7f
	offset := 1

	if has1stBit(f.Flags) {
		b[offset] = f.IMSILength
		offset++

		v, err := utils.StrToSwappedBytes(f.IMSI, "f")
		if err != nil {
			return err
		}
		if len(v) != int(f.IMSILength) {
			return ErrInvalidLength
		}
		copy(b[offset:offset+int(f.IMSILength)], v)
		offset += int(f.IMSILength)
	}

	if has2ndBit(f.Flags) {
		b[offset] = f.IMEILength
		offset++

		v, err := utils.StrToSwappedBytes(f.IMEI, "f")
		if err != nil {
			return err
		}
		if len(v) != int(f.IMEILength) {
			return ErrInvalidLength
		}
		copy(b[offset:offset+int(f.IMEILength)], v)
		offset += int(f.IMEILength)
	}

	if has3rdBit(f.Flags) {
		b[offset] = f.MSISDNLength
		offset++

		v, err := utils.StrToSwappedBytes(f.MSISDN, "f")
		if err != nil {
			return err
		}
		if len(v) != int(f.MSISDNLength) {
			return ErrInvalidLength
		}
		copy(b[offset:offset+int(f.MSISDNLength)], v)
		offset += int(f.MSISDNLength)
	}

	if has4thBit(f.Flags) {
		b[offset] = f.NAILength
		offset++

		if len(f.NAI) != int(f.NAILength) {
			return ErrInvalidLength
		}
		copy(b[offset:offset+int(f.NAILength)], f.NAI)
		offset += int(f.NAILength)
	}

	for _, v := range []struct {
		bit   uint8
		value string
	}{{0x10, f.SUPI}, {0x20, f.GPSI}, {0x40, f.PEI}} {
		if f.Flags&v.bit != 0 {
			b[offset] = uint8(len(v.value))
			offset++
			copy(b[offset:], v.value)
			offset += len(v.value)
		}
	}
	return nil
}

// MarshalLen returns field length in integer.
func (f *UserIDFields) MarshalLen() int {
	if f == nil {
		return 0
	}

	l := 1
	if has1stBit(f.Flags) {
		l += 1 + int(f.IMSILength)
	}
	if has2ndBit(f.Flags) {
		l += 1 + int(f.IMEILength)
	}
	if has3rdBit(f.Flags) {
		l += 1 + int(f.MSISDNLength)
	}
	if has4thBit(f.Flags) {
		l += 1 + int(f.NAILength)
	}

	for _, v := range []struct {
		bit   uint8
		value string
	}{{0x10, f.SUPI}, {0x20, f.GPSI}, {0x40, f.PEI}} {
		if f.Flags&v.bit != 0 {
			l += 1 + len(v.value)
		}
	}
	return l
}

// NewUserIDFromFields encodes all current fields of UserID.
func NewUserIDFromFields(fields *UserIDFields) *IE {
	b, err := fields.Marshal()
	if err != nil {
		return newInvalidIE(UserID, ErrInvalidValue)
	}
	return New(UserID, b)
}

// TS 29.244 V20.1.0 clause 8.2.101 limits these fields to GCI/GLI,
// External Identifier, and MAC/EUI respectively. Other identity kinds belong
// in the existing IMSI/NAI, MSISDN, or IMEI fields.
func (f *UserIDFields) validateIdentities() error {
	if f.Flags&0x80 != 0 || f.Flags&1 == 0 && (f.IMSI != "" || f.IMSILength != 0) || f.Flags&2 == 0 && (f.IMEI != "" || f.IMEILength != 0) || f.Flags&4 == 0 && (f.MSISDN != "" || f.MSISDNLength != 0) || f.Flags&8 == 0 && (f.NAI != "" || f.NAILength != 0) {
		return ErrInvalidValue
	}

	for _, v := range []struct {
		bit      uint8
		s        string
		prefixes []string
	}{{0x10, f.SUPI, []string{"gci-", "gli-"}}, {0x20, f.GPSI, []string{"extid-"}}, {0x40, f.PEI, []string{"mac-", "eui-"}}} {
		if f.Flags&v.bit == 0 {
			if v.s != "" {
				return ErrInvalidValue
			}
			continue
		}
		if len(v.s) > 255 {
			return ErrInvalidLength
		}
		if !utf8.ValidString(v.s) {
			return ErrInvalidValue
		}
		ok := false
		for _, p := range v.prefixes {
			if strings.HasPrefix(v.s, p) && len(v.s) > len(p) {
				ok = true
			}
		}
		if !ok {
			return ErrInvalidValue
		}
		if v.bit == 0x20 {
			local, domain, found := strings.Cut(strings.TrimPrefix(v.s, "extid-"), "@")
			if !found || local == "" || domain == "" || strings.Contains(domain, "@") {
				return ErrInvalidValue
			}
		}
		if v.bit == 0x40 && !validUserPEI(v.s) {
			return ErrInvalidValue
		}
	}
	return nil
}

// TS 29.571 V20.1.0 clause 5.3.2, Pei: MAC/EUI octets are hyphen-separated.
func validUserPEI(s string) bool {
	parts := strings.Split(s, "-")
	n := 6
	if len(parts) == 0 {
		return false
	}
	switch parts[0] {
	case "mac":
		if len(parts) == 8 && parts[7] == "untrusted" {
			parts = parts[:7]
		}
	case "eui":
		n = 8
	default:
		return false
	}
	if len(parts) != n+1 {
		return false
	}
	for _, part := range parts[1:] {
		if len(part) != 2 {
			return false
		}
		for _, c := range part {
			switch {
			case c >= '0' && c <= '9', c >= 'a' && c <= 'f', c >= 'A' && c <= 'F':
			default:
				return false
			}
		}
	}
	return true
}
