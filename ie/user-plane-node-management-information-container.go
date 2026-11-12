// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// TS 29.244 V20.1.0 clause 8.2.182.

// TS 29.244 V20.1.0 clause 8.2.182: the content is a management
// message defined by TS 24.539 clause 9.
// NewUserPlaneNodeManagementInformationContainer creates a new UserPlaneNodeManagementInformationContainer IE.
func NewUserPlaneNodeManagementInformationContainer(info []byte) *IE {
	return New(UserPlaneNodeManagementInformationContainer, info)
}

// UserPlaneNodeManagementInformationContainer returns UserPlaneNodeManagementInformationContainer as opaque bytes if the type of IE matches.
func (i *IE) UserPlaneNodeManagementInformationContainer() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	switch i.Type {
	case UserPlaneNodeManagementInformationContainer:
		if len(i.Payload) == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		return i.Payload, nil
	default:
		return nil, &InvalidTypeError{Type: i.Type}
	}
}
