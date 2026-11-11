// Command pfcpspec extracts a version-pinned TS 29.244 DOCX into a catalogue.
package main

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gomaja-labs/go-pfcp/internal/spec"
)

const docxSHA256 = "2f0fe3a9cf744ffb5a7f3e170f29238ddbca9a2ea5419da815def531db163ed1"
const archiveSHA256 = "7ef553959af29150f42b0f8cec025aa39addb52b481a12f6d60dee63f6b7a48b"

// node retains logical cells; expanding the grid alone loses content in tables
// that change cell spans halfway through. Only WordprocessingML is traversed.
type node struct {
	XMLName  xml.Name
	Attr     []xml.Attr `xml:",any,attr"`
	Text     string     `xml:",chardata"`
	Children []node     `xml:",any"`
}

func parseXML(r io.Reader) (*node, error) {
	var n node
	err := xml.NewDecoder(r).Decode(&n)
	return &n, err
}
func (n *node) child(name string) *node {
	if n != nil {
		for i := range n.Children {
			if n.Children[i].XMLName.Local == name {
				return &n.Children[i]
			}
		}
	}
	return nil
}
func (n *node) attr(name string) string {
	if n != nil {
		for _, a := range n.Attr {
			if a.Name.Local == name {
				return a.Value
			}
		}
	}
	return ""
}
func (n *node) walk(name string, fn func(*node)) {
	if n == nil {
		return
	}
	if n.XMLName.Local == name {
		fn(n)
	}
	for i := range n.Children {
		n.Children[i].walk(name, fn)
	}
}
func clean(s string) string { return strings.Join(strings.Fields(s), " ") }
func text(n *node) string {
	var paragraphs []string
	n.walk("p", func(p *node) {
		var b strings.Builder
		var visit func(*node)
		visit = func(v *node) {
			switch v.XMLName.Local {
			case "t":
				b.WriteString(v.Text)
			case "tab", "br", "cr":
				b.WriteByte(' ')
			}
			for i := range v.Children {
				visit(&v.Children[i])
			}
		}
		visit(p)
		paragraphs = append(paragraphs, b.String())
	})
	return clean(strings.Join(paragraphs, " "))
}
func expandedRows(tbl *node) ([][]string, [][]string, error) {
	previous := map[int]string{}
	rows := [][]string{}
	logical := [][]string{}
	for i := range tbl.Children {
		tr := &tbl.Children[i]
		if tr.XMLName.Local != "tr" {
			continue
		}
		col := 0
		if before := tr.child("trPr").child("gridBefore"); before != nil {
			var err error
			col, err = strconv.Atoi(before.attr("val"))
			if err != nil || col < 0 || col > 1000 {
				return nil, nil, fmt.Errorf("invalid leading grid")
			}
		}
		row := make([]string, col)
		cells := []string{}
		current := map[int]string{}
		for j := range tr.Children {
			tc := &tr.Children[j]
			if tc.XMLName.Local != "tc" {
				continue
			}
			v := text(tc)
			pr := tc.child("tcPr")
			width := 1
			if span := pr.child("gridSpan"); span != nil {
				var err error
				width, err = strconv.Atoi(span.attr("val"))
				if err != nil || width < 1 || width > 1000 {
					return nil, nil, fmt.Errorf("invalid grid span")
				}
			}
			if merge := pr.child("vMerge"); merge != nil && merge.attr("val") != "restart" {
				v = previous[col]
			}
			cells = append(cells, v)
			for range width {
				current[col] = v
				row = append(row, v)
				col++
			}
		}
		previous = current
		rows = append(rows, row)
		logical = append(logical, cells)
	}
	return rows, logical, nil
}

type table struct {
	ref, title, clause string
	rows, cells        [][]string
}

