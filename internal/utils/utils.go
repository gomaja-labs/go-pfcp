// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package utils

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// StrToSwappedBytes returns swapped bits from a byte.
// It is used for some values where some values are represented in swapped format.
//
// The second parameter is the hex character(0-f) to fill the last digit when
// handling a odd number. "f" is used In most cases.
func StrToSwappedBytes(s, filler string) ([]byte, error) {
	var raw []byte
	var err error
	if len(s)%2 == 0 {
		raw, err = hex.DecodeString(s)
	} else {
		raw, err = hex.DecodeString(s + filler)
	}
	if err != nil {
		return nil, err
	}

	return swap(raw), nil
}

// SwappedBytesToStr decodes raw swapped bytes into string.
// It is used for some values where some values are represented in swapped format.
//
// The second parameter is to decide whether to cut the last digit or not.
func SwappedBytesToStr(raw []byte, cutLastDigit bool) string {
	if len(raw) == 0 {
		return ""
	}

	s := hex.EncodeToString(swap(raw))
	if cutLastDigit {
		// TS 29.244 V20.1.0 clause 8.2.101 carries TBCD identities.
		// Only an actual high-nibble filler removes a final digit.
		s = strings.TrimSuffix(s, "f")
	}

	return s
}

func swap(raw []byte) []byte {
	swapped := make([]byte, len(raw))
	for n := range raw {
		t := ((raw[n] >> 4) & 0xf) + ((raw[n] << 4) & 0xf0)
		swapped[n] = t
	}
	return swapped
}

// Uint24To32 converts 24bits-length []byte value into the uint32 with 8bits of zeros as prefix.
// A slice whose length is not three returns an error.
func Uint24To32(b []byte) (uint32, error) {
	if len(b) != 3 {
		return 0, fmt.Errorf("24-bit integer requires 3 octets")
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2]), nil
}

// Uint32To24 converts the uint32 value into 24bits-length []byte. The values in 25-32 bit are cut off.
// This function is used for the fields with 3 octets.
func Uint32To24(n uint32) []byte {
	return []byte{uint8(n >> 16), uint8(n >> 8), uint8(n)}
}

// Uint40To64 converts 40bits-length []byte value into the uint64 with 8bits of zeros as prefix.
// A slice whose length is not five returns an error.
func Uint40To64(b []byte) (uint64, error) {
	if len(b) != 5 {
		return 0, fmt.Errorf("40-bit integer requires 5 octets")
	}
	return uint64(b[0])<<32 | uint64(b[1])<<24 | uint64(b[2])<<16 | uint64(b[3])<<8 | uint64(b[4]), nil
}

// Uint64To40 converts the uint64 value into 40bits-length []byte. The values in 25-64 bit are cut off.
// This function is used for the fields with 3 octets.
func Uint64To40(n uint64) []byte {
	return []byte{uint8(n >> 32), uint8(n >> 24), uint8(n >> 16), uint8(n >> 8), uint8(n)}
}

// EncodePLMN encodes MCC and MNC as BCD-encoded bytes.
func EncodePLMN(mcc, mnc string) ([]byte, error) {
	if len(mcc) != 3 || (len(mnc) != 2 && len(mnc) != 3) {
		return nil, fmt.Errorf("invalid PLMN length")
	}
	for _, digit := range mcc + mnc {
		if digit < '0' || digit > '9' {
			return nil, fmt.Errorf("invalid PLMN digit")
		}
	}

	c, err := StrToSwappedBytes(mcc, "f")
	if err != nil {
		return nil, err
	}
	n, err := StrToSwappedBytes(mnc, "f")
	if err != nil {
		return nil, err
	}

	// 2-digit
	b := make([]byte, 3)
	if len(mnc) == 2 {
		return append(c, n...), nil
	}

	// 3-digit
	b[0] = c[0]
	b[1] = (c[1] & 0x0f) | (n[1] << 4 & 0xf0)
	b[2] = n[0]

	return b, nil
}

// DecodePLMN decodes BCD-encoded bytes into MCC and MNC.
func DecodePLMN(b []byte) (mcc, mnc string, err error) {
	if len(b) < 3 {
		return "", "", io.ErrUnexpectedEOF
	}
	if len(b) != 3 {
		return "", "", fmt.Errorf("PLMN requires 3 octets")
	}
	// TS 29.244 V20.1.0 clause 8.2.108: decimal MCC/MNC digits;
	// only the third MNC digit may use the two-digit-MNC filler F.
	for pos := 0; pos < 6; pos++ {
		digit := b[pos/2] >> (4 * (pos % 2)) & 15
		if digit > 9 && (pos != 3 || digit != 15) {
			return "", "", fmt.Errorf("invalid PLMN digit")
		}
	}

	raw := hex.EncodeToString(b)
	mcc = string(raw[1]) + string(raw[0]) + string(raw[3])
	mnc = string(raw[5]) + string(raw[4])
	if string(raw[2]) != "f" {
		mnc += string(raw[2])
	}

	return
}

// ParseECI decodes ECI uint32 into e-NodeB ID and Cell ID.
func ParseECI(eci uint32) (enbID uint32, cellID uint8, err error) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, eci)
	cellID = buf[3]
	enbID = binary.BigEndian.Uint32([]byte{0, buf[0], buf[1], buf[2]})
	return
}

// EncodeFQDN encodes a PFCP length-delimited name without a root terminator.
// A trailing dot is accepted as an explicit root. Invalid names return nil:
// labels must contain 1..63 octets and the DNS name, including the implicit
// root octet, must fit in 255 octets (RFC 1035 section 3.1; RFC 2181 section 11).
func EncodeFQDN(fqdn string) []byte {
	if fqdn == "" || fqdn == "." {
		return []byte{0}
	}
	fqdn = strings.TrimSuffix(fqdn, ".")
	if len(fqdn)+2 > 255 {
		return nil
	}
	b := make([]byte, 0, len(fqdn)+1)
	for _, label := range strings.Split(fqdn, ".") {
		if len(label) == 0 || len(label) > 63 {
			return nil
		}
		b = append(b, byte(len(label)))
		b = append(b, label...)
	}
	return b
}

// DecodeFQDN decodes an uncompressed PFCP name. A malformed name returns the
// empty string, never a partial prefix. Compression pointers are not accepted.
// RFC 1035 section 3.1; RFC 2181 section 11.
func DecodeFQDN(b []byte) string {
	if len(b) > 255 {
		return ""
	}
	var labels []string
	for offset := 0; offset < len(b); {
		length := int(b[offset])
		offset++
		if length == 0 {
			if offset != len(b) {
				return ""
			}
			return strings.Join(labels, ".")
		}
		if length > 63 || length > len(b)-offset {
			return ""
		}
		labels = append(labels, string(b[offset:offset+length]))
		offset += length
	}
	if len(b) > 254 {
		return ""
	}
	return strings.Join(labels, ".")
}
