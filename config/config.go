package config

import (
	"github.com/Kaese72/sdup-lib/httpsdup"
)

//Config is the top level config structure
type Config struct {
	DebugLogging *bool           `json:"debug-logging"`
	Hue          HueConfig       `json:"hue"`
	SDUP         httpsdup.Config `json:"sdup-server"`
}

func (config *Config) PopulateExample() {
	config.Hue = HueConfig{}
	config.Hue.PopulateExample()

	config.SDUP = httpsdup.Config{}
	config.SDUP.PopulateExample()

	t := true
	config.DebugLogging = &t
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
