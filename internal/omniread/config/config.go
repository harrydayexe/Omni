package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config is a struct that holds the configuration for the OmniRead application.
type Config struct {
	Server struct {
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		DriverName      string `yaml:"driverName"`
		DataSourceName  string `yaml:"dataSourceName"`
		ConnMaxLifetime int    `yaml:"connMaxLifetime"`
		MaxOpenConns    int    `yaml:"maxOpenConns"`
		MaxIdleConns    int    `yaml:"maxIdleConns"`
	} `yaml:"database"`
}

func NewConfig(configFile string) (*Config, error) {
	f, err := os.Open(configFile)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
