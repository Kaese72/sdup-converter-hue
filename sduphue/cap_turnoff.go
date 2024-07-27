package sduphue

import (
	devicestoretemplates "github.com/Kaese72/device-store/rest/models"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[CapabilityDeactivate] = TriggerTurnOff
	gCapRegistry[CapabilityDeactivate] = GTriggerTurnOff
}

func TriggerTurnOff(id int, _ devicestoretemplates.DeviceCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: false})
	return err
}

func GTriggerTurnOff(id int, _ devicestoretemplates.GroupCapabilityArgs) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: false})
	return err
}
