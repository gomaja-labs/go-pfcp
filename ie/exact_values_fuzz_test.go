package ie_test

import (
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func FuzzExactValueEncodings(f *testing.F) {
	for _, seed := range [][]byte{nil, {0}, {1}, {0x45}, {0x28}, {0xff}, {1, 2, 3}, {1, 0, 0, 0, 0, 0, 0, 0, 1}} {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, b []byte) {
		if len(b) > 65535 {
			return
		}
		for _, tc := range timerCases() {
			value, err := tc.get(ie.New(tc.typ, b))
			if err != nil {
				continue
			}
			encoded := tc.new(value)
			if _, err := encoded.Marshal(); err != nil {
				t.Fatal(err)
			}
			again, err := tc.get(encoded)
			if err != nil || again != value {
				t.Fatal("timer changed its decoded duration")
			}
		}
		item := ie.New(ie.DownlinkDataServiceInformation, b)
		if value, err := item.PPI(); err == nil {
			if value > 63 {
				t.Fatal("PPI exceeds DSCP width")
			}
			fields, err := item.DownlinkDataServiceInformation()
			if err == nil && fields.PagingPolicyIndication != value {
				t.Fatal("PPI accessors disagree")
			}
		}
		_, err := ie.New(ie.DLBufferingSuggestedPacketCount, b).DLBufferingSuggestedPacketCount()
		if len(b) > 2 && !errors.Is(err, ie.ErrInvalidLength) {
			t.Fatal("excess packet count octets accepted")
		}
		fields, err := ie.ParseFSEIDFields(b)
		if err == nil {
			encoded, err := fields.Marshal()
			if err == nil {
				parsed, err := ie.ParseFSEIDFields(encoded)
				if err != nil || parsed.SEID != fields.SEID || parsed.Flags != fields.Flags || !parsed.IPv4Address.Equal(fields.IPv4Address) || !parsed.IPv6Address.Equal(fields.IPv6Address) {
					t.Fatal("F-SEID changed its decoded value")
				}
			}
		}
	})
}
