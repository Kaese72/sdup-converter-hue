package sduphue

import (
	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/openhue/openhue-go"
)

func init() {
	capRegistry[CapabilityDeactivate] = TriggerTurnOff
	gCapRegistry[CapabilityDeactivate] = GTriggerTurnOff
}

func TriggerTurnOff(target SDUPHueTarget, id string, _ ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	off := false
	err := target.home.UpdateLight(id, openhue.LightPut{On: &openhue.On{On: &off}})
	return adapterErrorFromErr(err)
}

func GTriggerTurnOff(target SDUPHueTarget, id string, _ ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	off := false
	err := target.home.UpdateGroupedLight(id, openhue.GroupedLightPut{On: &openhue.On{On: &off}})
	return adapterErrorFromErr(err)
}
