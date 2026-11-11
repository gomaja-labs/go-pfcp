// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewCumulativeRateRatioMeasurement creates a new CumulativeRateRatioMeasurement IE.
func NewCumulativeRateRatioMeasurement(measurement int32) *IE {
	return newUint32ValIE(CumulativeRateRatioMeasurement, uint32(measurement))
}

// CumulativeRateRatioMeasurement returns CumulativeRateRatioMeasurement in int32 if the type of IE matches.
func (i *IE) CumulativeRateRatioMeasurement() (int32, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case CumulativeRateRatioMeasurement:
		// TS 29.244 V20.1.0 clause 8.2.150: cumulativeRateRatio is Integer32.
		v, err := i.ValueAsUint32()
		return int32(v), err
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
