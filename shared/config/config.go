package config

type Config struct {
	DSN        string
	ServerPort uint16
}

func New(DSN string, ServerPort uint16) *Config {
	return &Config{DSN: DSN, ServerPort: ServerPort}
}
