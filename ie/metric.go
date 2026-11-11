// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewMetric creates a new Metric IE.
func NewMetric(metric uint8) *IE {
	return newUint8ValIE(Metric, metric)
}

// Metric returns Metric in uint8 if the type of IE matches.
func (i *IE) Metric() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	switch i.Type {
	case Metric:
		return i.ValueAsUint8()
	default:
		return 0, &InvalidTypeError{Type: i.Type}
	}
}
