# IPMux
![pipeline](https://github.com/optimus-hft/go-ipmux/actions/workflows/go-ci.yml/badge.svg)
[![codecov](https://codecov.io/gh/optimus-hft/go-ipmux/branch/main/graph/badge.svg)](#)
[![Go Report Card](https://goreportcard.com/badge/github.com/optimus-hft/go-ipmux)](https://goreportcard.com/report/github.com/optimus-hft/go-ipmux)
[![Go Reference](https://pkg.go.dev/badge/github.com/optimus-hft/go-ipmux.svg)](https://pkg.go.dev/github.com/optimus-hft/go-ipmux)

## GoLang Library for Multiplexing HTTP Clients based on Source IP
IPMux is an open-source GoLang library that provides a simple and efficient way to multiplex HTTP clients based on source IP addresses. This library is designed to handle scenarios where you need to make HTTP requests from specific network interfaces based on the source IP.

This library could be useful to bypass rate limit errors which operate on source IPs.

## Features
+ Multiplex HTTP clients based on source IP addresses.
+ Customizable options for creating HTTP clients.
+ Easy-to-use interface for seamless integration.


## Getting Started
### Installation
```
go get github.com/optimus-hft/go-ipmux
```

### Usage

```go
package main

import (
	"fmt"
	"io"
	"time"

	"github.com/optimus-hft/go-ipmux"
)

func main() {
	// both ips should be attached to your device for IPMux to operate properly
	ips := []string{"192.168.0.1", "192.168.0.2"}
	ipMux, err := ipmux.New(ips, ipmux.WithTimeout(10*time.Second))
	if err != nil {
		fmt.Println(err)
		return
	}

	client := ipMux.Client()
	response, err := client.Get("https://google.com")
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	fmt.Println(io.ReadAll(response.Body))
}

```

## API Documentation
### `func New(ips []string, options ...Option) (*IPMux, error)`
+ ips: List of IP addresses to create HTTP clients for.
+ options: Optional parameters to customize client creation.
Returns a new IPMux instance.

### `func (i *IPMux) Client() *http.Client`
Returns an HTTP client associated with one of the IPs provided in the New function. Returns `http.DefaultClient` if no clients are available.

### `func (i *IPMux) Clients() []*http.Client`
Returns a list of all clients created for the list of IPs given in `New`. Returns a list containing one client (`http.DefaultClient`) when there are no available clients.

### Options
### `WithBaseClient(baseClient *http.Client) Option`
Change the base client used to create clients for IPMux.

### `WithTimeout(timeout time.Duration) Option`
Set timeout on both client and dialer of the clients.

### `WithKeepAlive(keepalive time.Duration) Option`
Set keepalive on dialer of the clients.

### `WithBaseTransport(transport *http.Transport) Option`
Change the base transport (http.RoundTripper) used to create clients for IPMux.

### `WithDialer(dialer *net.Dialer) Option`
Change the base dialer used to create clients for IPMux.

## Contributing
Pull requests and bug reports are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
This project is licensed under the MIT License.
