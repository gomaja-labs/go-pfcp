// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie

// Cause definitions.
const (
	_ uint8 = 0
	// CauseRequestAccepted selects Request Accepted.
	CauseRequestAccepted uint8 = 1
	// CauseRequestRejected selects Request Rejected.
	CauseRequestRejected uint8 = 64
	// CauseSessionContextNotFound selects Session Context Not Found.
	CauseSessionContextNotFound uint8 = 65
	// CauseMandatoryIEMissing selects Mandatory IE Missing.
	CauseMandatoryIEMissing uint8 = 66
	// CauseConditionalIEMissing selects Conditional IE Missing.
	CauseConditionalIEMissing uint8 = 67
	// CauseInvalidLength selects Invalid Length.
	CauseInvalidLength uint8 = 68
	// CauseMandatoryIEIncorrect selects Mandatory IE Incorrect.
	CauseMandatoryIEIncorrect uint8 = 69
	// CauseInvalidForwardingPolicy selects Invalid Forwarding Policy.
	CauseInvalidForwardingPolicy uint8 = 70
	// CauseInvalidFTEIDAllocationOption selects Invalid FTEID Allocation Option.
	CauseInvalidFTEIDAllocationOption uint8 = 71
	// CauseNoEstablishedPFCPAssociation selects No Established PFCP Association.
	CauseNoEstablishedPFCPAssociation uint8 = 72
	// CauseRuleCreationModificationFailure selects Rule Creation Modification Failure.
	CauseRuleCreationModificationFailure uint8 = 73
	// CausePFCPEntityInCongestion selects PFCP Entity In Congestion.
	CausePFCPEntityInCongestion uint8 = 74
	// CauseNoResourcesAvailable selects No Resources Available.
	CauseNoResourcesAvailable uint8 = 75
	// CauseServiceNotSupported selects Service Not Supported.
	CauseServiceNotSupported uint8 = 76
	// CauseSystemFailure selects System Failure.
	CauseSystemFailure uint8 = 77
	// CauseRedirectionRequested selects Redirection Requested.
	CauseRedirectionRequested uint8 = 78
)

// NewCause creates a new Cause IE.
func NewCause(cause uint8) *IE {
	return newUint8ValIE(Cause, cause)
}

// Cause returns Cause in uint8 if the type of IE matches.
func (i *IE) Cause() (uint8, error) {
	if i == nil {
		return 0, ErrIENotFound
	}

	if i.Type != Cause {
		return 0, &InvalidTypeError{Type: i.Type}
	}
	return i.ValueAsUint8()
}

// TS 29.244 V20.1.0 Table 8.2.1-1.
const (
	// CauseMoreUsageReportsToSend selects More Usage Reports To Send.
	CauseMoreUsageReportsToSend uint8 = 2
	// CauseRequestPartiallyAccepted selects Request Partially Accepted.
	CauseRequestPartiallyAccepted uint8 = 3
	// CauseAllDynamicAddressesOccupied selects All Dynamic Addresses Occupied.
	CauseAllDynamicAddressesOccupied uint8 = 79
	// CauseUnknownPredefinedRule selects Unknown Predefined Rule.
	CauseUnknownPredefinedRule uint8 = 80
	// CauseUnknownApplicationID selects Unknown Application ID.
	CauseUnknownApplicationID uint8 = 81
	// CauseL2TPTunnelEstablishmentFailure selects L2 TP Tunnel Establishment Failure.
	CauseL2TPTunnelEstablishmentFailure uint8 = 82
	// CauseL2TPSessionEstablishmentFailure selects L2 TP Session Establishment Failure.
	CauseL2TPSessionEstablishmentFailure uint8 = 83
	// CauseL2TPTunnelRelease selects L2 TP Tunnel Release.
	CauseL2TPTunnelRelease uint8 = 84
	// CauseL2TPSessionRelease selects L2 TP Session Release.
	CauseL2TPSessionRelease uint8 = 85
	// CauseSessionRestorationResourceUnavailable selects Session Restoration Resource Unavailable.
	CauseSessionRestorationResourceUnavailable uint8 = 86
	// CauseL2TPTunnelAuthenticationFailure selects L2 TP Tunnel Authentication Failure.
	CauseL2TPTunnelAuthenticationFailure uint8 = 87
	// CauseL2TPSessionAuthenticationFailure selects L2 TP Session Authentication Failure.
	CauseL2TPSessionAuthenticationFailure uint8 = 88
	// CauseL2TPLNSUnreachable selects L2 TPLNS Unreachable.
	CauseL2TPLNSUnreachable uint8 = 89
	// CausePFDContentsSyntaxError selects PFD Contents Syntax Error.
	CausePFDContentsSyntaxError uint8 = 90
	// CausePFDContentsSemanticsError selects PFD Contents Semantics Error.
	CausePFDContentsSemanticsError uint8 = 91
	// CausePFDApplicationIDUnknown selects PFD Application ID Unknown.
	CausePFDApplicationIDUnknown uint8 = 92
)
