package capabilities

import (
	"io"
)

//CapabilityKey is the string identifier of a capability
type CapabilityKey string

const (
	//CapabilityActivate means the associated attribute can be activated
	CapabilityActivate CapabilityKey = "activate"
	//CapabilityDeactivate means the associated attribute can be deactivated
	CapabilityDeactivate CapabilityKey = "deactivate"
	//CapabilitySetAllKeyVal measn that you can change all attribute keyvals at the same time
	CapabilitySetAllKeyVal CapabilityKey = "setkeyval"
)

//Capability represents a capability
type Capability interface {
	TriggerCapability(capability io.ReadCloser) error
}
