package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

// Current fields in TS 29.244 V20.1.0 clauses 8.2.5, 8.2.101,
// 8.2.171, 8.2.70, 8.2.20, 8.2.63 and 8.2.156.
func TestCurrentLeafFieldPreservation(t *testing.T) {
	tests := []struct {
		name    string
		typ     ie.IEType
		payload []byte
		method  string
	}{
		{"IEA018", ie.SDFFilter, []byte{0x20, 0, 1, 0, 5, 1, 0, 0, 0, 9}, "SDFFilter"},
		{"IEA019", ie.UserID, append([]byte{0x20, 15}, []byte("extid-a@example")...), "UserID"},
		{"IEA020", ie.QoSMonitoringMeasurement, []byte{0xf0, 0x27, 0x10, 0, 1, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1}, "QoSMonitoringMeasurement"},
		{"IEA021", ie.RemoteGTPUPeer, []byte{0x12, 192, 0, 2, 1, 0xe9, 0x50, 0, 1}, "RemoteGTPUPeer"},
		{"IEA022", ie.RedirectInformation, []byte{5, 0x1f, 0x90}, "RedirectInformation"},
		{"IEA023", ie.PacketRate, []byte{5, 0, 0, 1, 1, 0, 9}, "PacketRate"},
		{"IEA039", ie.PMFControlInformation, []byte{7, 2, 9, 10}, "PMFControlInformation"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := reflect.ValueOf(ie.New(tt.typ, tt.payload)).MethodByName(tt.method).Call(nil)
			if !result[1].IsNil() {
				t.Fatal(result[1].Interface())
			}
			m, ok := result[0].Interface().(interface{ Marshal() ([]byte, error) })
			if !ok {
				t.Fatal("value helper does not expose structured fields")
			}
			got, err := m.Marshal()
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, tt.payload) {
				t.Fatal("current fields lost")
			}
		})
	}
}

func TestCurrentLeafValueRepresentation(t *testing.T) {
	for _, tt := range []struct {
		name   string
		i      *ie.IE
		method string
		kind   reflect.Kind
	}{
		{"IEA017", ie.New(ie.FiveGSUserPlaneNodeID, []byte{1, 1, 2, 3, 4, 5, 6, 7, 8}), "FiveGSUserPlaneNodeID", reflect.Struct},
		{"IEA024", ie.New(ie.RequestedQoSMonitoring, []byte{0, 3}), "RequestedQoSMonitoring", reflect.Uint16},
		{"IEA036", ie.NewDroppedDLTrafficThreshold(true, true, 123456, 654321), "DroppedDLTrafficThreshold", reflect.Struct},
	} {
		t.Run(tt.name, func(t *testing.T) {
			r := reflect.ValueOf(tt.i).MethodByName(tt.method).Call(nil)
			if !r[1].IsNil() {
				t.Fatal(r[1].Interface())
			}
			v := r[0]
			if v.Kind() == reflect.Pointer {
				v = v.Elem()
			}
			if v.Kind() != tt.kind {
				t.Fatalf("value kind %v; want %v", v.Kind(), tt.kind)
			}
		})
	}
}
func TestCurrentPMFFlags(t *testing.T) {
	v, err := ie.NewPMFControlInformation(3).PMFControlInformation()
	if err != nil || reflect.ValueOf(v).Kind() == reflect.Uint8 && reflect.ValueOf(v).Uint() != 3 {
		t.Fatal("DRTTI lost")
	}
}
func TestCurrentFQCSIDNodeType(t *testing.T) {
	i := ie.NewFQCSID("192.0.2.1", 2, 9)
	if len(i.Payload) != 8 {
		t.Fatal("mandatory Node Type absent")
	}
}
func TestCurrentDownlinkSize(t *testing.T) {
	v, err := ie.New(ie.DownlinkDataServiceInformation, []byte{4, 0, 9}).DLDataPacketsSize()
	if err != nil || v != 9 {
		t.Fatal("nested DL Data Packets Size unavailable", err)
	}
}

