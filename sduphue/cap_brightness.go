package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/mitchellh/mapstructure"
	"github.com/openhue/openhue-go"
)

func init() {
	capRegistry[CapabilitySetBrightness] = CTriggerSetBrightness
	gCapRegistry[CapabilitySetBrightness] = GTriggerSetBrightness
}

type CTBrightnessArgs struct {
	Value *float32 `mapstructure:"value"`
}

func CTriggerSetBrightness(target SDUPHueTarget, id string, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	var pArgs CTBrightnessArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Value == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "value must be set"}
	}
	// FIXME Limitations on value variable

	brightness := openhue.Brightness(*pArgs.Value)
	on := true
	err := target.home.UpdateLight(id, openhue.LightPut{On: &openhue.On{On: &on}, Dimming: &openhue.Dimming{Brightness: &brightness}})
	return adapterErrorFromErr(err)
}

func GTriggerSetBrightness(target SDUPHueTarget, id string, args ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	var pArgs CTBrightnessArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Value == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "value must be set"}
	}
	// FIXME Limitations on value variable

	brightness := openhue.Brightness(*pArgs.Value)
	on := true
	err := target.home.UpdateGroupedLight(id, openhue.GroupedLightPut{On: &openhue.On{On: &on}, Dimming: &openhue.Dimming{Brightness: &brightness}})
	return adapterErrorFromErr(err)
}
