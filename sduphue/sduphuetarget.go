package sduphue

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kaese72/device-store/ingestmodels"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-lib/adapter"
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

func parseHueDeviceID(id string) (HueDeviceID, *adapter.AdapterError) {
	res, err := base64.URLEncoding.DecodeString(string(id))
	if err != nil {
		log.Debug("Failed base64 decode id, resulting in no such device", map[string]interface{}{"VALUE": string(id)})
		return HueDeviceID{}, &adapter.AdapterError{Code: 404, Message: "No such device"}
	}

	splitString := strings.SplitN(string(res), "/", 2)
	if len(splitString) < 2 {
		log.Debug("Failed to split device ID on /, resulting in no such device", map[string]interface{}{"VALUE": string(res), "LEN": strconv.Itoa(len(splitString))})
		return HueDeviceID{}, &adapter.AdapterError{Code: 404, Message: "No such device"}
	}

	index, err := strconv.Atoi(splitString[1])
	if err != nil {
		log.Debug("Failed to parse hue index, resulting in no such device", map[string]interface{}{"VALUE": splitString[1]})
		return HueDeviceID{}, &adapter.AdapterError{Code: 404, Message: "No such device"}
	}

	return HueDeviceID{
		Index: index,
		Type:  HueDeviceIDType(splitString[0]),
	}, nil
}

type SDUPHueTarget struct {
}

func (target SDUPHueTarget) Initialize() (chan adapter.Update, error) {
	// Start updater
	channel := make(chan adapter.Update)
	go func() {
		timer := time.NewTicker(2 * time.Second)
		defer timer.Stop()
		for range timer.C {
			hueDevices, err := target.getAllDevices()
			if err != nil {
				log.Error("Error when fetching devices", map[string]interface{}{"error": err.Error()})
			} else {
				for _, newDevice := range hueDevices {
					channel <- adapter.Update{Device: &newDevice}
				}
			}
			hueGroups, err := target.getAllGroups()
			if err != nil {
				log.Error("Error when fetching groups", map[string]interface{}{"error": err.Error()})
			} else {
				for _, newGroup := range hueGroups {
					channel <- adapter.Update{Group: &newGroup}
				}
			}
			timer.Reset(2 * time.Second)
		}
	}()
	return channel, nil
}

func (target SDUPHueTarget) DeviceTriggerCapability(deviceID string, capabilityKey string, argument ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	capability, ok := capRegistry[capabilityKey]
	if !ok {
		// It might be worth looking into being able to differentiate between bridge not supporting and the capability truly not existing
		log.Debug("Could not find capability", map[string]interface{}{"device": string(deviceID), "capability": string(capabilityKey)})
		return &adapter.AdapterError{Code: 404, Message: "No such capability"}
	}

	hueID, err := parseHueDeviceID(deviceID)
	if err != nil {
		return err
	}

	switch hueID.Type {
	case LIGHT:
		return capability(hueID.Index, argument)

	default:
		return &adapter.AdapterError{Code: 404, Message: "No such device"}
	}
}

func (target SDUPHueTarget) GroupTriggerCapability(groupID string, capabilityKey string, argument ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	capability, ok := gCapRegistry[capabilityKey]
	if !ok {
		// It might be worth looking into being able to differentiate between bridge not supporting and the capability truly not existing
		log.Debug("Could not find capability", map[string]interface{}{"group": string(groupID), "capability": string(capabilityKey)})
		return &adapter.AdapterError{Code: 404, Message: "No such capability"}
	}

	groupIndex, err := strconv.Atoi(groupID)
	if err != nil {
		return &adapter.AdapterError{Code: 404, Message: "No such group"}
	}

	return capability(groupIndex, argument)
}

func InitSDUPHueTarget(URL, APIKey string) SDUPHueTarget {
	bridge = huego.New(URL, APIKey)
	target := SDUPHueTarget{}
	return target
}
