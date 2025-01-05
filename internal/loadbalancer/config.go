package loadbalancer

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Algorithm string   `yaml:"algorithm"`
	Paths     []string `yaml:"paths"`
}

// ReadConfig read configuration from `fileName` file
func ReadConfig(fileName string, logger *slog.Logger) (Config, error) {
	in, err := os.ReadFile(fileName)
	if err != nil {
		logger.Error("error reading file", slog.String("file", fileName), slog.String("error message", err.Error()))
		return Config{}, err
	}
	var config Config
	err = yaml.Unmarshal(in, &config)
	if err != nil {
		logger.Error("error unmarshalling yaml", slog.String("error message", err.Error()))
		return Config{}, err
	}

	config.Print(logger)
	return config, nil
}

func (c *Config) Print(logger *slog.Logger) {
	logger.Info("Algorithm", slog.String("algorithm", c.Algorithm))
	paths := make([]any, len(c.Paths))
	for i, path := range c.Paths {
		paths[i] = slog.String("path", path)
	}
	logger.Info("Paths", paths...)
}

func (c *Config) IsValid() error {
	if c.Algorithm != "round-robin" {
		return errors.New("the algorithm is unknown")
	}

	if len(c.Paths) == 0 {
		return errors.New("no paths are defined")
	}

	for _, path := range c.Paths {
		_, err := parsePattern(path)
		if err != nil {
			return fmt.Errorf("invalid path pattern found for %s when parsing: %w", path, err)
		}

	}

	return nil
}
