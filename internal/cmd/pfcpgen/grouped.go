package main

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"

	"github.com/gomaja-labs/go-pfcp/internal/spec"
)

type groupedData struct {
	Name, Tables string
	Fields       []field
}

// groupedDefinitions preserves the ordering constraints of every defining
// table. Reusing a row in another context does not
// make it a list. TS 29.244 V20.1.0 Tables 7.5.7.3-1 and 7.5.8.6-1 (type 214),
// Tables 7.5.2.9-4 and 7.5.4.18-2 (type 295).
func groupedDefinitions(c *spec.Catalogue, types map[int]spec.IEType) ([]groupedData, error) {
	var result []groupedData
	index := map[int]int{}
	contexts := map[int]int{}
	for _, g := range c.Grouped {
		contexts[g.Type]++
	}
	for _, g := range c.Grouped {
		n, ok := index[g.Type]
		if !ok {
			n = len(result)
			index[g.Type] = n
			result = append(result, groupedData{Name: constant(types[g.Type])})
		}
		d := &result[n]
		if d.Tables != "" {
			d.Tables += "; "
		}
		d.Tables += *g.Table + " (" + g.Context + ")"
		for _, f := range fields(g, types) {
			if contexts[g.Type] > 1 && !strings.Contains(f.Doc, "Context:") {
				f.Doc += " TS 29.244 V20.1.0 " + *g.Table + ". Context: " + g.Context + "."
			}
			found := false
			for j := range d.Fields {
				if d.Fields[j].Type == f.Type {
					d.Fields[j].Slice = d.Fields[j].Slice || f.Slice
					d.Fields[j].Doc += " " + f.Doc
					found = true
					break
				}
			}
			if !found {
				d.Fields = append(d.Fields, f)
			}
		}
	}
	for n := range result {
		g := &result[n]
		// A stable topological merge preserves both Session Report contexts:
		// SRR ID precedes Access Availability Report, which precedes QoS reports.
		edges := map[int]map[int]bool{}
		indegree := map[int]int{}
		for _, table := range c.Grouped {
			if constant(types[table.Type]) != g.Name {
				continue
			}
			seen := map[int]bool{}
			previous := 0
			for _, row := range table.IEs {
				if seen[row.IEType] {
					continue
				}
				seen[row.IEType] = true
				if previous != 0 {
					if edges[previous] == nil {
						edges[previous] = map[int]bool{}
					}
					if !edges[previous][row.IEType] {
						edges[previous][row.IEType] = true
						indegree[row.IEType]++
					}
				}
				previous = row.IEType
			}
		}
		ordered := make([]field, 0, len(g.Fields))
		used := map[int]bool{}
		for len(ordered) < len(g.Fields) {
			found := false
			for _, f := range g.Fields {
				if used[f.Type] || indegree[f.Type] != 0 {
					continue
				}
				ordered = append(ordered, f)
				used[f.Type] = true
				for next := range edges[f.Type] {
					indegree[next]--
				}
				found = true
				break
			}
			if !found {
				return nil, fmt.Errorf("inconsistent table order for %s", g.Name)
			}
		}
		g.Fields = ordered
	}
	return result, nil
}

func generateGrouped(c *spec.Catalogue, types map[int]spec.IEType, out map[string][]byte) error {
	tmpl, err := template.New("grouped").Parse(groupedTemplate)
	if err != nil {
		return err
	}
	groups, err := groupedDefinitions(c, types)
	if err != nil {
		return err
	}
	for _, g := range groups {
		var b bytes.Buffer
		b.WriteString(banner)
		if err := tmpl.Execute(&b, g); err != nil {
			return err
		}
		src, err := format.Source(b.Bytes())
		if err != nil {
			return fmt.Errorf("%s: %w", g.Name, err)
		}
		out["ie/"+fileName(g.Name+"Fields")] = src
	}
	var b bytes.Buffer
	b.WriteString(banner + "\npackage ie_test\n\nimport \"github.com/gomaja-labs/go-pfcp/ie\"\n\nfunc catalogueGrouped() []groupedCase {return []groupedCase{\n")
	for _, g := range groups {
		fmt.Fprintf(&b, "{typ:ie.%s, newIE:ie.New%s, empty:func() groupedFields {return &ie.%sFields{}}, access:func(i *ie.IE)(groupedFields,error){return i.%s()}, parse:func(b []byte)(groupedFields,error){return ie.Parse%sFields(b)}},\n", g.Name, g.Name, g.Name, g.Name, g.Name)
	}
	b.WriteString("}}\n")
	src, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}
	out["ie/grouped_catalogue_gen_test.go"] = src
	// A generated declaration's comment must remain a single logical line.
	for _, g := range groups {
		for _, f := range g.Fields {
			if strings.ContainsAny(f.Doc, "\n\r") {
				return fmt.Errorf("multiline grouped field comment: %s.%s", g.Name, f.Name)
			}
		}
	}
	return nil
}

