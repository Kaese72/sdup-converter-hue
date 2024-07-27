package sduphue

import (
	devicestoretemplates "github.com/Kaese72/device-store/rest/models"
)

// capRegistry contains device capability functions
var capRegistry = map[string]func(id int, _ devicestoretemplates.DeviceCapabilityArgs) error{}

// gGapRegistry contains group capability functions
var gCapRegistry = map[string]func(id int, _ devicestoretemplates.GroupCapabilityArgs) error{}
