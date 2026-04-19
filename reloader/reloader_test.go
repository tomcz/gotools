package reloader

import (
	"context"
	"crypto/tls"
	"log/slog"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestReloader_Compatibility(t *testing.T) {
	reloader := &Reloader{}
	_ = &tls.Config{
		GetCertificate:       reloader.GetCertificate,
		GetClientCertificate: reloader.GetClientCertificate,
	}
}

func TestReloader_GetCertificate(t *testing.T) {
	reloader, err := New(t.Context(), "testdata/dev.crt", "testdata/dev.key", time.Hour)
	assert.NilError(t, err)
	cert, err := reloader.GetCertificate(nil)
	assert.NilError(t, err)
	assert.Equal(t, "DevCert", cert.Leaf.Subject.CommonName)
}

func TestReloader_GetClientCertificate(t *testing.T) {
	reloader, err := New(t.Context(), "testdata/dev.crt", "testdata/dev.key", time.Hour)
	assert.NilError(t, err)
	cert, err := reloader.GetClientCertificate(nil)
	assert.NilError(t, err)
	assert.Equal(t, "DevCert", cert.Leaf.Subject.CommonName)
}

func TestReloader_reloadCertificate(t *testing.T) {
	reloader := &Reloader{
		certFile: "testdata/dev.crt",
		keyFile:  "testdata/dev.key",
		interval: 10 * time.Millisecond,
		log:      slog.New(slog.DiscardHandler),
	}
	assert.Assert(t, reloader.value.Load() == nil)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go reloader.reloadCertificate(ctx)

	<-ctx.Done()
	assert.Assert(t, reloader.value.Load() != nil)
}