func loadTables(path string) (result []table, err error) {
	z, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer func() { err = errors.Join(err, z.Close()) }()
	var root *node
	for _, f := range z.File {
		if f.Name == "word/document.xml" {
			r, e := f.Open()
			if e != nil {
				return nil, e
			}
			root, e = parseXML(r)
			closeErr := r.Close()
			if e != nil {
				return nil, e
			}
			if closeErr != nil {
				return nil, closeErr
			}
			break
		}
	}
	if root == nil || root.child("body") == nil {
		return nil, fmt.Errorf("missing document body")
	}
	heading := regexp.MustCompile(`^[78](?:\.[0-9A-Za-z]+)+\s`)
	captionRE := regexp.MustCompile(`Table\s+([78]\.[\d.]+-\d+)\.?\s*:\s*(.*)`)
	clause, ref, title := "", "", ""
	tables := []table{}
	for i := range root.child("body").Children {
		el := &root.child("body").Children[i]
		switch el.XMLName.Local {
		case "p":
			s := text(el)
			style := el.child("pPr").child("pStyle").attr("val")
			if heading.MatchString(s) && (strings.HasPrefix(style, "Heading") || style == "2" || style == "3" || style == "4" || style == "5" || style == "6") {
				clause = strings.Fields(s)[0]
			}
			if m := captionRE.FindStringSubmatch(s); m != nil {
				ref, title = "Table "+m[1], m[2]
			}
		case "tbl":
			if ref != "" {
				rows, cells, e := expandedRows(el)
				if e != nil {
					return nil, e
				}
				tables = append(tables, table{ref, title, clause, rows, cells})
				ref = ""
			}
		}
	}
	return tables, nil
}

func key(s string) string {
	return regexp.MustCompile(`[^a-z0-9]`).ReplaceAllString(strings.ToLower(s), "")
}
func goName(s string) string {
	for _, r := range [][2]string{{"Non-3GPP", "NonTGPP"}, {"Non 3GPP", "NonTGPP"}, {"3GPP", "TGPP"}, {"ID's", "IDs"}, {"5GS", "FiveGS"}} {
		s = strings.ReplaceAll(s, r[0], r[1])
	}
	parts := regexp.MustCompile(`[A-Za-z0-9]+`).FindAllString(s, -1)
	for i, p := range parts {
		parts[i] = strings.ToUpper(p[:1]) + p[1:]
	}
	s = strings.Join(parts, "")
	if s != "" && s[0] >= '0' && s[0] <= '9' {
		s = "N" + s
	}
	return s
}
func ptr[T any](v T) *T   { return &v }
func number(s string) int { v, _ := strconv.Atoi(s); return v }
func references(s string) []string {
	refs := []string{}
	for _, m := range regexp.MustCompile(`(Table|Clause)\s*([\d.]+(?:-\d+)?)`).FindAllStringSubmatch(s, -1) {
		refs = append(refs, clean(m[1]+" "+m[2]))
	}
	return refs
}
func registry(tables []table, surface *spec.Surface) ([]spec.IEType, error) {
	var source *table
	for i := range tables {
		if tables[i].ref == "Table 8.1.2-1" {
			source = &tables[i]
			break
		}
	}
	if source == nil {
		return nil, fmt.Errorf("missing registry")
	}
	out := []spec.IEType{}
	for _, r := range source.rows {
		if len(r) == 0 || !regexp.MustCompile(`^\d+(?: to \d+)?$`).MatchString(r[0]) {
			continue
		}
		if len(r) < 4 {
			return nil, fmt.Errorf("short registry row")
		}
		nums := regexp.MustCompile(`\d+`).FindAllString(r[0], -1)
		ty := number(nums[0])
		name, ref, fixed := r[1], r[2], r[3]
		kind := ""
		switch {
		case strings.Contains(name, "Reserved"):
			kind = "reserved"
		case strings.Contains(name, "Spare"):
			kind = "spare"
		case name == "Can be used":
			kind = "can-be-used"
		case strings.HasPrefix(ref, "Fixed"):
			kind = "fixed"
		case strings.HasPrefix(ref, "Variable"):
			kind = "variable"
		case strings.HasPrefix(ref, "Extendable"):
			kind = "extendable"
		default:
			return nil, fmt.Errorf("unknown IE kind %d", ty)
		}
		gn := surface.IENames[strconv.Itoa(ty)]
		var existing *string
		if gn != "" {
			existing = ptr(gn)
		} else {
			gn = goName(name)
		}
		if prefix := map[string]string{"reserved": "Reserved", "spare": "Spare", "can-be-used": "CanBeUsed"}[kind]; prefix != "" {
			gn = prefix + strconv.Itoa(ty)
		}
		gn = currentName(ty, gn)
		e := spec.IEType{Type: ty, Name: name, GoName: gn, Kind: kind, DefinedIn: references(ref), ExistingGoName: existing, Note: ref}
		if regexp.MustCompile(`^\d+$`).MatchString(fixed) {
			e.FixedOctets = ptr(number(fixed))
		}
		// TS 29.244 V20.1.0 actual clause headings resolve registry punctuation errors.
		if correction := map[int]string{55: "Clause 8.2.35", 56: "Clause 8.2.36", 57: "Clause 8.2.37", 275: "Clause 8.2.186", 299: "Clause 8.2.205"}[ty]; correction != "" {
			e.RegistryDefinedIn = e.DefinedIn
			e.DefinedIn = []string{correction}
			e.Note += "; definition appears under " + correction
		}
		if len(nums) > 1 {
			e.TypeEnd = ptr(number(nums[1]))
		}
		out = append(out, e)
	}
	return out, nil
}

