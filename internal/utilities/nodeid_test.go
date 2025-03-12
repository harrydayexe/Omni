package utilities

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"testing"
)

func TestGetNodeIDValid(t *testing.T) {
	nodename := "homelab1"
	expected := uint16(707)

	nodeID, err := GetNodeIDFromDeployment(slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})), nodename)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if nodeID != expected {
		t.Errorf("Expected %v, got %v", expected, nodeID)
	}
}

func TestGetNodeIDInvalid(t *testing.T) {
	hostname := "omnireadapi"

	_, err := GetNodeIDFromDeployment(slog.New(slog.NewTextHandler(io.Discard, &slog.HandlerOptions{})), hostname)
	if errors.Is(err, fmt.Errorf("unexpected hostname format: %s", hostname)) {
		t.Errorf("Did not receive expected error. Got: %v, want %v", err, fmt.Errorf("unexpected hostname format: %s", hostname))
	}
}
