// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewSuggestedBufferingPacketsCount creates a new SuggestedBufferingPacketsCount IE.
func NewSuggestedBufferingPacketsCount(count uint8) *IE {
	return newUint8ValIE(SuggestedBufferingPacketsCount, count)
}

// SuggestedBufferingPacketsCount returns SuggestedBufferingPacketsCount in uint8 if the type of IE matches.
func (i *IE) SuggestedBufferingPacketsCount() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case SuggestedBufferingPacketsCount:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
