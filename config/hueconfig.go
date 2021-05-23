package config

import "errors"

//HueConfig contains configuration related to the Hue bridge
type HueConfig struct {
	URL    string `json:"url"`
	APIKey string `json:"api-key"`
}

func (config *HueConfig) PopulateExample() {
	config.URL = "scheme://host:port/"
	config.APIKey = "Hue API key"
}

//Validate checks whether all fields are appropriately set
func (conf *HueConfig) Validate() error {
	if conf.URL == "" {
		return errors.New("must set hue.url")
	}

	if conf.APIKey == "" {
		return errors.New("must set hue.api-key")
	}
	return nil
}
