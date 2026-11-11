package ie

import (
	"bytes"
	"net/netip"
	"strings"
)

// NewURI creates a URI IE from its octet string, or a null IE for an empty value.
// TS 29.244 V20.1.0 clause 8.2.244; RFC 3986 sections 2–3 and appendix A.
// Generic URI syntax is validated without imposing scheme-specific structure
// (RFC 8820 sections 2–3). No normalization or percent-decoding is performed.
func NewURI(v []byte) *IE {
	if len(v) > 65535 || (len(v) > 0 && !validPFCPURI(string(v))) {
		return invalidRecentIE(URI, "URI syntax")
	}
	return New(URI, bytes.Clone(v))
}

// URI returns the URI octet string defined in TS 29.244 V20.1.0 clause 8.2.244.
func (i *IE) URI() ([]byte, error) {
	b, err := recentIEPayload(i, URI, 1, false)
	if err != nil {
		return nil, err
	}
	if !validPFCPURI(string(b)) {
		return nil, &DecodeError{Type: URI, Err: ErrInvalidValue}
	}
	return b, nil
}

// RFC 3986 appendix A: URI = scheme ":" hier-part [ "?" query ] [ "#" fragment ].
func validPFCPURI(s string) bool {
	scheme, rest, ok := strings.Cut(s, ":")
	if !ok || len(scheme) == 0 || !uriAlpha(scheme[0]) {
		return false
	}
	for j := 1; j < len(scheme); j++ {
		c := scheme[j]
		if !uriAlpha(c) && !uriDigit(c) && !strings.ContainsRune("+.-", rune(c)) {
			return false
		}
	}
	hier, fragment, hasFragment := strings.Cut(rest, "#")
	if hasFragment && !uriChars(fragment, "/:@?", true) {
		return false
	}
	hier, query, hasQuery := strings.Cut(hier, "?")
	if hasQuery && !uriChars(query, "/:@?", true) {
		return false
	}
	if strings.HasPrefix(hier, "//") {
		authority, path, _ := strings.Cut(hier[2:], "/")
		return uriAuthority(authority) && uriChars(path, "/:@", true)
	}
	return uriChars(hier, "/:@", true)
}

func uriAuthority(s string) bool {
	if user, host, ok := strings.Cut(s, "@"); ok {
		if !uriChars(user, ":", true) {
			return false
		}
		s = host
	}
	if strings.HasPrefix(s, "[") {
		literal, tail, ok := strings.Cut(s[1:], "]")
		if !ok {
			return false
		}
		if !uriIPLiteral(literal) {
			return false
		}
		if tail == "" {
			return true
		}
		if !strings.HasPrefix(tail, ":") {
			return false
		}
		return uriPort(tail[1:])
	}
	host, port, hasPort := strings.Cut(s, ":")
	return uriChars(host, "", true) && (!hasPort || uriPort(port))
}
func uriIPLiteral(s string) bool {
	if len(s) > 0 && (s[0] == 'v' || s[0] == 'V') {
		version, address, ok := strings.Cut(s[1:], ".")
		if !ok || len(version) == 0 || len(address) == 0 {
			return false
		}
		for j := 0; j < len(version); j++ {
			if !uriHex(version[j]) {
				return false
			}
		}
		return uriChars(address, ":", false)
	}
	ip, err := netip.ParseAddr(s)
	return err == nil && ip.Is6() && ip.Zone() == ""
}
func uriPort(s string) bool {
	for j := 0; j < len(s); j++ {
		if !uriDigit(s[j]) {
			return false
		}
	}
	return true
}
func uriChars(s, extra string, percent bool) bool {
	for j := 0; j < len(s); j++ {
		c := s[j]
		if uriAlpha(c) || uriDigit(c) || strings.ContainsRune("-._~!$&'()*+,;="+extra, rune(c)) {
			continue
		}
		if percent && c == '%' && j+2 < len(s) && uriHex(s[j+1]) && uriHex(s[j+2]) {
			j += 2
			continue
		}
		return false
	}
	return true
}
func uriAlpha(c byte) bool { return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' }
func uriDigit(c byte) bool { return c >= '0' && c <= '9' }
func uriHex(c byte) bool   { return uriDigit(c) || c >= 'A' && c <= 'F' || c >= 'a' && c <= 'f' }
