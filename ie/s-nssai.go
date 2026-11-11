// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

// NewSNSSAI creates a new SNSSAI IE.
func NewSNSSAI(sst uint8, sd uint32) *IE {
	i := New(SNSSAI, make([]byte, 4))
	i.Payload[0] = sst
	copy(i.Payload[1:4], utils.Uint32To24(sd))
	return i
}

// SNSSAI returns SNSSAI in []byte if the type of IE matches.
func (i *IE) SNSSAI() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case SNSSAI:
		if err := i.validateFixedValue(); err != nil {
			return nil, err
		}
		if len(i.Payload) < 4 {
			return nil, io.ErrUnexpectedEOF
		}

		return i.Payload[0:4], nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// SST returns SST in uint8 if the type of IE matches.
func (i *IE) SST() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	v, err := i.SNSSAI()
	if err != nil {
		return 0, err
	}

	return v[0], nil
}

// SD returns SD in uint32 if the type of IE matches.
func (i *IE) SD() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	v, err := i.SNSSAI()
	if err != nil {
		return 0, err
	}

	return utils.Uint24To32(v[1:4]), nil
}
