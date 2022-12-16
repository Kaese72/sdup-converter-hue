package sduphue

import (
	"github.com/Kaese72/sdup-lib/devicestoretemplates"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[CapabilityDeactivate] = TriggerTurnOff
	gCapRegistry[CapabilityDeactivate] = GTriggerTurnOff
}

func TriggerTurnOff(id int, _ devicestoretemplates.CapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: false})
	return err
}

func GTriggerTurnOff(id int, _ devicestoretemplates.CapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: false})
	return err
}
