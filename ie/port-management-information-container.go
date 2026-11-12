// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// TS 29.244 V20.1.0 clause 8.2.144: the content is a management
// message defined by TS 24.539 clause 8.
// NewPortManagementInformationContainer creates a new PortManagementInformationContainer IE.
func NewPortManagementInformationContainer(info []byte) *IE {
	return New(PortManagementInformationContainer, info)
}

// PortManagementInformationContainer returns PortManagementInformationContainer as opaque bytes if the type of IE matches.
func (i *IE) PortManagementInformationContainer() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case PortManagementInformationContainer:
		if len(i.Payload) == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		return i.Payload, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
