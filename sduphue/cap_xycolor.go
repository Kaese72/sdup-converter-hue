package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/amimof/huego"
	"github.com/mitchellh/mapstructure"
)

func init() {
	capRegistry[CapabilitySetColorXY] = TriggerSetXYColor
	gCapRegistry[CapabilitySetColorXY] = GTriggerSetXYColor
}

type XYColorArgs struct {
	X *float32 `mapstructure:"x"`
	Y *float32 `mapstructure:"y"`
}

func TriggerSetXYColor(id int, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs XYColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.X == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "x must be set"}
	}

	if pArgs.Y == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "y must be set"}
	}

	_, err := bridge.SetLightState(id, huego.State{On: true, Xy: []float32{*pArgs.X, *pArgs.Y}})
	return adapterErrorFromErr(err)
}

func GTriggerSetXYColor(id int, args ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs XYColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.X == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "x must be set"}
	}

	if pArgs.Y == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "y must be set"}
	}

	_, err := bridge.SetGroupState(id, huego.State{On: true, Xy: []float32{*pArgs.X, *pArgs.Y}})
	return adapterErrorFromErr(err)
}
