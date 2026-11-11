// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// TS 29.244 V20.1.0 clause 8.2.182.

// NewUserPlaneNodeManagementInformationContainer creates a new UserPlaneNodeManagementInformationContainer IE.
func NewUserPlaneNodeManagementInformationContainer(info string) *IE {
	return newStringIE(UserPlaneNodeManagementInformationContainer, info)
}

// UserPlaneNodeManagementInformationContainer returns UserPlaneNodeManagementInformationContainer in string if the type of IE matches.
func (i *IE) UserPlaneNodeManagementInformationContainer() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	switch i.Type {
	case UserPlaneNodeManagementInformationContainer:
		return i.ValueAsString()
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}
