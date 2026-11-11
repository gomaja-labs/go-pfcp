// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "github.com/gomaja-labs/go-pfcp/internal/utils"

// NewNetworkInstance creates a new NetworkInstance IE.
func NewNetworkInstance(instance string) *IE {
	return newStringIE(NetworkInstance, instance)
}

// NewNetworkInstanceFQDN creates a new NetworkInstance IE from the given
// FQDN string.
func NewNetworkInstanceFQDN(fqdn string) *IE {
	return newFQDNIE(NetworkInstance, fqdn)
}

// NetworkInstance returns NetworkInstance in string if the type of IE matches.
func (i *IE) NetworkInstance() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	switch i.Type {
	case NetworkInstance:
		return i.ValueAsString()
	default:
		return "", &InvalidTypeError{Type: i.Type}
	}
}

// NetworkInstanceFQDN returns NetworkInstance in string if the type of IE matches.
// This is for the case that NetworkInstance is encoded as a FQDN.
func (i *IE) NetworkInstanceFQDN() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	// for checking the type
	value, err := i.NetworkInstance()
	if err != nil {
		return "", err
	}

	return decodeFQDNValue([]byte(value))
}

// NetworkInstanceHeuristic assumes that the payload is encoded in Name Syntax
// and returns the decoded string if it looks meaningful. Otherwise returns a
// string just converted from []byte.
func (i *IE) NetworkInstanceHeuristic() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}

	value, err := i.NetworkInstance()
	if err != nil {
		return "", err
	}
	if name := utils.DecodeFQDN([]byte(value)); name != "" {
		return name, nil
	}
	return value, nil
}
