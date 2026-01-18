package sduphue

import (
	"strconv"
	"strings"

	"github.com/Kaese72/device-store/ingestmodels"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/amimof/huego"
)

func (target SDUPHueTarget) getAllDevices() ([]ingestmodels.IngestDevice, error) {
	specs := []ingestmodels.IngestDevice{}
	hueLights, err := bridge.GetLights()
	if err != nil {
		return nil, err
	}

	for _, light := range hueLights {
		hueLight := createLightDevice(light)
		specs = append(specs, hueLight)
	}
	return specs, nil
}

func (target *SDUPHueTarget) getAllGroups() (specs []ingestmodels.IngestGroup, err error) {
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

// Not lowercase according to docs, but it seems like strings are not properly speced
const (
	OnOffLight         = "on/off light"
	DimmableLight      = "dimmable light"
	ColorTempLight     = "color temperature light"
	ColorLight         = "color light"
	ExtendedColorLight = "extended color light"
)

// xyColorLights contains the different lights that support xy color control
// FIXME I might be able to tell if the light has uspport of xy color mode by checking the presence of "xy" in the state retrieved form the bridge
var xyColorLights = map[string]bool{
	ColorLight:         true,
	ExtendedColorLight: true,
}

var ctColorLights = map[string]bool{
	ColorTempLight:     true,
	ColorLight:         true,
	ExtendedColorLight: true,
}

const (
	//AttributeActive represents whether the device is currently on or off
	AttributeActive string = "active"
	//AttributeColorXY represents the primary color of the device, represented by xy coordinates
	AttributeColorX string = "colorx"
	AttributeColorY string = "colory"
	//AttributeColorTemp represents the primary color of the device, represented by xy coordinates
	AttributeColorTemp string = "colorct"
	//AttributeDescription is a readable description, presentable to a user. Should not be used to identify the device
	AttributeDescription string = "description"
	//AttributeUniqueID globally identifes a device across bridges. eg. MAC addresses
	AttributeUniqueID string = "uniqueID"
	//AttributeGroups globally identifies group names a device is part of (Groups generally do not have a unique ID, so strings is the best we have)
	AttributeGroups string = "groups"
)

const (
	//CapabilityActivate means the associated attribute can be activated
	CapabilityActivate string = "activate"
	//CapabilityDeactivate means the associated attribute can be deactivated
	CapabilityDeactivate string = "deactivate"
	//CapabilitySetColorXY means that you can change the x and y coordinates in color mode
	CapabilitySetColorXY string = "setcolorxy"
	//CapabilitySetColorTemp means that you can change the temperature in color mode
	CapabilitySetColorTemp string = "setcolorct"
	// CapabilitySetBrightness means that you can change the brightness of the light
	CapabilitySetBrightness string = "setbrightness"
	// CapabilityDim means that you can dim the light by some value
	CapabilityDim string = "dim"
)

func createLightDevice(light huego.Light) ingestmodels.IngestDevice {
	device := ingestmodels.IngestDevice{
		BridgeIdentifier: HueDeviceID{Type: LIGHT, Index: light.ID}.SDUPEncode(),
		Attributes: []ingestmodels.IngestAttribute{
			{
				Name:    AttributeActive,
				Boolean: &light.State.On,
			},
			{
				Name: AttributeDescription,
				Text: &light.Name,
			},
			{
				Name: AttributeUniqueID,
				Text: &light.UniqueID,
			},
		},
		Capabilities: []ingestmodels.IngestDeviceCapability{
			{
				Name:          CapabilityActivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
			{
				Name:          CapabilityDeactivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
			{
				Name: CapabilitySetBrightness,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{
					{
						Name: "value",
						Numeric: &ingestmodels.IngestNumericArgumentSpec{
							Min: 1,
							Max: 100,
						},
					},
				},
			},
			{
				Name: CapabilityDim,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{
					{
						Name: "inc",
						Numeric: &ingestmodels.IngestNumericArgumentSpec{
							Min: -100,
							Max: 100,
						},
					},
				},
			},
		},
	}
	// #################
	// # XY Color Mode #
	// #################
	if xyColorLights[strings.ToLower(light.Type)] {
		if len(light.State.Xy) == 2 {
			// If the XY is set, use it as an attribute
			device.Attributes = append(
				device.Attributes,
				ingestmodels.IngestAttribute{
					Name:    AttributeColorX,
					Numeric: &light.State.Xy[0],
				},
			)
			device.Attributes = append(
				device.Attributes,
				ingestmodels.IngestAttribute{
					Name:    AttributeColorY,
					Numeric: &light.State.Xy[1],
				},
			)

		} else {
			if len(light.State.Xy) != 0 {
				log.Error("Invalid length on XY array, assuming nil values")
			}
			device.Attributes = append(
				device.Attributes,
				ingestmodels.IngestAttribute{
					Name: AttributeColorX,
				},
			)
			device.Attributes = append(
				device.Attributes,
				ingestmodels.IngestAttribute{
					Name: AttributeColorY,
				},
			)
		}
		//Attach capability to change color with xy coordinates
		device.Capabilities = append(device.Capabilities, ingestmodels.IngestDeviceCapability{
			Name: CapabilitySetColorXY,
			ArgumentSpecs: []ingestmodels.IngestArgumentSpec{
				{
					Name: "x",
					Numeric: &ingestmodels.IngestNumericArgumentSpec{
						Min: 0,
						Max: 1,
					},
				},
				{
					Name: "y",
					Numeric: &ingestmodels.IngestNumericArgumentSpec{
						Min: 0,
						Max: 1,
					},
				},
			},
		})
	}
	// #################
	// # CT Color Mode #
	// #################
	if ctColorLights[strings.ToLower(light.Type)] {
		//Attach color temperature attrbiute
		ct := float32(light.State.Ct)
		device.Attributes = append(
			device.Attributes,
			ingestmodels.IngestAttribute{
				Name:    AttributeColorTemp,
				Numeric: &ct,
			},
		)
		//Attach capability to change color temperature
		device.Capabilities = append(device.Capabilities, ingestmodels.IngestDeviceCapability{
			Name: CapabilitySetColorTemp,
			ArgumentSpecs: []ingestmodels.IngestArgumentSpec{
				{
					Name: "ct",
					Numeric: &ingestmodels.IngestNumericArgumentSpec{
						Min: 153,
						Max: 500,
					},
				},
			},
		})
	}

	return device
}

func createDeviceGroup(group huego.Group) ingestmodels.IngestGroup {
	g := ingestmodels.IngestGroup{
		BridgeIdentifier: strconv.Itoa(group.ID),
		Name:             group.Name,
		Capabilities: []ingestmodels.IngestGroupCapability{
			{
				Name:          CapabilityActivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
			{
				Name:          CapabilityDeactivate,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{},
			},
			{
				Name: CapabilitySetBrightness,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{
					{
						Name: "value",
						Numeric: &ingestmodels.IngestNumericArgumentSpec{
							Min: 1,
							Max: 100,
						},
					},
				},
			},
			{
				Name: CapabilityDim,
				ArgumentSpecs: []ingestmodels.IngestArgumentSpec{
					{
						Name: "inc",
						Numeric: &ingestmodels.IngestNumericArgumentSpec{
							Min: -100,
							Max: 100,
						},
					},
				},
			},
		},
	}
	return g
}
