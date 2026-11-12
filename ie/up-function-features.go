// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

import "io"

// NewUPFunctionFeatures creates a UPFunctionFeatures IE from successive feature octets.
// TS 29.244 V20.1.0 clause 8.2.25 defines two-octet feature groups.
// An incomplete final group is padded with zero; all supplied octets are retained.
func NewUPFunctionFeatures(features ...uint8) *IE {
	length := max(2, (len(features)+1)&^1)
	payload := make([]byte, length)
	copy(payload, features)
	return New(UPFunctionFeatures, payload)
}

// UPFunctionFeatures returns UPFunctionFeatures in []byte if the type of IE matches.
func (i *IE) UPFunctionFeatures() ([]byte, error) {
	if i == nil {
		return nil, ErrIENotFound
	}

	if i.Type != UPFunctionFeatures {
		return nil, &InvalidTypeError{Type: i.Type}
	}

	// TS 29.244 V20.1.0 clause 8.2.25: first two-octet feature group.
	if len(i.Payload) < 2 {
		return nil, io.ErrUnexpectedEOF
	}
	return i.Payload, nil
}

// HasBUCP reports whether an IE has BUCP bit.
func (i *IE) HasBUCP() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has1stBit(i.Payload[0])
}

// HasDDND reports whether an IE has DDND bit.
func (i *IE) HasDDND() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has2ndBit(i.Payload[0])
}

// HasDLBD reports whether an IE has DLBD bit.
func (i *IE) HasDLBD() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has3rdBit(i.Payload[0])
}

// HasTRST reports whether an IE has TRST bit.
func (i *IE) HasTRST() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has4thBit(i.Payload[0])
}

// HasFTUP reports whether an IE has FTUP bit.
func (i *IE) HasFTUP() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has5thBit(i.Payload[0])
}

// HasPFDM reports whether an IE has PFDM bit.
func (i *IE) HasPFDM() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has6thBit(i.Payload[0])
}

// HasHEEU reports whether an IE has HEEU bit.
func (i *IE) HasHEEU() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has7thBit(i.Payload[0])
}

// HasTREU reports whether an IE has TREU bit.
func (i *IE) HasTREU() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 1 {
		return false
	}

	return has8thBit(i.Payload[0])
}

// HasEMPU reports whether an IE has EMPU bit.
func (i *IE) HasEMPU() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 2 {
		return false
	}

	return has1stBit(i.Payload[1])
}

// HasPDIU reports whether an IE has PDIU bit.
func (i *IE) HasPDIU() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 2 {
		return false
	}

	return has2ndBit(i.Payload[1])
}

// HasUDBC reports whether an IE has UDBC bit.
func (i *IE) HasUDBC() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 2 {
		return false
	}

	return has3rdBit(i.Payload[1])
}

// HasQUOAC reports whether an IE has QUOAC bit.
func (i *IE) HasQUOAC() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 2 {
		return false
	}

	return has4thBit(i.Payload[1])
}

// HasTRACE reports whether an IE has TRACE bit.
func (i *IE) HasTRACE() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 2 {
		return false
	}

	return has5thBit(i.Payload[1])
}

// HasFRRT reports whether an IE has FRRT bit.
func (i *IE) HasFRRT() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 2 {
		return false
	}

	return has6thBit(i.Payload[1])
}

// HasPFDE reports whether an IE has PFDE bit.
func (i *IE) HasPFDE() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 2 {
		return false
	}

	return has7thBit(i.Payload[1])
}

// HasEPFAR reports whether an IE has EPFAR bit.
func (i *IE) HasEPFAR() bool {
	if i == nil {
		return false
	}

	switch i.Type {
	case UPFunctionFeatures:
		if len(i.Payload) < 2 {
			return false
		}

		return has8thBit(i.Payload[1])
	case CPFunctionFeatures:
		if len(i.Payload) < 1 {
			return false
		}

		return has3rdBit(i.Payload[0])
	default:
		return false
	}
}

// HasDPDRA reports whether an IE has DPDRA  bit.
func (i *IE) HasDPDRA() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 3 {
		return false
	}

	return has1stBit(i.Payload[2])
}

