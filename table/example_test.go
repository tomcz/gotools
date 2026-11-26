package table_test

import (
	"fmt"

	"github.com/tomcz/gotools/table"
)

func Example() {
	tbl, err := table.NewTable(100, "IP", "ASN", "AS Country", "AS Description", "Geo Country", "Geo Region", "Geo City")
	if err != nil {
		panic(err)
	}

	tbl.AddRow("1.1.1.1", "13335", "US", "CLOUDFLARENET - Cloudflare, Inc.", "Australia", "Queensland", "Brisbane")
	tbl.AddRow("2.2.2.2", "", "", "", "United States of America", "Virginia", "Redwood")
	tbl.AddRow("8.8.8.8", "15169", "US", "GOOGLE - Google LLC", "United States of America", "California", "Mountain View")
	tbl.AddRow("9.9.9.9", "19281", "US", "QUAD9-AS-1 - Quad9", "United States of America", "California", "Berkeley")
	tbl.AddRow("127.0.0.1", "", "", "", "", "", "")

	tbl.BoldHeaders(false)
	tbl.Alignment("IP", table.Middle, table.Left)
	tbl.Alignment("ASN", table.Middle, table.Right)
	tbl.Alignment("AS Country", table.Middle, table.Middle)

	fmt.Println(tbl.String())

	// Output:
	// ┌───────────┬───────┬────────────┬─────────────────┬─────────────────┬────────────┬────────────────┐
	// │    IP     │  ASN  │ AS Country │ AS Description  │ Geo Country     │ Geo Region │ Geo City       │
	// ├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
	// │ 1.1.1.1   │ 13335 │     US     │ CLOUDFLARENET   │ Australia       │ Queensland │ Brisbane       │
	// │           │       │            │ - Cloudflare,   │                 │            │                │
	// │           │       │            │ Inc.            │                 │            │                │
	// ├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
	// │ 2.2.2.2   │       │            │                 │ United States   │ Virginia   │ Redwood        │
	// │           │       │            │                 │ of America      │            │                │
	// ├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
	// │ 8.8.8.8   │ 15169 │     US     │ GOOGLE -        │ United States   │ California │ Mountain View  │
	// │           │       │            │ Google LLC      │ of America      │            │                │
	// ├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
	// │ 9.9.9.9   │ 19281 │     US     │ QUAD9-AS-1 -    │ United States   │ California │ Berkeley       │
	// │           │       │            │ Quad9           │ of America      │            │                │
	// ├───────────┼───────┼────────────┼─────────────────┼─────────────────┼────────────┼────────────────┤
	// │ 127.0.0.1 │       │            │                 │                 │            │                │
	// └───────────┴───────┴────────────┴─────────────────┴─────────────────┴────────────┴────────────────┘
}
