package sduphue

import (
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[sduptemplates.CapabilityDeactivate] = TriggerTurnOff
	gCapRegistry[sduptemplates.CapabilityDeactivate] = GTriggerTurnOff
}

func TriggerTurnOff(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: false})
	return err
}

func GTriggerTurnOff(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: false})
	return err
}
