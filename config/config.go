package config

import (
	"errors"
	"fmt"

	"github.com/Kaese72/sdup-lib/deviceupdates"
)

type ListenConfig struct {
	ListenAddress string `mapstructure:"address"`
	ListenPort    int    `mapstructure:"port"`
}

func (conf ListenConfig) URL() string {
	return fmt.Sprintf("http://%s:%d", conf.ListenAddress, conf.ListenPort)
}

func (config ListenConfig) Validate() error {
	if config.ListenAddress == "" {
		return errors.New("empty listen address")
	}

	if config.ListenPort < 0 || config.ListenPort > 65665 {
		return fmt.Errorf("invalid port number, %d", config.ListenPort)
	}
	//FIXME Validate StoreEnrollmentConfig
	return nil
}

// Config is the top level config structure
type Config struct {
	DebugLogging      *bool                               `mapstructure:"debug-logging"`
	Hue               HueConfig                           `mapstructure:"hue"`
	HTTPServer        ListenConfig                        `mapstructure:"http-server"`
	EnrollDeviceStore deviceupdates.StoreEnrollmentConfig `mapstructure:"enroll"`
}

// Validate checks whether all fields are appropriately set
func (conf *Config) Validate() error {
	if err := conf.HTTPServer.Validate(); err != nil {
		return err
	}

	if err := conf.Hue.Validate(); err != nil {
		return err
	}

	// Ignore logging since there is not much to validate
	return nil
}
