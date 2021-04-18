package sduphue

import (
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

type TurnOnLight struct{}

func (cap TurnOnLight) Trigger(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: true})
	return err
}

func (cap TurnOnLight) Spec() sduptemplates.CapabilitySpec {
	return sduptemplates.CapabilitySpec{}
}
