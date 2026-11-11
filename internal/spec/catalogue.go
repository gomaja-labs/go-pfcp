// Package spec describes the version-pinned TS 29.244 catalogue used by pfcpgen.
package spec

import (
	_ "embed"
	"encoding/json"
)

//go:embed ts29244-v20.1.0.json
var catalogue string

//go:embed go-surface.json
var surface string

// Catalogue holds the published registry and message/grouped content tables.
// Field order preserves the canonical extraction's JSON representation.
type Catalogue struct {
	Spec     Document    `json:"spec"`
	IETypes  []IEType    `json:"ie_types"`
	Messages []Container `json:"messages"`
	Grouped  []Container `json:"grouped"`
}

// Document identifies the exact normative input.
type Document struct {
	Document      string `json:"document"`
	Version       string `json:"version"`
	Date          string `json:"date"`
	Archive       string `json:"archive"`
	ArchiveSHA256 string `json:"archive_sha256"`
	DocxSHA256    string `json:"docx_sha256"`
}

// IEType retains literal registry references alongside actual definitions.
type IEType struct {
	Type              int      `json:"type"`
	Name              string   `json:"name"`
	GoName            string   `json:"go_name"`
	Kind              string   `json:"kind"`
	DefinedIn         []string `json:"defined_in"`
	FixedOctets       *int     `json:"fixed_octets"`
	ExistingGoName    *string  `json:"existing_go_name"`
	Note              string   `json:"note"`
	RegistryDefinedIn []string `json:"registry_defined_in,omitempty"`
	TypeEnd           *int     `json:"type_end,omitempty"`
	Variants          *bool    `json:"variants,omitempty"`
}

// Container is a message or one context of a grouped IE. Message and grouped
// JSON members have different order, handled by MarshalJSON.
type Container struct {
	Type     int     `json:"type"`
	Name     string  `json:"name"`
	GoName   string  `json:"go_name"`
	Category string  `json:"category,omitempty"`
	Context  string  `json:"context,omitempty"`
	Table    *string `json:"table"`
	IEs      []Row   `json:"ies"`
	Clause   string  `json:"clause"`
}

// Row preserves a physical row, including conditional same-type variants.
type Row struct {
	Name             string        `json:"name"`
	Presence         string        `json:"presence"`
	IEType           int           `json:"ie_type"`
	IETypeName       string        `json:"ie_type_name"`
	Multiple         bool          `json:"multiple"`
	MultipleEvidence string        `json:"multiple_evidence"`
	Applicability    Applicability `json:"applicability"`
	ChildTable       *string       `json:"child_table"`
	GoField          string        `json:"go_field"`
	SharedTypeGroup  *string       `json:"shared_type_group"`
	SourceIETypeName *string       `json:"source_ie_type_name,omitempty"`
	ResolutionNote   string        `json:"resolution_note,omitempty"`
}

// Applicability is the union of interface marks within a source row.
type Applicability struct {
	Sxa  bool `json:"Sxa"`
	Sxb  bool `json:"Sxb"`
	Sxc  bool `json:"Sxc"`
	N4   bool `json:"N4"`
	N4mb bool `json:"N4mb"`
}

// Surface pins the original exported names; it is not derived from generated code.
type Surface struct {
	IENames  map[string]string `json:"ie_names"`
	Messages map[string]struct {
		Bindings map[string]string `json:"bindings"`
	} `json:"messages"`
}

// Load reads the committed catalogue.
func Load() (*Catalogue, error) {
	var c Catalogue
	err := json.Unmarshal([]byte(catalogue), &c)
	return &c, err
}

// LoadSurface reads the original Go name snapshot at revision 2fbe893.
func LoadSurface() (*Surface, error) {
	var s Surface
	err := json.Unmarshal([]byte(surface), &s)
	return &s, err
}

// MarshalJSON preserves the source extractor's grouped member ordering.
func (c Container) MarshalJSON() ([]byte, error) {
	if c.Category != "" {
		type plain Container
		return json.Marshal(plain(c))
	}
	return json.Marshal(struct {
		Type    int     `json:"type"`
		Name    string  `json:"name"`
		GoName  string  `json:"go_name"`
		Context string  `json:"context"`
		Table   *string `json:"table"`
		IEs     []Row   `json:"ies"`
		Clause  string  `json:"clause"`
	}{c.Type, c.Name, c.GoName, c.Context, c.Table, c.IEs, c.Clause})
}
