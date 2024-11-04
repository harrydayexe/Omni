package balancer

import (
	"net/url"
	"sync/atomic"
)

type RoundRobinBalancer struct {
	BaseBalancer
	current atomic.Uint64
}

func (r *RoundRobinBalancer) Balance() (url.URL, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.servers) == 0 {
		return url.URL{}, NoHealthyHostsError
	}

	i := r.current.Add(1) % uint64(len(r.servers))
	return r.servers[i], nil
}

func NewRoundRobinBalancer(servers []url.URL) Balancer {
	return &RoundRobinBalancer{
		BaseBalancer: BaseBalancer{
			servers: servers,
		},
		current: atomic.Uint64{},
	}
}

func init() {
	factories["round_robin"] = NewRoundRobinBalancer
}
