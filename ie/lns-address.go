package ie

import "net/netip"

// NewLNSAddress creates an LNS Address IE.
// TS 29.244 V20.1.0 clause 8.2.188: the payload length alone selects IPv4
// (four octets) or IPv6 (sixteen octets); no extension octets can be recognised.
func NewLNSAddress(address netip.Addr) *IE {
	if !address.IsValid() || address.Zone() != "" {
		return newInvalidIE(LNSAddress, ErrInvalidValue)
	}
	return New(LNSAddress, address.AsSlice())
}

// LNSAddress returns the address selected by the exact payload length.
// TS 29.244 V20.1.0 clause 8.2.188. Other lengths return ErrMalformed;
// raw IE parsing and encoding still preserve the complete payload.
func (i *IE) LNSAddress() (netip.Addr, error) {
	b, e := addedValuePayload(i, LNSAddress, 0, false)
	if e != nil {
		return netip.Addr{}, e
	}
	if len(b) != 4 && len(b) != 16 {
		return netip.Addr{}, &DecodeError{Type: LNSAddress, Err: ErrMalformed}
	}
	a, _ := netip.AddrFromSlice(b)
	return a, nil
}
