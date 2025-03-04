package config

// Config is a struct that holds the configuration for the Omni applications.
type Config struct {
	Host string `env:"HOST"`
	Port int    `env:"PORT" envDefault:"80"`
}

// DatabaseConfig is a struct that holds the configuration for connecting to a database.
type DatabaseConfig struct {
	Config
	DataSourceName  string `env:"DATA_SOURCE_NAME,required"`
	ConnMaxLifetime int    `env:"CONNECTION_MAX_LIFETIME" envDefault:"3"`
	MaxOpenConns    int    `env:"MAX_OPEN_CONNECTIONS" envDefault:"10"`
	MaxIdleConns    int    `env:"MAX_IDLE_CONNECTIONS" envDefault:"10"`
}

// AuthConfig is a struct that holds the configuration for Omni applications
// that require JWT token auth
type AuthConfig struct {
	DatabaseConfig
	JWTSecret string `env:"JWT_SECRET,required"`
}

// WriteConfig is a struct that holds the configuration for the OmniWrite application.
type WriteConfig struct {
	AuthConfig
	NodeName string `env:"NODE_NAME,required"`
}
