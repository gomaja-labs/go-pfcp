package ie

import "net/netip"

// NewDNSServerAddress creates a DNS Server Address IE.
// TS 29.244 V20.1.0 clause 8.2.193.
func NewDNSServerAddress(address netip.Addr) *IE {
	if !address.Is4() || address.Zone() != "" {
		return newInvalidIE(DNSServerAddress, ErrInvalidValue)
	}
	return New(DNSServerAddress, address.AsSlice())
}

// DNSServerAddress returns the IPv4 server address, ignoring trailing extension octets.
// TS 29.244 V20.1.0 clauses 8.2.193 and 8.1.2.
func (i *IE) DNSServerAddress() (netip.Addr, error) {
	b, err := addedValuePayload(i, DNSServerAddress, 4, false)
	if err != nil {
		return netip.Addr{}, err
	}
	return netip.AddrFrom4([4]byte(b[:4])), nil
}
