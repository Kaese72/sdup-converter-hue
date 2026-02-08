package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/amimof/huego"
	"github.com/mitchellh/mapstructure"
)

func init() {
	capRegistry[CapabilitySetBrightness] = CTriggerSetBrightness
	gCapRegistry[CapabilitySetBrightness] = GTriggerSetBrightness
}

type CTBrightnessArgs struct {
	Value *float32 `mapstructure:"value"`
}

func CTriggerSetBrightness(id int, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	var pArgs CTBrightnessArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Value == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "value must be set"}
	}
	// FIXME Limitations on value variable

	_, err := bridge.SetLightState(id, huego.State{On: true, Bri: uint8(*pArgs.Value)})
	return adapterErrorFromErr(err)
}

func GTriggerSetBrightness(id int, args ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	var pArgs CTBrightnessArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Value == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "value must be set"}
	}
	// FIXME Limitations on value variable

	_, err := bridge.SetGroupState(id, huego.State{On: true, Bri: uint8(*pArgs.Value)})
	return adapterErrorFromErr(err)
}
