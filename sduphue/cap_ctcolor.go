package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/mitchellh/mapstructure"
	"github.com/openhue/openhue-go"
)

func init() {
	capRegistry[CapabilitySetColorTemp] = TriggerSetCTColor
	gCapRegistry[CapabilitySetColorTemp] = GTriggerSetCTColor
}

type CTColorArgs struct {
	Ct *float32 `mapstructure:"ct"`
}

func TriggerSetCTColor(target SDUPHueTarget, id string, args ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs CTColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Ct == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "ct must be set"}
	}

	mirek := openhue.Mirek(int(*pArgs.Ct))
	err := target.home.UpdateLight(id, openhue.LightPut{ColorTemperature: &openhue.ColorTemperature{Mirek: &mirek}})
	return adapterErrorFromErr(err)
}

func GTriggerSetCTColor(target SDUPHueTarget, id string, args ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	//FIXME Is there anythig interesting in the huego response ?
	//FIXME Limitations on x and y variables
	var pArgs CTColorArgs
	if err := mapstructure.Decode(args, &pArgs); err != nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if pArgs.Ct == nil {
		return &adapter.AdapterError{Code: http.StatusBadRequest, Message: "ct must be set"}
	}

	mirek := openhue.Mirek(int(*pArgs.Ct))
	err := target.home.UpdateGroupedLight(id, openhue.GroupedLightPut{ColorTemperature: &openhue.ColorTemperature{Mirek: &mirek}})
	return adapterErrorFromErr(err)
}
