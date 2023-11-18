package ipmux

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/rs/dnscache"
)

const (
	defaultTimeout               = 30 * time.Second
	defaultKeepAlive             = 30 * time.Second
	defaultMaxIdleConns          = 100
	defaultIdleConnTimeout       = 90 * time.Second
	defaultTLSHandshakeTimeout   = 10 * time.Second
	defaultExpectContinueTimeout = 1 * time.Second
)

func getDefaultClientBaseOpts() *clientBaseOpts {
	dialer := &net.Dialer{
		Timeout:   defaultTimeout,
		KeepAlive: defaultKeepAlive,
	}
	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          defaultMaxIdleConns,
		IdleConnTimeout:       defaultIdleConnTimeout,
		TLSHandshakeTimeout:   defaultTLSHandshakeTimeout,
		ExpectContinueTimeout: defaultExpectContinueTimeout,
	}

	return &clientBaseOpts{
		client: &http.Client{
			Transport: transport,
		},
		dialer:          dialer,
		transport:       transport,
		refreshInterval: time.Second,
		ctx:             context.Background(),
	}
}

// nolint:containedctx
type clientBaseOpts struct {
	client          *http.Client
	transport       *http.Transport
	dialer          *net.Dialer
	resolver        *dnscache.Resolver
	refreshInterval time.Duration
	ctx             context.Context
}

type Option func(base *clientBaseOpts)

// WithBaseClient is used to change the base client that is used to create the clients for IPMux.
func WithBaseClient(baseClient *http.Client) Option {
	return func(base *clientBaseOpts) {
		base.client = baseClient
	}
}

// WithTimeout is used to set timeout on both client and dialer of the clients.
func WithTimeout(timeout time.Duration) Option {
	return func(base *clientBaseOpts) {
		base.client.Timeout = timeout
		base.dialer.Timeout = timeout
	}
}

// WithKeepAlive is used to set keepalive on dialer of the clients.
func WithKeepAlive(keepalive time.Duration) Option {
	return func(base *clientBaseOpts) {
		base.dialer.KeepAlive = keepalive
	}
}

// WithBaseTransport is used to change the base transport (http.RoundTripper) that is used to create the clients for IPMux.
func WithBaseTransport(transport *http.Transport) Option {
	return func(base *clientBaseOpts) {
		base.transport = transport
	}
}

// WithDialer is used to change the base dialer that is used to create the clients for IPMux.
func WithDialer(dialer *net.Dialer) Option {
	return func(base *clientBaseOpts) {
		base.dialer = dialer
	}
}

// WithDNSCache is used to enable dns cache for the clients.
func WithDNSCache(refreshInterval time.Duration) Option {
	return func(base *clientBaseOpts) {
		base.resolver = &dnscache.Resolver{}
		base.refreshInterval = refreshInterval
	}
}

// WithContext is used to change the context that is used to detect when to stop background goroutines like refreshing dns cache.
func WithContext(ctx context.Context) Option {
	return func(base *clientBaseOpts) {
		base.ctx = ctx
	}
}
