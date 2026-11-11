// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Rule ID Type definitions.
const (
	RuleIDTypePDR uint8 = 0 // 16
	RuleIDTypeFAR uint8 = 1 // 32
	RuleIDTypeQER uint8 = 2 // 32
	RuleIDTypeURR uint8 = 3 // 32
	RuleIDTypeBAR uint8 = 4 // 8
	RuleIDTypeMAR uint8 = 5 // 16
	RuleIDTypeSRR uint8 = 6 // 8 (clause 8.2.151)
)

// NewFailedRuleID creates a new FailedRuleID IE.
func NewFailedRuleID(typ uint8, id uint32) *IE {
	// TS 29.244 V20.1.0 clause 8.2.80, Table 8.2.80-1.
	if typ > RuleIDTypeSRR || ((typ == RuleIDTypePDR || typ == RuleIDTypeMAR) && id > 0xffff) || ((typ == RuleIDTypeBAR || typ == RuleIDTypeSRR) && id > 0xff) {
		return newInvalidIE(FailedRuleID, fmt.Errorf("%w: rule type %d or ID %d is not encodable", ErrInvalidValue, typ, id))
	}
	switch typ {
	case RuleIDTypePDR, RuleIDTypeMAR:
		b := make([]byte, 3)
		b[0] = typ
		binary.BigEndian.PutUint16(b[1:3], uint16(id))
		return New(FailedRuleID, b)
	case RuleIDTypeFAR, RuleIDTypeQER, RuleIDTypeURR:
		b := make([]byte, 5)
		b[0] = typ
		binary.BigEndian.PutUint32(b[1:5], id)
		return New(FailedRuleID, b)
	case RuleIDTypeBAR, RuleIDTypeSRR:
		return New(FailedRuleID, []byte{typ, uint8(id)})
	default:
		return New(FailedRuleID, []byte{typ})
	}
}

// RuleIDType returns RuleIDType in uint8 if the type of IE matches.
func (i *IE) RuleIDType() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != FailedRuleID {
		return 0, &InvalidTypeError{Type: i.Type}
	}

	v, err := i.valueAsUint8Masked(0x1f)
	if v > RuleIDTypeSRR {
		v = RuleIDTypeFAR
	}
	return v, err
}

// FailedRuleID returns FailedRuleID in uint32 if the type of IE matches.
func (i *IE) FailedRuleID() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != FailedRuleID {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	if len(i.Payload) < 1 {
		return 0, io.ErrUnexpectedEOF
	}

	typ, err := i.RuleIDType()
	if err != nil {
		return 0, err
	}
	switch typ {
	case RuleIDTypePDR, RuleIDTypeMAR:
		if len(i.Payload) < 3 {
			return 0, io.ErrUnexpectedEOF
		}
		return uint32(binary.BigEndian.Uint16(i.Payload[1:3])), nil
	case RuleIDTypeFAR, RuleIDTypeQER, RuleIDTypeURR:
		if len(i.Payload) < 5 {
			return 0, io.ErrUnexpectedEOF
		}
		return binary.BigEndian.Uint32(i.Payload[1:5]), nil
	case RuleIDTypeBAR, RuleIDTypeSRR:
		if len(i.Payload) < 2 {
			return 0, io.ErrUnexpectedEOF
		}
		return uint32(i.Payload[1]), nil
	default:
		return 0, nil
	}
}
