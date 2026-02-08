package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/amimof/huego"
	"github.com/mitchellh/mapstructure"
)

func init() {
	capRegistry[CapabilityDim] = CTriggerDim
	gCapRegistry[CapabilityDim] = GTriggerDim
}

type CTDimArgs struct {
	Inc *float32 `mapstructure:"inc"`
}

func CTriggerDim(id int, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	var pArgs CTDimArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Inc == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "inc must be set"}
	}
	// FIXME Limitations on inc variable

	_, err := bridge.SetLightState(id, huego.State{On: true, BriInc: int(*pArgs.Inc)})
	return adapterErrorFromErr(err)
}

func GTriggerDim(id int, args ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	var pArgs CTDimArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Inc == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "inc must be set"}
	}
	// FIXME Limitations on inc variable

	_, err := bridge.SetGroupState(id, huego.State{On: true, BriInc: int(*pArgs.Inc)})
	return adapterErrorFromErr(err)
}