const groupedTemplate = `
package ie

// {{.Name}}Fields represents the children of {{.Name}}.
// TS 29.244 V20.1.0 {{.Tables}}.
// Presence and conditional applicability remain the application's responsibility.
// IE builds children in table order, followed by OtherIEs; use the original IE
// to retain a received ordering. Repeated rows across contexts share one field.
type {{.Name}}Fields struct {
 {{range .Fields}}// {{.Name}} is the {{.Doc}}
 {{.Name}} {{if .Slice}}[]{{end}}*IE
 {{end}}// OtherIEs retains unknown types and repetitions beyond a singular field's
 // first occurrence. TS 29.244 V20.1.0 clauses 7.6.9 and 7.6.10.
 OtherIEs []*IE
}

// New{{.Name}} creates a {{.Name}} IE, skipping nil children.
func New{{.Name}}(ies ...*IE) *IE {return newGroupedIE({{.Name}},0,ies...)}

// {{.Name}} returns this IE's typed children. It accepts only {{.Name}}.
func (i *IE) {{.Name}}() (*{{.Name}}Fields,error) {
 if i==nil {return nil,ErrIENotFound}
 if i.Type!={{.Name}} {return nil,&InvalidTypeError{Type:i.Type}}
 children,err:=groupedChildren(i);if err!=nil {return nil,withDecodePath(i.Type,err)}
 f:=&{{.Name}}Fields{};f.addIEs(children);return f,nil
}

// IE builds a grouped IE in table order, followed by OtherIEs, skipping nils.
// A nil receiver builds an IE whose encoding reports ErrIENotFound.
func (f *{{.Name}}Fields) IE() *IE {
 if f==nil {return newInvalidIE({{.Name}},ErrIENotFound)}
 return New{{.Name}}(f.children()...)
}

// Parse{{.Name}}Fields decodes the value (concatenated child IEs, without a parent header).
// Decoding may borrow b; keep b unchanged while the result is in use.
func Parse{{.Name}}Fields(b []byte) (*{{.Name}}Fields,error) {
 f:=&{{.Name}}Fields{};if err:=f.UnmarshalBinary(b);err!=nil {return nil,err};return f,nil
}

// UnmarshalBinary replaces the fields with decoded children, accepting any order.
// A decoding error leaves the receiver unchanged.
// Decoding borrows b; keep b unchanged while the fields are in use.
func (f *{{.Name}}Fields) UnmarshalBinary(b []byte) error {
 if f==nil {return ErrIENotFound}
 children,err:=parseGroupedFields({{.Name}},b);if err!=nil {return withDecodePath({{.Name}},err)}
 *f={{.Name}}Fields{};f.addIEs(children);return nil
}

// MarshalLen returns the size of the concatenated child IEs, or zero for nil.
func (f *{{.Name}}Fields) MarshalLen() int {
 if f==nil {return 0};return f.IE().MarshalLen()-4
}

// MarshalTo encodes the value in table order without modifying the receiver.
func (f *{{.Name}}Fields) MarshalTo(b []byte) error {
 if f==nil {return ErrIENotFound}
 return marshalGroupedFields(f.IE(),b)
}

func (f *{{.Name}}Fields) children() []*IE {
 var children []*IE
 {{range .Fields}}children=append(children,f.{{.Name}}{{if .Slice}}...{{end}})
 {{end}}return append(children,f.OtherIEs...)
}

func (f *{{.Name}}Fields) addIEs(children []*IE) {
 for _,child:=range children {if child==nil {continue};switch child.Type {
 {{range .Fields}}case {{.Constant}}:
 {{if .Slice}}f.{{.Name}}=append(f.{{.Name}},child){{else}}if f.{{.Name}}==nil {f.{{.Name}}=child} else {f.OtherIEs=append(f.OtherIEs,child)}{{end}}
 {{end}}default:f.OtherIEs=append(f.OtherIEs,child)
 }}
}
`