func repetition(comment string) (bool, string) {
	patterns := []string{`Several IEs?\b[^.!?]*(?:[.!?]|$)`, `More than one IE\b[^.!?]*(?:[.!?]|$)`, `This IE may be present up to \d+ times\.`, `There may be several instances of this IE\.`, `This IE shall contain the list of all the (?:URR|QER) IDs[^.!?]*[.!?]`, `The CP function shall provide the full set of Additional Monitoring Times IE\(s\)\.`, `One IE with this type shall be included per PDR ID reported in the message\.`, `[^.!?]*\b(?:two|2) (?:IEs|instances)[^.!?]*(?:[.!?]|$)`, `[^.!?]*\bmay be (?:present|included) (?:several|multiple) times[^.!?]*(?:[.!?]|$)`}
	for _, p := range patterns {
		if m := regexp.MustCompile("(?i)" + p).FindString(comment); m != "" {
			words := strings.Fields(m)
			if len(words) > 25 {
				words = words[:25]
			}
			return true, strings.Join(words, " ")
		}
	}
	return false, "none"
}

type rawRow struct {
	name, presence, comment, ieName string
	blank                           bool
	applicability                   spec.Applicability
}

func contentRows(t table) ([]rawRow, error) {
	rows := []rawRow{}
	for pos, r := range t.cells {
		if len(r) < 2 || (r[1] != "M" && r[1] != "C" && r[1] != "CO" && r[1] != "O") {
			continue
		}
		if len(r) != 9 {
			return nil, fmt.Errorf("%s row %d: unexpected logical cell count %d", t.ref, pos+1, len(r))
		}
		name := r[0]
		if name == "" || (r[8] == "" && (t.ref != "Table 7.5.2.4-1" || name != "Ethernet Inactivity Timer")) {
			return nil, fmt.Errorf("%s row %d: missing required cell", t.ref, pos+1)
		}
		applies := func(i int) bool {
			for _, w := range strings.Fields(r[i]) {
				if w == "X" {
					return true
				}
			}
			return false
		}
		ieName := r[8]
		if ieName == "" {
			ieName = name
		}
		rows = append(rows, rawRow{name, r[1], r[2], ieName, r[8] == "", spec.Applicability{Sxa: applies(3), Sxb: applies(4), Sxc: applies(5), N4: applies(6), N4mb: applies(7)}})
	}
	return rows, nil
}

// Aliases resolve spelling from defining tables, never by fuzzy matching.
// TS 29.244 V20.1.0 Tables 7.4.4.1-1/-4, 7.4.7.1-2, 7.5.2.2-5,
// 7.5.3.1-6, 7.4.5.1.6-1; clauses 8.2.103, 8.2.117 and 8.2.285.
func aliases() map[string]int {
	return map[string]int{"Node level Mitigation Actions Information": 404, "Abnormal Traffic Pattern ID": 407, "Group-Id": 291, "Redundant Transmission Detection Parameters": 255, "MoQ Information": 386, "QoS Information": 240, "DNN": 159, "MAC Addressed Detected": 144}
}