func TestCurrentLeafSemanticValues(t *testing.T) {
	t.Run("IEA017", func(t *testing.T) {
		for _, id := range []uint64{0, 1, 0x0102030405060708, ^uint64(0)} {
			i := ie.NewFiveGSUserPlaneNodeID(true, id)
			f, e := i.FiveGSUserPlaneNodeID()
			if e != nil || !f.BID || f.Value != id || len(i.Payload) != 9 {
				t.Fatal("Unsigned64 lost", e)
			}
		}
		f, e := ie.NewFiveGSUserPlaneNodeID(false, 0).FiveGSUserPlaneNodeID()
		if e != nil || f.BID {
			t.Fatal("absent identifier", e)
		}
		for n := 1; n < 9; n++ {
			if _, e := ie.New(ie.FiveGSUserPlaneNodeID, append([]byte{1}, make([]byte, n-1)...)).FiveGSUserPlaneNodeID(); e == nil {
				t.Fatal("accepted truncated ID")
			}
		}
	})
	t.Run("IEA024", func(t *testing.T) {
		for bit := uint16(1); bit <= 0x200; bit <<= 1 {
			v, e := ie.NewRequestedQoSMonitoringFlags(bit).RequestedQoSMonitoring()
			if e != nil || v != bit {
				t.Fatal("request bit lost", e)
			}
		}
		v, e := ie.New(ie.RequestedQoSMonitoring, []byte{0xff, 0xff, 0xee}).RequestedQoSMonitoring()
		if e != nil || v != 0x3ff {
			t.Fatal("spares affected flags", e)
		}
	})
	t.Run("IEA025", func(t *testing.T) {
		f, e := ie.NewPMFControlInformation(3).PMFControlInformation()
		if e != nil || f.Flags != 3 {
			t.Fatal("DRTTI lost", e)
		}
	})
	t.Run("IEA028", func(t *testing.T) {
		for count := 0; count <= 15; count++ {
			ids := make([]uint16, count)
			for k := range ids {
				ids[k] = uint16(k + 1)
			}
			i := ie.NewFQCSID("192.0.2.1", 5, ids...)
			v, e := i.FQCSIDNodeType()
			if e != nil || v != 5 {
				t.Fatal("Node Type lost", e)
			}
			got, e := i.CSIDs()
			if e != nil || !reflect.DeepEqual(got, ids) {
				t.Fatal("CSID count changed", e)
			}
			i.Payload = i.Payload[:len(i.Payload)-1]
			if _, e := i.FQCSID(); e == nil {
				t.Fatal("missing Node Type accepted")
			}
		}
	})
	t.Run("IEA035", func(t *testing.T) {
		for flags := uint8(0); flags < 8; flags++ {
			f := &ie.DownlinkDataServiceInformationFields{Flags: flags}
			if flags&1 != 0 {
				f.PagingPolicyIndication = 7
			}
			if flags&2 != 0 {
				f.QFI = 63
			}
			if flags&4 != 0 {
				f.DLDataPacketsSize = 65535
			}
			i := ie.NewDownlinkDataServiceInformationFromFields(f)
			g, e := i.DownlinkDataServiceInformation()
			if e != nil || !reflect.DeepEqual(f, g) {
				t.Fatal("conditional field changed", e)
			}
			if flags&4 != 0 {
				v, e := i.DLDataPacketsSize()
				if e != nil || v != 65535 {
					t.Fatal("size lost", e)
				}
			}
		}
	})
	t.Run("IEA036", func(t *testing.T) {
		for flags := uint8(0); flags < 4; flags++ {
			p, b := uint64(0), uint64(0)
			if flags&1 != 0 {
				p = ^uint64(0)
			}
			if flags&2 != 0 {
				b = 0x0102030405060708
			}
			i := ie.NewDroppedDLTrafficThreshold(flags&1 != 0, flags&2 != 0, p, b)
			f, e := i.DroppedDLTrafficThreshold()
			if e != nil || f.DownlinkPackets != p || f.DownlinkDataBytes != b || f.Flags != flags {
				t.Fatal("counter lost", e)
			}
		}
	})
	t.Run("IEA039", func(t *testing.T) {
		i := ie.NewPMFControlInformation(7, 0, 9, 63)
		f, e := i.PMFControlInformation()
		if e != nil || !reflect.DeepEqual(f.QFIs, []uint8{0, 9, 63}) {
			t.Fatal("flow list lost", e)
		}
		if _, e := ie.NewPMFControlInformation(7, 64).Marshal(); e == nil {
			t.Fatal("oversize QFI accepted")
		}
	})
}

