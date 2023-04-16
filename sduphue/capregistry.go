package sduphue

import (
	devicestoretemplates "github.com/Kaese72/device-store/rest/models"
)

// capRegistry contains device capability functions
var capRegistry = map[devicestoretemplates.CapabilityKey]func(id int, _ devicestoretemplates.CapabilityArgs) error{}

// gGapRegistry contains group capability functions
var gCapRegistry = map[devicestoretemplates.CapabilityKey]func(id int, _ devicestoretemplates.CapabilityArgs) error{}
