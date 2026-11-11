package main

import (
	"strings"
	"testing"
)

func TestRepetition(t *testing.T) {
	for _, s := range []string{"Several IEs within the same IE type may be present to represent a list of Usage Reports.", "This IE may be present up to 16 times.", "More than one IE with this type may be included to represent multiple Time Domain Numbers.", "There may be several instances of this IE.", "This IE shall contain the list of all the URR IDs to be associated to the PDR.", "One IE with this type shall be included per PDR ID reported in the message."} {
		multi, evidence := repetition(s)
		if !multi || !strings.Contains(s, evidence) || len(strings.Fields(evidence)) > 25 {
			t.Fatal(s, multi, evidence)
		}
	}
	for _, s := range []string{"This IE shall contain multiple IP addresses.", "QERs of several PFCP sessions", "This IE may be present."} {
		multi, quote := repetition(s)
		if multi || quote != "none" {
			t.Fatal(s)
		}
	}
}

func TestLogicalRows(t *testing.T) {
	const source = `<tbl><tr><tc><tcPr><gridSpan val="2"/><vMerge val="restart"/></tcPr><p><r><t>Source IP</t></r></p></tc><tc><p><r><t>O</t></r></p></tc></tr><tr><tc><tcPr><gridSpan val="2"/><vMerge/></tcPr><p/></tc><tc><p><r><t>M</t></r></p></tc></tr><tr><trPr><gridBefore val="1"/></trPr><tc><p><r><t>tail</t></r></p></tc></tr></tbl>`
	n, err := parseXML(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}
	rows, cells, err := expandedRows(n)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows[0]) != 3 || len(cells[0]) != 2 || cells[1][0] != "Source IP" || rows[2][0] != "" || rows[2][1] != "tail" {
		t.Fatalf("rows=%v logical=%v", rows, cells)
	}
}