type leafFieldsCodec interface {
	Marshal() ([]byte, error)
	MarshalTo([]byte) error
	MarshalLen() int
}

func currentLeafCodec(typ uint8, b []byte) (leafFieldsCodec, error) {
	switch typ % 10 {
	case 0:
		return ie.ParseSDFFilterFields(b)
	case 1:
		return ie.ParseUserIDFields(b)
	case 2:
		return ie.ParseQoSMonitoringMeasurementFields(b)
	case 3:
		return ie.ParseRemoteGTPUPeerFields(b)
	case 4:
		return ie.ParseRedirectInformationFields(b)
	case 5:
		return ie.ParsePacketRateFields(b)
	case 6:
		return ie.ParsePMFControlInformationFields(b)
	case 7:
		return ie.ParseDownlinkDataServiceInformationFields(b)
	case 8:
		return ie.ParseDroppedDLTrafficThresholdFields(b)
	default:
		f := &ie.SDFMultiplexedMedia{}
		e := f.UnmarshalBinary(b)
		return f, e
	}
}

func TestCurrentLeafTruncation(t *testing.T) {
	vectors := [][]byte{
		{0x20, 0, 1, 0, 5, 1, 0, 0, 0, 9},
		append([]byte{0x20, 15}, []byte("extid-a@example")...),
		{0xf0, 0x27, 0x10, 0, 1, 1, 2, 3, 4, 5, 5, 4, 3, 2, 1},
		{0x12, 192, 0, 2, 1, 0xe9, 0x50, 0, 1},
		{5, 0x1f, 0x90},
		{7, 0, 0, 1, 1, 0, 2, 2, 0, 3, 3, 0, 4},
		{7, 3, 0, 9, 63},
		{7, 7, 63, 0xff, 0xff},
		{3, 1, 2, 3, 4, 5, 6, 7, 8, 8, 7, 6, 5, 4, 3, 2, 1},
		{31, 1, 2, 3, 4, 96, 1, 'x', 1, 200},
	}
	// Last vector is RTCP only: PTPI is forbidden with RPTPI.
	vectors[9] = []byte{29, 1, 2, 3, 4, 1, 'x', 1, 200}
	for typ, b := range vectors {
		f, e := currentLeafCodec(uint8(typ), b)
		if e != nil {
			t.Fatalf("type %d: %v", typ, e)
		}
		for n := 0; n < len(b); n++ {
			if _, e := currentLeafCodec(uint8(typ), b[:n]); e == nil {
				t.Fatalf("type %d accepted prefix %d", typ, n)
			}
		}
		for n := 0; n < f.MarshalLen(); n++ {
			if e := f.MarshalTo(make([]byte, n)); e == nil {
				t.Fatalf("type %d accepted buffer %d", typ, n)
			}
		}
		ext := append(append([]byte(nil), b...), 0xee, 0xff)
		if _, e := currentLeafCodec(uint8(typ), ext); e != nil {
			t.Fatalf("type %d rejected extensions: %v", typ, e)
		}
	}
}

func FuzzCurrentLeafFields(f *testing.F) {
	for _, b := range [][]byte{{}, {0}, {7, 3, 0, 9, 63}, {0x20, 0, 1, 0, 5, 1, 0, 0, 0, 9}, {5, 0x1f, 0x90}} {
		for typ := uint8(0); typ < 10; typ++ {
			f.Add(typ, b)
		}
	}
	f.Fuzz(func(t *testing.T, typ uint8, b []byte) {
		v, e := currentLeafCodec(typ, b)
		if e != nil {
			return
		}
		wire, e := v.Marshal()
		if e != nil {
			return
		}
		if len(wire) != v.MarshalLen() {
			t.Fatal("encoded length differs")
		}
		if _, e := currentLeafCodec(typ, wire); e != nil {
			t.Fatal("cannot decode serialized fields", e)
		}
		if len(wire) > 0 {
			if e := v.MarshalTo(wire[:len(wire)-1]); e == nil {
				t.Fatal("short buffer accepted")
			}
		}
	})
}

