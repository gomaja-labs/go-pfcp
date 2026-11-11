package ie_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestULPeriodicityCurrent(t *testing.T) {
	for _, v := range []time.Duration{0, time.Millisecond, 0x12345678 * time.Millisecond, 0xffffffff * time.Millisecond} {
		b := make([]byte, 5)
		binary.BigEndian.PutUint32(b[1:], uint32(v/time.Millisecond))
		assertAddedVector(t, ie.NewULPeriodicity(v), b, "ULPeriodicity", false)
		got, err := ie.NewULPeriodicity(v).ULPeriodicity()
		if err != nil || got != v {
			t.Fatalf("duration %v: %v %v", v, got, err)
		}
		b[0] = 255
		got, err = ie.New(ie.ULPeriodicity, b).ULPeriodicity()
		if err != nil || got != v {
			t.Fatal("spare octet changes value")
		}
	}
	for _, v := range []time.Duration{-1, time.Nanosecond, 0x100000000 * time.Millisecond} {
		assertAddedInvalid(t, ie.NewULPeriodicity(v))
	}
}
