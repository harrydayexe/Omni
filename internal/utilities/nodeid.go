package utilities

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

func GetNodeIDFromDeployment(logger *slog.Logger) (uint16, error) {
	hostname, err := os.Hostname()
	if err != nil {
		logger.Error("Failed to get hostname", slog.Any("error", err))
		return 0, fmt.Errorf("failed to get hostname: %w", err)
	}

	// Split by '-' and take the last two parts as a unique identifier
	parts := strings.Split(hostname, "-")
	if len(parts) < 2 {
		logger.Error("Unexpected hostname format", slog.String("hostname", hostname))
		return 0, fmt.Errorf("unexpected hostname format: %s", hostname)
	}
	uniquePart := parts[len(parts)-2] + "-" + parts[len(parts)-1]

	// Hash the extracted unique part
	hash := sha256.Sum256([]byte(uniquePart))
	// Convert first 2 bytes to uint16
	// The SHA256 hash is 32 bytes long but we only need a NodeId that is 10 bits
	// The first two bytes of the hash are just as random as any other part of the hash
	// therefore we can just use the first two bytes for our modulo operation
	hashedValue := binary.BigEndian.Uint16(hash[:2])

	// Ensure it's within 0-1023 range
	return hashedValue % 1024, nil
}
