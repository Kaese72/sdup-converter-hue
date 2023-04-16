package devicestoreupdater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Kaese72/sdup-converter-hue/config"
	"github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/subscription"
)

func InitDeviceStoreUpdater(config config.StoreEnrollmentConfig, subscriptions subscription.Subscriptions) error {
	logging.Info("Starting device store updater")
	sub := subscriptions.Subscribe()
	defer subscriptions.UnSubscribe(sub)
	for update := range sub.Updates() {
		dUpdate, err := update.GetDeviceUpdate()
		if err != nil {
			// This is a device update and not a group update
			continue
		}
		bPayload, err := json.Marshal(dUpdate.DeviceStorePatch())
		if err != nil {
			logging.Error("Failed to marshal struct to JSON to send to device store", map[string]interface{}{
				"error": err.Error(),
			})
			continue
		}
		logging.Error("Sending blob to device store", map[string]interface{}{"blob": string(bPayload)})
		devicePayload, err := http.NewRequest("POST", fmt.Sprintf("%s/rest/v0/devices", config.StoreURL), bytes.NewBuffer(bPayload))
		if err != nil {
			logging.Error("Failed to create request", map[string]interface{}{"error": err.Error()})
			continue
		}
		devicePayload.Header.Set("Bridge-Key", config.AdapterKey)
		resp, err := http.DefaultClient.Do(
			devicePayload,
		)
		if err != nil {
			logging.Error("Failed to http.Do request", map[string]interface{}{"error": err.Error()})
			continue
		}
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logging.Error("Failed to read response body on response", map[string]interface{}{"error": err.Error()})
			continue
		}

		logging.Info("Sent payload to device store", map[string]interface{}{"Response Code": resp.Status, "Response Body": string(respBody)})
	}
	return nil
}
