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
	hueTarget, err := sduphue.InitSDUPHueTarget(conf.Adapter.Hue.Host, conf.Adapter.Hue.APIKey, conf.Adapter.Hue.IgnoreTLSErrors)
	if err != nil {
		log.Error("Failed to initialize Hue target", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
	// Make sure we fulfil the trigger interfaces,
	var _ adapter.DeviceTriggerCapability = hueTarget
	var _ adapter.GroupTriggerCapability = hueTarget
	err = adapter.StartAdapter(hueTarget)
	if err != nil {
		log.Error("Failed to initiate device store updater", map[string]interface{}{"error": err.Error()})
		os.Exit(1)
	}
}
