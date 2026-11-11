package ie_test

import (
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestReportingSuggestionInfoCurrent(t *testing.T) {
	for _, v := range []ie.ReportingSuggestionInfoFields{{}, {ReportingTime: 0xffffffff * time.Second}, {Urgency: ie.ReportingNonDelayTolerant}} {
		b := []byte{byte(v.Urgency)}
		if v.Urgency == 0 {
			u := uint32(v.ReportingTime / time.Second)
			b = append(b, byte(u>>24), byte(u>>16), byte(u>>8), byte(u))
		}
		assertAddedFields(t, ie.NewReportingSuggestionInfo(v), &v, b, "ReportingSuggestionInfo", false)
	}

	var nilFields *ie.ReportingSuggestionInfoFields
	assertAddedNilFields(t, nilFields)
}

func TestReportingSuggestionInfoInvalidFields(t *testing.T) {
	for _, v := range []ie.ReportingSuggestionInfoFields{{Urgency: 2}, {Urgency: ie.ReportingNonDelayTolerant, ReportingTime: time.Second}, {ReportingTime: -1}, {ReportingTime: time.Nanosecond}, {ReportingTime: 0x100000000 * time.Second}} {
		assertAddedInvalid(t, ie.NewReportingSuggestionInfo(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}
