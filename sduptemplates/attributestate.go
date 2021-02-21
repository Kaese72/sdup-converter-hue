package sduptemplates

import "github.com/Kaese72/sdup-hue/log"

//KeyValContainer contains an unknown set keys with an unknown value type
// Exposes methods that collect and convert appropriately
type KeyValContainer interface{}

//AttributeStateMap defines the relationship between AttributeKeys and AttributeStates
type AttributeStateMap map[AttributeKey]AttributeState

//AttributeState defines how an attribute state is communicated over SDUP
type AttributeState struct {
	BooleanState *bool            `json:"boolean-state,omitempty"`
	KeyVal       *KeyValContainer `json:"keyval-state,omitempty"`
}

func (state AttributeState) Equivalent(other AttributeState) bool {
	if state.BooleanState != nil && other.BooleanState != nil {
		return *state.BooleanState == *other.BooleanState
	}
	//TODO Implement Equivalence for KeyVal state
	log.Log(log.Error, "Could not find common state", nil)
	return false
}
