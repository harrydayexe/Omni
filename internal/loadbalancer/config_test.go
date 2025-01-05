package loadbalancer

import (
	"io"
	"log/slog"
	"testing"
)

func TestReadConfig(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	config, err := ReadConfig("../../testdata/loadbalancer-config.yaml", logger)
	if err != nil {
		t.Fatalf("error reading config: %v", err)
	}

	if config.Algorithm != "round-robin" {
		t.Fatalf("expected algorithm %s, got algorithm %s\n", "round-robin", config.Algorithm)
	}

	if len(config.Paths) != 2 {
		t.Fatalf("expected 2 paths, got %d paths\n", len(config.Paths))
	}

	if config.Paths[0] != "GET /test1" {
		t.Fatalf("unexpected %s, got %s", "GET /test1\n", config.Paths[0])
	}
	if config.Paths[1] != "POST /test2" {
		t.Fatalf("unexpected %s, got %s", "POST /test2\n", config.Paths[1])
	}
}

func TestReadConfigInvalidPath(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	_, err := ReadConfig("invalid-path", logger)
	if err == nil {
		t.Fatalf("expected error reading config\n")
	}
}

func TestReadConfigInvalidYaml(t *testing.T) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelDebug}))
	_, err := ReadConfig("../../testdata/user-repo-no-posts.sql", logger)
	if err == nil {
		t.Fatalf("expected error reading config\n")
	}
}

func TestIsValidConfig(t *testing.T) {
	var cases = []struct {
		name          string
		config        Config
		expectedError bool
	}{
		{
			name: "valid config",
			config: Config{
				Algorithm: "round-robin",
				Paths:     []string{"GET /test1", "POST /test2"},
			},
			expectedError: false,
		},
		{
			name: "invalid algorithm",
			config: Config{
				Algorithm: "invalid-algorithm",
				Paths:     []string{"GET /test1", "POST /test2"},
			},
			expectedError: true,
		},
		{
			name: "invalid path",
			config: Config{
				Algorithm: "round-robin",
				Paths:     []string{"GET /test1", "POST /test2", "invalid-path"},
			},
			expectedError: true,
		},
		{
			name: "empty path",
			config: Config{
				Algorithm: "round-robin",
				Paths:     []string{},
			},
			expectedError: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := c.config.IsValid()
			if c.expectedError && err == nil {
				t.Fatalf("expected error, got nil\n")
			} else if !c.expectedError && err != nil {
				t.Fatalf("expected no error, got %v\n", err)
			}
		})
	}
}