func TestCurrentAdditionalValues(t *testing.T) {
	t.Run("SDF", func(t *testing.T) {
		entries := []ie.SDFMultiplexedMedia{
			{Flags: 1, SynchronizationSource: 0x01020304},
			{Flags: 2, PayloadType: 127},
			{Flags: 12, MediaIdentificationTag: "audio", RTPHeaderExtensionID: 255},
			{Flags: 20, MediaIdentificationTag: "video", RTCPPacketType: 200},
			{Flags: 31 &^ 2, SynchronizationSource: 9, MediaIdentificationTag: "x", RTPHeaderExtensionID: 1, RTCPPacketType: 201},
		}
		f := &ie.SDFFilterFields{Flags: 0x3f, FDLength: 1, FlowDescription: "a", ToSTrafficClass: []byte{0, 255}, SecurityParameterIndex: []byte{0, 1, 2, 3}, FlowLabel: []byte{0, 2, 3}, SDFFilterID: 9, MultiplexedMedia: entries}
		got, e := ie.NewSDFFilterFromFields(f).SDFFilter()
		if e != nil || !reflect.DeepEqual(got, f) {
			t.Fatal("SDF media or legacy fields lost", e)
		}
	})
	t.Run("UserID", func(t *testing.T) {
		f := &ie.UserIDFields{Flags: 0x70, SUPI: "gci-0123", GPSI: "extid-a@example", PEI: "mac-02-00-00-00-00-01"}
		got, e := ie.NewUserIDFromFields(f).UserID()
		if e != nil || !reflect.DeepEqual(got, f) {
			t.Fatal("identities lost", e)
		}
		for _, bad := range []*ie.UserIDFields{{Flags: 16, SUPI: "imsi-123"}, {Flags: 32, GPSI: "msisdn-123"}, {Flags: 64, PEI: "imei-123"}} {
			if _, e := ie.NewUserIDFromFields(bad).Marshal(); e == nil {
				t.Fatal("wrong identity kind accepted")
			}
		}
	})
	t.Run("QoS", func(t *testing.T) {
		f := &ie.QoSMonitoringMeasurementFields{Flags: 0xf7, DownlinkPacketDelay: 1, UplinkPacketDelay: 2, RoundTripPacketDelay: 3, DownlinkCongestionInformation: 10000, UplinkCongestionInformation: 9574, AverageDownlinkDataRate: 1<<40 - 1, AverageUplinkDataRate: 1<<32 + 1}
		got, e := ie.NewQoSMonitoringMeasurementFromFields(f).QoSMonitoringMeasurement()
		if e != nil || !reflect.DeepEqual(f, got) {
			t.Fatal("measurement changed", e)
		}
		f.DownlinkCongestionInformation = 10001
		if _, e := ie.NewQoSMonitoringMeasurementFromFields(f).Marshal(); e == nil {
			t.Fatal("congestion range ignored")
		}
		f.DownlinkCongestionInformation = 0
		f.AverageUplinkDataRate = 1 << 40
		if _, e := ie.NewQoSMonitoringMeasurementFromFields(f).Marshal(); e == nil {
			t.Fatal("bitrate truncated")
		}
	})
	t.Run("Redirect", func(t *testing.T) {
		for typ := uint8(0); typ <= 8; typ++ {
			f := &ie.RedirectInformationFields{RedirectAddressType: typ}
			if typ != 5 {
				f.RedirectServerAddress = "192.0.2.1"
				f.ServerAddrLength = 9
				if typ == 1 || typ == 7 {
					f.RedirectServerAddress = "2001:db8::1"
					f.ServerAddrLength = 11
				}
			}
			if typ == 4 || typ == 8 {
				f.OtherRedirectServerAddress = "2001:db8::1"
				f.OtherServerAddrLength = 11
			}
			if typ >= 5 {
				f.RedirectPort = 65535
			}
			got, e := ie.NewRedirectInformationFromFields(f).RedirectInformation()
			if e != nil || !reflect.DeepEqual(f, got) {
				t.Fatal("address/port changed", typ, e)
			}
		}
	})
	t.Run("PacketRate", func(t *testing.T) {
		f := &ie.PacketRateFields{Flags: 7, UplinkTimeUnit: 0, DownlinkTimeUnit: 1, UplinkPacketRate: 2, DownlinkPacketRate: 3, AdditionalUplinkTimeUnit: 3, AdditionalDownlinkTimeUnit: 4, AdditionalUplinkPacketRate: 65535, AdditionalDownlinkPacketRate: 5}
		got, e := ie.NewPacketRateFromFields(f).PacketRate()
		if e != nil || !reflect.DeepEqual(f, got) {
			t.Fatal("rate pairs changed", e)
		}
	})
}

