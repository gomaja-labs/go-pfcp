// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import (
	"io"

	"github.com/gomaja-labs/go-pfcp/internal/utils"
)

// NewGBR creates a new GBR IE.
func NewGBR(ul, dl uint64) *IE {
	i := New(GBR, make([]byte, 10))
	copy(i.Payload[0:5], utils.Uint64To40(ul))
	copy(i.Payload[5:10], utils.Uint64To40(dl))
	return i
}

// GBR returns GBR in []byte if the type of IE matches.
func (i *IE) GBR() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case GBR:
		if len(i.Payload) < 10 {
			return nil, io.ErrUnexpectedEOF
		}

		return i.Payload, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}

// GBRUL returns GBRUL in uint64 if the type of IE matches.
func (i *IE) GBRUL() (uint64, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	v, err := i.GBR()
	if err != nil {
		return 0, err
	}
	return utils.Uint40To64(v[0:5]), nil
}

// GBRDL returns GBRDL in uint64 if the type of IE matches.
func (i *IE) GBRDL() (uint64, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	v, err := i.GBR()
	if err != nil {
		return 0, err
	}
	return utils.Uint40To64(v[5:10]), nil
}
