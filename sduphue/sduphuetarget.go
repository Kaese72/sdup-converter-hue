package sduphue

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kaese72/sdup-lib/devicestoretemplates"
	log "github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/amimof/huego"
)

var bridge *huego.Bridge

type HueDeviceIDType string

const (
	LIGHT HueDeviceIDType = "light"
)

type HueDeviceID struct {
	Index int
	Type  HueDeviceIDType
}

// SDUPEncode converts the index and type into a DeviceID
// It is a base64 encoded string on the format "<type>/<index>"
func (id HueDeviceID) SDUPEncode() string {
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s/%d", id.Type, id.Index)))
}

func parseHueDeviceID(id string) (HueDeviceID, error) {
	res, err := base64.URLEncoding.DecodeString(string(id))
	if err != nil {
		log.Debug("Failed base64 decode id, resulting in no such device", map[string]interface{}{"VALUE": string(id)})
		return HueDeviceID{}, sduptemplates.NoSuchDevice
	}

	splitString := strings.SplitN(string(res), "/", 2)
	if len(splitString) < 2 {
		log.Debug("Failed to split device ID on /, resulting in no such device", map[string]interface{}{"VALUE": string(res), "LEN": strconv.Itoa(len(splitString))})
		return HueDeviceID{}, sduptemplates.NoSuchDevice
	}

	index, err := strconv.Atoi(splitString[1])
	if err != nil {
		log.Debug("Failed to parse hue index, resulting in no such device", map[string]interface{}{"VALUE": splitString[1]})
		return HueDeviceID{}, sduptemplates.NoSuchDevice
	}

	return HueDeviceID{
		Index: index,
		Type:  HueDeviceIDType(splitString[0]),
	}, nil
}

type HueGroupID struct {
	Index int
}

func parseHueGroupID(id sduptemplates.DeviceGroupID) (HueGroupID, error) {
	intId, err := strconv.Atoi(string(id))
	if err != nil {
		log.Debug("Failed to atoi group id, leading to error", map[string]interface{}{"VALUE": string(id)})
		return HueGroupID{}, err
	}
	return HueGroupID{
		Index: intId,
	}, nil
}

type SDUPHueTarget struct {
	updateChan  chan sduptemplates.Update
	initialized bool
}

func (target *SDUPHueTarget) Initialize() (channel chan sduptemplates.Update, err error) {
	if target.initialized {
		panic("Hue target already initialized")
	}

	// Start updater
	go func() {
		timer := time.NewTicker(2 * time.Second)
		//FIXME
		defer timer.Stop()
		for range timer.C {
			err := target.UpdateEverything()
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
	return target.getAllDevices()
}

func (target *SDUPHueTarget) Groups() (devices []sduptemplates.DeviceGroupSpec, err error) {
	return target.getAllGroups()
}

func (target *SDUPHueTarget) TriggerCapability(deviceID string, capabilityKey devicestoretemplates.CapabilityKey, argument devicestoretemplates.CapabilityArgs) error {
	capability, ok := capRegistry[capabilityKey]
	if !ok {
		// It might be worth looking into being able to differentiate between bridge not supporting and the capability truly not existing
		log.Debug("Could not find capability", map[string]interface{}{"device": string(deviceID), "capability": string(capabilityKey)})
		return sduptemplates.NoSuchCapability
	}

	hueID, err := parseHueDeviceID(deviceID)
	if err != nil {
		return err
	}

	switch hueID.Type {
	case LIGHT:
		return capability(hueID.Index, argument)

	default:
		return sduptemplates.NoSuchDevice
	}
}

func (target *SDUPHueTarget) GTriggerCapability(groupId sduptemplates.DeviceGroupID, capabilityKey devicestoretemplates.CapabilityKey, argument devicestoretemplates.CapabilityArgs) error {
	capability, ok := gCapRegistry[capabilityKey]
	if !ok {
		// It might be worth looking into being able to differentiate between bridge not supporting and the capability truly not existing
		log.Debug("Could not find group capability", map[string]interface{}{"capability": string(capabilityKey)})
		return sduptemplates.NoSuchCapability
	}

	hueGroupId, err := parseHueGroupID(groupId)
	if err != nil {
		return err
	}

	return capability(hueGroupId.Index, argument)

}

func InitSDUPHueTarget(URL, APIKey string) sduptemplates.SDUPTarget {
	bridge = huego.New(URL, APIKey)
	target := &SDUPHueTarget{
		updateChan: make(chan sduptemplates.Update, 10),
	}
	return target
}
