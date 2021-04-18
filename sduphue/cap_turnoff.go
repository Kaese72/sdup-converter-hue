package sduphue

import (
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

type TurnOffLight struct{}

func (cap TurnOffLight) Trigger(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: false})
	return err
}

func (cap TurnOffLight) Spec() sduptemplates.CapabilitySpec {
	return sduptemplates.CapabilitySpec{}
}
