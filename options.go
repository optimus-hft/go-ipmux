package ipmux

import (
	"net"
	"net/http"
	"time"
)

var defaultBaseClient = http.DefaultClient
var defaultBaseTransport = http.DefaultTransport.(*http.Transport)

var defaultClientBaseOpts = &clientBaseOpts{
	client:    defaultBaseClient,
	transport: defaultBaseTransport,
	dialer: &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	},
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
