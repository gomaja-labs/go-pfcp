package ie_test

import (
	"reflect"
	"testing"

	"github.com/gomaja-labs/go-pfcp/ie"
)

func TestCurrentFlagHelpers(t *testing.T) {
	for _, tc := range []struct {
		typ        ie.IEType
		octet, bit int
		name       string
	}{
		{178, 0, 2, "HasPPE"},
		{43, 0, 0, "HasBUCP"},
		{43, 0, 1, "HasDDND"},
		{43, 0, 2, "HasDLBD"},
		{43, 0, 3, "HasTRST"},
		{43, 0, 4, "HasFTUP"},
		{43, 0, 5, "HasPFDM"},
		{43, 0, 6, "HasHEEU"},
		{43, 0, 7, "HasTREU"},
		{43, 1, 0, "HasEMPU"},
		{43, 1, 1, "HasPDIU"},
		{43, 1, 2, "HasUDBC"},
		{43, 1, 3, "HasQUOAC"},
		{43, 1, 4, "HasTRACE"},
		{43, 1, 5, "HasFRRT"},
		{43, 1, 6, "HasPFDE"},
		{43, 1, 7, "HasEPFAR"},
		{43, 2, 0, "HasDPDRA"},
		{43, 2, 1, "HasADPDP"},
		{43, 2, 2, "HasUEIP"},
		{43, 2, 3, "HasSSET"},
		{43, 2, 4, "HasMNOP"},
		{43, 2, 5, "HasMTE"},
		{43, 2, 6, "HasBUNDL"},
		{43, 2, 7, "HasGCOM"},
		{43, 3, 0, "HasMPAS"},
		{43, 3, 1, "HasRTTL"},
		{43, 3, 2, "HasVTIME"},
		{43, 3, 3, "HasNORP"},
		{43, 3, 4, "HasIPTV"},
		{43, 3, 5, "HasIP6PL"},
		{43, 3, 6, "HasTSN"},
		{43, 3, 7, "HasMPTCP"},
		{43, 4, 0, "HasATSSSLL"},
		{43, 4, 1, "HasQFQM"},
		{43, 4, 2, "HasGPQM"},
		{43, 4, 3, "HasMTEDT"},
		{43, 4, 4, "HasCIOT"},
		{43, 4, 5, "HasETHAR"},
		{43, 4, 6, "HasDDDS"},
		{43, 4, 7, "HasRDS"},
		{43, 5, 0, "HasRTTWP"},
		{43, 5, 1, "HasQUASF"},
		{43, 5, 2, "HasNSPOC"},
		{43, 5, 3, "HasL2TP"},
		{43, 5, 4, "HasUPBER"},
		{43, 5, 5, "HasRESPS"},
		{43, 5, 6, "HasIPREP"},
		{43, 5, 7, "HasDNSTS"},
		{43, 6, 0, "HasDRQOS"},
		{43, 6, 1, "HasMBSN4"},
		{43, 6, 2, "HasPSUPRM"},
		{43, 6, 3, "HasEPPPI"},
		{43, 6, 4, "HasRATP"},
		{43, 6, 5, "HasUPIDP"},
		{43, 6, 6, "HasAFSFC"},
		{43, 6, 7, "HasMPQUICUDP"},
		{43, 7, 0, "HasREDSM"},
		{43, 7, 1, "HasDBDM"},
		{43, 7, 2, "HasTSCTS"},
		{43, 7, 3, "HasDRTSC"},
		{43, 7, 4, "HasN6JEDB"},
		{43, 7, 5, "HasQMCON"},
		{43, 7, 6, "HasDETNET"},
		{43, 7, 7, "HasEML4S"},
		{43, 8, 0, "HasPDUSM"},
		{43, 8, 1, "HasCNTL"},
		{43, 8, 2, "HasQMDRM"},
		{43, 8, 3, "HasEDBNC"},
		{43, 8, 4, "HasMTSDT"},
		{43, 8, 5, "HasUPSBIES"},
		{43, 8, 6, "HasUMN6IP"},
		{43, 8, 7, "HasUN6TU"},
		{43, 9, 0, "HasMBSCH"},
		{43, 9, 1, "HasUN6DM"},
		{43, 9, 2, "HasNATPUB"},
		{43, 9, 3, "HasUSHPH"},
		{43, 9, 4, "HasMPQUICIP"},
		{43, 9, 5, "HasMPQUICE"},
		{43, 9, 6, "HasDYNTR"},
		{43, 9, 7, "HasMUXMF"},
		{43, 10, 0, "HasCONUDP"},
		{43, 10, 1, "HasMOQ"},
		{43, 10, 2, "HasULM"},
		{43, 10, 3, "HasPSITLM"},
		{43, 10, 4, "HasUDPOPT"},
		{43, 10, 5, "HasQMABR"},
		{43, 10, 6, "HasPAPFD"},
		{43, 10, 7, "HasNLMAT"},
		{43, 11, 0, "HasUHRR"},
		{43, 11, 1, "HasSCONE"},
		{89, 0, 0, "HasLOAD"},
		{89, 0, 1, "HasOVRL"},
		{89, 0, 2, "HasEPFAR"},
		{89, 0, 3, "HasSSET"},
		{89, 0, 4, "HasBUNDL"},
		{89, 0, 5, "HasMPAS"},
		{89, 0, 6, "HasARDR"},
		{89, 0, 7, "HasUIAUR"},
		{89, 1, 0, "HasPSUCC"},
		{89, 1, 1, "HasRPGUR"},
		{89, 1, 2, "HasPAPFD"},
		{39, 0, 0, "HasDLDR"},
		{39, 0, 1, "HasUSAR"},
		{39, 0, 2, "HasERIR"},
		{39, 0, 3, "HasUPIR"},
		{39, 0, 4, "HasTMIR"},
		{39, 0, 5, "HasSESR"},
		{39, 0, 6, "HasUISR"},
		{49, 0, 0, "HasDROBU"},
		{49, 0, 1, "HasSNDEM"},
		{49, 0, 2, "HasQAURR"},
		{49, 0, 3, "HasSUMPC"},
		{49, 0, 4, "HasRUMUC"},
		{49, 0, 5, "HasDETEID"},
		{49, 0, 6, "HasHRSBOM"},
		{100, 0, 0, "HasMBQE"},
		{100, 0, 1, "HasINAM"},
		{100, 0, 2, "HasRADI"},
		{100, 0, 3, "HasISTM"},
		{100, 0, 4, "HasMNOP"},
		{100, 0, 5, "HasSSPOC"},
		{100, 0, 6, "HasASPOC"},
		{100, 0, 7, "HasCIAM"},
		{101, 0, 0, "HasUPFR"},
		{101, 0, 1, "HasUPRR"},
		{101, 0, 2, "HasCKDR"},
		{101, 0, 3, "HasGPQR"},
		{101, 0, 4, "HasPURR"},
		{101, 0, 5, "HasVSR"},
		{101, 0, 6, "HasN6DR"},
		{110, 0, 0, "HasAOCI"},
		{161, 0, 0, "HasPSDBU"},
		{161, 0, 1, "HasPPPCP"},
		{184, 0, 0, "HasPSREI"},
		{184, 0, 1, "HasUUPSI"},
		{186, 0, 0, "HasRESTI"},
		{186, 0, 1, "HasSUMPC"},
		{186, 0, 2, "HasHRSBOM"},
		{186, 0, 3, "HasEPRI"},
		{194, 0, 0, "HasBII"},
		{194, 0, 1, "HasRII"},
		{265, 0, 0, "HasMTAI"},
		{265, 0, 1, "HasMQAI"},
	} {
		t.Run(tc.typ.String()+"/"+tc.name, func(t *testing.T) {
			invoke := func(i *ie.IE) bool {
				t.Helper()
				m := reflect.ValueOf(i).MethodByName(tc.name)
				if !m.IsValid() {
					t.Fatalf("missing helper %s", tc.name)
				}
				return m.Call(nil)[0].Bool()
			}
			if invoke(nil) || invoke(ie.New(0xffff, []byte{255})) || invoke(ie.New(tc.typ, nil)) {
				t.Fatal("nil, wrong type or null IE reported set flag")
			}
			for bit := 0; bit < 8; bit++ {
				p := make([]byte, tc.octet+1)
				p[tc.octet] = 1 << bit
				if got := invoke(ie.New(tc.typ, p)); got != (bit == tc.bit) {
					t.Fatalf("one hot bit %d returned %v", bit, got)
				}
			}
			if tc.octet > 0 && invoke(ie.New(tc.typ, make([]byte, tc.octet))) {
				t.Fatal("truncated payload reported flag")
			}
		})
	}
}
