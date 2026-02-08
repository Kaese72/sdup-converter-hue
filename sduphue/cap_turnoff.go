package sduphue

import (
	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[CapabilityDeactivate] = TriggerTurnOff
	gCapRegistry[CapabilityDeactivate] = GTriggerTurnOff
}

func TriggerTurnOff(id int, _ ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: false})
	return adapterErrorFromErr(err)
}

func GTriggerTurnOff(id int, _ ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: false})
	return adapterErrorFromErr(err)
}
