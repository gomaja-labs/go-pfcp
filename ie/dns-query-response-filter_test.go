package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestDNSQueryResponseFilterCurrent(t *testing.T) {
	for _, value := range [][]byte{[]byte(`{"regex":""}`), []byte(`{"regex":".*example"}`)} {
		v := ie.DNSQueryResponseFilterFields{DomainNamePattern: string(value)}
		b := []byte{byte(len(value) >> 8), byte(len(value))}
		b = append(b, value...)
		assertAddedFields(t, ie.NewDNSQueryResponseFilter(v), &v, b, "DNSQueryResponseFilter", false)
	}

	for _, value := range []string{`{}`, `null`, `[]`, `{"regex":null}`, `{"regex":1}`, `{"regex":"a","stringMatchingRule":{}}`, `{"stringMatchingRule":{"stringMatchingConditions":[]}}`} {
		assertAddedInvalid(t, ie.NewDNSQueryResponseFilter(ie.DNSQueryResponseFilterFields{DomainNamePattern: value}))
	}
	var nilFields *ie.DNSQueryResponseFilterFields
	assertAddedNilFields(t, nilFields)
}

func TestDNSQueryResponseFilterInvalidUTF8(t *testing.T) {
	b := []byte{'{', '"', 'r', 'e', 'g', 'e', 'x', '"', ':', '"', 255, '"', '}'}
	assertAddedInvalid(t, ie.NewDNSQueryResponseFilter(ie.DNSQueryResponseFilterFields{DomainNamePattern: string(b)}))
	wire := append([]byte{0, byte(len(b))}, b...)
	if _, err := ie.New(ie.DNSQueryResponseFilter, wire).DNSQueryResponseFilter(); err == nil {
		t.Fatal("malformed UTF8 JSON accepted")
	}
}

func TestDNSQueryResponseFilterLengthBoundary(t *testing.T) {
	b := append([]byte(`{"regex":"`), bytes.Repeat([]byte{65}, 65521)...)
	b = append(b, []byte(`"}`)...)
	if len(b) != 65533 {
		t.Fatal("boundary vector length")
	}
	v := ie.DNSQueryResponseFilterFields{DomainNamePattern: string(b)}
	payload := append([]byte{255, 253}, b...)
	assertAddedFields(t, ie.NewDNSQueryResponseFilter(v), &v, payload, "DNSQueryResponseFilter", false)
	b = append([]byte(" "), b...)
	assertAddedInvalid(t, ie.NewDNSQueryResponseFilter(ie.DNSQueryResponseFilterFields{DomainNamePattern: string(b)}))
}

var _ = registerConstructorFixture(ie.NewDNSQueryResponseFilter, ie.DNSQueryResponseFilterFields{DomainNamePattern: `{"regex":"example"}`})
