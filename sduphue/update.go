package sduphue

import (
	"fmt"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/openhue/openhue-go"
)

func (target SDUPHueTarget) getAllDevices() ([]ingestmodels.IngestDevice, error) {
	specs := []ingestmodels.IngestDevice{}
	if target.home == nil {
		return nil, fmt.Errorf("home not initialized")
	}
	hueLights, err := target.home.GetLights()
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
	if target.home == nil {
		err = fmt.Errorf("home not initialized")
		return
	}
	groupNameByID := target.getGroupedLightNameMap()
	groupedLights, err := target.home.GetGroupedLights()
	if err != nil {
		return
	}

	for id, group := range groupedLights {
		name := groupNameByID[id]
		if name == "" {
			name = "Group"
		}
		hueGroup := createDeviceGroup(group, name)
		specs = append(specs, hueGroup)
	}
	return
}

const (
	//AttributeActive represents whether the device is currently on or off
	AttributeActive string = "active"
	//AttributeBrightness represents the brightness of the device
	AttributeBrightness string = "brightness"
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

func createLightDevice(light openhue.LightGet) ingestmodels.IngestDevice {
	lightID := safeString(light.Id)
	device := ingestmodels.IngestDevice{
		BridgeIdentifier: lightID,
		Attributes:       []ingestmodels.IngestAttribute{},
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
	if light.On != nil && light.On.On != nil {
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name:    AttributeActive,
			Boolean: light.On.On,
		})
	}
	if light.Dimming != nil && light.Dimming.Brightness != nil {
		brightness := float32(*light.Dimming.Brightness)
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name:    AttributeBrightness,
			Numeric: &brightness,
		})
	}
	if light.Metadata != nil && light.Metadata.Name != nil {
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name: AttributeDescription,
			Text: light.Metadata.Name,
		})
	}
	if lightID != "" {
		device.Attributes = append(device.Attributes, ingestmodels.IngestAttribute{
			Name: AttributeUniqueID,
			Text: &lightID,
		})
	}
	// #################
	// # XY Color Mode #
	// #################
	if light.Color != nil && light.Color.Xy != nil {
		if light.Color.Xy.X != nil && light.Color.Xy.Y != nil {
			x := *light.Color.Xy.X
			y := *light.Color.Xy.Y
			device.Attributes = append(
				device.Attributes,
				ingestmodels.IngestAttribute{
					Name:    AttributeColorX,
					Numeric: &x,
				},
			)
			device.Attributes = append(
				device.Attributes,
				ingestmodels.IngestAttribute{
					Name:    AttributeColorY,
					Numeric: &y,
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
	if light.ColorTemperature != nil && light.ColorTemperature.Mirek != nil {
		ct := float32(*light.ColorTemperature.Mirek)
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

func createDeviceGroup(group openhue.GroupedLightGet, name string) ingestmodels.IngestGroup {
	g := ingestmodels.IngestGroup{
		BridgeIdentifier: safeString(group.Id),
		Name:             name,
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

func safeString(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func (target *SDUPHueTarget) getGroupedLightNameMap() map[string]string {
	result := map[string]string{}
	rooms, err := target.home.GetRooms()
	if err == nil {
		for _, room := range rooms {
			name := "Room"
			if room.Metadata != nil && room.Metadata.Name != nil {
				name = *room.Metadata.Name
			}
			for serviceID, serviceType := range room.GetServices() {
				if serviceType == openhue.ResourceIdentifierRtypeGroupedLight {
					result[serviceID] = name
				}
			}
		}
	}

	bridgeHome, err := target.home.GetBridgeHome()
	if err == nil && bridgeHome != nil && bridgeHome.Services != nil {
		for _, service := range *bridgeHome.Services {
			if service.Rtype == nil || service.Rid == nil {
				continue
			}
			if *service.Rtype == openhue.ResourceIdentifierRtypeGroupedLight {
				if _, ok := result[*service.Rid]; !ok {
					result[*service.Rid] = "Bridge Home"
				}
			}
		}
	}

	return result
}
