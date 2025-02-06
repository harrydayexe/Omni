package config

// Config is a struct that holds the configuration for the OmniRead application.
type Config struct {
	Host            string `env:"HOST"`
	Port            int    `env:"PORT" envDefault:"80"`
	DataSourceName  string `env:"DATA_SOURCE_NAME,required"`
	ConnMaxLifetime int    `env:"CONNECTION_MAX_LIFETIME" envDefault:"3"`
	MaxOpenConns    int    `env:"MAX_OPEN_CONNECTIONS" envDefault:"10"`
	MaxIdleConns    int    `env:"MAX_IDLE_CONNECTIONS" envDefault:"10"`
	JWTSecret       string `env:"JWT_SECRET,required"`
}
