// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/google/go-cmp/cmp"
)

func TestDurationIEs(t *testing.T) {
	cases := []struct {
		description string
		structured  *ie.IE
		decoded     time.Duration
		decoderFunc func(*ie.IE) (time.Duration, error)
	}{
		{
			description: "AveragePacketDelay",
			structured:  ie.NewAveragePacketDelay(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.AveragePacketDelay() },
		}, {
			description: "AveragePacketDelay/GTPUPathQoSControlInformation",
			structured: ie.NewGTPUPathQoSControlInformation(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewAveragePacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.GTPUPathQoSControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.AveragePacketDelayThreshold.AveragePacketDelay()
			},
		}, {
			description: "AveragePacketDelay/GTPUPathQoSReport",
			structured: ie.NewGTPUPathQoSReport(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewAveragePacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.GTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				return f0.OtherIEs[0].AveragePacketDelay()
			},
		}, {
			description: "AveragePacketDelay/QoSInformationInGTPUPathQoSReport",
			structured: ie.NewQoSInformationInGTPUPathQoSReport(
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
				ie.NewAveragePacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.QoSInformationInGTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				return f0.AveragePacketDelay.AveragePacketDelay()
			},
		}, {
			description: "DLBufferingDuration/20hr",
			structured:  ie.NewDLBufferingDuration(20 * time.Hour),
			decoded:     20 * time.Hour,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.DLBufferingDuration() },
		}, {
			description: "DLBufferingDuration/15min",
			structured:  ie.NewDLBufferingDuration(15 * time.Minute),
			decoded:     15 * time.Minute,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.DLBufferingDuration() },
		}, {
			description: "DLBufferingDuration/30sec",
			structured:  ie.NewDLBufferingDuration(30 * time.Second),
			decoded:     30 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.DLBufferingDuration() },
		}, {
			description: "DLBufferingDuration/UpdateBARWithinSessionReportResponse",
			structured: ie.NewUpdateBARWithinSessionReportResponse(
				ie.NewBARID(0xff),
				ie.NewDLBufferingDuration(30*time.Second),
			),
			decoded: 30 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateBARWithinSessionReportResponse()
				if err != nil {
					return 0, err
				}
				return f0.DLBufferingDuration.DLBufferingDuration()
			},
		}, {
			description: "DownlinkDataNotificationDelay",
			structured:  ie.NewDownlinkDataNotificationDelay(100 * time.Millisecond),
			decoded:     100 * time.Millisecond,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.DownlinkDataNotificationDelay() },
		}, {
			description: "DownlinkDataNotificationDelay/CreateBAR",
			structured: ie.NewCreateBAR(
				ie.NewBARID(0xff),
				ie.NewDownlinkDataNotificationDelay(100*time.Millisecond),
			),
			decoded: 100 * time.Millisecond,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateBAR()
				if err != nil {
					return 0, err
				}
				return f0.DownlinkDataNotificationDelay.DownlinkDataNotificationDelay()
			},
		}, {
			description: "DownlinkDataNotificationDelay/UpdateBARWithinSessionReportResponse",
			structured: ie.NewUpdateBARWithinSessionReportResponse(
				ie.NewBARID(0xff),
				ie.NewDownlinkDataNotificationDelay(100*time.Millisecond),
			),
			decoded: 100 * time.Millisecond,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateBARWithinSessionReportResponse()
				if err != nil {
					return 0, err
				}
				return f0.DownlinkDataNotificationDelay.DownlinkDataNotificationDelay()
			},
		}, {
			description: "DownlinkDataNotificationDelay/UpdateBARWithinSessionModificationRequest",
			structured: ie.NewUpdateBARWithinSessionModificationRequest(
				ie.NewBARID(0xff),
				ie.NewDownlinkDataNotificationDelay(100*time.Millisecond),
			),
			decoded: 100 * time.Millisecond,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateBARWithinSessionModificationRequest()
				if err != nil {
					return 0, err
				}
				return f0.DownlinkDataNotificationDelay.DownlinkDataNotificationDelay()
			},
		}, {
			description: "DurationMeasurement",
			structured:  ie.NewDurationMeasurement(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.DurationMeasurement() },
		}, {
			description: "DurationMeasurement/UsageReportWithinSessionModificationResponse",
			structured: ie.NewUsageReportWithinSessionModificationResponse(
				ie.NewURRID(0xffffffff),
				ie.NewDurationMeasurement(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UsageReportWithinSessionModificationResponse()
				if err != nil {
					return 0, err
				}
				return f0.DurationMeasurement.DurationMeasurement()
			},
		}, {
			description: "DurationMeasurement/UsageReportWithinSessionDeletionResponse",
			structured: ie.NewUsageReportWithinSessionDeletionResponse(
				ie.NewURRID(0xffffffff),
				ie.NewDurationMeasurement(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UsageReportWithinSessionDeletionResponse()
				if err != nil {
					return 0, err
				}
				return f0.DurationMeasurement.DurationMeasurement()
			},
		}, {
			description: "DurationMeasurement/UsageReportWithinSessionReportRequest",
			structured: ie.NewUsageReportWithinSessionReportRequest(
				ie.NewURRID(0xffffffff),
				ie.NewDurationMeasurement(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UsageReportWithinSessionReportRequest()
				if err != nil {
					return 0, err
				}
				return f0.DurationMeasurement.DurationMeasurement()
			},
		}, {
			description: "EthernetInactivityTimer",
			structured:  ie.NewEthernetInactivityTimer(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.EthernetInactivityTimer() },
		}, {
			description: "EthernetInactivityTimer/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewEthernetInactivityTimer(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.EthernetInactivityTimer.EthernetInactivityTimer()
			},
		}, {
			description: "EthernetInactivityTimer/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewEthernetInactivityTimer(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.EthernetInactivityTimer.EthernetInactivityTimer()
			},
		}, {
			description: "GracefulReleasePeriod/20hr",
			structured:  ie.NewGracefulReleasePeriod(20 * time.Hour),
			decoded:     20 * time.Hour,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.GracefulReleasePeriod() },
		}, {
			description: "GracefulReleasePeriod/15min",
			structured:  ie.NewGracefulReleasePeriod(15 * time.Minute),
			decoded:     15 * time.Minute,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.GracefulReleasePeriod() },
		}, {
			description: "GracefulReleasePeriod/30sec",
			structured:  ie.NewGracefulReleasePeriod(30 * time.Second),
			decoded:     30 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.GracefulReleasePeriod() },
		}, {
			description: "MaximumPacketDelay",
			structured:  ie.NewMaximumPacketDelay(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.MaximumPacketDelay() },
		}, {
			description: "MaximumPacketDelay/GTPUPathQoSControlInformation",
			structured: ie.NewGTPUPathQoSControlInformation(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewMaximumPacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.GTPUPathQoSControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.MaximumPacketDelayThreshold.MaximumPacketDelay()
			},
		}, {
			description: "MaximumPacketDelay/GTPUPathQoSReport",
			structured: ie.NewGTPUPathQoSReport(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewMaximumPacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.GTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				return f0.OtherIEs[0].MaximumPacketDelay()
			},
		}, {
			description: "MaximumPacketDelay/QoSInformationInGTPUPathQoSReport",
			structured: ie.NewQoSInformationInGTPUPathQoSReport(
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
				ie.NewMaximumPacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.QoSInformationInGTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				return f0.MaximumPacketDelay.MaximumPacketDelay()
			},
		}, {
			description: "MeasurementPeriod",
			structured:  ie.NewMeasurementPeriod(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.MeasurementPeriod() },
		}, {
			description: "MeasurementPeriod/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewMeasurementPeriod(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.MeasurementPeriod.MeasurementPeriod()
			},
		}, {
			description: "MeasurementPeriod/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewMeasurementPeriod(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.MeasurementPeriod.MeasurementPeriod()
			},
		}, {
			description: "MeasurementPeriod/QoSMonitoringPerQoSFlowControlInformation",
			structured: ie.NewQoSMonitoringPerQoSFlowControlInformation(
				ie.NewQFI(0x01),
				ie.NewMeasurementPeriod(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.QoSMonitoringPerQoSFlowControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.MeasurementPeriod.MeasurementPeriod()
			},
		}, {
			description: "MinimumPacketDelay",
			structured:  ie.NewMinimumPacketDelay(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.MinimumPacketDelay() },
		}, {
			description: "MinimumPacketDelay/GTPUPathQoSControlInformation",
			structured: ie.NewGTPUPathQoSControlInformation(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewMinimumPacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.GTPUPathQoSControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.MinimumPacketDelayThreshold.MinimumPacketDelay()
			},
		}, {
			description: "MinimumPacketDelay/GTPUPathQoSReport",
			structured: ie.NewGTPUPathQoSReport(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewMinimumPacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.GTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				return f0.OtherIEs[0].MinimumPacketDelay()
			},
		}, {
			description: "MinimumPacketDelay/QoSInformationInGTPUPathQoSReport",
			structured: ie.NewQoSInformationInGTPUPathQoSReport(
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
				ie.NewMinimumPacketDelay(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.QoSInformationInGTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				return f0.MinimumPacketDelay.MinimumPacketDelay()
			},
		}, {
			description: "MinimumWaitTime",
			structured:  ie.NewMinimumWaitTime(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.MinimumWaitTime() },
		}, {
			description: "MinimumWaitTime/QoSMonitoringPerQoSFlowControlInformation",
			structured: ie.NewQoSMonitoringPerQoSFlowControlInformation(
				ie.NewQFI(0x01),
				ie.NewMinimumWaitTime(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.QoSMonitoringPerQoSFlowControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.MinimumWaitTime.MinimumWaitTime()
			},
		}, {
			description: "QuotaHoldingTime",
			structured:  ie.NewQuotaHoldingTime(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.QuotaHoldingTime() },
		}, {
			description: "QuotaHoldingTime/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewQuotaHoldingTime(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.QuotaHoldingTime.QuotaHoldingTime()
			},
		}, {
			description: "QuotaHoldingTime/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewQuotaHoldingTime(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.QuotaHoldingTime.QuotaHoldingTime()
			},
		}, {
			description: "QuotaValidityTime",
			structured:  ie.NewQuotaValidityTime(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.QuotaValidityTime() },
		}, {
			description: "QuotaValidityTime/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewQuotaValidityTime(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.QuotaValidityTime.QuotaValidityTime()
			},
		}, {
			description: "QuotaValidityTime/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewQuotaValidityTime(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.QuotaValidityTime.QuotaValidityTime()
			},
		}, {
			description: "SubsequentTimeQuota",
			structured:  ie.NewSubsequentTimeQuota(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.SubsequentTimeQuota() },
		}, {
			description: "SubsequentTimeQuota/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewSubsequentTimeQuota(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.SubsequentTimeQuota.SubsequentTimeQuota()
			},
		}, {
			description: "SubsequentTimeQuota/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewSubsequentTimeQuota(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.SubsequentTimeQuota.SubsequentTimeQuota()
			},
		}, {
			description: "SubsequentTimeQuota/AdditionalMonitoringTime",
			structured: ie.NewAdditionalMonitoringTime(
				ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
				ie.NewSubsequentTimeQuota(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.AdditionalMonitoringTime()
				if err != nil {
					return 0, err
				}
				return f0.SubsequentTimeQuota.SubsequentTimeQuota()
			},
		}, {
			description: "SubsequentTimeThreshold",
			structured:  ie.NewSubsequentTimeThreshold(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.SubsequentTimeThreshold() },
		}, {
			description: "SubsequentTimeThreshold/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewSubsequentTimeThreshold(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.SubsequentTimeThreshold.SubsequentTimeThreshold()
			},
		}, {
			description: "SubsequentTimeThreshold/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewSubsequentTimeThreshold(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.SubsequentTimeThreshold.SubsequentTimeThreshold()
			},
		}, {
			description: "SubsequentTimeThreshold/AdditionalMonitoringTime",
			structured: ie.NewAdditionalMonitoringTime(
				ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
				ie.NewSubsequentTimeThreshold(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.AdditionalMonitoringTime()
				if err != nil {
					return 0, err
				}
				return f0.SubsequentTimeThreshold.SubsequentTimeThreshold()
			},
		}, {
			description: "TimeOffsetMeasurement",
			structured:  ie.NewTimeOffsetMeasurement(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.TimeOffsetMeasurement() },
		}, {
			description: "TimeOffsetThreshold",
			structured:  ie.NewTimeOffsetThreshold(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.TimeOffsetThreshold() },
		}, {
			description: "TimeOffsetThreshold/ClockDriftControlInformation",
			structured: ie.NewClockDriftControlInformation(
				ie.NewRequestedClockDriftInformation(1, 1),
				ie.NewTimeOffsetThreshold(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.ClockDriftControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.TimeOffsetThreshold.TimeOffsetThreshold()
			},
		}, {
			description: "TimeOffsetThreshold/ClockDriftReport",
			structured: ie.NewClockDriftReport(
				ie.NewTimeDomainNumber(255),
				ie.NewTimeOffsetThreshold(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.ClockDriftReport()
				if err != nil {
					return 0, err
				}
				return f0.OtherIEs[0].TimeOffsetThreshold()
			},
		}, {
			description: "TimeQuota",
			structured:  ie.NewTimeQuota(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.TimeQuota() },
		}, {
			description: "TimeQuota/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewTimeQuota(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.TimeQuota.TimeQuota()
			},
		}, {
			description: "TimeQuota/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewTimeQuota(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.TimeQuota.TimeQuota()
			},
		}, {
			description: "TimeThreshold",
			structured:  ie.NewTimeThreshold(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.TimeThreshold() },
		}, {
			description: "TimeThreshold/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewTimeThreshold(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.TimeThreshold.TimeThreshold()
			},
		}, {
			description: "TimeThreshold/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewTimeThreshold(10*time.Second),
			),
			decoded: 10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.TimeThreshold.TimeThreshold()
			},
		}, {
			description: "Timer/20hr",
			structured:  ie.NewTimer(20 * time.Hour),
			decoded:     20 * time.Hour,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.Timer() },
		}, {
			description: "Timer/15min",
			structured:  ie.NewTimer(15 * time.Minute),
			decoded:     15 * time.Minute,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.Timer() },
		}, {
			description: "Timer/30sec",
			structured:  ie.NewTimer(30 * time.Second),
			decoded:     30 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.Timer() },
		}, {
			description: "Timer/OverloadControlInformation",
			structured: ie.NewOverloadControlInformation(
				ie.NewSequenceNumber(0xffffffff),
				ie.NewTimer(30*time.Second),
			),
			decoded: 30 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.OverloadControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.PeriodOfValidity.Timer()
			},
		}, {
			description: "Timer/GTPUPathQoSControlInformation",
			structured: ie.NewGTPUPathQoSControlInformation(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewTimer(30*time.Second),
			),
			decoded: 30 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) {
				f0, err := i.GTPUPathQoSControlInformation()
				if err != nil {
					return 0, err
				}
				return f0.MinimumWaitingTime.Timer()
			},
		}, {
			description: "UserPlaneInactivityTimer",
			structured:  ie.NewUserPlaneInactivityTimer(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.UserPlaneInactivityTimer() },
		}, {
			description: "ValidityTimer",
			structured:  ie.NewValidityTimer(10 * time.Second),
			decoded:     10 * time.Second,
			decoderFunc: func(i *ie.IE) (time.Duration, error) { return i.ValidityTimer() },
		},
	}

	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			got, err := c.decoderFunc(c.structured)
			if err != nil {
				t.Fatal(err)
			}

			if diff := cmp.Diff(got, c.decoded); diff != "" {
				t.Error(diff)
			}
		})
	}
}
