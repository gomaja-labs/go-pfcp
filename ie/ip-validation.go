package ie

import "net"

// Address presence and family follow the selecting flag, never the slice length.
// TS 29.244 V20.1.0 clause 8.1.2 and each caller's address-field clause.
func ipv4FieldValid(ip net.IP, present bool) bool {
	if !present {
		return ip == nil
	}
	return ip.To4() != nil
}
func ipv6FieldValid(ip net.IP, present bool) bool {
	if !present {
		return ip == nil
	}
	return ip.To16() != nil && ip.To4() == nil
}

// Keep a malformed nonempty argument distinguishable from an absent argument,
// so the fields encoder and IE constructor both report ErrInvalidValue.
func parseIPArgument(value string) net.IP {
	if value == "" {
		return nil
	}
	if ip := net.ParseIP(value); ip != nil {
		if v4 := ip.To4(); v4 != nil {
			return v4
		}
		return ip
	}
	return net.IP{}
}
