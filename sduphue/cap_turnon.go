package sduphue

import (
	devicestoretemplates "github.com/Kaese72/device-store/rest/models"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[CapabilityActivate] = TriggerTurnOn
	gCapRegistry[CapabilityActivate] = GTriggerTurnOn
}

func TriggerTurnOn(id int, _ devicestoretemplates.DeviceCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: true})
	return err
}

func GTriggerTurnOn(id int, _ devicestoretemplates.GroupCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: true})
	return err
}
