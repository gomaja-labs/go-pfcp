// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"encoding/binary"
	"net"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/google/go-cmp/cmp"
)

func TestUint16IEs(t *testing.T) {
	cases := []struct {
		description string
		structured  *ie.IE
		decoded     uint16
		decoderFunc func(*ie.IE) (uint16, error)
	}{
		{
			description: "AdditionalUsageReportsInformation",
			structured:  ie.NewAdditionalUsageReportsInformation(0x00ff),
			decoded:     0x00ff,
			decoderFunc: func(i *ie.IE) (uint16, error) { return i.AdditionalUsageReportsInformation() },
		}, {
			description: "DLBufferingSuggestedPacketCount",
			structured:  ie.NewDLBufferingSuggestedPacketCount(0xffff),
			decoded:     0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) { return i.DLBufferingSuggestedPacketCount() },
		}, {
			description: "DLBufferingSuggestedPacketCount/UpdateBARWithinSessionReportResponse",
			structured: ie.NewUpdateBARWithinSessionReportResponse(
				ie.NewBARID(0xff),
				ie.NewDownlinkDataNotificationDelay(100*time.Millisecond),
				ie.NewDLBufferingDuration(30*time.Second),
				ie.NewDLBufferingSuggestedPacketCount(0xffff),
				ie.NewSuggestedBufferingPacketsCount(0x01),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UpdateBARWithinSessionReportResponse()
				if err != nil {
					return 0, err
				}
				return f0.DLBufferingSuggestedPacketCount.DLBufferingSuggestedPacketCount()
			},
		}, {
			description: "DLDataPacketsSize",
			structured:  ie.NewDLDataPacketsSize(0xffff),
			decoded:     0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) { return i.DLDataPacketsSize() },
		}, {
			description: "DLDataPacketsSize/DownlinkDataReport",
			structured: ie.NewDownlinkDataReport(
				ie.NewPDRID(0xffff),
				ie.NewDLDataPacketsSize(0xffff),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.DownlinkDataReport()
				if err != nil {
					return 0, err
				}
				return f0.DLDataPacketsSize.DLDataPacketsSize()
			},
		}, {
			description: "Ethertype",
			structured:  ie.NewEthertype(0xffff),
			decoded:     0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) { return i.Ethertype() },
		}, {
			description: "Ethertype/PDI",
			structured: ie.NewPDI(
				ie.NewSourceInterface(ie.SrcInterfaceAccess),
				ie.NewEthertype(0xffff),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.PDI()
				if err != nil {
					return 0, err
				}
				return f0.OtherIEs[0].Ethertype()
			},
		}, {
			description: "Ethertype/EthernetPacketFilter",
			structured: ie.NewEthernetPacketFilter(
				ie.NewEthernetFilterID(0xffffffff),
				ie.NewEthertype(0xffff),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.EthernetPacketFilter()
				if err != nil {
					return 0, err
				}
				return f0.Ethertype.Ethertype()
			},
		}, {
			description: "MARID",
			structured:  ie.NewMARID(0x1111),
			decoded:     0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) { return i.MARID() },
		}, {
			description: "MARID/CreatePDR",
			structured: ie.NewCreatePDR(
				ie.NewPDRID(0xffff),
				ie.NewMARID(0x1111),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.CreatePDR()
				if err != nil {
					return 0, err
				}
				return f0.MARID.MARID()
			},
		}, {
			description: "MARID/CreateMAR",
			structured: ie.NewCreateMAR(
				ie.NewMARID(0x1111),
				ie.NewSteeringFunctionality(ie.SteeringFunctionalityATSSSLL),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.CreateMAR()
				if err != nil {
					return 0, err
				}
				return f0.MARID.MARID()
			},
		}, {
			description: "MARID/RemoveMAR",
			structured: ie.NewRemoveMAR(
				ie.NewMARID(0x1111),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.RemoveMAR()
				if err != nil {
					return 0, err
				}
				return f0.MARID.MARID()
			},
		}, {
			description: "MARID/UpdateMAR",
			structured: ie.NewUpdateMAR(
				ie.NewMARID(0x1111),
				ie.NewSteeringFunctionality(ie.SteeringFunctionalityATSSSLL),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UpdateMAR()
				if err != nil {
					return 0, err
				}
				return f0.MARID.MARID()
			},
		}, {
			description: "NumberOfReports",
			structured:  ie.NewNumberOfReports(0x1111),
			decoded:     0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) { return i.NumberOfReports() },
		}, {
			description: "NumberOfReports/CreateURR",
			structured: ie.NewCreateURR(
				ie.NewURRID(0xffffffff),
				ie.NewNumberOfReports(0x1111),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.CreateURR()
				if err != nil {
					return 0, err
				}
				return f0.NumberOfReports.NumberOfReports()
			},
		}, {
			description: "NumberOfReports/UpdateURR",
			structured: ie.NewUpdateURR(
				ie.NewURRID(0xffffffff),
				ie.NewNumberOfReports(0x1111),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UpdateURR()
				if err != nil {
					return 0, err
				}
				return f0.NumberOfReports.NumberOfReports()
			},
		}, {
			description: "PDRID",
			structured:  ie.NewPDRID(0xffff),
			decoded:     0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) { return i.PDRID() },
		}, {
			description: "PDRID/CreatePDR",
			structured: ie.NewCreatePDR(
				ie.NewPDRID(0xffff),
				ie.NewPrecedence(0x11111111),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.CreatePDR()
				if err != nil {
					return 0, err
				}
				return f0.PDRID.PDRID()
			},
		}, {
			description: "PDRID/UpdatePDR",
			structured: ie.NewUpdatePDR(
				ie.NewPDRID(0xffff),
				ie.NewPrecedence(0x11111111),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UpdatePDR()
				if err != nil {
					return 0, err
				}
				return f0.PDRID.PDRID()
			},
		}, {
			description: "PDRID/RemovePDR",
			structured: ie.NewRemovePDR(
				ie.NewPDRID(0xffff),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.RemovePDR()
				if err != nil {
					return 0, err
				}
				return f0.PDRID.PDRID()
			},
		}, {
			description: "PDRID/CreatedPDR",
			structured: ie.NewCreatedPDR(
				ie.NewPDRID(0xffff),
				ie.NewFTEID(0x01, 0x11111111, net.ParseIP("127.0.0.1"), nil, 0),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.CreatedPDR()
				if err != nil {
					return 0, err
				}
				return f0.PDRID.PDRID()
			},
		}, {
			description: "PDRID/ApplicationDetectionInformation",
			structured: ie.NewApplicationDetectionInformation(
				ie.NewApplicationID([]byte("https://example.com/apps/go-pfcp/")),
				ie.NewPDRID(0xffff),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.ApplicationDetectionInformation()
				if err != nil {
					return 0, err
				}
				return f0.PDRID.PDRID()
			},
		}, {
			description: "PDRID/UsageReportWithinSessionReportRequest",
			structured: ie.NewUsageReportWithinSessionReportRequest(
				ie.NewURRID(0xffffffff),
				ie.NewApplicationDetectionInformation(
					ie.NewApplicationID([]byte("https://example.com/apps/go-pfcp/")),
					ie.NewPDRID(0xffff),
				),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UsageReportWithinSessionReportRequest()
				if err != nil {
					return 0, err
				}
				f1, err := f0.ApplicationDetectionInformation.ApplicationDetectionInformation()
				if err != nil {
					return 0, err
				}
				return f1.PDRID.PDRID()
			},
		}, {
			description: "PDRID/DownlinkDataReport",
			structured: ie.NewDownlinkDataReport(
				ie.NewPDRID(0xffff),
				ie.NewDownlinkDataServiceInformation(true, true, 7, 63),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.DownlinkDataReport()
				if err != nil {
					return 0, err
				}
				return f0.PDRID[0].PDRID()
			},
		}, {
			description: "PDRID/UpdatedPDR",
			structured: ie.NewUpdatedPDR(
				ie.NewPDRID(0xffff),
				ie.NewFTEID(0x01, 0x11111111, net.ParseIP("127.0.0.1"), nil, 0),
			),
			decoded: 0xffff,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UpdatedPDR()
				if err != nil {
					return 0, err
				}
				return f0.PDRID.PDRID()
			},
		}, {
			description: "TransportLevelMarking",
			structured:  ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			decoded:     0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				v, e := i.TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
			},
		}, {
			description: "TransportLevelMarking/ForwardingParameters",
			structured: ie.NewForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.ForwardingParameters()
				if err != nil {
					return 0, err
				}
				v, e := f0.TransportLevelMarking.TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
			},
		}, {
			description: "TransportLevelMarking/UpdateForwardingParameters",
			structured: ie.NewUpdateForwardingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UpdateForwardingParameters()
				if err != nil {
					return 0, err
				}
				v, e := f0.TransportLevelMarking.TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
			},
		}, {
			description: "TransportLevelMarking/DuplicatingParameters",
			structured: ie.NewDuplicatingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.DuplicatingParameters()
				if err != nil {
					return 0, err
				}
				v, e := f0.TransportLevelMarking.TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
			},
		}, {
			description: "TransportLevelMarking/UpdateDuplicatingParameters",
			structured: ie.NewUpdateDuplicatingParameters(
				ie.NewDestinationInterface(ie.DstInterfaceAccess),
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.UpdateDuplicatingParameters()
				if err != nil {
					return 0, err
				}
				v, e := f0.TransportLevelMarking.TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
			},
		}, {
			description: "TransportLevelMarking/GTPUPathQoSControlInformation",
			structured: ie.NewGTPUPathQoSControlInformation(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.GTPUPathQoSControlInformation()
				if err != nil {
					return 0, err
				}
				v, e := f0.DSCP[0].TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
			},
		}, {
			description: "TransportLevelMarking/GTPUPathQoSReport",
			structured: ie.NewGTPUPathQoSReport(
				ie.NewRemoteGTPUPeer(0x0e, "127.0.0.1", "", ie.DstInterfaceAccess, []byte("some.instance.example"), time.Time{}),
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.GTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				v, e := f0.OtherIEs[0].TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
			},
		}, {
			description: "TransportLevelMarking/QoSInformationInGTPUPathQoSReport",
			structured: ie.NewQoSInformationInGTPUPathQoSReport(
				ie.NewAveragePacketDelay(10*time.Second),
				ie.NewTransportLevelMarking([]byte{0x11, 0x11}),
			),
			decoded: 0x1111,
			decoderFunc: func(i *ie.IE) (uint16, error) {
				f0, err := i.QoSInformationInGTPUPathQoSReport()
				if err != nil {
					return 0, err
				}
				v, e := f0.DSCP.TransportLevelMarking()
				if e != nil {
					return 0, e
				}
				return binary.BigEndian.Uint16(v), nil
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
