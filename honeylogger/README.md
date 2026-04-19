# honeylogger

Simple [log/slog.Handler](https://pkg.go.dev/log/slog#Handler) adaptor for [honeycombio/libhoney-go](https://github.com/honeycombio/libhoney-go).

This adaptor sends logs as [Events](https://docs.honeycomb.io/send-data/go/libhoney). If you are looking for something with more features please use [Honeycomb's OpenTelemetry](https://docs.honeycomb.io/send-data/opentelemetry) integration instead.