// HasADPDP reports whether an IE has ADPDP  bit.
func (i *IE) HasADPDP() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 3 {
		return false
	}

	return has2ndBit(i.Payload[2])
}

// HasUEIP reports whether an IE has UEIP  bit.
func (i *IE) HasUEIP() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 3 {
		return false
	}

	return has3rdBit(i.Payload[2])
}

// HasSSET reports whether an IE has SSET  bit.
func (i *IE) HasSSET() bool {
	if i == nil {
		return false
	}

	switch i.Type {
	case UPFunctionFeatures:
		if len(i.Payload) < 3 {
			return false
		}

		return has4thBit(i.Payload[2])
	case CPFunctionFeatures:
		if len(i.Payload) < 1 {
			return false
		}

		return has4thBit(i.Payload[0])
	default:
		return false
	}
}

// HasMNOP reports whether an IE has MNOP  bit.
func (i *IE) HasMNOP() bool {
	if i == nil {
		return false
	}

	switch i.Type {
	case UPFunctionFeatures:
		if len(i.Payload) < 3 {
			return false
		}

		return has5thBit(i.Payload[2])
	case MeasurementInformation:
		if len(i.Payload) < 1 {
			return false
		}

		return has5thBit(i.Payload[0])
	default:
		return false
	}
}

// HasMTE reports whether an IE has MTE  bit.
func (i *IE) HasMTE() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 3 {
		return false
	}

	return has6thBit(i.Payload[2])
}

// HasBUNDL reports whether an IE has BUNDL  bit.
func (i *IE) HasBUNDL() bool {
	if i == nil {
		return false
	}

	switch i.Type {
	case UPFunctionFeatures:
		if len(i.Payload) < 3 {
			return false
		}

		return has7thBit(i.Payload[2])
	case CPFunctionFeatures:
		if len(i.Payload) < 1 {
			return false
		}

		return has5thBit(i.Payload[0])
	default:
		return false
	}
}

// HasGCOM reports whether an IE has GCOM  bit.
func (i *IE) HasGCOM() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 3 {
		return false
	}

	return has8thBit(i.Payload[2])
}

// HasMPAS reports whether an IE has MPAS bit.
func (i *IE) HasMPAS() bool {
	if i == nil {
		return false
	}

	switch i.Type {
	case UPFunctionFeatures:
		if len(i.Payload) < 4 {
			return false
		}

		return has1stBit(i.Payload[3])
	case CPFunctionFeatures:
		if len(i.Payload) < 1 {
			return false
		}

		return has6thBit(i.Payload[0])
	default:
		return false
	}
}

// HasRTTL reports whether an IE has RTTL bit.
func (i *IE) HasRTTL() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 4 {
		return false
	}

	return has2ndBit(i.Payload[3])
}

// HasVTIME reports whether an IE has VTIME bit.
func (i *IE) HasVTIME() bool {
	if i == nil {
		return false
	}

	if i.Type != UPFunctionFeatures {
		return false
	}
	if len(i.Payload) < 4 {
		return false
	}

	return has3rdBit(i.Payload[3])
}

// HasNORP reports whether the NORP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 8 bit 4.
func (i *IE) HasNORP() bool { return payloadFlag(i, UPFunctionFeatures, 3, 0x08) }

// HasIPTV reports whether the IPTV flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 8 bit 5.
func (i *IE) HasIPTV() bool { return payloadFlag(i, UPFunctionFeatures, 3, 0x10) }

// HasTSN reports whether the TSN flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 8 bit 7.
func (i *IE) HasTSN() bool { return payloadFlag(i, UPFunctionFeatures, 3, 0x40) }

// HasMPTCP reports whether the MPTCP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 8 bit 8.
func (i *IE) HasMPTCP() bool { return payloadFlag(i, UPFunctionFeatures, 3, 0x80) }

// HasATSSSLL reports whether the ATSSS-LL flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 1.
func (i *IE) HasATSSSLL() bool { return payloadFlag(i, UPFunctionFeatures, 4, 0x01) }

