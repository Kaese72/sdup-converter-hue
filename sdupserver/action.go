package sdupserver

import (
	"fmt"
	"net/http"

	"github.com/Kaese72/sdup-hue/attributes"
	"github.com/Kaese72/sdup-hue/capabilities"
	"github.com/Kaese72/sdup-hue/log"
	"github.com/gorilla/mux"
)

//CapabilityTrigger triggers a capability for a attribute
func (sdup *SDUPServer) CapabilityTrigger(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)

	//These variables must exist
	deviceID := vars["device_id"]
	attributeKey := vars["attribute_key"]
	capabilityKey := vars["capability_key"]

	if device, ok := sdup.DeviceContainer.Devices[deviceID]; ok {
		if attribute, ok := device.Attributes[attributes.AttributeKey(attributeKey)]; ok {
			if capInterface, ok := attribute.Capabilities[capabilities.CapabilityKey(capabilityKey)]; ok {
				if err := capInterface.TriggerCapability(reader.Body); err != nil {
					log.Log(log.Error, err.Error(), nil)
					http.Error(writer, err.Error(), http.StatusInternalServerError)
				} else {
					log.Log(log.Info, fmt.Sprintf("Set capability %s on attribute %s on device %s", capabilityKey, attributeKey, deviceID), nil)
					http.Error(writer, fmt.Sprintf("Set capability %s on attribute %s on device %s", capabilityKey, attributeKey, deviceID), http.StatusOK)
				}

			} else {
				http.Error(writer, fmt.Sprintf("No known capability %s for attribute with ID %s for Device with ID %s", capabilityKey, attributeKey, deviceID), http.StatusNotFound)
			}
		} else {
			http.Error(writer, fmt.Sprintf("No known attribute with ID %s for Device with ID %s", attributeKey, deviceID), http.StatusNotFound)
		}
	} else {
		http.Error(writer, fmt.Sprintf("No known device with ID %s", deviceID), http.StatusNotFound)
	}
}
