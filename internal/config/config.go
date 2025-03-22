package config

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/viper"

	"kdb/internal/utils"
)

type AppConfig struct {
	Data *Config
}

var errInvalidLogger error = errors.New("invalid logger")

func NewConfig() (*AppConfig, error) {
	return &AppConfig{}, nil
}

func (a *AppConfig) Init(ctx context.Context, env string) error {
	viper.SetConfigName(fmt.Sprintf("%s.config", env))
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		wErr := fmt.Errorf("viper trying read in config: %w", err)
		return wErr
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		wErr := fmt.Errorf("viper trying unmarshall config: %w", err)
		return wErr
	}

	a.Data = &config

	a.overrideByFlags()

	return nil
}

func (a *AppConfig) String() string {
	return utils.StructToString(a)
}
