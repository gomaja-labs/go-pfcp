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

func TestTimeIEs(t *testing.T) {
	cases := []struct {
		description string
		structured  *ie.IE
		decoded     time.Time
		decoderFunc func(*ie.IE) (time.Time, error)
	}{
		{
			description: "ActivationTime",
			structured:  ie.NewActivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.ActivationTime() },
		}, {
			description: "ActivationTime/CreatePDR",
			structured: ie.NewCreatePDR(
				ie.NewPDRID(0xffff),
				ie.NewActivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.CreatePDR()
				if err != nil {
					return time.Time{}, err
				}
				return f0.ActivationTime.ActivationTime()
			},
		}, {
			description: "ActivationTime/UpdatePDR",
			structured: ie.NewUpdatePDR(
				ie.NewPDRID(0xffff),
				ie.NewActivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UpdatePDR()
				if err != nil {
					return time.Time{}, err
				}
				return f0.ActivationTime.ActivationTime()
			},
		}, {
			description: "DeactivationTime",
			structured:  ie.NewDeactivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.DeactivationTime() },
		}, {
			description: "DeactivationTime/CreatePDR",
			structured: ie.NewCreatePDR(
				ie.NewPDRID(0xffff),
				ie.NewDeactivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.CreatePDR()
				if err != nil {
					return time.Time{}, err
				}
				return f0.DeactivationTime.DeactivationTime()
			},
		}, {
			description: "DeactivationTime/UpdatePDR",
			structured: ie.NewUpdatePDR(
				ie.NewPDRID(0xffff),
				ie.NewDeactivationTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UpdatePDR()
				if err != nil {
					return time.Time{}, err
				}
				return f0.DeactivationTime.DeactivationTime()
			},
		}, {
			description: "EndTime",
			structured:  ie.NewEndTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.EndTime() },
		}, {
			description: "EndTime/UsageReportWithinSessionModificationResponse",
			structured: ie.NewUsageReportWithinSessionModificationResponse(
				ie.NewURRID(0xffffffff),
				ie.NewEndTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionModificationResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.EndTime.EndTime()
			},
		}, {
			description: "EndTime/UsageReportWithinSessionDeletionResponse",
			structured: ie.NewUsageReportWithinSessionDeletionResponse(
				ie.NewURRID(0xffffffff),
				ie.NewEndTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionDeletionResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.EndTime.EndTime()
			},
		}, {
			description: "EndTime/UsageReportWithinSessionReportRequest",
			structured: ie.NewUsageReportWithinSessionReportRequest(
				ie.NewURRID(0xffffffff),
				ie.NewEndTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionReportRequest()
				if err != nil {
					return time.Time{}, err
				}
				return f0.EndTime.EndTime()
			},
		}, {
			description: "EventTimeStamp",
			structured:  ie.NewEventTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.EventTimeStamp() },
		}, {
			description: "EventTimeStamp/UsageReportWithinSessionReportRequest",
			structured: ie.NewUsageReportWithinSessionReportRequest(
				ie.NewURRID(0xffffffff),
				ie.NewEventTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionReportRequest()
				if err != nil {
					return time.Time{}, err
				}
				return f0.EventTimeStamp[0].EventTimeStamp()
			},
		}, {
			description: "EventTimeStamp/GTPUPathQoSReport",
			structured: ie.NewGTPUPathQoSReport(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, "some.instance.example"),
				ie.NewEventTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.GTPUPathQoSReport()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeStamp.EventTimeStamp()
			},
		}, {
			description: "EventTimeStamp/QoSMonitoringReport",
			structured: ie.NewQoSMonitoringReport(
				ie.NewQFI(0x01),
				ie.NewEventTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.QoSMonitoringReport()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeStamp.EventTimeStamp()
			},
		}, {
			description: "EventTimeStamp/ClockDriftReport",
			structured: ie.NewClockDriftReport(
				ie.NewTimeDomainNumber(255),
				ie.NewEventTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.ClockDriftReport()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeStamp.EventTimeStamp()
			},
		}, {
			description: "MonitoringTime",
			structured:  ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.MonitoringTime() },
		}, {
			description: "MonitoringTime/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return time.Time{}, err
				}
				return f0.MonitoringTime.MonitoringTime()
			},
		}, {
			description: "MonitoringTime/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return time.Time{}, err
				}
				return f0.MonitoringTime.MonitoringTime()
			},
		}, {
			description: "MonitoringTime/AdditionalMonitoringTime",
			structured: ie.NewAdditionalMonitoringTime(
				ie.NewSubsequentVolumeThreshold(0x01, 0x1111111111111111, 0, 0),
				ie.NewMonitoringTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.AdditionalMonitoringTime()
				if err != nil {
					return time.Time{}, err
				}
				return f0.MonitoringTime.MonitoringTime()
			},
		}, {
			description: "RecoveryTimeStamp",
			structured:  ie.NewRecoveryTimeStamp(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.RecoveryTimeStamp() },
		}, {
			description: "StartTime",
			structured:  ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.StartTime() },
		}, {
			description: "StartTime/UsageReportWithinSessionModificationResponse",
			structured: ie.NewUsageReportWithinSessionModificationResponse(
				ie.NewURRID(0xffffffff),
				ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionModificationResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.StartTime.StartTime()
			},
		}, {
			description: "StartTime/UsageReportWithinSessionDeletionResponse",
			structured: ie.NewUsageReportWithinSessionDeletionResponse(
				ie.NewURRID(0xffffffff),
				ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionDeletionResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.StartTime.StartTime()
			},
		}, {
			description: "StartTime/UsageReportWithinSessionReportRequest",
			structured: ie.NewUsageReportWithinSessionReportRequest(
				ie.NewURRID(0xffffffff),
				ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionReportRequest()
				if err != nil {
					return time.Time{}, err
				}
				return f0.StartTime.StartTime()
			},
		}, {
			description: "StartTime/GTPUPathQoSReport",
			structured: ie.NewGTPUPathQoSReport(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, "some.instance.example"),
				ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.GTPUPathQoSReport()
				if err != nil {
					return time.Time{}, err
				}
				return f0.StartTime.StartTime()
			},
		}, {
			description: "StartTime/QoSMonitoringReport",
			structured: ie.NewQoSMonitoringReport(
				ie.NewQFI(0x01),
				ie.NewStartTime(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.QoSMonitoringReport()
				if err != nil {
					return time.Time{}, err
				}
				return f0.StartTime.StartTime()
			},
		}, {
			description: "TimeOfFirstPacket",
			structured:  ie.NewTimeOfFirstPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.TimeOfFirstPacket() },
		}, {
			description: "TimeOfFirstPacket/UsageReportWithinSessionModificationResponse",
			structured: ie.NewUsageReportWithinSessionModificationResponse(
				ie.NewURRID(0xffffffff),
				ie.NewTimeOfFirstPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionModificationResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeOfFirstPacket.TimeOfFirstPacket()
			},
		}, {
			description: "TimeOfFirstPacket/UsageReportWithinSessionDeletionResponse",
			structured: ie.NewUsageReportWithinSessionDeletionResponse(
				ie.NewURRID(0xffffffff),
				ie.NewTimeOfFirstPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionDeletionResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeOfFirstPacket.TimeOfFirstPacket()
			},
		}, {
			description: "TimeOfFirstPacket/UsageReportWithinSessionReportRequest",
			structured: ie.NewUsageReportWithinSessionReportRequest(
				ie.NewURRID(0xffffffff),
				ie.NewTimeOfFirstPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionReportRequest()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeOfFirstPacket.TimeOfFirstPacket()
			},
		}, {
			description: "TimeOfLastPacket",
			structured:  ie.NewTimeOfLastPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			decoded:     time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) { return i.TimeOfLastPacket() },
		}, {
			description: "TimeOfLastPacket/UsageReportWithinSessionModificationResponse",
			structured: ie.NewUsageReportWithinSessionModificationResponse(
				ie.NewURRID(0xffffffff),
				ie.NewTimeOfLastPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionModificationResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeOfLastPacket.TimeOfLastPacket()
			},
		}, {
			description: "TimeOfLastPacket/UsageReportWithinSessionDeletionResponse",
			structured: ie.NewUsageReportWithinSessionDeletionResponse(
				ie.NewURRID(0xffffffff),
				ie.NewTimeOfLastPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionDeletionResponse()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeOfLastPacket.TimeOfLastPacket()
			},
		}, {
			description: "TimeOfLastPacket/UsageReportWithinSessionReportRequest",
			structured: ie.NewUsageReportWithinSessionReportRequest(
				ie.NewURRID(0xffffffff),
				ie.NewTimeOfLastPacket(time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)),
			),
			decoded: time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC),
			decoderFunc: func(i *ie.IE) (time.Time, error) {
				f0, err := i.UsageReportWithinSessionReportRequest()
				if err != nil {
					return time.Time{}, err
				}
				return f0.TimeOfLastPacket.TimeOfLastPacket()
			},
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
