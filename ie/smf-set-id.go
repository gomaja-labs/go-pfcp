// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewSMFSetID creates a new SMFSetID IE.
func NewSMFSetID(id string) *IE {
	// TS 29.244 V20.1.0 clause 8.2.131: a spare octet precedes the FQDN.
	i := newFQDNIE(SMFSetID, id)
	if i.err == nil {
		i.Payload = append([]byte{0}, i.Payload...)
	}
	return i
}

// SMFSetID returns SMFSetID in string if the type of IE matches.
func (i *IE) SMFSetID() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	if i.Type != SMFSetID {
		return "", &InvalidTypeError{Type: i.Type}
	}

	if len(i.Payload) < 2 {
		return "", io.ErrUnexpectedEOF
	}
	return decodeFQDNValue(i.Payload[1:])
}
