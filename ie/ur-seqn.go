// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewURSEQN creates a new URSEQN IE.
func NewURSEQN(seq uint32) *IE {
	return newUint32ValIE(URSEQN, seq)
}

// URSEQN returns URSEQN in uint32 if the type of IE matches.
func (i *IE) URSEQN() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case URSEQN:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
