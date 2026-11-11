// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewCumulativeRateRatioThreshold creates a new CumulativeRateRatioThreshold IE.
func NewCumulativeRateRatioThreshold(threshold uint32) *IE {
	return newUint32ValIE(CumulativeRateRatioThreshold, threshold)
}

// CumulativeRateRatioThreshold returns CumulativeRateRatioThreshold in uint32 if the type of IE matches.
func (i *IE) CumulativeRateRatioThreshold() (uint32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case CumulativeRateRatioThreshold:
		return i.ValueAsUint32()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
