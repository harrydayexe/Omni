package balancer

import (
	"errors"
	"net/url"
	"sync"
)

var (
	NoHealthyHostsError        = errors.New("no healthy hosts available")
	AlgorithmNotSupportedError = errors.New("algorithm not supported")
)

// A Balancer keeps track of a pool of servers and selects the next server to use
// based on its load balancing algorithm
type Balancer interface {
	Add(*url.URL)               // Add a new server to the load balancer
	Remove(*url.URL)            // Remove a server from the load balancer
	Balance() (*url.URL, error) // Return the next server to use
	Len() int                   // Return the number of servers in the load balancer
}

var factories = make(map[string]func([]*url.URL) Balancer)

func BuildBalancer(algorithm string, servers []*url.URL) (Balancer, error) {
	fac, ok := factories[algorithm]
	if !ok {
		return nil, AlgorithmNotSupportedError
	}

	return fac(servers), nil
}

type BaseBalancer struct {
	sync.RWMutex
	servers []*url.URL
}

func (b *BaseBalancer) Add(server *url.URL) {
	b.Lock()
	defer b.Unlock()
	for _, s := range b.servers {
		if s.String() == server.String() {
			return
		}
	}
	b.servers = append(b.servers, server)
}

func (b *BaseBalancer) Remove(server *url.URL) {
	b.Lock()
	defer b.Unlock()
	for i, s := range b.servers {
		if s.String() == server.String() {
			b.servers = append(b.servers[:i], b.servers[i+1:]...)
			return
		}
	}
}

func (b *BaseBalancer) Balance() (*url.URL, error) {
	return &url.URL{}, NoHealthyHostsError
}

func (b *BaseBalancer) Len() int {
	b.RLock()
	defer b.RUnlock()
	return len(b.servers)
}
