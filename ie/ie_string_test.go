// Copyright go-pfcp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package ie_test

import (
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
	"github.com/google/go-cmp/cmp"
)

func TestStringIEs(t *testing.T) {
	cases := []struct {
		description string
		structured  *ie.IE
		decoded     string
		decoderFunc func(*ie.IE) (string, error)
	}{
		{
			description: "ActivatePredefinedRules",
			structured:  ie.NewActivatePredefinedRules("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ActivatePredefinedRules() },
		}, {
			description: "APNDNN",
			structured:  ie.NewAPNDNN("some.apn.example"),
			decoded:     "some.apn.example",
			decoderFunc: func(i *ie.IE) (string, error) { return i.APNDNN() },
		}, {
			description: "ApplicationID",
			structured:  ie.NewApplicationID("https://example.com/apps/go-pfcp/"),
			decoded:     "https://example.com/apps/go-pfcp/",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ApplicationID() },
		}, {
			description: "ApplicationInstanceID",
			structured:  ie.NewApplicationInstanceID("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ApplicationInstanceID() },
		}, {
			description: "UserPlaneNodeManagementInformationContainer",
			structured:  ie.NewUserPlaneNodeManagementInformationContainer("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.UserPlaneNodeManagementInformationContainer() },
		}, {
			description: "UserPlaneNodeManagementInformationContainer/TSCManagementInformationWithinSessionModificationRequest",
			structured: ie.NewTSCManagementInformationWithinSessionModificationRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewUserPlaneNodeManagementInformationContainer("go-pfcp"),
			),
			decoded: "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) {
				f0, err := i.TSCManagementInformationWithinSessionModificationRequest()
				if err != nil {
					return "", err
				}
				return f0.UserPlaneNodeManagementInformationContainer.UserPlaneNodeManagementInformationContainer()
			},
		}, {
			description: "UserPlaneNodeManagementInformationContainer/TSCManagementInformationWithinSessionModificationResponse",
			structured: ie.NewTSCManagementInformationWithinSessionModificationResponse(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewUserPlaneNodeManagementInformationContainer("go-pfcp"),
			),
			decoded: "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) {
				f0, err := i.TSCManagementInformationWithinSessionModificationResponse()
				if err != nil {
					return "", err
				}
				return f0.UserPlaneNodeManagementInformationContainer.UserPlaneNodeManagementInformationContainer()
			},
		}, {
			description: "UserPlaneNodeManagementInformationContainer/TSCManagementInformationWithinSessionReportRequest",
			structured: ie.NewTSCManagementInformationWithinSessionReportRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewUserPlaneNodeManagementInformationContainer("go-pfcp"),
			),
			decoded: "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) {
				f0, err := i.TSCManagementInformationWithinSessionReportRequest()
				if err != nil {
					return "", err
				}
				return f0.UserPlaneNodeManagementInformationContainer.UserPlaneNodeManagementInformationContainer()
			},
		}, {
			description: "DataNetworkAccessIdentifier",
			structured:  ie.NewDataNetworkAccessIdentifier("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.DataNetworkAccessIdentifier() },
		}, {
			description: "DeactivatePredefinedRules",
			structured:  ie.NewDeactivatePredefinedRules("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.DeactivatePredefinedRules() },
		}, {
			description: "FlowInformation/FlowDescription",
			structured:  ie.NewFlowInformation(ie.FlowDirectionDownlink, "go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.FlowDescription() },
		}, {
			description: "ForwardingPolicyIdentifier",
			structured:  ie.NewForwardingPolicy("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.ForwardingPolicyIdentifier() },
		}, {
			description: "FramedIPv6Route",
			structured:  ie.NewFramedIPv6Route("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.FramedIPv6Route() },
		}, {
			description: "FramedRoute",
			structured:  ie.NewFramedRoute("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.FramedRoute() },
		}, {
			description: "NetworkInstance",
			structured:  ie.NewNetworkInstance("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstance() },
		}, {
			description: "NetworkInstance/Heuristic",
			structured:  ie.NewNetworkInstance("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstanceHeuristic() },
		}, {
			description: "NetworkInstance/FQDN",
			structured:  ie.NewNetworkInstanceFQDN("go-pfcp.com"),
			decoded:     "go-pfcp.com",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstanceFQDN() },
		}, {
			description: "NetworkInstance/FQDN/Heuristic",
			structured:  ie.NewNetworkInstanceFQDN("go-pfcp.com"),
			decoded:     "go-pfcp.com",
			decoderFunc: func(i *ie.IE) (string, error) { return i.NetworkInstanceHeuristic() },
		}, {
			description: "PortManagementInformationContainer",
			structured:  ie.NewPortManagementInformationContainer("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.PortManagementInformationContainer() },
		}, {
			description: "PortManagementInformationContainer/TSCManagementInformationWithinSessionModificationRequest",
			structured: ie.NewTSCManagementInformationWithinSessionModificationRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewUserPlaneNodeManagementInformationContainer("go-pfcp"),
			),
			decoded: "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) {
				f0, err := i.TSCManagementInformationWithinSessionModificationRequest()
				if err != nil {
					return "", err
				}
				return f0.PortManagementInformationContainer.PortManagementInformationContainer()
			},
		}, {
			description: "PortManagementInformationContainer/TSCManagementInformationWithinSessionModificationResponse",
			structured: ie.NewTSCManagementInformationWithinSessionModificationResponse(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewUserPlaneNodeManagementInformationContainer("go-pfcp"),
			),
			decoded: "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) {
				f0, err := i.TSCManagementInformationWithinSessionModificationResponse()
				if err != nil {
					return "", err
				}
				return f0.PortManagementInformationContainer.PortManagementInformationContainer()
			},
		}, {
			description: "PortManagementInformationContainer/TSCManagementInformationWithinSessionReportRequest",
			structured: ie.NewTSCManagementInformationWithinSessionReportRequest(
				ie.NewPortManagementInformationContainer("go-pfcp"),
				ie.NewUserPlaneNodeManagementInformationContainer("go-pfcp"),
			),
			decoded: "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) {
				f0, err := i.TSCManagementInformationWithinSessionReportRequest()
				if err != nil {
					return "", err
				}
				return f0.PortManagementInformationContainer.PortManagementInformationContainer()
			},
		}, {
			description: "SMFSetID",
			structured:  ie.NewSMFSetID("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.SMFSetID() },
		}, {
			description: "UEIPAddressPoolIdentity",
			structured:  ie.NewUEIPAddressPoolIdentity("go-pfcp"),
			decoded:     "go-pfcp",
			decoderFunc: func(i *ie.IE) (string, error) { return i.UEIPAddressPoolIdentityString() },
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
