package ie_test

import (
	"errors"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestHelperSemanticAllocation(t *testing.T) {
	// TS 29.244 V20.1.0 clause 8.2.54: zero denotes dynamic CP
	// provisioning; one denotes a rule predefined in the UP function.
	for _, tc := range []struct {
		name string
		call func(*ie.IE) bool
		cp   bool
	}{
		{"IsAllocatedByCPFunction", (*ie.IE).IsAllocatedByCPFunction, true},
		{"IsAllocatedByUPFunction", (*ie.IE).IsAllocatedByUPFunction, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for label, input := range map[string]*ie.IE{"nil": nil, "wrong-type": ie.NewCause(1), "null": ie.New(ie.URRID, nil)} {
				t.Run(label, func(t *testing.T) {
					if tc.call(input) {
						t.Fatal("reported allocation on absent or wrong value")
					}
				})
			}
			for n := 1; n < 4; n++ {
				if tc.call(ie.New(ie.URRID, make([]byte, n))) {
					t.Errorf("reported allocation on truncated length %d", n)
				}
			}
			for _, high := range []bool{false, true} {
				id := uint32(0x12345678)
				if high {
					id |= 0x80000000
				}
				if got, want := tc.call(ie.NewURRID(id)), tc.cp != high; got != want {
					t.Errorf("high bit=%v result=%v; want %v", high, got, want)
				}
			}
		})
	}
}

func TestHelperSemanticMustAPNDNN(t *testing.T) {
	if got := ie.NewAPNDNN("abc.example").MustAPNDNN(); got != "abc.example" {
		t.Fatalf("value=%q", got)
	}
	for label, input := range map[string]*ie.IE{"nil": nil, "wrong-type": ie.NewCause(1), "null": ie.New(ie.APNDNN, nil), "truncated": ie.New(ie.APNDNN, []byte{3, 'a'})} {
		t.Run(label, func(t *testing.T) {
			if got := input.MustAPNDNN(); got != "" {
				t.Fatalf("value=%q; want empty", got)
			}
		})
	}
}

func TestHelperSemanticValidateApplyAction(t *testing.T) {
	// TS 29.244 V20.1.0 clause 8.2.26: FORW is a valid action;
	// null has no action and cannot supply a value to the validator.
	if err := ie.NewApplyAction(2).ValidateApplyAction(); err != nil {
		t.Fatal(err)
	}
	var absent *ie.IE
	if err := absent.ValidateApplyAction(); !errors.Is(err, ie.ErrIENotFound) {
		t.Fatal(err)
	}
	var wrong *ie.InvalidTypeError
	if err := ie.NewCause(1).ValidateApplyAction(); !errors.As(err, &wrong) {
		t.Fatal(err)
	}
	if err := ie.New(ie.ApplyAction, nil).ValidateApplyAction(); err == nil {
		t.Fatal("accepted null action")
	}
}
