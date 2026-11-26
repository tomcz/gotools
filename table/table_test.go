package table

import (
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
	assert.DeepEqual(t, []string{"Virginia"}, rw.WrapCell(rw.cell("Virginia"), 10))
	assert.DeepEqual(t, []string{"Queensland"}, rw.WrapCell(rw.cell("Queensland"), 10))

	actual := rw.WrapCell(rw.cell("CLOUDFLARENET - Cloudflare, Inc."), 15)
	assert.DeepEqual(t, []string{"CLOUDFLARENET", "- Cloudflare,", "Inc."}, actual)

	actual = rw.WrapCell(rw.cell("2001:0db8:85a3:0000:0000:8a2e:0370:7334"), 5)
	assert.DeepEqual(t, []string{"2001:", "0db8:", "85a3:", "0000:", "0000:", "8a2e:", "0370:", "7334"}, actual)
}
