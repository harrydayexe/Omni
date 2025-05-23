package loadbalancer

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"time"
)

func CheckHealth(serviceURL *url.URL) bool {
	client := http.Client{
		Timeout: 5 * time.Second, // Set a timeout to avoid hanging indefinitely
	}

	healthEndpoint := url.URL{
		Scheme: serviceURL.Scheme,
		Host:   serviceURL.Host,
		Path:   "/readyz",
	}

	resp, err := client.Get(healthEndpoint.String())
	if err != nil {
		log.Printf("Error reaching health endpoint: %v", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func (p *LoadBalancerProxy) ReadAliveMap(server *url.URL) bool {
	p.RLock()
	defer p.RUnlock()
	return p.isAliveMap[server.Host]
}

func (p *LoadBalancerProxy) StartHealthCheck(ctx context.Context, interval time.Duration) {
	for server := range p.serviceMap {
		go p.runHealthCheck(ctx, server, interval)
	}
}

func (p *LoadBalancerProxy) runHealthCheck(ctx context.Context, server *url.URL, interval time.Duration) {
	ticker := time.NewTicker(interval * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			p.healthCheck(ctx, server)
		}
	}
}

func (p *LoadBalancerProxy) healthCheck(ctx context.Context, server *url.URL) {
	if CheckHealth(server) && !p.ReadAliveMap(server) {
		p.Lock()
		p.isAliveMap[server.Host] = true
		p.Unlock()
		p.balancer.Add(server)
	} else if !CheckHealth(server) && p.ReadAliveMap(server) {
		p.Lock()
		p.isAliveMap[server.Host] = false
		p.Unlock()
		p.balancer.Remove(server)
	}
}
