package sduphue

import (
	"github.com/Kaese72/sdup-lib/devicestoretemplates"
)

// capRegistry contains device capability functions
var capRegistry = map[devicestoretemplates.CapabilityKey]func(id int, _ devicestoretemplates.CapabilityArgs) error{}

// gGapRegistry contains group capability functions
var gCapRegistry = map[devicestoretemplates.CapabilityKey]func(id int, _ devicestoretemplates.CapabilityArgs) error{}
