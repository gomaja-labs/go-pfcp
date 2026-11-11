// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// NewDataNetworkAccessIdentifier creates a new DataNetworkAccessIdentifier IE.
func NewDataNetworkAccessIdentifier(id string) *IE {
	return newStringIE(DataNetworkAccessIdentifier, id)
}

// DataNetworkAccessIdentifier returns DataNetworkAccessIdentifier in string if the type of IE matches.
func (i *IE) DataNetworkAccessIdentifier() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	switch i.Type {
	case DataNetworkAccessIdentifier:
		return i.ValueAsString()
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}
