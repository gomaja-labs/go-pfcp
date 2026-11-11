package spec

import (
	"encoding/json"
	"reflect"
	"testing"
)

func readCatalogue(t *testing.T) *Catalogue {
	t.Helper()
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestCatalogue(t *testing.T) {
	c := readCatalogue(t)
	if err := Validate(c); err != nil {
		t.Fatal(err)
	}
	if len(c.Messages) != 25 || len(c.Grouped) != 123 {
		t.Fatal("incomplete containers")
	}
	rows, multi, shared := 0, 0, 0
	for _, g := range append(c.Messages, c.Grouped...) {
		for _, r := range g.IEs {
			rows++
			if r.Multiple {
				multi++
			}
			if r.SharedTypeGroup != nil {
				shared++
			}
		}
	}
	if rows != 927 || multi != 200 || shared != 28 {
		t.Fatalf("rows=%d multiple=%d shared=%d", rows, multi, shared)
	}
}

func TestCatalogueRejectsMutations(t *testing.T) {
	cases := map[string]func(*Catalogue){
		"missing_reserved_type":    func(c *Catalogue) { c.IETypes = c.IETypes[1:] },
		"unaccounted_range":        func(c *Catalogue) { n := 65534; c.IETypes[len(c.IETypes)-1].TypeEnd = &n },
		"missing_type":             func(c *Catalogue) { c.IETypes = append(c.IETypes[:2], c.IETypes[3:]...) },
		"duplicate_type":           func(c *Catalogue) { c.IETypes = append(c.IETypes, c.IETypes[2]) },
		"unknown_row_type":         func(c *Catalogue) { c.Messages[0].IEs[0].IEType = 10000 },
		"missing_group_definition": func(c *Catalogue) { c.Grouped = c.Grouped[1:] },
		"missing_child_table":      func(c *Catalogue) { s := "Table 0.0-1"; c.Grouped[0].IEs[1].ChildTable = &s },
		"wrong_child_type":         func(c *Catalogue) { s := "Table 7.5.2.2-1"; c.Grouped[0].IEs[1].ChildTable = &s },
		"duplicate_go_name":        func(c *Catalogue) { c.IETypes[2].GoName = c.IETypes[1].GoName },
		"duplicate_go_field":       func(c *Catalogue) { c.Messages[0].IEs[1].GoField = c.Messages[0].IEs[0].GoField },
		"unproven_multiple":        func(c *Catalogue) { c.Messages[0].IEs[0].Multiple = true },
		"long_quote": func(c *Catalogue) {
			r := &c.Messages[0].IEs[0]
			r.Multiple = true
			r.MultipleEvidence = "a a a a a a a a a a a a a a a a a a a a a a a a a a"
		},
		"missing_shared_marker": func(c *Catalogue) {
			for i := range c.Messages {
				for j := range c.Messages[i].IEs {
					if c.Messages[i].IEs[j].SharedTypeGroup != nil {
						c.Messages[i].IEs[j].SharedTypeGroup = nil
						return
					}
				}
			}
		},
	}
	for name, mutate := range cases {
		t.Run(name, func(t *testing.T) {
			c := readCatalogue(t)
			mutate(c)
			if err := Validate(c); err == nil {
				t.Fatal("accepted corrupt catalogue")
			}
		})
	}
}

func TestCatalogueNamesAndContexts(t *testing.T) {
	c := readCatalogue(t)
	s, err := LoadSurface()
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range c.IETypes {
		want := s.IENames[itoa(e.Type)]
		if e.ExistingGoName == nil {
			if want != "" {
				t.Fatal(e.Type)
			}
		} else if *e.ExistingGoName != want {
			t.Fatal(e.Type)
		}
	}
	for _, m := range c.Messages {
		for _, r := range m.IEs {
			old := s.Messages[m.GoName].Bindings[itoa(r.IEType)]
			if current := map[int]string{194: "CreateBridgeRouterInfo", 195: "CreatedBridgeRouterInfo", 201: "TSCManagementInformation"}[r.IEType]; current != "" {
				old = current
			}
			if old != "" && r.GoField != old {
				t.Fatal(m.GoName, r.GoField)
			}
		}
	}
	for _, v := range []struct {
		m   int
		row string
		ty  int
	}{{52, "Update BAR", 86}, {57, "Update BAR", 12}, {53, "Usage Report", 78}, {55, "Usage Report", 79}, {56, "Usage Report", 80}, {52, "TSC Management Information", 199}, {53, "TSC Management Information", 200}, {56, "TSC Management Information", 201}, {53, "Packet Rate Status Report", 264}} {
		found := false
		for _, m := range c.Messages {
			if m.Type == v.m {
				for _, r := range m.IEs {
					if r.Name == v.row {
						found = r.IEType == v.ty
					}
				}
			}
		}
		if !found {
			t.Fatal(v)
		}
	}
	// The ordered model must reproduce canonical JSON, independently of the DOCX.
	b, err := CanonicalJSON(c)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != catalogue {
		t.Fatal("canonical serialization changed catalogue bytes")
	}
	var generic any
	if err := json.Unmarshal(b, &generic); err != nil {
		t.Fatal(err)
	}
}

func TestCatalogueMergedAndInheritedRows(t *testing.T) {
	c := readCatalogue(t)
	messages, grouped := map[int]Container{}, map[int]Container{}
	messageRows, groupedRows := 0, 0
	for _, m := range c.Messages {
		messages[m.Type] = m
		messageRows += len(m.IEs)
	}
	for _, g := range c.Grouped {
		grouped[g.Type] = g
		groupedRows += len(g.IEs)
	}
	inherited := len(grouped[167].IEs) + len(grouped[176].IEs)
	if messageRows != 268 || groupedRows != 659 || messageRows+groupedRows-inherited != 917 {
		t.Fatal("source/inherited row census")
	}
	if !reflect.DeepEqual(grouped[166].IEs, grouped[167].IEs) || !reflect.DeepEqual(grouped[175].IEs, grouped[176].IEs) {
		t.Fatal("incorporated tables changed")
	}
	if len(messages[50].IEs) != 39 || messages[50].IEs[38].IEType != 385 || len(grouped[2].IEs) != 25 || grouped[2].IEs[24].IEType != 395 {
		t.Fatal("lost rows after changed grid spans")
	}
	r := grouped[188].IEs
	if len(r) != 3 || r[1].Name != "Source IP Address" || r[2].Name != r[1].Name || r[1].Presence != "O" || r[2].Presence != "M" || !r[2].Applicability.N4mb || r[2].Applicability.N4 {
		t.Fatal("lost merged alternatives")
	}
}
