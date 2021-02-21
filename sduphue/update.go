package sduphue

import (
	"errors"
	"fmt"

	"github.com/Kaese72/sdup-hue/log"
	"github.com/Kaese72/sdup-hue/sduptemplates"
	"github.com/amimof/huego"
)

func (target *SDUPHueTarget) UpdateAllDevices() error {
	hueLights, err := bridge.GetLights()
	if err != nil {
		return errors.New("Failed to enumerate devices (lights) on Hue bridge")
	}
	for _, light := range hueLights {
		newDevice := createLightDevice(light)
		deviceUpdates := sduptemplates.DeviceUpdate{
			ID:             newDevice.ID.SDUPEncode(),
			AttributesDiff: sduptemplates.AttributeStateMap{},
		}

		if oldLight, ok := target.devices[newDevice.ID.SDUPEncode()]; ok {
			log.Log(log.Debug, "Updating light", map[string]string{"light": fmt.Sprint(newDevice.ID.Index)})
			//Previously known device
			for key, newAttrVal := range newDevice.Attributes {
				if oldAttrVal, ok := oldLight.Attributes[key]; ok {
					if !oldAttrVal.State.Equivalent(newAttrVal.State) {
						deviceUpdates.AttributesDiff[key] = newAttrVal.State
					}
					// TODO Might encounter diff in capabilities here
				} else {
					// TODO Added attributes?
					log.Log(log.Error, "Found more capabilities?", nil)
				}
			}
			// TODO Lost attributes?

		} else {
			log.Log(log.Info, "Adding light", map[string]string{"light": fmt.Sprint(newDevice.ID.Index)})

			// TODO capability updates
			for id, attribute := range newDevice.Attributes {
				deviceUpdates.AttributesDiff[id] = attribute.State
			}
		}
		// TODO What about lost devices?
		// Just update all detected lights
		target.devices[newDevice.ID.SDUPEncode()] = newDevice
		if len(deviceUpdates.AttributesDiff) > 0 {
			target.updateChan <- deviceUpdates
		}
	}

	return nil
}

// ###########
// # TurnOn #
// ###########
type TurnOnLight struct{}

func (cap TurnOnLight) Trigger(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: true})
	return err
}

func (cap TurnOnLight) Spec() sduptemplates.CapabilitySpec {
	return sduptemplates.CapabilitySpec{}
}

// ###########
// # TurnOff #
// ###########
type TurnOffLight struct{}

func (cap TurnOffLight) Trigger(id int, _ sduptemplates.CapabilityArgument) error {
	//FIXME Is there anythig interesting in the huego response ?
	_, err := bridge.SetLightState(id, huego.State{On: false})
	return err
}

func (cap TurnOffLight) Spec() sduptemplates.CapabilitySpec {
	return sduptemplates.CapabilitySpec{}
}

func createLightDevice(light huego.Light) HueDevice {
	return HueDevice{
		ID: HueDeviceID{Type: "light", Index: light.ID},
		Attributes: map[sduptemplates.AttributeKey]HueAttribute{
			sduptemplates.AttributeActive: {
				State: sduptemplates.AttributeState{BooleanState: &light.State.On},
			},
		},
		Capabilities: map[sduptemplates.CapabilityKey]HueCapability{
			sduptemplates.CapabilityActivate:   TurnOnLight{},
			sduptemplates.CapabilityDeactivate: TurnOffLight{},
		},
	}
}
