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
			URL    string `mapstructure:"url"`
			APIKey string `mapstructure:"api-key"`
		} `mapstructure:"hue"`
	} `mapstructure:"adapter"`
}

func ReadConfig() (Config, error) {
	myVip := viper.New()
	// Set replaces to allow keys like "database.mongodb.connection-string"
	// WARNING. Overriding any of these may hav unintended consequences.
	myVip.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	// # Hue server configuration
	// URL to server, eg. http://
	// ADAPTER_HUE_URL
	myVip.BindEnv("adapter.hue.url")
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
	if conf.Adapter.Hue.URL == "" {
		return Config{}, errors.New("must provide adapter.hue.url")
	}
	return conf, nil
}
