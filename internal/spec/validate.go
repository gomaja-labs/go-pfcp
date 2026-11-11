package spec

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func itoa(n int) string { return strconv.Itoa(n) }

// CanonicalJSON encodes the catalogue with the original extractor's member order,
// indentation, Unicode and final newline.
func CanonicalJSON(c *Catalogue) ([]byte, error) {
	var b bytes.Buffer
	e := json.NewEncoder(&b)
	e.SetEscapeHTML(false)
	e.SetIndent("", "  ")
	if err := e.Encode(c); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// Validate checks coverage, source identities, list evidence and Go identifiers.
// TS 29.244 V20.1.0 Tables 8.1.2-1 and 7.3-1; clauses 7.2.3.4 and 7.6.10.
func Validate(c *Catalogue) error {
	if c == nil || c.Spec.Version != "20.1.0" {
		return fmt.Errorf("unexpected catalogue version")
	}
	fail := func(s string) error { return fmt.Errorf("invalid catalogue: %s", s) }
	ident := regexp.MustCompile(`^[A-Z][A-Za-z0-9]*$`)
	coverage := make([]int, 65536)
	byType := map[int]IEType{}
	names := map[string]bool{}
	for _, e := range c.IETypes {
		hi := e.Type
		if e.TypeEnd != nil {
			hi = *e.TypeEnd
		}
		if e.Type < 0 || hi < e.Type || hi > 65535 {
			return fail("invalid IE range")
		}
		for t := e.Type; t <= hi; t++ {
			coverage[t]++
		}
		if names[e.GoName] || !ident.MatchString(e.GoName) {
			return fail("invalid or duplicate IE identifier")
		}
		names[e.GoName] = true
		if e.TypeEnd == nil {
			byType[e.Type] = e
		}
		if e.FixedOctets != nil && *e.FixedOctets < 0 {
			return fail("invalid fixed_octets")
		}
		switch e.Kind {
		case "grouped", "fixed", "variable", "extendable", "spare", "can-be-used", "reserved":
		default:
			return fail("invalid kind")
		}
	}
	for _, n := range coverage {
		if n != 1 {
			return fail("missing or overlapping IE type")
		}
	}
	msgs := map[int]bool{}
	names = map[string]bool{}
	for _, m := range c.Messages {
		if (m.Type < 1 || m.Type > 17) && (m.Type < 50 || m.Type > 57) || msgs[m.Type] || names[m.GoName] || !ident.MatchString(m.GoName) {
			return fail("invalid message coverage or name")
		}
		msgs[m.Type] = true
		names[m.GoName] = true
	}
	if len(msgs) != 25 {
		return fail("incomplete message coverage")
	}
	groups := map[int][]Container{}
	tables := map[string]Container{}
	for _, g := range c.Grouped {
		if g.Table == nil || tables[*g.Table].Table != nil || byType[g.Type].Kind != "grouped" || byType[g.Type].GoName != g.GoName || len(g.IEs) == 0 {
			return fail("invalid grouped definition")
		}
		tables[*g.Table] = g
		groups[g.Type] = append(groups[g.Type], g)
	}
	for _, e := range c.IETypes {
		if e.Kind != "grouped" {
			continue
		}
		gs := groups[e.Type]
		if len(gs) == 0 || len(gs) != len(e.DefinedIn) {
			return fail("incomplete defining tables")
		}
		refs := map[string]bool{}
		for _, g := range gs {
			refs[*g.Table] = true
		}
		for _, ref := range e.DefinedIn {
			if !refs[ref] {
				return fail("wrong defining table")
			}
		}
		if len(gs) > 1 {
			different := false
			for _, g := range gs[1:] {
				different = different || !reflect.DeepEqual(gs[0].IEs, g.IEs)
			}
			if e.Variants == nil || *e.Variants != different {
				return fail("wrong variants flag")
			}
		}
	}
	all := append(append([]Container{}, c.Messages...), c.Grouped...)
	for _, g := range all {
		fields := map[string][]Row{}
		types := map[int][]Row{}
		for _, r := range g.IEs {
			if r.Presence != "M" && r.Presence != "C" && r.Presence != "CO" && r.Presence != "O" {
				return fail("invalid presence")
			}
			e, ok := byType[r.IEType]
			if !ok || e.Kind == "reserved" || e.Kind == "spare" || e.Kind == "can-be-used" {
				return fail("unresolved or unusable row IE")
			}
			if r.MultipleEvidence == "" || len(strings.Fields(r.MultipleEvidence)) > 25 || (r.MultipleEvidence != "none") != r.Multiple {
				return fail("invalid repetition evidence")
			}
			if !ident.MatchString(r.GoField) {
				return fail("invalid field identifier")
			}
			if e.Kind == "grouped" {
				if r.ChildTable == nil || tables[*r.ChildTable].Type != r.IEType {
					return fail("missing or wrong grouped child")
				}
			} else if r.ChildTable != nil {
				return fail("scalar has child table")
			}
			fields[r.GoField] = append(fields[r.GoField], r)
			types[r.IEType] = append(types[r.IEType], r)
		}
		for field, rs := range fields {
			if len(rs) > 1 {
				for _, r := range rs {
					if r.IEType != rs[0].IEType || r.SharedTypeGroup == nil || *r.SharedTypeGroup != field {
						return fail("unintended field collision")
					}
				}
			}
		}
		for _, rs := range types {
			for _, r := range rs {
				if len(rs) > 1 {
					if r.GoField != rs[0].GoField || r.SharedTypeGroup == nil || *r.SharedTypeGroup != r.GoField {
						return fail("same type needs one shared field")
					}
				} else if r.SharedTypeGroup != nil {
					return fail("spurious shared group")
				}
			}
		}
	}
	return nil
}
