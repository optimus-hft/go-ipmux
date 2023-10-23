package ipmux

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	ipMux, err := New([]string{
		"127.0.0.1",
		"1.1.1.1",
	}, WithBaseClient(http.DefaultClient))
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ErrInvalidIP)
	assert.NotNil(t, ipMux)
	assert.Len(t, ipMux.clients, 1)
}

func TestIPMux_Clients(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ipMux := IPMux{}

		clients := ipMux.Clients()
		assert.Len(t, clients, 1)
		assert.Contains(t, clients, http.DefaultClient)
	})
	t.Run("non-empty", func(t *testing.T) {
		secondClient := http.DefaultClient
		secondClient.Timeout = time.Second

		ipMux := IPMux{
			counter: atomic.Uint64{},
			clients: []*http.Client{http.DefaultClient, secondClient},
		}

		clients := ipMux.Clients()
		assert.Len(t, clients, 2)
		assert.Contains(t, clients, secondClient)
		assert.Contains(t, clients, http.DefaultClient)
	})
}

func TestIPMux_Client(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ipMux := IPMux{}

		firstClient := ipMux.Client()
		assert.Equal(t, firstClient, http.DefaultClient)
		secondClient := ipMux.Client()
		assert.Equal(t, secondClient, http.DefaultClient)
	})
	t.Run("non-empty", func(t *testing.T) {
		secondClient := http.DefaultClient
		secondClient.Timeout = time.Second

		ipMux := IPMux{
			counter: atomic.Uint64{},
			clients: []*http.Client{http.DefaultClient, secondClient},
		}

		firstClient := ipMux.Client()
		assert.Equal(t, firstClient, http.DefaultClient)

		second := ipMux.Client()
		assert.Equal(t, second, secondClient)

		assert.Equal(t, ipMux.counter.Load(), uint64(2))
	})
}
