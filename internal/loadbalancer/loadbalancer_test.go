package loadbalancer

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/harrydayexe/Omni/internal/loadbalancer/balancer"
)

func TestReadyz(t *testing.T) {
	cases := []struct {
		name    string
		proxies map[string]*LoadBalancerProxy
		want    int
	}{
		{
			name:    "no paths",
			proxies: make(map[string]*LoadBalancerProxy),
			want:    http.StatusServiceUnavailable,
		},
		{
			name: "proxy with no backends",
			proxies: map[string]*LoadBalancerProxy{
				"GET /": {
					serviceMap: make(map[*url.URL]*httputil.ReverseProxy),
					isAliveMap: make(map[string]bool),
					balancer:   &balancer.RoundRobinBalancer{},
				},
			},
			want: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
			rr := httptest.NewRecorder()

			lb := &loadBalancer{
				Config:  Config{},
				Logger:  slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug})),
				Proxies: tt.proxies,
				Mux:     http.NewServeMux(),
			}

			lb.readyz(rr, req)
			resp := rr.Result()
			if resp.StatusCode != tt.want {
				t.Errorf("got %d, want %d", resp.StatusCode, tt.want)
			}
		})
	}

}
