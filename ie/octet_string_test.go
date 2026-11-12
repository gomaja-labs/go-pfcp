package ie_test

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestIdentifierStringAPIs(t *testing.T) {
	cases := []struct {
		name  string
		newIE any
		get   any
	}{
		{"NetworkInstance", ie.NewNetworkInstance, (*ie.IE).NetworkInstance},
		{"ApplicationID", ie.NewApplicationID, (*ie.IE).ApplicationID},
		{"ApplicationInstanceID", ie.NewApplicationInstanceID, (*ie.IE).ApplicationInstanceID},
		{"ForwardingPolicyIdentifier", ie.NewForwardingPolicy, (*ie.IE).ForwardingPolicyIdentifier},
		{"ActivatePredefinedRules", ie.NewActivatePredefinedRules, (*ie.IE).ActivatePredefinedRules},
		{"DeactivatePredefinedRules", ie.NewDeactivatePredefinedRules, (*ie.IE).DeactivatePredefinedRules},
		{"PredefinedRulesName", ie.NewPredefinedRulesName, (*ie.IE).PredefinedRulesName},
		{"DataNetworkAccessIdentifier", ie.NewDataNetworkAccessIdentifier, (*ie.IE).DataNetworkAccessIdentifier},
		{"GroupID", ie.NewGroupID, (*ie.IE).GroupID},
		{"UEIPAddressPoolIdentity", ie.NewUEIPAddressPoolIdentity, (*ie.IE).UEIPAddressPoolIdentity},
		{"URI", ie.NewURI, (*ie.IE).URI},
		{"EventNotificationURI", ie.NewEventNotificationURI, (*ie.IE).EventNotificationURI},
		{"NotificationCorrelationID", ie.NewNotificationCorrelationID, (*ie.IE).NotificationCorrelationID},
		{"Metadata", ie.NewMetadata, (*ie.IE).Metadata},
		{"HeaderDetectionReference", ie.NewHeaderDetectionReference, (*ie.IE).HeaderDetectionReference},
		{"HeaderDetectionSupportInformation", ie.NewHeaderDetectionSupportInformation, (*ie.IE).HeaderDetectionSupportInformation},
		{"HeaderHandlingControlReference", ie.NewHeaderHandlingControlReference, (*ie.IE).HeaderHandlingControlReference},
		{"HeaderInformation", ie.NewHeaderInformation, (*ie.IE).HeaderInformation},
		{"HeaderValue", ie.NewHeaderValue, (*ie.IE).HeaderValue},
		{"BindingIndication", ie.NewBindingIndication, (*ie.IE).BindingIndication},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			newIE, ok := tc.newIE.(func(string) *ie.IE)
			if !ok {
				t.Fatalf("constructor accepts %T, want string", tc.newIE)
			}
			get, ok := tc.get.(func(*ie.IE) (string, error))
			if !ok {
				t.Fatalf("accessor is %T, want string", tc.get)
			}
			values := []string{"a", string([]byte{0, 255, 128, '.', 0xc0, 0xaf}), "a\x00b"}
			switch tc.name {
			case "URI", "EventNotificationURI":
				values = []string{"https://example.test/%2f?q=%FF#frag", "urn:example:a"}
			case "Metadata":
				values = []string{"AAEA/w==", "YQ=="}
			}
			for _, value := range values {
				item := newIE(value)
				wire, e := item.Marshal()
				if e != nil {
					t.Fatal(e)
				}
				parsed, e := ie.Parse(wire)
				if e != nil {
					t.Fatal(e)
				}
				got, e := get(parsed)
				if e != nil || got != value {
					t.Fatalf("octet string changed, error %v", e)
				}
				replay, e := parsed.Marshal()
				if e != nil || !bytes.Equal(wire, replay) {
					t.Fatal("wire changed", e)
				}
				for j := range wire {
					wire[j] = 0
				}
				if got != value {
					t.Fatal("returned string aliases receive buffer")
				}
			}
			if _, e := get(nil); !errors.Is(e, ie.ErrIENotFound) {
				t.Fatal("nil receiver", e)
			}
		})
	}
}

func TestStructuredIdentifierStringFields(t *testing.T) {
	for _, tc := range []struct {
		typ   reflect.Type
		field string
	}{
		{reflect.TypeFor[ie.RemoteGTPUPeerFields](), "NetworkInstance"},
		{reflect.TypeFor[ie.DNSQueryResponseFilterFields](), "DomainNamePattern"},
	} {
		f, ok := tc.typ.FieldByName(tc.field)
		if !ok || f.Type != reflect.TypeFor[string]() {
			t.Fatalf("%s.%s must use string", tc.typ, tc.field)
		}
	}
	value := string([]byte{0, 255, 128, 1})
	ctor, ok := any(ie.NewRemoteGTPUPeer).(func(uint8, string, string, uint8, string, time.Time) *ie.IE)
	if !ok {
		t.Fatal("peer network instance must accept string")
	}
	item := ctor(10, "192.0.2.1", "", 0, value, time.Time{})
	wire, e := item.Marshal()
	if e != nil {
		t.Fatal(e)
	}
	parsed, e := ie.Parse(wire)
	if e != nil {
		t.Fatal(e)
	}
	fields, e := parsed.RemoteGTPUPeer()
	if e != nil {
		t.Fatal(e)
	}
	if reflect.ValueOf(fields).Elem().FieldByName("NetworkInstance").String() != value {
		t.Fatal("peer network instance changed")
	}
}
