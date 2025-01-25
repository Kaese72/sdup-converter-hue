package sduphue

import (
	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[CapabilityActivate] = TriggerTurnOn
	gCapRegistry[CapabilityActivate] = GTriggerTurnOn
}

func TriggerTurnOn(id int, _ ingestmodels.DeviceCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: true})
	return err
}

func GTriggerTurnOn(id int, _ ingestmodels.GroupCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: true})
	return err
}
