// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewTSCManagementInformation creates a new TSCManagementInformation IE.
func NewTSCManagementInformation(typ IEType, ies ...*IE) *IE {
	return newGroupedIE(typ, 0, ies...)
}

// NewTSCManagementInformationWithinSessionModificationRequest creates a new TSCManagementInformationWithinSessionModificationRequest IE.
func NewTSCManagementInformationWithinSessionModificationRequest(ies ...*IE) *IE {
	return newGroupedIE(TSCManagementInformationWithinSessionModificationRequest, 0, ies...)
}

// NewTSCManagementInformationWithinSessionModificationResponse creates a new TSCManagementInformationWithinSessionModificationResponse IE.
func NewTSCManagementInformationWithinSessionModificationResponse(ies ...*IE) *IE {
	return newGroupedIE(TSCManagementInformationWithinSessionModificationResponse, 0, ies...)
}

// NewTSCManagementInformationWithinSessionReportRequest creates a new TSCManagementInformationWithinSessionReportRequest IE.
func NewTSCManagementInformationWithinSessionReportRequest(ies ...*IE) *IE {
	return newGroupedIE(TSCManagementInformationWithinSessionReportRequest, 0, ies...)
}

// TSCManagementInformation returns the IEs above TSCManagementInformation if the type of IE matches.
func (i *IE) TSCManagementInformation() ([]*IE, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case TSCManagementInformationWithinSessionModificationRequest,
		TSCManagementInformationWithinSessionModificationResponse,
		TSCManagementInformationWithinSessionReportRequest:

		return i.ValueAsGrouped()
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
