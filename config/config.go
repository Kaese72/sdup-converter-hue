package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"
)

// Config is the top level config structure
type Config struct {
	Adapter struct {
		Hue struct {
			Host   string `mapstructure:"host"`
			APIKey string `mapstructure:"api-key"`
			IgnoreTLSErrors bool `mapstructure:"ignore-tls-errors"`
		} `mapstructure:"hue"`
	} `mapstructure:"adapter"`
}

func ReadConfig() (Config, error) {
	myVip := viper.New()
	// Set replaces to allow keys like "database.mongodb.connection-string"
	// WARNING. Overriding any of these may hav unintended consequences.
	myVip.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// # Hue server configuration
	// Hostname or IP to bridge (no scheme), eg. 192.168.1.10
	// ADAPTER_HUE_HOST
	myVip.BindEnv("adapter.hue.host")
	// Whether to ignore TLS errors when connecting to Hue bridge
	// ADAPTER_HUE_IGNORE_TLS_ERRORS
	myVip.SetDefault("adapter.hue.ignore-tls-errors", true)
	myVip.BindEnv("adapter.hue.ignore-tls-errors")
	// preconfigured api key string
	// ADAPTER_HUE_API_KEY
	myVip.BindEnv("adapter.hue.api-key")

	var conf Config
	err := myVip.Unmarshal(&conf)
	if err != nil {
		return Config{}, err
	}
	if conf.Adapter.Hue.APIKey == "" {
		return Config{}, errors.New("must provide adapter.hue.api-key")
	}
	if conf.Adapter.Hue.Host == "" {
		return Config{}, errors.New("must provide adapter.hue.host")
	}
	return conf, nil
}
