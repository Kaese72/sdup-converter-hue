package config

import (
	"github.com/Kaese72/sdup-lib/httpsdup"
)

type StoreEnrollmentConfig struct {
	StoreURL  string `json:"store-url"`
	EnrollURL string `json:"enroll-url"`
}

func (config *StoreEnrollmentConfig) ShouldEnroll() bool {
	return len(config.StoreURL) != 0 && len(config.EnrollURL) != 0
}

func (config *StoreEnrollmentConfig) PopulateExample() {
	config.EnrollURL = "127.0.0.1:8086"
	config.StoreURL = "127.0.0.1:8080"
}

//Config is the top level config structure
type Config struct {
	DebugLogging      *bool                 `json:"debug-logging"`
	Hue               HueConfig             `json:"hue"`
	SDUP              httpsdup.Config       `json:"sdup-server"`
	EnrollDeviceStore StoreEnrollmentConfig `json:"store-enrollment"`
}

func (config *Config) PopulateExample() {
	config.Hue = HueConfig{}
	config.Hue.PopulateExample()

	config.SDUP = httpsdup.Config{}
	config.SDUP.PopulateExample()

	t := true
	config.DebugLogging = &t
	config.EnrollDeviceStore.PopulateExample()
}

//Validate checks whether all fields are appropriately set
func (conf *Config) Validate() error {
	if err := conf.SDUP.Validate(); err != nil {
		return err
	}

	if err := conf.Hue.Validate(); err != nil {
		return err
	}

	// Ignore logging since there is not much to validate
	return nil
}
