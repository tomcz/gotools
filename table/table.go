package table

import (
	"errors"
	"fmt"
	"math"
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

type columnAlignment struct {
	Header Alignment
	Values Alignment
}

type columnAlignments map[int]columnAlignment

type runeWidths map[string]int

func (r runeWidths) StringWidth(txt string) int {
	if width, ok := r[txt]; ok {
		return width
	}
	width := runewidth.StringWidth(txt)
	r[txt] = width
	return width
}

type tblCell struct {
	Value string
	Width int
}

func tblCells(values []string, runes runeWidths) []tblCell {
	return slices.Map(values, func(s string) tblCell {
		return tblCell{Value: s, Width: runes.StringWidth(s)}
	})
}

type Table struct {
	headers     []tblCell
	rows        [][]tblCell
	cols        columnAlignments
	runes       runeWidths
	boldHeaders bool
	maxWidth    int
}

// NewTable constructor. Invalid maxWidth and empty headers result in errors.
func NewTable(maxWidth int, headers ...string) (*Table, error) {
	if maxWidth < 1 {
		return nil, fmt.Errorf("invalid maxWidth: %d", maxWidth)
	}
	if len(headers) == 0 {
		return nil, errors.New("headers cannot be empty")
	}
	runes := make(runeWidths)
	return &Table{
		boldHeaders: true,
		maxWidth:    maxWidth,
		headers:     tblCells(headers, runes),
		runes:       runes,
	}, nil
}

// BoldHeaders overrides default display style for every header.
func (t *Table) BoldHeaders(bold bool) {
	t.boldHeaders = bold
}

// Alignment sets the alignment for an individual table column.
func (t *Table) Alignment(column int, header Alignment, values Alignment) {
	if t.cols == nil {
		t.cols = make(columnAlignments)
	}
	t.cols[column] = columnAlignment{
		Header: header,
		Values: values,
	}
}

// AddRow to table. The number of values must match the previously set number of headers.
func (t *Table) AddRow(values ...string) error {
	if len(values) != len(t.headers) {
		return fmt.Errorf("len(values) %d != len(headers) %d", len(values), len(t.headers))
	}
	t.rows = append(t.rows, tblCells(values, t.runes))
	return nil
}

// String renders a table that will not be longer than the previously set maximum width.
func (t *Table) String() string {
	colWidths := t.columnWidths()
	sb := new(strings.Builder)
	writeLine(sb, colWidths, "┌", "┬", "┐")
	t.writeRow(sb, t.headers, colWidths, true, t.boldHeaders)
	for _, row := range t.rows {
		sb.WriteString("\n")
		writeLine(sb, colWidths, "├", "┼", "┤")
		t.writeRow(sb, row, colWidths, false, false)
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

func (t *Table) writeRow(sb *strings.Builder, row []tblCell, colWidths []int, isHeader bool, isBold bool) {
	cells := make([][]string, len(row))
	for i, cell := range row {
		cells[i] = t.wrapText(cell, colWidths[i])
	}
	maxLines := slices.Reduce(cells, 0, func(acc int, cell []string) int {
		return maths.Max(acc, len(cell))
	})
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
				alignment := maths.Choose(isHeader, t.cols[i].Header, t.cols[i].Values)
				text := t.alignText(cell[line], colWidths[i], alignment)
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

func (t *Table) alignText(text string, width int, alignment Alignment) string {
	if text == "" {
		return strings.Repeat(" ", width)
	}
	textWidth := t.runes.StringWidth(text)
	if textWidth == width {
		return text
	}
	diff := width - textWidth
	if alignment == Left {
		return text + strings.Repeat(" ", diff)
	}
	if alignment == Right {
		return strings.Repeat(" ", diff) + text
	}
	padLeft := diff / 2
	remainder := diff % 2
	return strings.Repeat(" ", padLeft) + text + strings.Repeat(" ", padLeft+remainder)
}

func (t *Table) wrapText(cell tblCell, width int) []string {
	if cell.Width <= width {
		return []string{cell.Value}
	}
	var line string
	var lines []string
	for word := range strings.SplitSeq(cell.Value, " ") {
		if t.runes.StringWidth(word) > width {
			wrapped := runewidth.Wrap(cell.Value, width)
			return strings.Split(wrapped, "\n")
		}
		if line == "" {
			line = word
			continue
		}
		next := line + " " + word
		if t.runes.StringWidth(next) < width {
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
	paddingCount := 2 * len(t.headers)
	dividerCount := 1 + len(t.headers)
	overhead := paddingCount + dividerCount

	columns := make([]*columnInfo, len(t.headers))
	for i, col := range t.headers {
		widths := slices.Map(t.rows, func(row []tblCell) int { return row[i].Width })
		columns[i] = &columnInfo{index: i, width: maths.MaxOf(col.Width, widths...)}
	}

	tableWidth := slices.Reduce(columns, overhead, func(acc int, c *columnInfo) int { return acc + c.width })
	if tableWidth <= t.maxWidth {
		return columnWidths(columns)
	}

	fairWidth := int(math.Floor(float64(t.maxWidth-overhead) / float64(len(columns))))
	longCols := slices.Select(columns, func(c *columnInfo) bool { return c.width > fairWidth })
	shortCols := slices.Select(columns, func(c *columnInfo) bool { return c.width < fairWidth })
	diff := slices.Reduce(shortCols, 0, func(acc int, c *columnInfo) int { return acc + (fairWidth - c.width) })
	longWidth := fairWidth + int(math.Floor(float64(diff)/float64(len(longCols))))

	for _, col := range longCols {
		col.width = longWidth
	}
	return columnWidths(columns)
}
