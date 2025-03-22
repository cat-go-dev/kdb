package config

type Config struct {
	Engine  Engine  `mapstructure:"engine"`
	Network Network `mapstructure:"network"`
	Logging Logging `mapstructure:"logging"`
}

type Engine struct {
	Type string `mapstructure:"type"`
}

type Network struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	MaxConnections int    `mapstructure:"max_connections"`
	MaxMessageSize string `mapstructure:"max_message_size"`
	IdleTimeout    string `mapstructure:"idle_timeout"`
}

type Logging struct {
	Level  string `mapstructure:"level"`
	OutputDir string `mapstructure:"output_dir"`
}
