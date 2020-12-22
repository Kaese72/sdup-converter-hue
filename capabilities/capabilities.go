package capabilities

//CapabilityKey is the string identifier of a capability
type CapabilityKey string

const (
	//CapabilityActivate means the associated attribute can be activated
	CapabilityActivate CapabilityKey = "activate"
	//CapabilityDeactivate means the associated attribute can be deactivated
	CapabilityDeactivate CapabilityKey = "deactivate"
)

//Capability represents a capability
type Capability interface {
}
