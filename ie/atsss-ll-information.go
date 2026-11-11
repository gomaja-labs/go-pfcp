// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewATSSSLLInformation creates a new ATSSSLLInformation IE.
func NewATSSSLLInformation(lli uint8) *IE {
	return newUint8ValIE(ATSSSLLInformation, lli&0x01)
}

// ATSSSLLInformation returns ATSSSLLInformation in uint8 if the type of IE matches.
func (i *IE) ATSSSLLInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case ATSSSLLInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
