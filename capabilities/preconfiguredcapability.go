package capabilities

import "io"

//PreConfiguredCapability represents a capability that does not take a ny input
type PreConfiguredCapability struct {
	CapabilityCallback func() error `json:"-"`
}

//TriggerCapability triggers the SimpleCapability callback function
func (cap PreConfiguredCapability) TriggerCapability(_ io.ReadCloser) error {
	return cap.CapabilityCallback()
}
