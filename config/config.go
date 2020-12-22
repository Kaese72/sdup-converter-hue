package config

import "errors"

//HueConfig contains configuration related to the Hue bridge
type HueConfig struct {
	URL    string `json:"url"`
	APIKey string `json:"api-key"`
}

//SDUPConfig contains SDUP related configs
type SDUPConfig struct {
	ListenAddress string `json:"ip"`
	ListenPort    int    `json:"port"`
}

//Config is the top level config structure
type Config struct {
	Hue  HueConfig  `json:"hue"`
	SDUP SDUPConfig `json:"sdup"`
}

//Validate checks whether all fields are appropriately set
func (conf *Config) Validate() error {
	if conf.Hue.URL == "" {
		return errors.New("Must set hue.url")
	}

	if conf.Hue.APIKey == "" {
		return errors.New("Must set hue.api-key")
	}

	if conf.SDUP.ListenAddress == "" {
		//Simply default to localhost, why noy
		conf.SDUP.ListenAddress = "localhost"
	}

	if conf.SDUP.ListenPort == 0 {
		//Simply default to 8080, why not.
		conf.SDUP.ListenPort = 8080
	}

	return nil
}

//NewExampleConfig generates an example config
func NewExampleConfig() *Config {
	return &Config{
		Hue: HueConfig{
			URL:    "http://localhost:8081/",
			APIKey: "some api key string here",
		},
		SDUP: SDUPConfig{
			ListenAddress: "localhost",
			ListenPort:    8080,
		},
	}
}
