package table

import (
	"cmp"
	"errors"
	"fmt"
	sdk "slices"
	"strings"

	"github.com/mattn/go-runewidth"

	"github.com/tomcz/gotools/maths"
	"github.com/tomcz/gotools/slices"
)

type Alignment int

const (
	Left Alignment = iota
	Right
	Middle
)

type tblCell struct {
	Value string
	Width int
}

type runeWidths map[string]int

func (r runeWidths) StringWidth(txt string) int {
	if width, ok := r[txt]; ok {
		return width
	}
	width := runewidth.StringWidth(txt)
	r[txt] = width
	return width
}

func (r runeWidths) cell(value string) tblCell {
	return tblCell{Value: value, Width: r.StringWidth(value)}
}

func (r runeWidths) Cells(values []string) []tblCell {
	return slices.Map(values, r.cell)
}

func (r runeWidths) Format(value string, width int, align Alignment) string {
	if value == "" {
		return strings.Repeat(" ", width)
	}
	valueWidth := r.StringWidth(value)
	if valueWidth == width {
		return value
	}
	diff := width - valueWidth
	if align == Left {
		return value + strings.Repeat(" ", diff)
	}
	if align == Right {
		return strings.Repeat(" ", diff) + value
	}
	padLeft := diff / 2
	remainder := diff % 2
	return strings.Repeat(" ", padLeft) + value + strings.Repeat(" ", padLeft+remainder)
}

func (r runeWidths) wrap(value string, width int) []string {
	wrapped := runewidth.Wrap(value, width)
	return strings.Split(wrapped, "\n")
}

func (r runeWidths) WrapCell(cell tblCell, width int) []string {
	if cell.Width <= width {
		return []string{cell.Value}
	}
	var line string
	var lines []string
	for word := range strings.SplitSeq(cell.Value, " ") {
		if r.StringWidth(word) > width {
			return r.wrap(cell.Value, width)
		}
		if line == "" {
			line = word
			continue
		}
		next := line + " " + word
		if r.StringWidth(next) < width {
			line = next
			continue
		}
		lines = append(lines, line)
		line = word
	}
	if line != "" {
		lines = append(lines, line)
	}
	return lines
}

type Table struct {
	headers     []tblCell
	rows        [][]tblCell
	headerAlign map[int]Alignment
	valuesAlign map[int]Alignment
	runes       runeWidths
	boldHeaders bool
	maxWidth    int
}

// NewTable constructor.
func NewTable(maxWidth int, headers ...string) (*Table, error) {
	if len(headers) == 0 {
		return nil, errors.New("headers cannot be empty")
	}
	// minimum width is overhead + one character per header
	minWidth := formatOverhead(len(headers)) + len(headers)
	if maxWidth < minWidth {
		return nil, fmt.Errorf("invalid maxWidth: %d, need at least %d", maxWidth, minWidth)
	}
	runes := make(runeWidths)
	return &Table{
		maxWidth:    maxWidth,
		headers:     runes.Cells(headers),
		headerAlign: make(map[int]Alignment),
		valuesAlign: make(map[int]Alignment),
		runes:       runes,
		boldHeaders: true,
	}, nil
}

// BoldHeaders overrides default display style for every header.
func (t *Table) BoldHeaders(bold bool) {
	t.boldHeaders = bold
}

// Alignment sets the alignment for an individual table column identified by the given header.
func (t *Table) Alignment(header string, headerAlign Alignment, valuesAlign Alignment) {
	column := slices.IndexOf(t.headers, func(cell tblCell) bool {
		return cell.Value == header
	})
	if column < 0 {
		panic(fmt.Sprintf("unknown header %q", header))
	}
	t.headerAlign[column] = headerAlign
	t.valuesAlign[column] = valuesAlign
}

// AddRow to table. The count of values must match the count of headers.
func (t *Table) AddRow(values ...string) {
	if len(values) != len(t.headers) {
		panic(fmt.Sprintf("len(values) %d != len(headers) %d", len(values), len(t.headers)))
	}
	t.rows = append(t.rows, t.runes.Cells(values))
}

