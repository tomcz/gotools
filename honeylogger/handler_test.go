package honeylogger

import (
	"log/slog"
	"testing"

	"gotest.tools/v3/assert"
)

func TestHoneyLogger_Enabled(t *testing.T) {
	h := &Handler{Level: slog.LevelInfo}
	assert.Assert(t, h.Enabled(t.Context(), slog.LevelInfo))
	assert.Assert(t, h.Enabled(t.Context(), slog.LevelWarn))
	assert.Assert(t, !h.Enabled(t.Context(), slog.LevelDebug))
}
