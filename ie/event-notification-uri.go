package ie

import (
	"net/netip"
	"strings"
)

// NewEventNotificationURI creates an Event Notification URI IE from the exact
// ASCII URI octets, or a null IE for an empty value.
// TS 29.244 V20.1.0 clause 8.2.202; RFC 3986 sections 2, 3 and Appendix A.
func NewEventNotificationURI(value string) *IE {
	if len(value) > 65535 || (len(value) > 0 && !addedURIValid([]byte(value))) {
		return newInvalidIE(EventNotificationURI, ErrInvalidValue)
	}
	return New(EventNotificationURI, []byte(value))
}

// EventNotificationURI returns the URI octets without normalizing their spelling.
// TS 29.244 V20.1.0 clause 8.2.202; RFC 3986 Appendix A.
func (i *IE) EventNotificationURI() (string, error) {
	b, err := addedValuePayload(i, EventNotificationURI, 1, false)
	if err != nil {
		return "", err
	}
	if !addedURIValid(b) {
		return "", &DecodeError{Type: EventNotificationURI, Err: ErrInvalidValue}
	}
	return string(b), nil
}

func addedURIAlpha(c byte) bool { return c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' }
func addedURIHex(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F'
}
func addedURIChars(s, extra string) bool {
	for n := 0; n < len(s); n++ {
		c := s[n]
		if addedURIAlpha(c) || c >= '0' && c <= '9' || strings.ContainsRune("-._~!$&'()*+,;="+extra, rune(c)) {
			continue
		}
		if c == '%' && n+2 < len(s) && addedURIHex(s[n+1]) && addedURIHex(s[n+2]) {
			n += 2
			continue
		}
		return false
	}
	return true
}
func addedURIAuthority(s string) bool {
	if at := strings.IndexByte(s, '@'); at >= 0 {
		if !addedURIChars(s[:at], ":") {
			return false
		}
		s = s[at+1:]
	}
	if strings.HasPrefix(s, "[") {
		end := strings.IndexByte(s, ']')
		if end < 0 {
			return false
		}
		host := s[1:end]
		rest := s[end+1:]
		if strings.HasPrefix(host, "v") || strings.HasPrefix(host, "V") {
			dot := strings.IndexByte(host, '.')
			if dot < 2 || dot == len(host)-1 {
				return false
			}
			for n := 1; n < dot; n++ {
				if !addedURIHex(host[n]) {
					return false
				}
			}
			// IPvFuture does not permit pct-encoded sequences.
			if strings.ContainsRune(host[dot+1:], '%') || !addedURIChars(host[dot+1:], ":") {
				return false
			}
		} else {
			a, err := netip.ParseAddr(host)
			if err != nil || !a.Is6() || a.Zone() != "" {
				return false
			}
		}
		if rest == "" {
			return true
		}
		if rest[0] != ':' {
			return false
		}
		s = rest[1:]
	} else {
		if colon := strings.IndexByte(s, ':'); colon >= 0 {
			if !addedURIChars(s[:colon], "") {
				return false
			}
			s = s[colon+1:]
		} else {
			return addedURIChars(s, "")
		}
	}
	for n := 0; n < len(s); n++ {
		if s[n] < '0' || s[n] > '9' {
			return false
		}
	}
	return true
}
func addedURIValid(value []byte) bool {
	s := string(value)
	colon := strings.IndexByte(s, ':')
	if colon < 1 || !addedURIAlpha(s[0]) {
		return false
	}
	for n := 1; n < colon; n++ {
		c := s[n]
		if !addedURIAlpha(c) && (c < '0' || c > '9') && c != '+' && c != '-' && c != '.' {
			return false
		}
	}
	rest := s[colon+1:]
	if at := strings.IndexByte(rest, '#'); at >= 0 {
		if !addedURIChars(rest[at+1:], ":@/?") {
			return false
		}
		rest = rest[:at]
	}
	if at := strings.IndexByte(rest, '?'); at >= 0 {
		if !addedURIChars(rest[at+1:], ":@/?") {
			return false
		}
		rest = rest[:at]
	}
	if strings.HasPrefix(rest, "//") {
		authority, path := rest[2:], ""
		if at := strings.IndexByte(authority, '/'); at >= 0 {
			path = authority[at:]
			authority = authority[:at]
		}
		return addedURIAuthority(authority) && addedURIChars(path, ":@/")
	}
	return addedURIChars(rest, ":@/")
}