func TestCurrentLeafRejectsUnflaggedValues(t *testing.T) {
	cases := []struct {
		name string
		i    *ie.IE
	}{
		{"SDF", ie.NewSDFFilterFromFields(&ie.SDFFilterFields{FlowDescription: "a"})},
		{"SDFMedia", ie.NewSDFFilterFromFields(&ie.SDFFilterFields{Flags: 32, MultiplexedMedia: []ie.SDFMultiplexedMedia{{Flags: 1, SynchronizationSource: 1, PayloadType: 2}}})},
		{"UserID", ie.NewUserID(0, "123", "", "", "")},
		{"QoS", ie.NewQoSMonitoringMeasurementFromFields(&ie.QoSMonitoringMeasurementFields{Flags: 7, DownlinkCongestionInformation: 1})},
		{"QoSDelay", ie.NewQoSMonitoringMeasurement(8, 1, 0, 0)},
		{"RemoteAddress", ie.NewRemoteGTPUPeer(2, "192.0.2.1", "2001:db8::1", 0, "", time.Time{})},
		{"RemoteInterface", ie.NewRemoteGTPUPeer(2, "192.0.2.1", "", 1, "", time.Time{})},
		{"RemoteNetwork", ie.NewRemoteGTPUPeer(2, "192.0.2.1", "", 0, "network", time.Time{})},
		{"RemoteTimestamp", ie.NewRemoteGTPUPeerFromFields(&ie.RemoteGTPUPeerFields{Flags: 2, IPv4Address: []byte{192, 0, 2, 1}, RecoveryTimestamp: time.Unix(1700000000, 0)})},
		{"PacketRate", ie.NewPacketRate(1, 1, 1, 2, 2)},
		{"PacketRateAdditional", ie.NewPacketRateFromFields(&ie.PacketRateFields{Flags: 1, UplinkPacketRate: 1, AdditionalUplinkPacketRate: 2})},
		{"Downlink", ie.NewDownlinkDataServiceInformation(false, true, 1, 9)},
		{"DownlinkSize", ie.NewDownlinkDataServiceInformationFromFields(&ie.DownlinkDataServiceInformationFields{Flags: 2, QFI: 9, DLDataPacketsSize: 1})},
		{"PMF", ie.NewPMFControlInformation(1, 9)},
		{"Dropped", ie.NewDroppedDLTrafficThreshold(false, true, 1, 9)},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if tt.i == nil {
				t.Fatal("constructor returned nil")
			}
			if _, e := tt.i.Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
				t.Fatal("unflagged constructor input accepted", e)
			}
		})
	}
}

