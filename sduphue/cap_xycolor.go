package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/mitchellh/mapstructure"
	"github.com/openhue/openhue-go"
)

func init() {
	capRegistry[CapabilitySetColorXY] = TriggerSetXYColor
	gCapRegistry[CapabilitySetColorXY] = GTriggerSetXYColor
}

type XYColorArgs struct {
	X *float32 `mapstructure:"x"`
	Y *float32 `mapstructure:"y"`
}

func TriggerSetXYColor(target SDUPHueTarget, id string, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
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

	color := openhue.Color{Xy: &openhue.GamutPosition{X: pArgs.X, Y: pArgs.Y}}
	err := target.home.UpdateLight(id, openhue.LightPut{Color: &color})
	return adapterErrorFromErr(err)
}

func GTriggerSetXYColor(target SDUPHueTarget, id string, args ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
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

	color := openhue.Color{Xy: &openhue.GamutPosition{X: pArgs.X, Y: pArgs.Y}}
	err := target.home.UpdateGroupedLight(id, openhue.GroupedLightPut{Color: &color})
	return adapterErrorFromErr(err)
}
