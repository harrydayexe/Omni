package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"

	"github.com/harrydayexe/Omni/internal/loadbalancer/balancer"
)

var (
	XRealIP       = http.CanonicalHeaderKey("X-Real-IP")
	XProxy        = http.CanonicalHeaderKey("X-Proxy")
	XForwardedFor = http.CanonicalHeaderKey("X-Forwarded-For")
	ReverseProxy  = "Omni-LoadBalancer-Proxy"
)

// A LoadBalancerProxy is an HTTP handler that forwards requests to a pool of
// backend servers. It uses a Balancer to select the next server to use.
type LoadBalancerProxy struct {
	serviceMap map[*url.URL]*httputil.ReverseProxy
	balancer   balancer.Balancer

	sync.RWMutex // Protect isAliveMap
	isAliveMap   map[string]bool
}

func customRewrite(rf func(*httputil.ProxyRequest)) func(*httputil.ProxyRequest) {
	return func(r *httputil.ProxyRequest) {
		rf(r)
		r.Out.Header.Set(XRealIP, r.In.RemoteAddr)
		r.Out.Header.Set(XProxy, ReverseProxy)
	}
}

func NewLoadBalancerProxy(algorithm string) (*LoadBalancerProxy, error) {
	services := make(map[*url.URL]*httputil.ReverseProxy)
	isAlive := make(map[string]bool)

	bal, err := balancer.BuildBalancer(algorithm)
	if err != nil {
		return nil, err
	}

	lb := LoadBalancerProxy{
		serviceMap: services,
		isAliveMap: isAlive,
		balancer:   bal,
	}

	return &lb, nil
}

func (p *LoadBalancerProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			// log.Printf("proxy causes panic :%s", err)
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte(err.(error).Error()))
		}
	}()

	host, err := p.balancer.Balance()
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte(fmt.Sprintf("balance error: %s", err.Error())))
		return
	}

	p.serviceMap[host].ServeHTTP(w, r)
}

func (p *LoadBalancerProxy) Add(server *url.URL) {
	p.Lock()
	defer p.Unlock()

	proxy := httputil.NewSingleHostReverseProxy(server)

	rewrite := proxy.Rewrite
	proxy.Rewrite = customRewrite(rewrite)

	// Initially set to not alive
	p.isAliveMap[server.Host] = false
	p.serviceMap[server] = proxy
	p.balancer.Add(server)
}

func (p *LoadBalancerProxy) Remove(server *url.URL) {
	p.Lock()
	defer p.Unlock()

	delete(p.serviceMap, server)
	delete(p.isAliveMap, server.Host)
	p.balancer.Remove(server)
}
