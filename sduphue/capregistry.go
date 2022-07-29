package sduphue

import "github.com/Kaese72/sdup-lib/sduptemplates"

//capRegistry contains device capability functions
var capRegistry = map[sduptemplates.CapabilityKey]func(id int, _ sduptemplates.CapabilityArgument) error{}

//gGapRegistry contains group capability functions
var gCapRegistry = map[sduptemplates.CapabilityKey]func(id int, _ sduptemplates.CapabilityArgument) error{}
