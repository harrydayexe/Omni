package loadbalancer

import (
	"context"
	"net/http/httputil"
	"net/url"
	"testing"

	"github.com/harrydayexe/Omni/internal/loadbalancer/balancer"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestReadAliveMap(t *testing.T) {
	// Create a new LoadBalancerProxy
	lb := LoadBalancerProxy{
		serviceMap: make(map[*url.URL]*httputil.ReverseProxy),
		isAliveMap: make(map[string]bool),
		balancer:   &balancer.RoundRobinBalancer{},
	}

	testUrl, _ := url.Parse("https://127.0.0.1/")

	lb.Lock()
	lb.isAliveMap[testUrl.Host] = true
	lb.Unlock()

	if !lb.ReadAliveMap(testUrl) {
		t.Errorf("Expected true, got false")
	}
}

func createHealthCheckTester(t *testing.T, envVars map[string]string) (testcontainers.Container, error) {
	ctx := context.Background()
	req := testcontainers.ContainerRequest{
		Image:        "harrydayexe/healthcheck-tester",
		ExposedPorts: []string{"80/tcp"},
		WaitingFor:   wait.ForHTTP("/testcontainersz"),
		Env:          envVars,
	}

	htc, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	testcontainers.CleanupContainer(t, htc)
	if err != nil {
		return nil, err
	}

	return htc, nil
}

func TestCheckHealth(t *testing.T) {
	cases := []struct {
		name     string
		envVars  map[string]string
		expected bool
	}{
		{
			name: "healthy",
			envVars: map[string]string{
				"READYZ": "TRUE",
			},
			expected: true,
		},
		{
			name: "unhealthy",
			envVars: map[string]string{
				"READYZ": "FALSE",
			},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			htc, err := createHealthCheckTester(t, c.envVars)
			if err != nil {
				t.Errorf("Error starting container: %v", err)
			}
			endpoint, err := htc.Endpoint(context.Background(), "")
			if err != nil {
				t.Errorf("Error getting endpoint: %v", err)
			}

			url := &url.URL{
				Scheme: "http",
				Host:   endpoint,
			}

			if CheckHealth(url) != c.expected {
				t.Errorf("Expected %v, got %v", c.expected, !c.expected)
			}
		})
	}
}

func TestCheckHealthUnreachable(t *testing.T) {
	if CheckHealth(&url.URL{
		Scheme: "http",
		Host:   "localhost:1234",
	}) {
		t.Errorf("Expected %v, got %v", false, true)
	}
}

func TestHealthCheck(t *testing.T) {
	cases := []struct {
		name          string
		envVars       map[string]string
		initialHealth bool
		loadBalancer  *LoadBalancerProxy
		expectedAlive bool
	}{
		{
			name: "currently healthy, is still healthy",
			envVars: map[string]string{
				"READYZ": "TRUE",
			},
			initialHealth: true,
			loadBalancer: &LoadBalancerProxy{
				isAliveMap: make(map[string]bool),
				serviceMap: make(map[*url.URL]*httputil.ReverseProxy),
				balancer:   balancer.NewRoundRobinBalancer(),
			},
			expectedAlive: true,
		},
		{
			name: "currently unhealthy, is still unhealthy",
			envVars: map[string]string{
				"READYZ": "FALSE",
			},
			initialHealth: false,
			loadBalancer: &LoadBalancerProxy{
				isAliveMap: make(map[string]bool),
				serviceMap: make(map[*url.URL]*httputil.ReverseProxy),
				balancer:   balancer.NewRoundRobinBalancer(),
			},
			expectedAlive: false,
		},
		{
			name: "currently healthy, is now unhealthy",
			envVars: map[string]string{
				"READYZ": "FALSE",
			},
			initialHealth: true,
			loadBalancer: &LoadBalancerProxy{
				isAliveMap: make(map[string]bool),
				serviceMap: make(map[*url.URL]*httputil.ReverseProxy),
				balancer:   balancer.NewRoundRobinBalancer(),
			},
			expectedAlive: false,
		},
		{
			name: "currently unhealthy, is now healthy",
			envVars: map[string]string{
				"READYZ": "TRUE",
			},
			initialHealth: false,
			loadBalancer: &LoadBalancerProxy{
				isAliveMap: make(map[string]bool),
				serviceMap: make(map[*url.URL]*httputil.ReverseProxy),
				balancer:   balancer.NewRoundRobinBalancer(),
			},
			expectedAlive: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			htc, err := createHealthCheckTester(t, c.envVars)
			if err != nil {
				t.Errorf("Error starting container: %v", err)
			}
			endpoint, err := htc.Endpoint(context.Background(), "")
			if err != nil {
				t.Errorf("Error getting endpoint: %v", err)
			}

			url := &url.URL{
				Scheme: "http",
				Host:   endpoint,
			}

			c.loadBalancer.balancer.Add(url)
			c.loadBalancer.isAliveMap[url.Host] = c.initialHealth

			c.loadBalancer.healthCheck(context.Background(), url)

			if c.loadBalancer.isAliveMap[url.Host] != c.expectedAlive {
				t.Errorf("Expected %v, got %v", c.expectedAlive, !c.expectedAlive)
			}
		})
	}
}
