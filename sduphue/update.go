package sduphue

import (
	"fmt"
	"strings"

	"github.com/Kaese72/sdup-converter-hue/log"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

func (target *SDUPHueTarget) getAllDevices() (specs []HueDevice, err error) {
	hueLights, err := bridge.GetLights()
	if err != nil {
		return
	}
	for _, light := range hueLights {
		hueLight := createLightDevice(light)
		specs = append(specs, hueLight)
	}
	return
}

func (target *SDUPHueTarget) UpdateAllDevices() error {
	hueDevices, err := target.getAllDevices()
	if err != nil {
		return err
	}
	for _, newDevice := range hueDevices {
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

//Not lowercase according to docs, but it seems like strings are not properly speced
const (
	OnOffLight         = "on/off light"
	DimmableLight      = "dimmable light"
	ColorTempLight     = "color temperature light"
	ColorLight         = "color light"
	ExtendedColorLight = "extended color light"
)

//xyColorLights contains the different lights that support xy color control
//FIXME I might be able to tell if the light has uspport of xy color mode by checking the presence of "xy" in the state retrieved form the bridge
var xyColorLights = map[string]bool{
	ColorLight:         true,
	ExtendedColorLight: true,
}

var ctColorLights = map[string]bool{
	ColorTempLight:     true,
	ColorLight:         true,
	ExtendedColorLight: true,
}

//var ctColorLights = map[string]bool{}

func createLightDevice(light huego.Light) HueDevice {
	device := HueDevice{
		ID: HueDeviceID{Type: "light", Index: light.ID},
		Attributes: map[sduptemplates.AttributeKey]HueAttribute{
			sduptemplates.AttributeActive: {
				State: sduptemplates.AttributeState{Boolean: &light.State.On},
			},
		},
		Capabilities: map[sduptemplates.CapabilityKey]HueCapability{
			sduptemplates.CapabilityActivate:   TurnOnLight{},
			sduptemplates.CapabilityDeactivate: TurnOffLight{},
		},
	}
	// #########################
	// # Description of device #
	// #########################
	device.Attributes[sduptemplates.AttributeDescription] = HueAttribute{
		State: sduptemplates.AttributeState{
			Text: &light.Name,
		},
	}
	// ############
	// # UniqueID #
	// ############
	device.Attributes[sduptemplates.AttributeUniqueID] = HueAttribute{
		State: sduptemplates.AttributeState{
			Text: &light.UniqueID,
		},
	}

	// #################
	// # XY Color Mode #
	// #################
	if xyColorLights[strings.ToLower(light.Type)] {
		if len(light.State.Xy) == 2 {
			// If the XY is set, use it as an attribute
			device.Attributes[sduptemplates.AttributeColorXY] = HueAttribute{
				State: sduptemplates.AttributeState{
					KeyVal: &sduptemplates.KeyValContainer{
						"x": light.State.Xy[0],
						"y": light.State.Xy[1],
					},
				},
			}

		} else {
			if len(light.State.Xy) != 0 {
				log.Log(log.Error, "Invalid length on XY array, assuming nil values", nil)
			}
			//Attach attribute with nil color xy coordinates
			//This happens when the colormode is not set to xy but rather eg. ct
			device.Attributes[sduptemplates.AttributeColorXY] = HueAttribute{
				State: sduptemplates.AttributeState{
					KeyVal: &sduptemplates.KeyValContainer{
						"x": nil,
						"y": nil,
					},
				},
			}
		}
		//Attach capability to change color with xy coordinates
		device.Capabilities[sduptemplates.CapabilitySetColorXY] = XYColor{}
	}
	// #################
	// # CT Color Mode #
	// #################
	if ctColorLights[strings.ToLower(light.Type)] {
		//Attach color temperature attrbiute
		ct := float32(light.State.Ct)
		device.Attributes[sduptemplates.AttributeColorTemp] = HueAttribute{
			State: sduptemplates.AttributeState{
				Numeric: &ct,
			},
		}
		//Attach capability to change color temperature
		device.Capabilities[sduptemplates.CapabilitySetColorTemp] = CTColor{}
	}

	return device
}
