package ie

// RATTypeValue encodes RAT Type; TS 29.244 V20.1.0 clause 8.2.186.
type RATTypeValue uint8

const (
	// RATTypeUTRAN is the UTRAN value defined in clause 8.2.186.
	RATTypeUTRAN RATTypeValue = 1
	// RATTypeGERAN is the GERAN value defined in clause 8.2.186.
	RATTypeGERAN RATTypeValue = 2
	// RATTypeWLAN is the WLAN value defined in clause 8.2.186.
	RATTypeWLAN RATTypeValue = 3
	// RATTypeGAN is the GAN value defined in clause 8.2.186.
	RATTypeGAN RATTypeValue = 4
	// RATTypeHSPAEvolution is the HSPAEvolution value defined in clause 8.2.186.
	RATTypeHSPAEvolution RATTypeValue = 5
	// RATTypeEUTRAN is the EUTRAN value defined in clause 8.2.186.
	RATTypeEUTRAN RATTypeValue = 6
	// RATTypeVirtual is the Virtual value defined in clause 8.2.186.
	RATTypeVirtual RATTypeValue = 7
	// RATTypeEUTRANNBIoT is the EUTRANNBIoT value defined in clause 8.2.186.
	RATTypeEUTRANNBIoT RATTypeValue = 8
	// RATTypeLTEM is the LTEM value defined in clause 8.2.186.
	RATTypeLTEM RATTypeValue = 9
	// RATTypeNR is the NR value defined in clause 8.2.186.
	RATTypeNR RATTypeValue = 10
	// RATTypeEUTRANLEO is the EUTRANLEO value defined in clause 8.2.186.
	RATTypeEUTRANLEO RATTypeValue = 11
	// RATTypeEUTRANMEO is the EUTRANMEO value defined in clause 8.2.186.
	RATTypeEUTRANMEO RATTypeValue = 12
	// RATTypeEUTRANGEO is the EUTRANGEO value defined in clause 8.2.186.
	RATTypeEUTRANGEO RATTypeValue = 13
	// RATTypeEUTRANOtherSatellite is the EUTRANOtherSatellite value defined in clause 8.2.186.
	RATTypeEUTRANOtherSatellite RATTypeValue = 14
	// RATTypeNBIoTLEO is the NBIoTLEO value defined in clause 8.2.186.
	RATTypeNBIoTLEO RATTypeValue = 15
	// RATTypeNBIoTMEO is the NBIoTMEO value defined in clause 8.2.186.
	RATTypeNBIoTMEO RATTypeValue = 16
	// RATTypeNBIoTGEO is the NBIoTGEO value defined in clause 8.2.186.
	RATTypeNBIoTGEO RATTypeValue = 17
	// RATTypeNBIoTOtherSatellite is the NBIoTOtherSatellite value defined in clause 8.2.186.
	RATTypeNBIoTOtherSatellite RATTypeValue = 18
	// RATTypeLTEMLEO is the LTEMLEO value defined in clause 8.2.186.
	RATTypeLTEMLEO RATTypeValue = 19
	// RATTypeLTEMMEO is the LTEMMEO value defined in clause 8.2.186.
	RATTypeLTEMMEO RATTypeValue = 20
	// RATTypeLTEMGEO is the LTEMGEO value defined in clause 8.2.186.
	RATTypeLTEMGEO RATTypeValue = 21
	// RATTypeLTEMOtherSatellite is the LTEMOtherSatellite value defined in clause 8.2.186.
	RATTypeLTEMOtherSatellite RATTypeValue = 22
	// RATTypeNRLEO is the NRLEO value defined in clause 8.2.186.
	RATTypeNRLEO RATTypeValue = 23
	// RATTypeNRMEO is the NRMEO value defined in clause 8.2.186.
	RATTypeNRMEO RATTypeValue = 24
	// RATTypeNRGEO is the NRGEO value defined in clause 8.2.186.
	RATTypeNRGEO RATTypeValue = 25
	// RATTypeNROtherSatellite is the NROtherSatellite value defined in clause 8.2.186.
	RATTypeNROtherSatellite RATTypeValue = 26
	// RATTypeNRRedCap is the NRRedCap value defined in clause 8.2.186.
	RATTypeNRRedCap RATTypeValue = 27
	// RATTypeNReRedCap is the NReRedCap value defined in clause 8.2.186.
	RATTypeNReRedCap RATTypeValue = 28
)

// NewRATType creates a RAT Type IE.
// TS 29.244 V20.1.0 clause 8.2.186.
func NewRATType(value RATTypeValue) *IE {
	if value < 1 || value > 28 {
		return newInvalidIE(RATType, ErrInvalidValue)
	}
	return New(RATType, []byte{byte(value)})
}

// RATType returns the RAT Type value, ignoring spare bits.
// TS 29.244 V20.1.0 clauses 8.2.186 and 8.1.2.
func (i *IE) RATType() (RATTypeValue, error) {
	b, err := addedValuePayload(i, RATType, 1, false)
	if err != nil {
		return 0, err
	}
	return RATTypeValue(b[0] & 255), nil
}