// HasQFQM reports whether the QFQM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 2.
func (i *IE) HasQFQM() bool { return payloadFlag(i, UPFunctionFeatures, 4, 0x02) }

// HasGPQM reports whether the GPQM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 3.
func (i *IE) HasGPQM() bool { return payloadFlag(i, UPFunctionFeatures, 4, 0x04) }

// HasMTEDT reports whether the MT-EDT flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 4.
func (i *IE) HasMTEDT() bool { return payloadFlag(i, UPFunctionFeatures, 4, 0x08) }

// HasCIOT reports whether the CIOT flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 5.
func (i *IE) HasCIOT() bool { return payloadFlag(i, UPFunctionFeatures, 4, 0x10) }

// HasETHAR reports whether the ETHAR flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 6.
func (i *IE) HasETHAR() bool { return payloadFlag(i, UPFunctionFeatures, 4, 0x20) }

// HasDDDS reports whether the DDDS flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 9 bit 7.
func (i *IE) HasDDDS() bool { return payloadFlag(i, UPFunctionFeatures, 4, 0x40) }

// HasRTTWP reports whether the RTTWP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 1.
func (i *IE) HasRTTWP() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x01) }

// HasQUASF reports whether the QUASF flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 2.
func (i *IE) HasQUASF() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x02) }

// HasNSPOC reports whether the NSPOC flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 3.
func (i *IE) HasNSPOC() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x04) }

// HasL2TP reports whether the L2TP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 4.
func (i *IE) HasL2TP() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x08) }

// HasUPBER reports whether the UPBER flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 5.
func (i *IE) HasUPBER() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x10) }

// HasRESPS reports whether the RESPS flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 6.
func (i *IE) HasRESPS() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x20) }

// HasIPREP reports whether the IPREP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 7.
func (i *IE) HasIPREP() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x40) }

// HasDNSTS reports whether the DNSTS flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 10 bit 8.
func (i *IE) HasDNSTS() bool { return payloadFlag(i, UPFunctionFeatures, 5, 0x80) }

// HasDRQOS reports whether the DRQOS flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 1.
func (i *IE) HasDRQOS() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x01) }

// HasMBSN4 reports whether the MBSN4 flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 2.
func (i *IE) HasMBSN4() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x02) }

// HasPSUPRM reports whether the PSUPRM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 3.
func (i *IE) HasPSUPRM() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x04) }

// HasEPPPI reports whether the EPPPI flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 4.
func (i *IE) HasEPPPI() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x08) }

// HasRATP reports whether the RATP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 5.
func (i *IE) HasRATP() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x10) }

// HasUPIDP reports whether the UPIDP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 6.
func (i *IE) HasUPIDP() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x20) }

// HasAFSFC reports whether the AFSFC flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 7.
func (i *IE) HasAFSFC() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x40) }

// HasMPQUICUDP reports whether the MPQUIC-UDP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 11 bit 8.
func (i *IE) HasMPQUICUDP() bool { return payloadFlag(i, UPFunctionFeatures, 6, 0x80) }

// HasREDSM reports whether the REDSM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 1.
func (i *IE) HasREDSM() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x01) }

// HasDBDM reports whether the DBDM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 2.
func (i *IE) HasDBDM() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x02) }

// HasTSCTS reports whether the TSCTS flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 3.
func (i *IE) HasTSCTS() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x04) }

// HasDRTSC reports whether the DRTSC flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 4.
func (i *IE) HasDRTSC() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x08) }

// HasN6JEDB reports whether the N6JEDB flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 5.
func (i *IE) HasN6JEDB() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x10) }

// HasQMCON reports whether the QMCON flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 6.
func (i *IE) HasQMCON() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x20) }

// HasDETNET reports whether the DETNET flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 7.
func (i *IE) HasDETNET() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x40) }

// HasEML4S reports whether the EML4S flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 12 bit 8.
func (i *IE) HasEML4S() bool { return payloadFlag(i, UPFunctionFeatures, 7, 0x80) }

