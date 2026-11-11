# TS 29.244 V20.1.0 catalogue

`ts29244-v20.1.0.json` is the validated catalogue extracted from 3GPP TS
29.244 V20.1.0 (2026-09). It includes the approved current IE and message field names; SHA-256:
`82c8de9f31082d39d3d3e6c0d778dabeddf8cd3c577589dff5d1a7cd772830c6`.

The [3GPP archive](https://www.3gpp.org/ftp/Specs/archive/29_series/29.244/)
was checked on 2026-11-11. The newest publication listed was
[29244-k10.zip](https://www.3gpp.org/ftp/Specs/archive/29_series/29.244/29244-k10.zip).

| Input | SHA-256 |
| --- | --- |
| Archive | `7ef553959af29150f42b0f8cec025aa39addb52b481a12f6d60dee63f6b7a48b` |
| `29244-k10.docx` | `2f0fe3a9cf744ffb5a7f3e170f29238ddbca9a2ea5419da815def531db163ed1` |

The DOCX and specification text are not distributed in this repository.
`go-surface.json` pins existing names at revision
`2fbe893b27b0ccc38232a1312300fa7ebdb6411b`, including source hashes. It is a
naming input, not current-checkout validation. Extraction does not inspect its
own generated output to decide names.

From the repository root:

```sh
go run ./internal/cmd/pfcpspec -docx /path/to/29244-k10.docx -output /tmp/catalogue.json
cmp internal/spec/ts29244-v20.1.0.json /tmp/catalogue.json
go test ./internal/spec ./internal/cmd/pfcpspec ./internal/cmd/pfcpgen -count=1
go generate ./...
git diff --exit-code
```

The Go extractor uses only the standard library. It checks the DOCX hash,
retains logical cells independently of expanded grid coordinates, inherits
vertical merges by column, expands incorporated tables and records list evidence.
Committed-catalogue tests require no DOCX, network access or Python.

For a new specification version, first verify the newest publication against
the archive listing and retain download hashes. Update the extractor's pinned
input identities and document metadata, then review every changed table header,
registry allocation, child reference, repetition excerpt and contextual variant.
Capture intentional naming changes separately from numeric identity changes.
Regenerate to a new versioned JSON file and compare the complete diff before
updating the embedded catalogue filename and generator provenance banner.
Update the census assertions and run the extractor, generator and library gates.
A new hash by itself is not acceptance of changed wire semantics.

## Interpretation decisions

These decisions apply to this release; literal conflicting source references
remain in the unchanged JSON for provenance.

- Defining table headers determine numeric identity where the registry conflicts:
  detection parameters are 255 (Table 7.5.2.2-5), forwarding parameters are 270
  (Table 7.5.2.3-4), Update SRR is 213 (Table 7.5.4.20-1), Session Report is 214
  (Tables 7.5.7.3-1 and 7.5.8.6-1), and MoQ Information is 386 (Table 7.5.3.1-6).
- The blank IE Type cell for Ethernet Inactivity Timer in Table 7.5.2.4-1 maps to
  146. The row name and clause 8.2.105 agree; V16.7.0 preserves that cell and the
  current Table 7.5.4.4-1 also identifies it. The JSON's historical confirmation
  note remains unchanged; this paragraph records the approved resolution.
- Session Report (214) and Direct Reporting Information (295) retain separate
  table contexts and child selections. A future grouped-fields API uses one
  union struct per numeric type and documents each row's context.
- Same-type rows share one slice field, including the FQ-CSID roles. This does
  not infer arbitrary multiplicity or remove interface-specific conditions.
  Other lists require the explicit bounded repetition evidence in the JSON.
- Current IE names and message fields are emitted from `go_name` and `go_field`.
  `existing_go_name` and `go-surface.json` retain the historical API record.
- Types 116 (reserved, clause 8.2.82 is Void), 273 (spare), and 381 (can be
  used) have no generated constants. A received IE of these types is retained
  as an unknown IE (clause 7.6.9).

Seven row-comment references are resolved against the actual child definitions:

| Parent | Row | Printed reference | Actual child table |
| --- | --- | --- | --- |
| 7.5.4.1-1 | Remove PDR | 7.5.4-6 | 7.5.4.6-1 |
| 7.5.4.1-1 | Remove FAR | 7.5.4-7 | 7.5.4.7-1 |
| 7.5.4.1-1 | Remove URR | 7.5.4-8 | 7.5.4.8-1 |
| 7.5.4.1-1 | Remove QER | 7.5.4-9 | 7.5.4.9-1 |
| 7.5.4.1-1 | Remove SRR | 7.5.4-19 | 7.5.4.19-1 |
| 7.5.8.1-1 | Packet Rate Status Report | 7.5.7.1-1 | 7.5.7.1-2 |
| 7.5.4.20-1 | Traffic Parameter Measurement Control Information | 7.5.2.9-4 | 7.5.2.9-5 |

The catalogue contains 25 messages, 268 message rows, 121 grouped types across
123 tables, and 659 grouped rows. It records presence and the union of interface
marks, not the full applicability predicates or per-flag conditions. The
registry's `fixed_octets` is the fixed portion of an extendable IE, not its full
length and not a prefix to prepend to grouped children.
