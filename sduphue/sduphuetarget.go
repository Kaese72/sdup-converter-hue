package sduphue

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kaese72/device-store/rest/models"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-lib/deviceupdates"
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

type SDUPHueTarget struct {
}

func (target SDUPHueTarget) Initialize() (chan deviceupdates.Update, error) {
	// Start updater
	channel := make(chan deviceupdates.Update)
	go func() {
		timer := time.NewTicker(2 * time.Second)
		defer timer.Stop()
		for range timer.C {
			hueDevices, err := target.getAllDevices()
			if err != nil {
				log.Error("Error when fetching devices", map[string]interface{}{"error": err.Error()})
			} else {
				for _, newDevice := range hueDevices {
					channel <- deviceupdates.Update{Device: newDevice}
				}
			}
			timer.Reset(2 * time.Second)
		}
	}()
	return channel, nil
}

func (target SDUPHueTarget) TriggerCapability(deviceID string, capabilityKey string, argument models.CapabilityArgs) error {
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

func InitSDUPHueTarget(URL, APIKey string) SDUPHueTarget {
	bridge = huego.New(URL, APIKey)
	target := SDUPHueTarget{}
	return target
}
