package balancer

import (
	"net/url"
	"testing"
)

func TestBuildBalancer(t *testing.T) {
	// Test that the BuildBalancer function returns an error when the algorithm is not supported
	_, err := BuildBalancer("unsupported", nil)
	if err != AlgorithmNotSupportedError {
		t.Errorf("expected AlgorithmNotSupportedError, got %v", err)
	}
}

func TestBaseBalancer_Len(t *testing.T) {
	// Test that the Len method returns the correct number of servers
	b := &BaseBalancer{servers: []*url.URL{{}, {}}}
	if b.Len() != 2 {
		t.Errorf("expected 2, got %v", b.Len())
	}
}
