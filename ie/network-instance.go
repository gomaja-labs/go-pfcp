// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "github.com/gomaja-labs/go-pfcp/internal/utils"

// NewNetworkInstance creates a new NetworkInstance IE.
// The value is an octet string carried unchanged; no character set is imposed or checked.
// TS 29.244 V20.1.0 clause 8.2.4: no text repertoire is specified for this identifier.
func NewNetworkInstance(instance string) *IE {
	return newStringIE(NetworkInstance, instance)
}

// NewNetworkInstanceFQDN creates a new NetworkInstance IE from the given
// FQDN string.
func NewNetworkInstanceFQDN(fqdn string) *IE {
	return newFQDNIE(NetworkInstance, fqdn)
}

// NetworkInstance returns the unchanged octet string.
func (i *IE) NetworkInstance() (string, error) {
	if i == nil {
		return "", ErrIENotFound
	}
	if i.Type != NetworkInstance {
		return "", &InvalidTypeError{Type: i.Type}
	}
	b, err := i.payloadAtLeast(1)
	return string(b), err
}

// NetworkInstanceFQDN returns NetworkInstance as a decoded domain name.
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
