package sduphue

import (
	"github.com/Kaese72/device-store/ingestmodels"
)

// capRegistry contains device capability functions
var capRegistry = map[string]func(id int, _ ingestmodels.DeviceCapabilityArgs) error{}

// gGapRegistry contains group capability functions
var gCapRegistry = map[string]func(id int, _ ingestmodels.GroupCapabilityArgs) error{}
