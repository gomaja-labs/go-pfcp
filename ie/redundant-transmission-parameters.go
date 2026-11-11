// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewRedundantTransmissionParametersInPDI creates a new RedundantTransmissionParameters IE.
func NewRedundantTransmissionParametersInPDI(fteid, ni *IE) *IE {
	return newGroupedIE(RedundantTransmissionParameters, 0, fteid, ni)
}

// NewRedundantTransmissionParametersInFAR creates a new RedundantTransmissionParameters IE.
func NewRedundantTransmissionParametersInFAR(ohc, ni *IE) *IE {
	return newGroupedIE(RedundantTransmissionParameters, 0, ohc, ni)
}
