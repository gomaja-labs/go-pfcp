package ie

import "io"

// valueAsUint8Masked interprets a field while preserving the original encoding.
// TS 29.244 V20.1.0 clause 8.1.2: spare bits are never evaluated on receipt.
func (i *IE) valueAsUint8Masked(mask uint8) (uint8, error) {
	value, err := i.ValueAsUint8()
	return value & mask, err
}

// validateFixedValue applies Table 8.1.2-1 to typed values. Opaque IE framing
// retains unknown contents, including null-length removal instructions (8.1.2).
func (i *IE) validateFixedValue() error {
	var size int
	switch i.Type {
	case Cause, Metric, Weight:
		size = 1
	case OffendingIE, NumberOfReports:
		size = 2
	case SequenceNumber, URSEQN, AggregatedURRID, FramedRouting, PortNumber, NWTTPortNumber, SNSSAI:
		size = 4
	case Multiplier:
		size = 12
	case NFInstanceID:
		size = 16
	default:
		return nil
	}
	if len(i.Payload) < size {
		return &DecodeError{Type: i.Type, Err: io.ErrUnexpectedEOF}
	}
	if len(i.Payload) > size {
		return &DecodeError{Type: i.Type, Err: ErrInvalidLength}
	}
	return nil
}

// payloadFlag tests a known flag without interpreting spare or extension bits.
func payloadFlag(i *IE, typ IEType, offset int, mask uint8) bool {
	return i != nil && i.Type == typ && len(i.Payload) > offset && i.Payload[offset]&mask != 0
}
