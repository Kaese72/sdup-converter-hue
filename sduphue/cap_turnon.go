package sduphue

import (
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

func init() {
	capRegistry[sduptemplates.CapabilityActivate] = TriggerTurnOn
	gCapRegistry[sduptemplates.CapabilityActivate] = GTriggerTurnOn
}

func TriggerTurnOn(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: true})
	return err
}

func GTriggerTurnOn(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetGroupState(id, huego.State{On: true})
	return err
}
