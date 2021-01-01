package capabilities

import (
	"encoding/json"
	"fmt"
	"io"
)

//KeyValCapability represents a capability that takes a map with values as input
type KeyValCapability struct {
	CapabilityCallback func(KeyValContainer) error `json:"-"`
}

//KeyValContainer contains an unknown set keys with an unknown value type
// Exposes methods that collect and convert appropriately
type KeyValContainer interface {
	//IntKey(string) (int, error)
	Float32Key(string) (float32, error)
}

//RawKeyValContainer contains an unknown set keys with an unknown value type
type RawKeyValContainer map[string]interface{}

//Float32Key returns a float32 if the key is actually a float, otherwise error
func (container RawKeyValContainer) Float32Key(key string) (result float32, err error) {
	if rawValue, ok := container[key]; ok {
		switch ConvertedValue := rawValue.(type) {
		case float32:
			result = ConvertedValue
		case float64:
			result = float32(ConvertedValue)
		default:
			err = fmt.Errorf("Key value %s not float", key)
		}
	} else {
		err = fmt.Errorf("Missing value %s", key)
	}
	return
}

//TriggerCapability triggers the SimpleCapability callback function
func (cap KeyValCapability) TriggerCapability(bodyReader io.ReadCloser) error {
	var bodyMap RawKeyValContainer
	if err := json.NewDecoder(bodyReader).Decode(&bodyMap); err != nil {
		return err
	}
	return cap.CapabilityCallback(bodyMap)
}