// HasPDUSM reports whether the PDUSM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 1.
func (i *IE) HasPDUSM() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x01) }

// HasCNTL reports whether the CN-TL flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 2.
func (i *IE) HasCNTL() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x02) }

// HasQMDRM reports whether the QMDRM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 3.
func (i *IE) HasQMDRM() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x04) }

// HasEDBNC reports whether the EDBNC flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 4.
func (i *IE) HasEDBNC() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x08) }

// HasMTSDT reports whether the MT-SDT flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 5.
func (i *IE) HasMTSDT() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x10) }

// HasUPSBIES reports whether the UPSBIES flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 6.
func (i *IE) HasUPSBIES() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x20) }

// HasUMN6IP reports whether the UMN6IP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 7.
func (i *IE) HasUMN6IP() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x40) }

// HasUN6TU reports whether the UN6TU flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 13 bit 8.
func (i *IE) HasUN6TU() bool { return payloadFlag(i, UPFunctionFeatures, 8, 0x80) }

// HasMBSCH reports whether the MBSCH flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 1.
func (i *IE) HasMBSCH() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x01) }

// HasUN6DM reports whether the UN6DM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 2.
func (i *IE) HasUN6DM() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x02) }

// HasNATPUB reports whether the NATPUB flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 3.
func (i *IE) HasNATPUB() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x04) }

// HasUSHPH reports whether the USHPH flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 4.
func (i *IE) HasUSHPH() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x08) }

// HasMPQUICIP reports whether the MPQUIC-IP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 5.
func (i *IE) HasMPQUICIP() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x10) }

// HasMPQUICE reports whether the MPQUIC-E flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 6.
func (i *IE) HasMPQUICE() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x20) }

// HasDYNTR reports whether the DYNTR flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 7.
func (i *IE) HasDYNTR() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x40) }

// HasMUXMF reports whether the MUXMF flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 14 bit 8.
func (i *IE) HasMUXMF() bool { return payloadFlag(i, UPFunctionFeatures, 9, 0x80) }

// HasCONUDP reports whether the CONUDP flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 1.
func (i *IE) HasCONUDP() bool { return payloadFlag(i, UPFunctionFeatures, 10, 0x01) }

// HasMOQ reports whether the MOQ flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 2.
func (i *IE) HasMOQ() bool { return payloadFlag(i, UPFunctionFeatures, 10, 0x02) }

// HasULM reports whether the ULM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 3.
func (i *IE) HasULM() bool { return payloadFlag(i, UPFunctionFeatures, 10, 0x04) }

// HasPSITLM reports whether the PSITLM flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 4.
func (i *IE) HasPSITLM() bool { return payloadFlag(i, UPFunctionFeatures, 10, 0x08) }

// HasUDPOPT reports whether the UDPOPT flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 5.
func (i *IE) HasUDPOPT() bool { return payloadFlag(i, UPFunctionFeatures, 10, 0x10) }

// HasQMABR reports whether the QMABR flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 6.
func (i *IE) HasQMABR() bool { return payloadFlag(i, UPFunctionFeatures, 10, 0x20) }

// HasPAPFD reports whether the PAPFD flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 7.
func (i *IE) HasPAPFD() bool {
	if payloadFlag(i, CPFunctionFeatures, 1, 0x04) {
		return true
	}
	return payloadFlag(i, UPFunctionFeatures, 10, 0x40)
}

// HasNLMAT reports whether the NLMAT flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 15 bit 8.
func (i *IE) HasNLMAT() bool { return payloadFlag(i, UPFunctionFeatures, 10, 0x80) }

// HasUHRR reports whether the UHRR flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 16 bit 1.
func (i *IE) HasUHRR() bool { return payloadFlag(i, UPFunctionFeatures, 11, 0x01) }

// HasSCONE reports whether the SCONE flag is set.
// TS 29.244 V20.1.0 clause 8.2.25, octet 16 bit 2.
func (i *IE) HasSCONE() bool { return payloadFlag(i, UPFunctionFeatures, 11, 0x02) }
