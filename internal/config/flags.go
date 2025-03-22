package config

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	flagEngineType = "engine_type"

	flagHost           = "host"
	flagPort           = "port"
	flagMaxConnections = "max_connections"
	flagMaxMessageSize = "max_message_size"
	flagIdleTimeout    = "idle_timeout"

	flagLogLevel     = "log_level"
	flagLogOutputDir = "output_dir"
)

func (a *AppConfig) overrideByFlags() {
	a.overideEngine()
	a.overideNetwork()
	a.overideLogging()
}

func (a *AppConfig) overideEngine() {
	pflag.String(flagEngineType, "", "engine type")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	engineType := viper.GetString(flagEngineType)
	if engineType != "" {
		a.Data.Engine.Type = engineType
	}
}

func (a *AppConfig) overideNetwork() {
	pflag.String(flagHost, "", "server host")
	pflag.Int(flagPort, 0, "server port")
	pflag.Int(flagMaxConnections, 0, "max connections")
	pflag.String(flagMaxMessageSize, "", "max message size")
	pflag.String(flagIdleTimeout, "", "idle timeout")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	host := viper.GetString(flagHost)
	if host != "" {
		a.Data.Network.Host = host
	}

	port := viper.GetInt(flagPort)
	if port != 0 {
		a.Data.Network.Port = port
	}

	maxConn := viper.GetInt(flagMaxConnections)
	if maxConn != 0 {
		a.Data.Network.MaxConnections = maxConn
	}

	maxMesSize := viper.GetString(flagMaxMessageSize)
	if maxMesSize != "" {
		a.Data.Network.MaxMessageSize = maxMesSize
	}

	idleTimeout := viper.GetString(flagIdleTimeout)
	if idleTimeout != "" {
		a.Data.Network.IdleTimeout = idleTimeout
	}
}

func (a *AppConfig) overideLogging() {
	pflag.String(flagLogLevel, "", "logging level")
	pflag.String(flagLogOutputDir, "", "logs dir")

	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	level := viper.GetString(flagLogLevel)
	if level != "" {
		a.Data.Logging.Level = level
	}

	dir := viper.GetString(flagLogOutputDir)
	if dir != "" {
		a.Data.Logging.OutputDir = dir
	}
}