func TestCurrentLeafReviewBoundaries(t *testing.T) {
	t.Run("RemoteMalformed", func(t *testing.T) {
		for _, i := range []*ie.IE{ie.NewRemoteGTPUPeer(2, "junk", "", 0, "", time.Time{}), ie.NewRemoteGTPUPeer(1, "", "abcdefghijklmnop", 0, "", time.Time{})} {
			if _, e := i.Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
				t.Fatal("malformed address accepted", e)
			}
		}
	})
	t.Run("RemotePresence", func(t *testing.T) {
		if _, e := ie.NewRemoteGTPUPeer(0, "", "", 0, "", time.Time{}).Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
			t.Fatal("missing peer address accepted", e)
		}
	})
	t.Run("PacketRatePresence", func(t *testing.T) {
		for _, flags := range []uint8{0} {
			if _, e := ie.NewPacketRate(flags, 0, 0, 0, 0).Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
				t.Fatal("missing packet rate flags accepted", e)
			}
		}
	})
	t.Run("PacketRateReservedUnits", func(t *testing.T) {
		for unit := uint8(5); unit <= 7; unit++ {
			f, e := ie.ParsePacketRateFields([]byte{7, unit | 0xf8, 0, 1, unit, 0, 2, unit, 0, 3, unit, 0, 4})
			if e != nil || f.UplinkTimeUnit != 0 || f.DownlinkTimeUnit != 0 || f.AdditionalUplinkTimeUnit != 0 || f.AdditionalDownlinkTimeUnit != 0 {
				t.Fatal("reserved unit not interpreted as minute", e)
			}
			for k := 0; k < 4; k++ {
				f := &ie.PacketRateFields{Flags: 7}
				switch k {
				case 0:
					f.UplinkTimeUnit = unit
				case 1:
					f.DownlinkTimeUnit = unit
				case 2:
					f.AdditionalUplinkTimeUnit = unit
				case 3:
					f.AdditionalDownlinkTimeUnit = unit
				}
				if _, e := ie.NewPacketRateFromFields(f).Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
					t.Fatal("constructor accepted reserved unit", unit, k, e)
				}
			}
		}
	})
	t.Run("SDFFlowLabelSpare", func(t *testing.T) {
		if _, e := ie.NewSDFFilter("", nil, nil, []byte{0xf1, 2, 3}, 0).Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
			t.Fatal("constructor accepted spare bits", e)
		}
		i := ie.New(ie.SDFFilter, []byte{8, 0, 0xf1, 2, 3})
		before := append([]byte(nil), i.Payload...)
		f, e := i.SDFFilter()
		if e != nil || !bytes.Equal(f.FlowLabel, []byte{1, 2, 3}) || !bytes.Equal(i.Payload, before) {
			t.Fatal("spares evaluated or original payload mutated", e)
		}
	})
}

func TestCurrentRedirectAddressValidation(t *testing.T) {
	for _, tt := range []struct {
		typ  uint8
		addr string
	}{{0, ""}, {0, "not-an-ip"}, {0, string([]byte{255})}, {0, "2001:db8::1"}, {1, "192.0.2.1"}, {1, "2001:0db8:0:0:0:0:0:1"}, {2, ""}, {3, string([]byte{255})}} {
		if _, e := ie.NewRedirectInformation(tt.typ, tt.addr).Marshal(); !errors.Is(e, ie.ErrInvalidValue) {
			t.Fatal("invalid sender address accepted", tt.typ, e)
		}
	}
	for _, addr := range []string{"", "junk", string([]byte{255}), "2001:db8::1"} {
		b := append([]byte{0, 0, byte(len(addr))}, []byte(addr)...)
		if _, e := ie.ParseRedirectInformationFields(b); e == nil {
			t.Fatal("invalid received IPv4 address accepted")
		}
	}
	f := &ie.RedirectInformationFields{RedirectAddressType: 1, ServerAddrLength: 23, RedirectServerAddress: "2001:0db8:0:0:0:0:0:1"}
	f.ServerAddrLength = uint16(len(f.RedirectServerAddress))
	b := append([]byte{1, 0, byte(f.ServerAddrLength)}, []byte(f.RedirectServerAddress)...)
	got, e := ie.ParseRedirectInformationFields(b)
	if e != nil || got.RedirectServerAddress != f.RedirectServerAddress {
		t.Fatal("legitimate received IPv6 spelling rejected/changed", e)
	}
}
