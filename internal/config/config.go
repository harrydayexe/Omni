package config

// Config is a struct that holds the configuration for the OmniRead application.
type Config struct {
	Port            int    `env:"PORT"`
	DataSourceName  string `env:"DATA_SOURCE_NAME"`
	ConnMaxLifetime int    `env:"CONNECTION_MAX_LIFETIME"`
	MaxOpenConns    int    `env:"MAX_OPEN_CONNECTIONS"`
	MaxIdleConns    int    `env:"MAX_IDLE_CONNECTIONS"`
}
