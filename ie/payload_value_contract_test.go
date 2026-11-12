package ie_test

import (
	"bytes"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestPayloadValueLengths(t *testing.T) {
	for _, tc := range []struct {
		name   string
		value  *ie.IE
		helper func(*ie.IE) ([]byte, error)
	}{
		{"downlink", ie.NewDownlinkDataServiceInformation(true, true, 1, 2), func(i *ie.IE) ([]byte, error) {
			if _, err := i.DownlinkDataServiceInformation(); err != nil {
				return nil, err
			}
			return i.Payload, nil
		}},
		{"flow", ie.NewFlowInformation(ie.FlowDirectionUplink, "abc"), (*ie.IE).FlowInformation},
		{"policy", ie.NewForwardingPolicy("abc"), (*ie.IE).ForwardingPolicy},
		{"time quota", ie.NewTimeQuotaMechanism(ie.BTITCTP, time.Second), (*ie.IE).TimeQuotaMechanism},
		{"pool", ie.NewUEIPAddressPoolIdentity("abc"), (*ie.IE).UEIPAddressPoolIdentity},
		{"features", ie.NewUPFunctionFeatures(1, 2), (*ie.IE).UPFunctionFeatures},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.helper(tc.value)
			if err != nil || !bytes.Equal(got, tc.value.Payload) {
				t.Fatalf("valid: %v", err)
			}
			for n := 0; n < len(tc.value.Payload); n++ {
				_, err := tc.helper(ie.New(tc.value.Type, tc.value.Payload[:n]))
				if !errors.Is(err, io.ErrUnexpectedEOF) {
					t.Errorf("prefix %d: %v", n, err)
				}
			}
			extension := append(append([]byte(nil), tc.value.Payload...), 0xaa)
			got, err = tc.helper(ie.New(tc.value.Type, extension))
			if err != nil || !bytes.Equal(got, extension) {
				t.Fatalf("extension: %v", err)
			}
		})
	}
}

func TestFlowDirectionReservedValues(t *testing.T) {
	// TS 29.244 V20.1.0 Table 8.2.61-1 maps received spare values 4..7 to 0.
	for raw := uint8(0); raw < 8; raw++ {
		want := raw
		if want > 3 {
			want = ie.FlowDirectionUnspecified
		}
		got, err := ie.New(ie.FlowInformation, []byte{raw | 0xf8, 0, 0}).FlowDirection()
		if err != nil || got != want {
			t.Errorf("direction %d: got %d, %v; want %d", raw, got, err, want)
		}
	}
}
