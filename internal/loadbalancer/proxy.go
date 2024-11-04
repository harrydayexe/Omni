package loadbalancer

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
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
	serviceMap map[url.URL]*httputil.ReverseProxy
	balancer   Balancer

	sync.RWMutex // Protect isAliveMap
	isAliveMap   map[url.URL]bool
}

func NewLoadBalancerProxy(algorithm string, targetServers []url.URL) (*LoadBalancerProxy, error) {
	services := make(map[url.URL]*httputil.ReverseProxy)
	isAlive := make(map[url.URL]bool)

	for _, server := range targetServers {
		proxy := httputil.NewSingleHostReverseProxy(&server)

		rewrite := proxy.Rewrite
		proxy.Rewrite = func(r *httputil.ProxyRequest) {
			rewrite(r)
			r.Out.Header.Set(XRealIP, r.In.RemoteAddr)
			r.Out.Header.Set(XProxy, ReverseProxy)
		}

		isAlive[server] = false
		services[server] = proxy
	}

	bal, err := BuildBalancer(algorithm, targetServers)
	if err != nil {
		return nil, err
	}

	return &LoadBalancerProxy{
		serviceMap: services,
		balancer:   bal,
		isAliveMap: isAlive,
	}, nil
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
