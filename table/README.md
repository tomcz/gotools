# table

Simple table that wraps table cells to respect the set maximum width. Text wrapping will first attempt to break on spaces, but if that is not possible it will break words to ensure a fair cell width.

```go
import (
	"fmt"
	"log"
	"os"
	"github.com/tomcz/gotools/table"
	"golang.org/x/term"
)

func main() {
	var results []*Result // initialisation omitted
	
	boldHeader := true
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// not running in a terminal, use fixed width and plain headers
		boldHeader = false
		width = 100
	}
	tbl, err := table.NewTable(width, "IP", "ASN", "AS Country", "AS Description", "Geo Country", "Geo Region", "Geo City")
	if err != nil {
		log.Fatalln(err)
	}
	tbl.BoldHeader(boldHeader)
	for _, res := range results {
		row := make([]string, 7)
		row[0] = res.IP
		if res.ASN != nil {
			row[1] = strconv.Itoa(res.ASN.ASN)
			row[2] = res.ASN.CountryCode
			row[3] = res.ASN.Description
		}
		if res.Location != nil {
			row[4] = res.Location.Country
			row[5] = res.Location.Region
			row[6] = res.Location.City
		}
		if err = tbl.AddRow(row...); err != nil {
			log.Fatalln(err)
		}
	}
	fmt.Println(tbl.String())
}
```

```
┌───────────┬───────┬────────────┬────────────────┬────────────────┬────────────┬────────────────┐
│ IP        │ ASN   │ AS Country │ AS Description │ Geo Country    │ Geo Region │ Geo City       │
├───────────┼───────┼────────────┼────────────────┼────────────────┼────────────┼────────────────┤
│ 1.1.1.1   │ 13335 │ US         │ CLOUDFLARENET  │ Australia      │ Queensland │ Brisbane       │
│           │       │            │ - Cloudflare,  │                │            │                │
│           │       │            │ Inc.           │                │            │                │
├───────────┼───────┼────────────┼────────────────┼────────────────┼────────────┼────────────────┤
│ 2.2.2.2   │       │            │                │ United States  │ Virginia   │ Redwood        │
│           │       │            │                │ of America     │            │                │
├───────────┼───────┼────────────┼────────────────┼────────────────┼────────────┼────────────────┤
│ 8.8.8.8   │ 15169 │ US         │ GOOGLE -       │ United States  │ California │ Mountain View  │
│           │       │            │ Google LLC     │ of America     │            │                │
├───────────┼───────┼────────────┼────────────────┼────────────────┼────────────┼────────────────┤
│ 9.9.9.9   │ 19281 │ US         │ QUAD9-AS-1 -   │ United States  │ California │ Berkeley       │
│           │       │            │ Quad9          │ of America     │            │                │
├───────────┼───────┼────────────┼────────────────┼────────────────┼────────────┼────────────────┤
│ 127.0.0.1 │       │            │                │                │            │                │
└───────────┴───────┴────────────┴────────────────┴────────────────┴────────────┴────────────────┘
```
