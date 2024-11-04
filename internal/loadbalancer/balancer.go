package loadbalancer

import (
	"errors"
	"net/url"
)

var (
	NoHealthyHostsError        = errors.New("no healthy hosts available")
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

// A Balancer keeps track of a pool of servers and selects the next server to use
// based on its load balancing algorithm
type Balancer interface {
	Add(url.URL)               // Add a new server to the load balancer
	Remove(url.URL)            // Remove a server from the load balancer
	Balance() (url.URL, error) // Return the next server to use
}

var factories = make(map[string]func([]url.URL) Balancer)

func BuildBalancer(algorithm string, servers []url.URL) (Balancer, error) {
	fac, ok := factories[algorithm]
	if !ok {
		return nil, AlgorithmNotSupportedError
	}

	return fac(servers), nil
}
