package honeylogger

import (
	"log/slog"
	"testing"

	"github.com/honeycombio/libhoney-go"
	"gotest.tools/v3/assert"
)

func TestHandler_Enabled(t *testing.T) {
	h := &Handler{Level: slog.LevelInfo}
	assert.Assert(t, h.Enabled(t.Context(), slog.LevelInfo))
	assert.Assert(t, h.Enabled(t.Context(), slog.LevelWarn))
	assert.Assert(t, !h.Enabled(t.Context(), slog.LevelDebug))
}

func TestHandler_Copy(t *testing.T) {
	h1 := &Handler{
		Level:   slog.LevelWarn,
		Builder: new(libhoney.Builder),
		Reject:  func(record slog.Record) bool { return true },
	}
	attrs := map[string]slog.Value{"foor": slog.StringValue("bar")}
	groups := []string{"wibble"}
	h2 := h1.copy(attrs, groups)
	assert.Equal(t, h1.Level, h2.Level)
	assert.Equal(t, h1.Builder, h2.Builder)
	assert.Equal(t, h1.Reject(slog.Record{}), h2.Reject(slog.Record{}))
	assert.DeepEqual(t, attrs, h2.attrs)
	assert.DeepEqual(t, groups, h2.groups)
}
