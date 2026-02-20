package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/mitchellh/mapstructure"
	"github.com/openhue/openhue-go"
)

func init() {
	capRegistry[CapabilityDim] = CTriggerDim
	gCapRegistry[CapabilityDim] = GTriggerDim
}

type CTDimArgs struct {
	Inc *float32 `mapstructure:"inc"`
}

func CTriggerDim(target SDUPHueTarget, id string, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	var pArgs CTDimArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Inc == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "inc must be set"}
	}
	// FIXME Limitations on inc variable

	brightnessDelta := float32(absFloat32(*pArgs.Inc))
	var action *openhue.DimmingDeltaAction
	if *pArgs.Inc > 0 {
		up := openhue.DimmingDeltaActionUp
		action = &up
	} else if *pArgs.Inc < 0 {
		down := openhue.DimmingDeltaActionDown
		action = &down
	}
	err := target.home.UpdateLight(id, openhue.LightPut{DimmingDelta: &openhue.DimmingDelta{Action: action, BrightnessDelta: &brightnessDelta}})
	return adapterErrorFromErr(err)
}

func GTriggerDim(target SDUPHueTarget, id string, args ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	var pArgs CTDimArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Inc == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "inc must be set"}
	}
	// FIXME Limitations on inc variable

	brightnessDelta := float32(absFloat32(*pArgs.Inc))
	var action *openhue.DimmingDeltaAction
	if *pArgs.Inc > 0 {
		up := openhue.DimmingDeltaActionUp
		action = &up
	} else if *pArgs.Inc < 0 {
		down := openhue.DimmingDeltaActionDown
		action = &down
	}
	err := target.home.UpdateGroupedLight(id, openhue.GroupedLightPut{DimmingDelta: &openhue.DimmingDelta{Action: action, BrightnessDelta: &brightnessDelta}})
	return adapterErrorFromErr(err)
}

func absFloat32(value float32) float32 {
	if value < 0 {
		return -value
	}
	return value
}
