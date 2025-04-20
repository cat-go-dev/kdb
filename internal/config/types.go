package config

type Config struct {
	Engine  Engine  `mapstructure:"engine"`
	Network Network `mapstructure:"network"`
	Logging Logging `mapstructure:"logging"`
	WAL     WAL     `mapstructure:"wal"`
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
	Level     string `mapstructure:"level"`
	OutputDir string `mapstructure:"output_dir"`
}

type WAL struct {
	Flushing       Flushing `mapstructure:"flushing"`
	MaxSegmentSize string   `mapstructure:"max_segment_size"`
	DataDitrectory string   `mapstructure:"data_directory"`
}

type Flushing struct {
	BatchSize    int    `mapstructure:"batch_size"`
	BatchTimeout string `mapstructure:"batch_timeout"`
}
