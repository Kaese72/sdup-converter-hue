package sduphue

import (
	"strconv"
	"strings"

	log "github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

func (target *SDUPHueTarget) getAllDevices() (specs []sduptemplates.DeviceSpec, err error) {
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

func (target *SDUPHueTarget) getAllGroups() (specs []sduptemplates.DeviceGroupSpec, err error) {
	hueGroups, err := bridge.GetGroups()
	if err != nil {
		return
	}

	for _, group := range hueGroups {
		hueGroup := createDeviceGroup(group)
		specs = append(specs, hueGroup)
	}
	return
}

func (target *SDUPHueTarget) UpdateEverything() error {
	hueDevices, err := target.getAllDevices()
	if err != nil {
		return err
	}
	for _, newDevice := range hueDevices {
		deviceUpdate := sduptemplates.UpdateFromDeviceUpdate(newDevice.SpecToInitialUpdate())

		// TODO What about lost devices?
		// Just update all detected lights
		target.updateChan <- deviceUpdate
	}

	hueGroups, err := target.getAllGroups()
	if err != nil {
		return err
	}
	for _, newGroup := range hueGroups {
		groupUpdate := sduptemplates.UpdateFromDeviceGroupUpdate(newGroup.SpecToInitialUpdate())

		target.updateChan <- groupUpdate
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

func createLightDevice(light huego.Light) sduptemplates.DeviceSpec {
	device := sduptemplates.DeviceSpec{
		ID: HueDeviceID{Type: LIGHT, Index: light.ID}.SDUPEncode(),
		Attributes: sduptemplates.AttributeStateMap{
			sduptemplates.AttributeActive: {
				Boolean: &light.State.On,
			},
		},
		Capabilities: sduptemplates.CapabilitySpecMap{
			sduptemplates.CapabilityActivate:   sduptemplates.CapabilitySpec{},
			sduptemplates.CapabilityDeactivate: sduptemplates.CapabilitySpec{},
		},
	}
	// #########################
	// # Description of device #
	// #########################
	device.Attributes[sduptemplates.AttributeDescription] = sduptemplates.AttributeState{
		Text: &light.Name,
	}
	// ############
	// # UniqueID #
	// ############
	device.Attributes[sduptemplates.AttributeUniqueID] = sduptemplates.AttributeState{
		Text: &light.UniqueID,
	}

	// #################
	// # XY Color Mode #
	// #################
	if xyColorLights[strings.ToLower(light.Type)] {
		if len(light.State.Xy) == 2 {
			// If the XY is set, use it as an attribute
			device.Attributes[sduptemplates.AttributeColorX] = sduptemplates.AttributeState{Numeric: &light.State.Xy[0]}
			device.Attributes[sduptemplates.AttributeColorY] = sduptemplates.AttributeState{Numeric: &light.State.Xy[1]}

		} else {
			if len(light.State.Xy) != 0 {
				log.Error("Invalid length on XY array, assuming nil values")
			}
			//Attach attribute with nil color xy coordinates
			//This happens when the colormode is not set to xy but rather eg. ct
			device.Attributes[sduptemplates.AttributeColorX] = sduptemplates.AttributeState{}
			device.Attributes[sduptemplates.AttributeColorY] = sduptemplates.AttributeState{}
		}
		//Attach capability to change color with xy coordinates
		device.Capabilities[sduptemplates.CapabilitySetColorXY] = sduptemplates.CapabilitySpec{}
	}
	// #################
	// # CT Color Mode #
	// #################
	if ctColorLights[strings.ToLower(light.Type)] {
		//Attach color temperature attrbiute
		ct := float32(light.State.Ct)
		device.Attributes[sduptemplates.AttributeColorTemp] = sduptemplates.AttributeState{
			Numeric: &ct,
		}
		//Attach capability to change color temperature
		device.Capabilities[sduptemplates.CapabilitySetColorTemp] = sduptemplates.CapabilitySpec{}
	}

	return device
}

func createDeviceGroup(group huego.Group) sduptemplates.DeviceGroupSpec {
	g := sduptemplates.DeviceGroupSpec{
		ID:        sduptemplates.DeviceGroupID(strconv.Itoa(group.ID)),
		Name:      group.Name,
		DeviceIDs: []sduptemplates.DeviceID{},
	}
	for _, lightId := range group.Lights {
		lid, err := strconv.Atoi(lightId)
		if err != nil {
			// FIXME Do not panic; ignore or something
			panic(err)
		}
		g.DeviceIDs = append(g.DeviceIDs, HueDeviceID{Type: LIGHT, Index: lid}.SDUPEncode())
	}
	return g
}
