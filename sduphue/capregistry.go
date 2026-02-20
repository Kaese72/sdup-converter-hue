package sduphue

import (
	"net/http"

	"github.com/Kaese72/device-store/ingestmodels"
	"github.com/Kaese72/sdup-lib/adapter"
)

// capRegistry contains device capability functions
var capRegistry = map[string]func(target SDUPHueTarget, id string, _ ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError{}

// gGapRegistry contains group capability functions
var gCapRegistry = map[string]func(target SDUPHueTarget, id string, _ ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError{}

func adapterErrorFromErr(err error) *adapter.AdapterError {
	if err == nil {
		return nil
	}
	return &adapter.AdapterError{Code: http.StatusInternalServerError, Message: err.Error()}
}
