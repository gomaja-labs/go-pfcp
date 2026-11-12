package ie_test

import (
	"errors"
	"fmt"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func ExampleEncodeError() {
	// Four-octet timestamps cannot represent this date in the package's era window.
	invalid := ie.NewRecoveryTimeStamp(time.Date(1969, time.December, 31, 0, 0, 0, 0, time.UTC))
	_, err := invalid.Marshal()
	var encoding *ie.EncodeError
	if !errors.As(err, &encoding) {
		fmt.Println("expected encoding error:", err)
		return
	}
	fmt.Println("type:", encoding.Type)
	fmt.Println("invalid value:", errors.Is(err, ie.ErrInvalidValue))
	// Output:
	// type: RecoveryTimeStamp
	// invalid value: true
}

func ExampleIE_ValueAsGrouped() {
	vendor := ie.NewVendorSpecificGroupedIE(0x8001, 32473, ie.NewPDRID(7))
	wire, err := vendor.Marshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	decoded, err := ie.Parse(wire)
	if err != nil {
		fmt.Println(err)
		return
	}
	children, err := decoded.ValueAsGrouped()
	if err != nil {
		fmt.Println(err)
		return
	}
	id, err := children[0].PDRID()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("enterprise:", decoded.EnterpriseID, "PDR:", id)
	// Output: enterprise: 32473 PDR: 7
}
