package sduphue

import (
	"encoding/base64"
	"fmt"
	"sync"
	"time"

	log "github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

var bridge *huego.Bridge

type HueCapability interface {
	Trigger(int, sduptemplates.CapabilityArgument) error
	Spec() sduptemplates.CapabilitySpec
}

type HueAttribute struct {
	State sduptemplates.AttributeState
}

func (attr HueAttribute) Spec() sduptemplates.AttributeSpec {
	return sduptemplates.AttributeSpec{
		AttributeState: attr.State,
	}
}

type HueDeviceID struct {
	Index int
	Type  string
}

//SDUPEncode converts the index and type into a DeviceID
// It is a base64 encoded string on the format "<type>/<index>"
func (id HueDeviceID) SDUPEncode() sduptemplates.DeviceID {
	return sduptemplates.DeviceID(base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s/%d", id.Type, id.Index))))
}

type HueDevice struct {
	ID           HueDeviceID
	Attributes   map[sduptemplates.AttributeKey]HueAttribute
	Capabilities map[sduptemplates.CapabilityKey]HueCapability
}

func (device HueDevice) Spec() sduptemplates.DeviceSpec {
	attrMap := sduptemplates.AttributeSpecMap{}
	for attrKey, attrVal := range device.Attributes {
		attrMap[attrKey] = attrVal.Spec()
	}
	capMap := sduptemplates.CapabilitySpecMap{}
	for capKey, capVal := range device.Capabilities {
		capMap[capKey] = capVal.Spec()
	}
	return sduptemplates.DeviceSpec{
		ID:           device.ID.SDUPEncode(),
		Attributes:   attrMap,
		Capabilities: capMap,
	}
}

type SDUPHueTarget struct {
	devices     map[sduptemplates.DeviceID]HueDevice
	updateChan  chan sduptemplates.DeviceUpdate
	lock        sync.RWMutex
	initialized bool
}

func (target *SDUPHueTarget) Initialize() (channel chan sduptemplates.DeviceUpdate, err error) {
	if target.initialized {
		panic("Hue target already initialized")
	}
	// Fetch all devices currently present on bridge
	devices, err := target.getAllDevices()
	if err != nil {
		return
	}
	// Register all devices
	for _, device := range devices {
		log.Info("Initializeing bridge with light", map[string]string{"light": fmt.Sprint(device.ID.Index)})
		target.devices[device.ID.SDUPEncode()] = device
	}

	// Start updater
	go func() {
		timer := time.NewTicker(2 * time.Second)
		//FIXME
		defer timer.Stop()
		for range timer.C {
			err := target.UpdateAllDevices()
			if err != nil {
				log.Error(err.Error())
			}
		}
	}()
	target.initialized = true
	channel = target.updateChan
	return
}

func (target *SDUPHueTarget) Devices() (devices []sduptemplates.DeviceSpec, err error) {
	target.lock.RLock()
	defer target.lock.RUnlock()
	for _, device := range target.devices {
		devices = append(devices, device.Spec())
	}
	return
}

func (target *SDUPHueTarget) TriggerCapability(deviceID sduptemplates.DeviceID, capabilityKey sduptemplates.CapabilityKey, argument sduptemplates.CapabilityArgument) error {
	target.lock.RLock()
	device, ok := target.devices[deviceID]
	if !ok {
		log.Debug("Could not find device", map[string]string{"device": string(deviceID)})
		target.lock.RUnlock()
		return sduptemplates.NoSuchDevice
	}

	capability, ok := device.Capabilities[capabilityKey]
	if !ok {
		log.Debug("Could not find capability", map[string]string{"device": string(deviceID), "capability": string(capabilityKey)})
		target.lock.RUnlock()
		return sduptemplates.NoSuchCapability
	}
	//We are unlocking before since the triggering can take quite a while
	target.lock.RUnlock()
	return capability.Trigger(device.ID.Index, argument)

}

func InitSDUPHueTarget(URL, APIKey string) sduptemplates.SDUPTarget {
	bridge = huego.New(URL, APIKey)
	target := &SDUPHueTarget{
		devices:    map[sduptemplates.DeviceID]HueDevice{},
		updateChan: make(chan sduptemplates.DeviceUpdate, 10),
	}
	return target
}
