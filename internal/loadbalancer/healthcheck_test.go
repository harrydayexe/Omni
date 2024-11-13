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
			t.Log(endpoint)

			url := &url.URL{
				Scheme: "http",
				Host:   endpoint,
			}
			t.Logf("URL: %v", url)

			if CheckHealth(url) != c.expected {
				t.Errorf("Expected %v, got %v", c.expected, !c.expected)
			}
		})
	}
}
