package devicecontainer

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"

	"github.com/Kaese72/sdup-hue/attributes"
	"github.com/Kaese72/sdup-hue/capabilities"
)

//SDUPDeviceID describes what is necessary for a ID to be valid in SDUP
type SDUPDeviceID interface {
	Stringify() string
}

//HueSDUPDeviceID is a Hue specific Device ID struct
type HueSDUPDeviceID struct {
	DeviceType string
	DeviceID   int
}

var idRegex = regexp.MustCompile(`^(\w+)/(\d+)$`)

//HueSDUPDeviceIDFromSDUPID decodes and SDUP ID to bridge specific identifier
func HueSDUPDeviceIDFromSDUPID(encoded string) (*HueSDUPDeviceID, error) {
	bytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	matchGroups := idRegex.FindSubmatch(bytes)
	if len(matchGroups) == 0 {
		return nil, fmt.Errorf("Invalid decoded ID format")
	}

	deviceID, _ := strconv.Atoi(string(matchGroups[2]))
	return &HueSDUPDeviceID{
		DeviceType: string(matchGroups[1]),
		DeviceID:   deviceID,
	}, nil
}

//NewHueSDUPDeviceID creates a new HueSDUPDeviceID based on inputs
func NewHueSDUPDeviceID(deviceType string, deviceID int) *HueSDUPDeviceID {
	return &HueSDUPDeviceID{
		DeviceType: deviceType,
		DeviceID:   deviceID,
	}
}

//Stringify encodes a HueSDUPDeviceID to a base64 url encoded string
func (id *HueSDUPDeviceID) Stringify() string {
	return base64.URLEncoding.EncodeToString([]byte(fmt.Sprintf("%s/%d", id.DeviceType, id.DeviceID)))
}

//SDUPDevice is the SDUP representation of any SDUP enabled device
type SDUPDevice struct {
	ID string `json:"id"`
	//Name       string                                    `json:"-"`
	Attributes map[attributes.AttributeKey]SDUPAttribute `json:"attributes"`
}

//HowChanged figures out what capabilities have been added, changed and lost
func (device *SDUPDevice) HowChanged(otherDevice *SDUPDevice) ([]SDUPAttribute, []SDUPAttribute, []SDUPAttribute, error) {
	var added, changed, lost []SDUPAttribute
	for attributeKey, oldAttribute := range device.Attributes {
		newAttribute, ok := otherDevice.Attributes[attributeKey]
		if ok {
			equal, err := newAttribute.Equal(&oldAttribute)
			if err != nil {
				// Something went wrong, abort
				return nil, nil, nil, err
			}

			if !equal {
				changed = append(changed, newAttribute)
			}
			// If there is no diff, the attribute exist but nothing has changed

		} else {
			lost = append(lost, oldAttribute)
		}
	}
	for attributeKey, newAttribute := range otherDevice.Attributes {
		if _, ok := device.Attributes[attributeKey]; !ok {
			added = append(added, newAttribute)
		}
	}
	return added, changed, lost, nil
}

//SDUPAttribute describes the state and capabilities of an SDUP attribute
type SDUPAttribute struct {
	Name         attributes.AttributeKey                                `json:"-"`
	BooleanState *bool                                                  `json:"boolean-state,omitempty"`
	KeyVal       capabilities.RawKeyValContainer                        `json:"keyval-state,omitempty"`
	Capabilities map[capabilities.CapabilityKey]capabilities.Capability `json:"capabilities,omitempty"`
}

//Equal determines if two attribute values of the same type are equal
func (attr *SDUPAttribute) Equal(other *SDUPAttribute) (bool, error) {
	if attr.BooleanState != nil && other.BooleanState != nil {
		return *attr.BooleanState == *other.BooleanState, nil

	} else if attr.KeyVal != nil && other.KeyVal != nil {
		return reflect.DeepEqual(attr.KeyVal, other.KeyVal), nil
	}

	return false, errors.New("Could not find common descriptor")
}
