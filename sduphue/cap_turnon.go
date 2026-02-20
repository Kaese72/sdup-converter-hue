package sduphue

import (
	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/openhue/openhue-go"
)

func init() {
	capRegistry[CapabilityActivate] = TriggerTurnOn
	gCapRegistry[CapabilityActivate] = GTriggerTurnOn
}

func TriggerTurnOn(target SDUPHueTarget, id string, _ ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	on := true
	err := target.home.UpdateLight(id, openhue.LightPut{On: &openhue.On{On: &on}})
	return adapterErrorFromErr(err)
}

func GTriggerTurnOn(target SDUPHueTarget, id string, _ ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	if target.home == nil {
		return &adapter.AdapterError{Code: 500, Message: "home not initialized"}
	}
	on := true
	err := target.home.UpdateGroupedLight(id, openhue.GroupedLightPut{On: &openhue.On{On: &on}})
	return adapterErrorFromErr(err)
}