type resolver struct {
	surface *spec.Surface
	byType  map[int]*spec.IEType
	lookup  map[string]int
	defs    map[int][]string
	tables  map[string]table
}

func (v *resolver) rows(t table, msg string, visiting map[string]bool) ([]spec.Row, error) {
	if visiting[t.ref] {
		return nil, fmt.Errorf("cyclic table inheritance at %s", t.ref)
	}
	visiting[t.ref] = true
	defer delete(visiting, t.ref)
	raw, err := contentRows(t)
	if err != nil {
		return nil, err
	}
	if len(raw) == 0 {
		var lines []string
		for _, r := range t.rows {
			lines = append(lines, strings.Join(r, " "))
		}
		if m := regexp.MustCompile(`Same IEs and requirements as defined in (Table\s+[\d.]+-\d+)`).FindStringSubmatch(strings.Join(lines, " ")); m != nil {
			parent, ok := v.tables[m[1]]
			if !ok {
				return nil, fmt.Errorf("missing inherited table %s", m[1])
			}
			return v.rows(parent, msg, visiting)
		}
	}
	out := []spec.Row{}
	for _, r := range raw {
		ty, ok := v.lookup[key(r.ieName)]
		// Context-specific numeric identities: TS 29.244 V20.1.0 table headers.
		contexts := map[[2]string]int{{"Table 7.5.4.1-1", "Update BAR"}: 86, {"Table 7.5.9.1-1", "Update BAR"}: 12, {"Table 7.5.5.1-1", "Usage Report"}: 78, {"Table 7.5.7.1-1", "Usage Report"}: 79, {"Table 7.5.8.1-1", "Usage Report"}: 80, {"Table 7.5.4.1-1", "TSC Management Information"}: 199, {"Table 7.5.5.1-1", "TSC Management Information"}: 200, {"Table 7.5.8.1-1", "TSC Management Information"}: 201, {"Table 7.5.5.1-1", "Packet Rate Status Report"}: 264}
		if contextual, exists := contexts[[2]string{t.ref, r.ieName}]; exists {
			ty, ok = contextual, true
		}
		if !ok {
			return nil, fmt.Errorf("%s: unresolved IE %q", t.ref, r.ieName)
		}
		multi, evidence := repetition(r.comment)
		children := v.defs[ty]
		explicit := map[string]bool{}
		for _, m := range regexp.MustCompile(`Table\s+([\d.]+-\d+)`).FindAllStringSubmatch(r.comment, -1) {
			for _, child := range children {
				if child == "Table "+m[1] {
					explicit[child] = true
				}
			}
		}
		var child *string
		if len(explicit) == 1 {
			for c := range explicit {
				child = ptr(c)
			}
		} else if len(children) == 1 {
			child = ptr(children[0])
		}
		// Tables 7.5.2.9-4 and 7.5.4.18-2 differ by Reporting Flags.
		if ty == 295 {
			if c := map[string]string{"Table 7.5.2.9-1": "Table 7.5.2.9-4", "Table 7.5.4.18-1": "Table 7.5.4.18-2"}[t.ref]; c != "" {
				child = ptr(c)
			}
		}
		if len(children) > 1 && child == nil {
			context := strings.ToLower(t.title + " " + r.name)
			candidates := []string{}
			for _, c := range children {
				title := strings.ToLower(v.tables[c].title)
				for _, word := range []string{"modification", "deletion", "report request"} {
					if strings.Contains(context, word) && strings.Contains(title, word) {
						candidates = append(candidates, c)
						break
					}
				}
			}
			if len(candidates) == 1 {
				child = ptr(candidates[0])
			}
		}
		field := goName(r.name)
		if old := v.surface.Messages[msg].Bindings[strconv.Itoa(ty)]; old != "" {
			field = currentField(ty, old)
		}
		row := spec.Row{Name: r.name, Presence: r.presence, IEType: ty, IETypeName: r.ieName, Multiple: multi, MultipleEvidence: evidence, Applicability: r.applicability, ChildTable: child, GoField: field}
		if r.blank {
			row.SourceIETypeName = ptr("")
			row.ResolutionNote = "Blank IE Type cell; row name matches type 146, clause 8.2.105; confirm before generation."
		}
		out = append(out, row)
	}
	counts := map[int]int{}
	for _, r := range out {
		counts[r.IEType]++
	}
	for i := range out {
		r := &out[i]
		if counts[r.IEType] < 2 {
			continue
		}
		field := v.surface.Messages[msg].Bindings[strconv.Itoa(r.IEType)]
		if field == "" {
			field = v.byType[r.IEType].GoName
		}
		field = currentField(r.IEType, field)
		r.GoField = field
		r.SharedTypeGroup = ptr(field)
	}
	return out, nil
}