// String renders a table that will not be longer than the previously set maximum width.
func (t *Table) String() string {
	colWidths := t.columnWidths()
	sb := new(strings.Builder)
	writeLine(sb, colWidths, "┌", "┬", "┐")
	t.writeRow(sb, t.headers, colWidths, t.headerAlign, t.boldHeaders)
	for _, row := range t.rows {
		sb.WriteString("\n")
		writeLine(sb, colWidths, "├", "┼", "┤")
		t.writeRow(sb, row, colWidths, t.valuesAlign, false)
	}
	sb.WriteString("\n")
	writeLine(sb, colWidths, "└", "┴", "┘")
	return sb.String()
}

func writeLine(sb *strings.Builder, colWidths []int, left, middle, right string) {
	sb.WriteString(left)
	for i, width := range colWidths {
		if i > 0 {
			sb.WriteString(middle)
		}
		sb.WriteString(strings.Repeat("─", width+2))
	}
	sb.WriteString(right)
}

var (
	boldStart = fmt.Sprintf("\u001b[%dm", 1)
	boldEnd   = fmt.Sprintf("\u001b[%dm", 22)
)

func (t *Table) writeRow(sb *strings.Builder, row []tblCell, colWidths []int, colAligns map[int]Alignment, isBold bool) {
	cells := make([][]string, len(row))
	for i, cell := range row {
		cells[i] = t.runes.WrapCell(cell, colWidths[i])
	}
	maxLines := slices.Reduce(cells, 1, func(acc int, cell []string) int { return max(acc, len(cell)) })
	for line := range maxLines {
		for i, cell := range cells {
			if i == 0 {
				sb.WriteString("\n")
				sb.WriteString("│")
				sb.WriteString(" ")
			} else {
				sb.WriteString(" ")
			}
			if line < len(cell) {
				text := t.runes.Format(cell[line], colWidths[i], colAligns[i])
				if isBold {
					text = boldStart + text + boldEnd
				}
				sb.WriteString(text)
			} else {
				sb.WriteString(strings.Repeat(" ", colWidths[i]))
			}
			sb.WriteString(" ")
			sb.WriteString("│")
		}
	}
}

type columnInfo struct {
	index int
	width int
}

func columnWidths(columns []*columnInfo) []int {
	result := make([]int, len(columns))
	for _, col := range columns {
		result[col.index] = col.width
	}
	return result
}

func (t *Table) columnWidths() []int {
	overhead := formatOverhead(len(t.headers))

	columns := make([]*columnInfo, len(t.headers))
	for i, h := range t.headers {
		cw := slices.Reduce(t.rows, h.Width, func(acc int, row []tblCell) int { return max(acc, row[i].Width) })
		columns[i] = &columnInfo{index: i, width: cw}
	}

	tableWidth := slices.Reduce(columns, overhead, func(acc int, c *columnInfo) int { return acc + c.width })
	if tableWidth <= t.maxWidth {
		return columnWidths(columns)
	}

	equalWidth := (t.maxWidth - overhead) / len(columns)
	wideCols := slices.Select(columns, func(c *columnInfo) bool { return c.width > equalWidth })
	sdk.SortStableFunc(wideCols, func(a, b *columnInfo) int { return cmp.Compare(b.width, a.width) })
	for _, col := range wideCols {
		col.width = equalWidth
	}

	tableWidth = slices.Reduce(columns, overhead, func(acc int, c *columnInfo) int { return acc + c.width })
	remainder := t.maxWidth - tableWidth

	var col int
	for range remainder {
		wideCols[col].width += 1
		col = maths.Choose(col < len(wideCols)-1, col+1, 0)
	}

	return columnWidths(columns)
}

func formatOverhead(numHeaders int) int {
	paddingCount := 2 * numHeaders
	dividerCount := 1 + numHeaders
	return paddingCount + dividerCount
}
