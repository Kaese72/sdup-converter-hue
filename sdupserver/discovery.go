package sdupserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Kaese72/sdup-hue/attributes"
	"github.com/Kaese72/sdup-hue/capabilities"
	"github.com/Kaese72/sdup-hue/devicecontainer"
	"github.com/Kaese72/sdup-hue/log"
	"github.com/amimof/huego"
)

//Discovery handles discovery requests from a client
func (sdup *SDUPServer) Discovery(writer http.ResponseWriter, reader *http.Request) {
	devices, err := sdup.discovery()
	if err != nil {
		log.Log(log.Error, err.Error(), nil)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonEncoded, err := json.MarshalIndent(devices, "", "   ")
	if err != nil {
		log.Log(log.Error, err.Error(), nil)
		http.Error(writer, "Failed to JSON encode SDUPDevices", http.StatusInternalServerError)
	}
	writer.Write(jsonEncoded)
}

//discovery handles internal discovery procedures
func (sdup *SDUPServer) discovery() ([]devicecontainer.SDUPDevice, error) {
	SDUPDevices := []devicecontainer.SDUPDevice{}
	sdup.deviceContainerLock.Lock()
	defer sdup.deviceContainerLock.Unlock()

	if sdup.RecentlyRefreshed() {
		fmt.Println("Using cache for discovery")
		for _, device := range sdup.DeviceContainer.Devices {
			SDUPDevices = append(SDUPDevices, *device)
		}

	} else {
		hueLights, err := sdup.Bridge.GetLights()
		if err != nil {
			return nil, errors.New("Failed to enumerate devices (lights) on Hue bridge")
		}
		for _, light := range hueLights {
			device, _ := createDeviceFromLight(sdup, light)
			//Update the converters state
			//FIXME What to do with lost devices?
			added, changed, _, err := sdup.DeviceContainer.Update(&device)
			if err != nil {
				return nil, err
			}
			attributeMap := map[attributes.AttributeKey]devicecontainer.SDUPAttribute{}
			for _, attribute := range append(added, changed...) {
				attribute.Capabilities = nil
				attributeMap[attribute.Name] = attribute
			}
			//FIXME Blocking
			if len(attributeMap) > 0 {
				sdup.Subscriptions.EventChan <- devicecontainer.SDUPDevice{
					ID:         device.ID,
					Attributes: attributeMap,
				}
			}
			SDUPDevices = append(SDUPDevices, device)
		}
		sdup.lastDiscovered = time.Now()
	}
	return SDUPDevices, nil
}

func createDeviceFromLight(sdup *SDUPServer, light huego.Light) (devicecontainer.SDUPDevice, error) {
	device := devicecontainer.SDUPDevice{
		ID: devicecontainer.NewHueSDUPDeviceID("light", light.ID).Stringify(),
		Attributes: map[attributes.AttributeKey]devicecontainer.SDUPAttribute{
			attributes.AttributeActive: {
				Name:         attributes.AttributeActive,
				BooleanState: &light.State.On,
				Capabilities: map[capabilities.CapabilityKey]capabilities.Capability{
					capabilities.CapabilityActivate: capabilities.PreConfiguredCapability{
						CapabilityCallback: func(lightID int) func() error {
							return func() error {
								log.Log(log.Info, fmt.Sprintf("Setting state on=true on light %d", lightID), nil)
								_, err := sdup.Bridge.SetLightState(lightID, huego.State{On: true})
								return err
							}
						}(light.ID),
					},
					capabilities.CapabilityDeactivate: capabilities.PreConfiguredCapability{
						CapabilityCallback: func(lightID int) func() error {
							return func() error {
								log.Log(log.Info, fmt.Sprintf("Setting state on=false on light %d", lightID), nil)
								_, err := sdup.Bridge.SetLightState(lightID, huego.State{On: false})
								return err
							}
						}(light.ID),
					},
				},
			},
		},
	}
	switch strings.ToLower(light.Type) {
	// Types fetched from Phillips Hue documentation
	// https://developers.meethue.com/develop/hue-api/supported-devices/
	// Lowercased just to be sure
	case "dimmable light":
	case "color temperature light":
	// 	device.Attributes[attributes.AttributeColor] = devicecontainer.SDUPAttribute{
	// 		KeyVal: map[string]interface{}{"x": light.State.Xy[0], "y": light.State.Xy[1]},
	// 		Capabilities: map[capabilities.CapabilityKey]capabilities.Capability{capabilities.CapabilitySetAllKeyVal: capabilities.KeyValCapability{
	// 			CapabilityCallback: GenerateXYFunction(sdup, light.ID),
	// 		},
	// 		},
	// 	}
	case "color light":
		device.Attributes[attributes.AttributeColor] = devicecontainer.SDUPAttribute{
			Name:   attributes.AttributeColor,
			KeyVal: capabilities.RawKeyValContainer{"x": light.State.Xy[0], "y": light.State.Xy[1]},
			Capabilities: map[capabilities.CapabilityKey]capabilities.Capability{capabilities.CapabilitySetAllKeyVal: capabilities.KeyValCapability{
				CapabilityCallback: GenerateXYFunction(sdup, light.ID),
			},
			},
		}
	case "extended color light":
		device.Attributes[attributes.AttributeColor] = devicecontainer.SDUPAttribute{
			Name:   attributes.AttributeColor,
			KeyVal: capabilities.RawKeyValContainer{"x": light.State.Xy[0], "y": light.State.Xy[1]},
			Capabilities: map[capabilities.CapabilityKey]capabilities.Capability{capabilities.CapabilitySetAllKeyVal: capabilities.KeyValCapability{
				CapabilityCallback: GenerateXYFunction(sdup, light.ID),
			},
			},
		}
	default:
		log.Log(log.Warning, "Unknown light mode", map[string]string{"mode": light.Type})
	}

	return device, nil
}

//GenerateXYFunction Given a sdupserver and an id, generate a function that takes a KeyValContainer that changes color with XY
func GenerateXYFunction(sdup *SDUPServer, lightID int) func(capabilities.KeyValContainer) error {
	return func(input capabilities.KeyValContainer) (err error) {
		var x, y float32
		if x, err = input.Float32Key("x"); err != nil {
			return
		}
		if y, err = input.Float32Key("y"); err != nil {
			return
		}
		_, err = sdup.Bridge.SetLightState(lightID, huego.State{On: true, Xy: []float32{x, y}})
		return
	}
}
