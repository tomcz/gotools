package table

import (
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestFormat(t *testing.T) {
	rw := make(runeWidths)
	assert.Equal(t, "  ", rw.Format("", 2, Left))
	assert.Equal(t, "hi", rw.Format("hi", 2, Left))
	assert.Equal(t, "hi   ", rw.Format("hi", 5, Left))
	assert.Equal(t, "   hi", rw.Format("hi", 5, Right))
	assert.Equal(t, " hi  ", rw.Format("hi", 5, Middle))
	assert.Equal(t, "  hi  ", rw.Format("hi", 6, Middle))
}

func TestWrapCell(t *testing.T) {
	rw := make(runeWidths)
	assert.DeepEqual(t, []string{"Virginia"}, rw.WrapCell(rw.Cell("Virginia"), 10))
	assert.DeepEqual(t, []string{"Queensland"}, rw.WrapCell(rw.Cell("Queensland"), 10))

	actual := rw.WrapCell(rw.Cell("CLOUDFLARENET - Cloudflare, Inc."), 15)
	assert.DeepEqual(t, []string{"CLOUDFLARENET", "- Cloudflare,", "Inc."}, actual)

	actual = rw.WrapCell(rw.Cell("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), 5)
	assert.DeepEqual(t, []string{"2001:", "0db8:", "85a3:", "0000:", "0000:", "8a2e:", "0370:", "7334"}, actual)
}

const expectedTable = `
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
	assert.Equal(t, strings.TrimSpace(expectedTable), table.String())
}