func extract(path string) (*spec.Catalogue, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if fmt.Sprintf("%x", sha256.Sum256(b)) != docxSHA256 {
		return nil, fmt.Errorf("unexpected DOCX SHA-256")
	}
	tables, err := loadTables(path)
	if err != nil {
		return nil, err
	}
	surface, err := spec.LoadSurface()
	if err != nil {
		return nil, err
	}
	ies, err := registry(tables, surface)
	if err != nil {
		return nil, err
	}
	v := resolver{surface: surface, byType: map[int]*spec.IEType{}, lookup: map[string]int{}, defs: map[int][]string{}, tables: map[string]table{}}
	for i := range ies {
		e := &ies[i]
		if e.TypeEnd == nil {
			v.byType[e.Type] = e
		}
		if e.Kind != "reserved" && e.Kind != "spare" && e.Kind != "can-be-used" {
			v.lookup[key(e.Name)] = e.Type
		}
	}
	for _, e := range ies {
		names := []string{regexp.MustCompile(`(?i)\s+(?:IE )?within\s+`).Split(e.Name, 2)[0], regexp.MustCompile(`\s*\([^)]*\)\s*$`).ReplaceAllString(e.Name, "")}
		for _, name := range names {
			k := key(name)
			if _, exists := v.lookup[k]; !exists {
				v.lookup[k] = e.Type
			}
		}
	}
	for name, ty := range aliases() {
		v.lookup[key(name)] = ty
	}
	grouped := []spec.Container{}
	messageTables := []table{}
	rowSources := []table{}
	for _, t := range tables {
		if strings.HasPrefix(t.ref, "Table 7.") {
			v.tables[t.ref] = t
		}
		if !strings.HasPrefix(t.ref, "Table 7.4.") && !strings.HasPrefix(t.ref, "Table 7.5.") {
			continue
		}
		condition := false
		for _, r := range t.rows {
			for _, c := range r {
				condition = condition || strings.Contains(c, "Condition")
			}
		}
		if !condition {
			continue
		}
		if len(t.rows) == 0 {
			continue
		}
		header := strings.Join(t.rows[0], " ")
		m := regexp.MustCompile(`(?:IE(?: Type)?|Type)\s*=\s*(\d+)`).FindStringSubmatch(header)
		if m == nil {
			m = regexp.MustCompile(`=\s*(\d+)\s*\(decimal\)`).FindStringSubmatch(header)
		}
		if m != nil {
			ty := number(m[1])
			e := v.byType[ty]
			if e == nil {
				return nil, fmt.Errorf("unregistered group %d", ty)
			}
			grouped = append(grouped, spec.Container{Type: ty, Name: e.Name, GoName: e.GoName, Context: t.title, Table: ptr(t.ref), IEs: []spec.Row{}, Clause: t.clause})
			rowSources = append(rowSources, t)
			v.defs[ty] = append(v.defs[ty], t.ref)
		} else {
			messageTables = append(messageTables, t)
		}
	}
	for ty, refs := range v.defs {
		e := v.byType[ty]
		e.Kind = "grouped"
		e.RegistryDefinedIn = e.DefinedIn
		e.DefinedIn = refs
		if !reflect.DeepEqual(e.RegistryDefinedIn, refs) {
			e.Note += "; defining table headers identify " + strings.Join(refs, ", ")
		}
	}
	for i := range grouped {
		grouped[i].IEs, err = v.rows(rowSources[i], "", map[string]bool{})
		if err != nil {
			return nil, err
		}
	}
	for ty, refs := range v.defs {
		if len(refs) > 1 {
			var first []spec.Row
			different := false
			for _, g := range grouped {
				if g.Type == ty {
					if first == nil {
						first = g.IEs
					} else {
						different = different || !reflect.DeepEqual(first, g.IEs)
					}
				}
			}
			v.byType[ty].Variants = ptr(different)
		}
	}
	messages := []spec.Container{}
	for _, t := range tables {
		if t.ref != "Table 7.3-1" {
			continue
		}
		for _, r := range t.rows {
			if len(r) < 2 || !regexp.MustCompile(`^\d+$`).MatchString(r[0]) {
				continue
			}
			ty := number(r[0])
			if (ty < 1 || ty > 17) && (ty < 50 || ty > 57) {
				continue
			}
			name := r[1]
			gn := strings.ReplaceAll(goName(strings.TrimPrefix(name, "PFCP ")), "Modifcation", "Modification")
			if ty == 11 {
				gn = "VersionNotSupportedResponse"
			}
			corrected := strings.TrimPrefix(strings.ReplaceAll(name, "Modifcation", "Modification"), "PFCP ")
			candidates := []table{}
			for _, candidate := range messageTables {
				if strings.Contains(candidate.title, corrected) {
					candidates = append(candidates, candidate)
				}
			}
			expected := 1
			if ty == 11 {
				expected = 0
			}
			if len(candidates) != expected {
				return nil, fmt.Errorf("message %d: %d candidate tables", ty, len(candidates))
			}
			category := "node"
			if ty >= 50 {
				category = "session"
			}
			m := spec.Container{Type: ty, Name: name, GoName: gn, Category: category, IEs: []spec.Row{}, Clause: "7.4.4.7"}
			if len(candidates) > 0 {
				source := candidates[0]
				m.Table = ptr(source.ref)
				m.Clause = source.clause
				m.IEs, err = v.rows(source, gn, map[string]bool{})
				if err != nil {
					return nil, err
				}
			}
			messages = append(messages, m)
		}
	}
	c := &spec.Catalogue{Spec: spec.Document{Document: "3GPP TS 29.244", Version: "20.1.0", Date: "2026-09", Archive: "29244-k10.zip", ArchiveSHA256: archiveSHA256, DocxSHA256: docxSHA256}, IETypes: ies, Messages: messages, Grouped: grouped}
	if err := spec.Validate(c); err != nil {
		return nil, err
	}
	return c, nil
}
func run() error {
	docx := flag.String("docx", "", "input V20.1.0 DOCX")
	output := flag.String("output", "", "output catalogue JSON")
	flag.Parse()
	if *docx == "" || *output == "" {
		return fmt.Errorf("both -docx and -output are required")
	}
	c, err := extract(*docx)
	if err != nil {
		return err
	}
	b, err := spec.CanonicalJSON(c)
	if err != nil {
		return err
	}
	return os.WriteFile(*output, b, 0o644)
}
func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// currentName applies the approved TS 29.244 V20.1.0 Table 8.1.2-1 names.
func currentName(typ int, name string) string {
	switch typ {
	case 194:
		return "CreateBridgeRouterInfo"
	case 195:
		return "CreatedBridgeRouterInfo"
	case 196:
		return "PortNumber"
	case 198:
		return "FiveGSUserPlaneNodeID"
	case 206:
		return "TimeDomainNumber"
	case 229:
		return "LinkSpecificMultipathIPAddress"
	case 265:
		return "MultipathApplicableIndication"
	case 266:
		return "UserPlaneNodeManagementInformationContainer"
	default:
		return name
	}
}

func currentField(typ int, name string) string {
	// TS 29.244 V20.1.0 Table 7.5.8.1-1.
	if typ == 201 {
		return "TSCManagementInformation"
	}
	return currentName(typ, name)
}
