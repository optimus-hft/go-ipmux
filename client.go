package ipmux

import (
	"context"
	"net"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/hashicorp/go-multierror"
	"github.com/rs/dnscache"
)

var ErrInvalidIP = errors.New("invalid source ip, not found in device network interfaces")

type IPMux struct {
	counter   atomic.Uint64
	clients   []*http.Client
	ctxCancel context.CancelFunc
	dnsCache  *dnscache.Resolver
}

// Client returns one of the clients that is associated with one of the IPs given in New.
// the function is safe to use without error handling of the constructor. it returns http.DefaultClient when there are no available clients.
func (i *IPMux) Client() *http.Client {
	length := uint64(len(i.clients))
	if length == 0 {
		return http.DefaultClient
	}
	defer i.counter.Add(1)

	return i.clients[i.counter.Load()%length]
}

// Clients returns a list of all clients created for the list of ips given in New.
// the function is safe to use without error handling of the constructor. it returns a list containing ob client (http.DefaultClient) when there are no available clients.
func (i *IPMux) Clients() []*http.Client {
	if len(i.clients) == 0 {
		return []*http.Client{http.DefaultClient}
	}

	return i.clients
}

// refreshDNSCache runs in each refreshInterval and refreshes the dns cache.
func (i *IPMux) refreshDNSCache(ctx context.Context, refreshInterval time.Duration) {
	ticker := time.NewTicker(refreshInterval)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			i.dnsCache.Refresh(true)
		}
	}

}

// Stop stops the ipMux.
func (i *IPMux) Stop() {
	i.ctxCancel()
}

// New is the constructor of IPMux. It creates a http.Client for each of the ips given. If there are any errors for one of the ips, the client will not be created for that ip but the other clients will be created.
// you can customize the created clients with Option functions.
func New(ips []string, options ...Option) (*IPMux, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return &IPMux{}, errors.Wrap(err, "could not list network interface addresses")
	}

	ipmux := &IPMux{
		clients: make([]*http.Client, 0),
	}

	clientAddrs := make([]net.Addr, 0, len(addrs))

	var resultErr error
	for _, ip := range ips {
		if addr, exists := ipExistingInAddrs(addrs, ip); !exists {
			resultErr = multierror.Append(errors.WithDetailf(ErrInvalidIP, "IP: %s", ip))
		} else {
			clientAddrs = append(clientAddrs, addr)
		}
	}

	clientOpts := getDefaultClientBaseOpts()
	for _, option := range options {
		option(clientOpts)
	}

	ctx, cancel := context.WithCancel(clientOpts.ctx)
	ipmux.ctxCancel = cancel
	ipmux.dnsCache = clientOpts.resolver

	if ipmux.dnsCache != nil {
		go ipmux.refreshDNSCache(ctx, clientOpts.refreshInterval)
	}

	for _, addr := range clientAddrs {
		ipmux.clients = append(ipmux.clients, createAddressedClient(addr, clientOpts))
	}

	if len(ipmux.clients) == 0 {
		ipmux.clients = append(ipmux.clients, createDefaultClient(clientOpts))
	}

	return ipmux, resultErr
}

func ipExistingInAddrs(addrs []net.Addr, ip string) (net.Addr, bool) {
	for _, addr := range addrs {
		if strings.Contains(addr.String(), ip) {
			localAddr, err := net.ResolveIPAddr("ip", ip)
			if err != nil {
				continue
			}

			return &net.TCPAddr{
				IP: localAddr.IP,
			}, true
		}
	}

	return nil, false
}

func createAddressedClient(addr net.Addr, opts *clientBaseOpts) *http.Client {
	dialer := opts.dialer
	dialer.LocalAddr = addr

	return createDefaultClient(opts)
}

func createDefaultClient(opts *clientBaseOpts) *http.Client {
	client := opts.client
	dialer := opts.dialer
	transport := opts.transport

	transport.DialContext = dialer.DialContext

	if opts.resolver != nil {
		//nolint:nonamedreturns
		transport.DialContext = func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			host, port, err := net.SplitHostPort(addr)
			if err != nil {
				//nolint:wrapcheck
				return nil, err
			}
			ips, err := opts.resolver.LookupHost(ctx, host)
			if err != nil {
				//nolint:wrapcheck
				return nil, err
			}

			for _, ip := range ips {
				conn, err = dialer.DialContext(ctx, network, net.JoinHostPort(ip, port))
				if err == nil {
					break
				}
			}

			return
		}
	}

	client.Transport = transport

	return client
}
