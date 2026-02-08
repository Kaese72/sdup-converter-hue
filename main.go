package main

import (
	"os"

	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-converter-hue/config"
	"github.com/Kaese72/sdup-converter-hue/sduphue"
	"github.com/Kaese72/sdup-lib/adapter"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		log.Error("Failed to read config", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
	hueTarget := sduphue.InitSDUPHueTarget(conf.Adapter.Hue.URL, conf.Adapter.Hue.APIKey)
	// Make sure we fulfil the trigger interfaces,
	var _ adapter.DeviceTriggerCapability = hueTarget
	var _ adapter.GroupTriggerCapability = hueTarget
	err = adapter.StartAdapter(hueTarget)
	if err != nil {
		log.Error("Failed to initiate device store updater", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}
