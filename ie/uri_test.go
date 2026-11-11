package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestURI(t *testing.T) {
	// RFC 3986 sections 1.1.2, 3 and appendix A, not scheme-specific normalization.
	for _, uri := range []string{"foo://example.com:8042/over/there?name=ferret#nose", "urn:example:animal:ferret:nose", "ftp://ftp.is.co.za/rfc/rfc1808.txt", "http://www.ietf.org/rfc/rfc2396.txt", "ldap://[2001:db8::7]/c=GB?objectClass?one", "mailto:John.Doe@example.com", "news:comp.infosystems.www.servers.unix", "tel:+1-816-555-1212", "telnet://192.0.2.16:80/", "urn:oasis:names:specification:docbook:dtd:xml:4.1.2", "https://u:p@host:/a%20b?#", "a:", "a:/", "a://", "a://[vF.ab:c]/", "a://[::ffff:192.0.2.1]/", "a:?x/?:@!$&'()*+,;=#y/?:@!$&'()*+,;="} {
		i := ie.NewURI([]byte(uri))
		checkNewIEWire(t, i, []byte(uri))
		got, err := i.URI()
		if err != nil || !bytes.Equal(got, []byte(uri)) {
			t.Fatal(uri, err)
		}
	}
	for _, uri := range []string{"relative", "//host/path", "1a:x", "a:%", "a:%xy", "a:b c", "a:b\\c", "a:é", "a://[zone%25eth0]/", "a://[xyz]/", "a://[::1", "a://[::1]extra/", "a://user@@host/", "a://host:abc/", "a:a#b#c", "a:a[b]", "a://[v1.]/", "a://a^b/"} {
		assertNewIEInvalid(t, ie.NewURI([]byte(uri)), ie.URI)
		if _, err := ie.New(ie.URI, []byte(uri)).URI(); err == nil {
			t.Fatal("invalid URI accepted", uri)
		}
	}
	checkNewIEFailures(t, ie.URI, func(i *ie.IE) (any, error) { return i.URI() }, 1)
	checkNewIEWire(t, ie.NewURI(nil), nil)
	assertNewIEInvalid(t, ie.NewURI(bytes.Repeat([]byte{'a'}, 65536)), ie.URI)
}
