package ipmux

import (
	"net"
	"net/http"
	"time"
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
		dialer:    dialer,
		transport: transport,
	}
}

type clientBaseOpts struct {
	client    *http.Client
	transport *http.Transport
	dialer    *net.Dialer
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
