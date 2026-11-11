// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewPacketReplicationAndDetectionCarryOnInformation creates a new PacketReplicationAndDetectionCarryOnInformation IE.
func NewPacketReplicationAndDetectionCarryOnInformation(flag uint8) *IE {
	return newUint8ValIE(PacketReplicationAndDetectionCarryOnInformation, flag)
}

// PacketReplicationAndDetectionCarryOnInformation returns PacketReplicationAndDetectionCarryOnInformation in uint8 if the type of IE matches.
func (i *IE) PacketReplicationAndDetectionCarryOnInformation() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case PacketReplicationAndDetectionCarryOnInformation:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}

// HasPRIUEAI reports whether an IE has PRIUEAI bit.
func (i *IE) HasPRIUEAI() bool {
	if i == nil {
		return false
	}

	v, err := i.PacketReplicationAndDetectionCarryOnInformation()
	if err != nil {
		return false
	}

	return has1stBit(v)
}

// HasPRINT19I reports whether an IE has PRINT19I bit.
func (i *IE) HasPRINT19I() bool {
	if i == nil {
		return false
	}

	v, err := i.PacketReplicationAndDetectionCarryOnInformation()
	if err != nil {
		return false
	}

	return has2ndBit(v)
}

// HasPRIN6I reports whether an IE has PRIN6I bit.
func (i *IE) HasPRIN6I() bool {
	if i == nil {
		return false
	}

	v, err := i.PacketReplicationAndDetectionCarryOnInformation()
	if err != nil {
		return false
	}

	return has3rdBit(v)
}

// HasDCARONI reports whether an IE has DCARONI bit.
func (i *IE) HasDCARONI() bool {
	if i == nil {
		return false
	}

	v, err := i.PacketReplicationAndDetectionCarryOnInformation()
	if err != nil {
		return false
	}

	return has4thBit(v)
}
