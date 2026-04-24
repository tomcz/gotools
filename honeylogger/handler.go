package honeylogger

import (
	"context"
	"fmt"
	"log/slog"
	"maps"
	"slices"
	"strings"

	"github.com/honeycombio/libhoney-go"
)

// Handler implements [slog.Handler] to send log records to [honeycomb.io] as events.
//
// Both Level and Builder are optional. The default [slog] level is [slog.LevelInfo].
// This handler will use [libhoney.NewEvent] to create events if a [libhoney.Builder]
// has not been provided.
//
// Reject is an optional filter for log records that should not be sent as events.
// The default behaviour is to forward all records that are enabled for each record's
// log level. This can be further extended by providing a filter that can perform
// additional inspection of records and returns true if the record should not be
// sent as an event.
//
// [honeycomb.io]: https://www.honeycomb.io
type Handler struct {
	Level   slog.Level             // Optional
	Builder *libhoney.Builder      // Optional
	Reject  func(slog.Record) bool // Optional
	attrs   map[string]slog.Value
	groups  []string
}

var _ slog.Handler = (*Handler)(nil)

func (h *Handler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.Level
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	if !h.Enabled(ctx, record.Level) {
		return nil
	}
	if h.Reject != nil && h.Reject(record) {
		return nil
	}
	params := make(map[string]any)
	for key, value := range h.attrs {
		params[key] = eventValue(value)
	}
	record.Attrs(func(attr slog.Attr) bool {
		key := h.groupKey(attr.Key)
		params[key] = eventValue(attr.Value)
		return true
	})
	params["event"] = record.Message
	params["level"] = record.Level.String()
	var event *libhoney.Event
	if h.Builder != nil {
		event = h.Builder.NewEvent()
	} else {
		event = libhoney.NewEvent()
	}
	err := event.Add(params)
	if err != nil {
		return err
	}
	return event.Send()
}

func eventValue(value slog.Value) any {
	obj := value.Any()
	if value.Kind() == slog.KindAny {
		if str, ok := obj.(fmt.Stringer); ok {
			return str.String()
		}
		if err, ok := obj.(error); ok {
			return err.Error()
		}
	}
	return obj
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var newAttrs map[string]slog.Value
	if len(h.attrs) > 0 {
		newAttrs = maps.Clone(h.attrs)
	} else {
		newAttrs = make(map[string]slog.Value)
	}
	for _, attr := range attrs {
		key := h.groupKey(attr.Key)
		newAttrs[key] = attr.Value
	}
	var newGroups []string
	if len(h.groups) > 0 {
		newGroups = slices.Clone(h.groups)
	}
	return &Handler{
		Level:   h.Level,
		Builder: h.Builder,
		attrs:   newAttrs,
		groups:  newGroups,
	}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	var newAttrs map[string]slog.Value
	if len(h.attrs) > 0 {
		newAttrs = maps.Clone(h.attrs)
	}
	newGroups := []string{name}
	if len(h.groups) > 0 {
		newGroups = slices.Concat(h.groups, newGroups)
	}
	return &Handler{
		Level:   h.Level,
		Builder: h.Builder,
		attrs:   newAttrs,
		groups:  newGroups,
	}
}

func (h *Handler) groupKey(key string) string {
	switch len(h.groups) {
	case 0:
		return key
	case 1:
		return h.groups[0] + "_" + key
	default:
		return strings.Join(h.groups, "_") + "_" + key
	}
}
