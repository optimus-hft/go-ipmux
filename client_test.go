package ipmux

import (
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
	nonDefaultHttpClient := http.DefaultClient
	nonDefaultHttpClient.Timeout = time.Second

	testCases := []struct {
		name string
		expectedLen int
		firstClient *http.Client
		secondClient *http.Client
		ipMuxClient IPMux
	} {
		{
			name: "Empty",
			expectedLen: 1,
			firstClient: http.DefaultClient,
			secondClient: http.DefaultClient,
			ipMuxClient: IPMux{},
		},
		{
			name: "NonEmpty",
			expectedLen: 2,
			firstClient: http.DefaultClient,
			secondClient: nonDefaultHttpClient,
			ipMuxClient: IPMux{
				counter: atomic.Uint64{},
				clients: []*http.Client{http.DefaultClient, nonDefaultHttpClient},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			clients := tc.ipMuxClient.Clients()
			assert.Len(t, clients, tc.expectedLen)
			assert.Contains(t, clients, tc.firstClient)
			assert.Contains(t, clients, tc.secondClient)
		})
	}
}

func TestIPMux_Client(t *testing.T) {

	nonDefaultHttpClient := http.DefaultClient
	nonDefaultHttpClient.Timeout = time.Second

	testCases := []struct {
		name string
		firstClient *http.Client
		secondClient *http.Client
		ipMuxClient IPMux
	} {
		{
			name: "Empty",
			firstClient: http.DefaultClient,
			secondClient: http.DefaultClient,
			ipMuxClient: IPMux{},
		},
		{
			name: "NonEmpty",
			firstClient: http.DefaultClient,
			secondClient: nonDefaultHttpClient,
			ipMuxClient: IPMux{
				counter: atomic.Uint64{},
				clients: []*http.Client{http.DefaultClient, nonDefaultHttpClient},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			firstClient := tc.ipMuxClient.Client()
			secondClient := tc.ipMuxClient.Client()
			assert.Equal(t, firstClient, tc.firstClient)
			assert.Equal(t, secondClient, tc.secondClient)
		})
	}
}
