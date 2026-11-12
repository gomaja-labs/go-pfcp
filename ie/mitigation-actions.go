package ie

// MitigationActions flags follow TS 29.244 V20.1.0 Figure 8.2.283-1.
const (
	// MitigationActionsDROP selects DROP.
	MitigationActionsDROP uint8 = 1
	// MitigationActionsTHRO selects THRO.
	MitigationActionsTHRO uint8 = 2
	// MitigationActionsIDDN selects IDDN.
	MitigationActionsIDDN uint8 = 4
)

// NewMitigationActions creates a Mitigation Actions IE with at least one action.
// TS 29.244 V20.1.0 clause 8.2.283 assigns DROP and THRO to bits 1 and 2.
// Figure 8.2.283-1 places IDDN at bit 3 and the prose reserves bits 4–8;
// the prose's second "Bit 1" label for IDDN is therefore an editorial error.
// Wireshark's hf_pfcp_mitigation_actions_b2_iddn also uses mask 0x04.
func NewMitigationActions(v uint8) *IE {
	if v == 0 || v > 7 {
		return invalidRecentIE(MitigationActions, "Mitigation Actions: at least one action and no spare bits required")
	}
	return New(MitigationActions, []byte{v})
}

// MitigationActions returns the action flags from TS 29.244 V20.1.0 clause
// 8.2.283, ignoring spare bits and extension octets as required by clause 8.1.2.
func (i *IE) MitigationActions() (uint8, error) {
	b, err := recentIEPayload(i, MitigationActions, 1, false)
	if err != nil {
		return 0, err
	}
	v := b[0] & 7
	if v == 0 {
		return 0, &DecodeError{Type: MitigationActions, Err: ErrInvalidValue}
	}
	return v, nil
}
