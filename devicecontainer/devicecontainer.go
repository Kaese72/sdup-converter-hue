package devicecontainer

import (
	"fmt"
	"time"

	"github.com/Kaese72/sdup-hue/log"
)

//HueDeviceContainer manages all assets related to this SDUP converter
type HueDeviceContainer struct {
	Devices  map[string]*SDUPDevice
	LastSync time.Time
}

//Update updates a already known device or adds it to the known devices
func (container *HueDeviceContainer) Update(device *SDUPDevice) (added []SDUPAttribute, changed []SDUPAttribute, lost []SDUPAttribute, err error) {
	oldDevice, ok := container.Devices[device.ID]
	if !ok {
		//TODO What to do when new device is added?
		// Registration event?
		log.Log(log.Info, "Added device", map[string]string{
			"device": device.ID,
		})
	} else {
		added, changed, lost, err = oldDevice.HowChanged(device)
		if err != nil {
			return

		}
		if len(added) == 0 && len(changed) == 0 && len(lost) == 0 {
			log.Log(log.Debug, "Updated device", map[string]string{
				"device":  device.ID,
				"added":   fmt.Sprint(len(added)),
				"changed": fmt.Sprint(len(changed)),
				"lost":    fmt.Sprint(len(lost)),
			})

		} else {
			log.Log(log.Info, "Updated device", map[string]string{
				"device":  device.ID,
				"added":   fmt.Sprint(len(added)),
				"changed": fmt.Sprint(len(changed)),
				"lost":    fmt.Sprint(len(lost)),
			})
		}

	}
	container.Devices[device.ID] = device
	return
}

//NewHueDeviceContainer initializes a new HueDeviceContainer
func NewHueDeviceContainer() HueDeviceContainer {
	return HueDeviceContainer{
		Devices: map[string]*SDUPDevice{},
	}
}
