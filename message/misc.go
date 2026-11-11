// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package message

func uint24To32(b []byte) uint32 {
	if len(b) != 3 {
		return 0
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

func uint32To24(n uint32) []byte {
	return []byte{uint8(n >> 16), uint8(n >> 8), uint8(n)}
}
