package ie

// IsGrouped reports whether the type is a grouped IE defined by
// TS 29.244 V20.1.0 Table 8.1.2-1 (types 1 through 271).
// Vendor-specific IEs require explicit interpretation with ValueAsGrouped.
func (i *IE) IsGrouped() bool {
	if i == nil {
		return false
	}
	switch i.Type {
	case IPMulticastAddressingInfo, CreatePDR,
		PDI,
		CreateFAR,
		ForwardingParameters,
		DuplicatingParameters,
		CreateURR,
		CreateQER,
		CreatedPDR,
		UpdatePDR,
		UpdateFAR,
		UpdateForwardingParameters,
		UpdateBARWithinSessionReportResponse,
		UpdateURR,
		UpdateQER,
		RemovePDR,
		RemoveFAR,
		RemoveURR,
		RemoveQER,
		LoadControlInformation,
		OverloadControlInformation,
		ApplicationIDsPFDs,
		PFDContext,
		ApplicationDetectionInformation,
		QueryURR,
		UsageReportWithinSessionModificationResponse,
		UsageReportWithinSessionDeletionResponse,
		UsageReportWithinSessionReportRequest,
		DownlinkDataReport,
		CreateBAR,
		UpdateBARWithinSessionModificationRequest,
		RemoveBAR,
		ErrorIndicationReport,
		UserPlanePathFailureReport,
		UpdateDuplicatingParameters,
		AggregatedURRs,
		CreateTrafficEndpoint,
		CreatedTrafficEndpoint,
		UpdateTrafficEndpoint,
		RemoveTrafficEndpoint,
		EthernetPacketFilter,
		EthernetTrafficInformation,
		AdditionalMonitoringTime,
		CreateMAR,
		TGPPAccessForwardingActionInformation,
		NonTGPPAccessForwardingActionInformation,
		RemoveMAR,
		UpdateMAR,
		UpdateTGPPAccessForwardingActionInformation,
		UpdateNonTGPPAccessForwardingActionInformation,
		PFCPSessionRetentionInformation,
		UserPlanePathRecoveryReport,
		JoinIPMulticastInformationWithinUsageReport,
		LeaveIPMulticastInformationWithinUsageReport,
		CreatedBridgeRouterInfo,
		TSCManagementInformationWithinSessionModificationRequest,
		TSCManagementInformationWithinSessionModificationResponse,
		TSCManagementInformationWithinSessionReportRequest,
		ClockDriftControlInformation,
		ClockDriftReport,
		RemoveSRR,
		CreateSRR,
		UpdateSRR,
		SessionReport,
		AccessAvailabilityControlInformation,
		AccessAvailabilityReport,
		ProvideATSSSControlInformation,
		ATSSSControlParameters,
		MPTCPParameters,
		ATSSSLLParameters,
		PMFParameters,
		UEIPAddressPoolInformation,
		GTPUPathQoSControlInformation,
		GTPUPathQoSReport,
		QoSInformationInGTPUPathQoSReport,
		QoSMonitoringPerQoSFlowControlInformation,
		QoSMonitoringReport,
		PacketRateStatusReport,
		EthernetContextInformation,
		RedundantTransmissionParameters,
		UpdatedPDR,
		ProvideRDSConfigurationInformation,
		QueryPacketRateStatusWithinSessionModificationRequest,
		PacketRateStatusReportWithinSessionModificationResponse,
		UEIPAddressUsageInformation,
		RedundantTransmissionForwardingParameters,
		TransportDelayReporting:
		return true
	default:
		return false
	}
}

// Add appends non-nil children to a grouped IE. A raw payload is parsed first.
// If that payload is malformed, it returns an error and leaves the IE unchanged.
func (i *IE) Add(ies ...*IE) error {
	if i == nil {
		return ErrIENotFound
	}

	if !i.IsGrouped() && i.ChildIEs == nil {
		return ErrInvalidType
	}
	children, err := i.ValueAsGrouped()
	if err != nil {
		return err
	}
	if children == nil {
		children = make([]*IE, 0)
	}
	for _, child := range ies {
		if child != nil {
			children = append(children, child)
		}
	}
	i.ChildIEs, i.Payload = children, nil
	return nil
}

// Remove removes all children of typ from a grouped IE. A raw payload is parsed
// first; a malformed payload returns an error and leaves the IE unchanged.
func (i *IE) Remove(typ IEType) error {
	if i == nil {
		return ErrIENotFound
	}

	if !i.IsGrouped() && i.ChildIEs == nil {
		return ErrInvalidType
	}
	children, err := i.ValueAsGrouped()
	if err != nil {
		return err
	}
	kept := make([]*IE, 0, len(children))
	for _, child := range children {
		if child != nil && child.Type != typ {
			kept = append(kept, child)
		}
	}
	i.ChildIEs, i.Payload = kept, nil
	return nil
}

// FindByType returns the first child with the requested type.
func (i *IE) FindByType(typ IEType) (*IE, error) {
	if i == nil {
		return nil, ErrIENotFound
	}
	if !i.IsGrouped() && i.ChildIEs == nil {
		return nil, ErrInvalidType
	}
	children, err := i.ValueAsGrouped()
	if err != nil {
		return nil, err
	}
	for _, child := range children {
		if child != nil && child.Type == typ {
			return child, nil
		}
	}
	return nil, ErrIENotFound
}
