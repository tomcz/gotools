package reloader

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"sync/atomic"
	"time"
)

// Reloader for [tls.Config] certificate and private key files.
type Reloader struct {
	certFile string
	keyFile  string
	interval time.Duration
	log      *slog.Logger
	value    atomic.Value
}

// New creates a new [Reloader] that will use a background goroutine to
// periodically reload the given certificate and private key files.
// The goroutine will keep on running until the given context is done.
func New(ctx context.Context, certificateFile string, privateKeyFile string, reloadInterval time.Duration) (*Reloader, error) {
	reloader := &Reloader{
		certFile: certificateFile,
		keyFile:  privateKeyFile,
		interval: reloadInterval,
		log:      slog.New(slog.DiscardHandler),
	}
	err := reloader.loadCertificate()
	if err != nil {
		return nil, fmt.Errorf("loadCertificate: %w", err)
	}
	go reloader.reloadCertificate(ctx)
	return reloader, nil
}

// SetLogger overrides the default logger created with a [slog.DiscardHandler].
// The logger is used to record errors found during the background certificate
// reload at warn level, and successful certificate load events at debug level.
func (r *Reloader) SetLogger(log *slog.Logger) {
	r.log = log
}

// GetCertificate returns the managed TLS certificate to be used for server authentication.
// See [tls.Config] for more details.
//
//goland:noinspection GoUnusedParameter
func (r *Reloader) GetCertificate(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return r.value.Load().(*tls.Certificate), nil
}

// GetClientCertificate returns the managed TLS certificate to be used for client authentication.
// See [tls.Config] for more details.
//
//goland:noinspection GoUnusedParameter
func (r *Reloader) GetClientCertificate(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
	return r.value.Load().(*tls.Certificate), nil
}

func (r *Reloader) loadCertificate() error {
	cert, err := tls.LoadX509KeyPair(r.certFile, r.keyFile)
	if err != nil {
		return err
	}
	r.log.Debug("loadCertificate successful")
	r.value.Store(&cert)
	return nil
}

func (r *Reloader) reloadCertificate(ctx context.Context) {
	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := r.loadCertificate(); err != nil {
				r.log.Warn("loadCertificate", "err", err)
			}
		}
	}
}
