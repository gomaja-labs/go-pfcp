package ie_test

import (
	"bytes"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestVendorSpecificNodeReportTypeCurrent(t *testing.T) {
	for _, v := range []ie.VendorSpecificNodeReportTypeFields{{EnterpriseID: 0, ReportTypes: []byte{1}}, {EnterpriseID: 65535, ReportTypes: []byte{0, 128}}} {
		b := []byte{byte(v.EnterpriseID >> 8), byte(v.EnterpriseID)}
		b = append(b, v.ReportTypes...)
		assertAddedFields(t, ie.NewVendorSpecificNodeReportType(v), &v, b, "VendorSpecificNodeReportType", false)
	}

	var nilFields *ie.VendorSpecificNodeReportTypeFields
	assertAddedNilFields(t, nilFields)
}

func TestVendorSpecificNodeReportTypeInvalidFields(t *testing.T) {
	for _, v := range []ie.VendorSpecificNodeReportTypeFields{{}, {ReportTypes: []byte{0}}, {ReportTypes: []byte{0, 0}}} {
		assertAddedInvalid(t, ie.NewVendorSpecificNodeReportType(v))
		if _, err := v.Marshal(); err == nil {
			t.Fatal("invalid fields accepted")
		}
	}
}

func TestVendorSpecificNodeReportTypeLengthBoundary(t *testing.T) {
	v := ie.VendorSpecificNodeReportTypeFields{EnterpriseID: 65535, ReportTypes: bytes.Repeat([]byte{1}, 65533)}
	b := append([]byte{255, 255}, v.ReportTypes...)
	assertAddedFields(t, ie.NewVendorSpecificNodeReportType(v), &v, b, "VendorSpecificNodeReportType", false)
	v.ReportTypes = append(v.ReportTypes, 1)
	assertAddedInvalid(t, ie.NewVendorSpecificNodeReportType(v))
}
