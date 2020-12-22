package sdupserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
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
			device := devicecontainer.SDUPDevice{
				ID: devicecontainer.NewHueSDUPDeviceID("light", light.ID).Stringify(),
				Attributes: map[attributes.AttributeKey]devicecontainer.SDUPAttribute{
					attributes.AttributeActive: {
						Name:         attributes.AttributeActive,
						BooleanState: &light.State.On,
						Capabilities: map[capabilities.CapabilityKey]capabilities.Capability{
							capabilities.CapabilityActivate: capabilities.SimpleCapability{
								Callable: func(lightID int) func() error {
									return func() error {
										log.Log(log.Info, fmt.Sprintf("Setting state on=true on light %d", lightID), nil)
										_, err = sdup.Bridge.SetLightState(lightID, huego.State{On: true})
										return err
									}
								}(light.ID),
							},
							capabilities.CapabilityDeactivate: capabilities.SimpleCapability{
								Callable: func(lightID int) func() error {
									return func() error {
										log.Log(log.Info, fmt.Sprintf("Setting state on=false on light %d", lightID), nil)
										_, err = sdup.Bridge.SetLightState(lightID, huego.State{On: false})
										return err
									}
								}(light.ID),
							},
						},
					},
				},
				LastSeen: time.Now(),
			}
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
