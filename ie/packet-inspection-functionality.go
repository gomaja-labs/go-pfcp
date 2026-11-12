package ie

// PacketInspectionFunctionality values are defined in TS 29.244 V20.1.0 clause 8.2.250.
const (
	// PacketInspectionFunctionalityIPBD selects IPBD.
	PacketInspectionFunctionalityIPBD uint8 = 1
	// PacketInspectionFunctionalityMACD selects MACD.
	PacketInspectionFunctionalityMACD uint8 = 2
	// PacketInspectionFunctionalityDPID selects DPID.
	PacketInspectionFunctionalityDPID uint8 = 4
)

// NewPacketInspectionFunctionality creates a Packet Inspection functionality IE.
// TS 29.244 V20.1.0 clause 8.2.250.
func NewPacketInspectionFunctionality(v uint8) *IE {
	if v > 7 {
		return invalidRecentIE(PacketInspectionFunctionality, "Packet Inspection functionality")
	}
	return New(PacketInspectionFunctionality, []byte{v})
}

// PacketInspectionFunctionality returns the Packet Inspection functionality value.
// TS 29.244 V20.1.0 clause 8.2.250; clause 8.1.2 governs spare bits and extensions.
func (i *IE) PacketInspectionFunctionality() (uint8, error) {
	b, err := recentIEPayload(i, PacketInspectionFunctionality, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 7

	if v > 7 {
		return 0, &DecodeError{Type: PacketInspectionFunctionality, Err: ErrInvalidValue}
	}
	return v, nil
}
