package utilities

import (
	"crypto/sha256"
	"encoding/binary"
	"log/slog"
)

// GetNodeIDFromDeployment generates a unique NodeId from the name of the
// Kubernetes node it is running on
func GetNodeIDFromDeployment(logger *slog.Logger, nodeName string) (uint16, error) {
	hash := sha256.Sum256([]byte(nodeName))
	// Convert first 2 bytes to uint16
	// The SHA256 hash is 32 bytes long but we only need a NodeId that is 10 bits
	// The first two bytes of the hash are just as random as any other part of the hash
	// therefore we can just use the first two bytes for our modulo operation
	hashedValue := binary.BigEndian.Uint16(hash[:2])
	// Ensure it's within 0-1023 range
	nodeid := hashedValue % 1024

	logger.Info("Node ID", slog.Int("nodeid", int(nodeid)))

	return nodeid, nil
}
