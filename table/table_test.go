package table

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

const expected = `
┌───────────┬───────┬────────────┬─────────────────┬─────────────────┬────────────┬────────────────┐
│    IP     │  ASN  │ AS Country │ AS Description  │ Geo Country     │ Geo Region │ Geo City       │
├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
│ 1.1.1.1   │ 13335 │     US     │ CLOUDFLARENET   │ Australia       │ Queensland │ Brisbane       │
│           │       │            │ - Cloudflare,   │                 │            │                │
│           │       │            │ Inc.            │                 │            │                │
├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
│ 2.2.2.2   │       │            │                 │ United States   │ Virginia   │ Redwood        │
│           │       │            │                 │ of America      │            │                │
├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
│ 8.8.8.8   │ 15169 │     US     │ GOOGLE -        │ United States   │ California │ Mountain View  │
│           │       │            │ Google LLC      │ of America      │            │                │
├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
│ 9.9.9.9   │ 19281 │     US     │ QUAD9-AS-1 -    │ United States   │ California │ Berkeley       │
│           │       │            │ Quad9           │ of America      │            │                │
├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
│ 127.0.0.1 │       │            │                 │                 │            │                │
└───────────┴───────┴────────────┴─────────────────┴─────────────────┴────────────┴────────────────┘
`

func TestTable(t *testing.T) {
	table, err := NewTable(100, "IP", "ASN", "AS Country", "AS Description", "Geo Country", "Geo Region", "Geo City")
	assert.NilError(t, err)

	table.AddRow("1.1.1.1", "13335", "US", "CLOUDFLARENET - Cloudflare, Inc.", "Australia", "Queensland", "Brisbane")
	table.AddRow("2.2.2.2", "", "", "", "United States of America", "Virginia", "Redwood")
	table.AddRow("8.8.8.8", "15169", "US", "GOOGLE - Google LLC", "United States of America", "California", "Mountain View")
	table.AddRow("9.9.9.9", "19281", "US", "QUAD9-AS-1 - Quad9", "United States of America", "California", "Berkeley")
	table.AddRow("127.0.0.1", "", "", "", "", "", "")

	table.BoldHeaders(false)
	table.Alignment("IP", Middle, Left)
	table.Alignment("ASN", Middle, Right)
	table.Alignment("AS Country", Middle, Middle)
	assert.Equal(t, strings.TrimSpace(expected), table.String())
}
