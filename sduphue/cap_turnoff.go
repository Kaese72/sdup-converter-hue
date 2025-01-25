package sduphue

import (
	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[CapabilityDeactivate] = TriggerTurnOff
	gCapRegistry[CapabilityDeactivate] = GTriggerTurnOff
}

func TriggerTurnOff(id int, _ ingestmodels.DeviceCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: false})
	return err
}

func GTriggerTurnOff(id int, _ ingestmodels.GroupCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: false})
	return err
}
