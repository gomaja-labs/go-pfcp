package ie

import "net/netip"

// NewNBNSServerAddress creates a NBNS Server Address IE.
// TS 29.244 V20.1.0 clause 8.2.194.
func NewNBNSServerAddress(address netip.Addr) *IE {
	if !address.Is4() || address.Zone() != "" {
		return newInvalidIE(NBNSServerAddress, ErrInvalidValue)
	}
	return New(NBNSServerAddress, address.AsSlice())
}

// NBNSServerAddress returns the IPv4 server address, ignoring trailing extension octets.
// TS 29.244 V20.1.0 clauses 8.2.194 and 8.1.2.
func (i *IE) NBNSServerAddress() (netip.Addr, error) {
	b, err := addedValuePayload(i, NBNSServerAddress, 4, false)
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.AddrFrom4([4]byte(b[:4])), nil
}
