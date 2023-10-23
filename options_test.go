package ipmux

import (
	"github.com/stretchr/testify/assert"
	"net"
	"net/http"
	"testing"
	"time"
)

func TestWithBaseClient(t *testing.T) {
	defaultBase := defaultClientBaseOpts
	client := http.DefaultClient
	client.Timeout = time.Second
	WithBaseClient(client)(defaultBase)
	assert.Equal(t, defaultBase.client.Timeout, time.Second)
}

func TestWithTimeout(t *testing.T) {
	defaultBase := defaultClientBaseOpts
	WithTimeout(time.Second)(defaultBase)
	assert.Equal(t, defaultBase.client.Timeout, time.Second)
	assert.Equal(t, defaultBase.dialer.Timeout, time.Second)
}

func TestWithKeepAlive(t *testing.T) {
	defaultBase := defaultClientBaseOpts
	WithKeepAlive(time.Second)(defaultBase)
	assert.Equal(t, defaultBase.dialer.KeepAlive, time.Second)
}

func TestWithBaseTransport(t *testing.T) {
	defaultBase := defaultClientBaseOpts
	transport := &http.Transport{
		DialContext:           defaultBase.dialer.DialContext,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          110,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	WithBaseTransport(transport)(defaultBase)
	assert.Equal(t, defaultBase.transport, transport)
}

func TestWithDialer(t *testing.T) {
	defaultBase := defaultClientBaseOpts
	dialer := &net.Dialer{
		Timeout: time.Second,
	}
	WithDialer(dialer)(defaultBase)
	assert.Equal(t, defaultBase.dialer, dialer)
}
