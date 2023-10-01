package devicestoreupdater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-converter-hue/config"
	"github.com/Kaese72/sdup-lib/sduptemplates"
	"github.com/Kaese72/sdup-lib/subscription"
)

func pushDeviceUpdate(config config.StoreEnrollmentConfig, deviceUpdate sduptemplates.DeviceUpdate) error {
	bPayload, err := json.Marshal(deviceUpdate.DeviceStorePatch())
	if err != nil {
		logging.Error("Failed to marshal struct to JSON to send to device store", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	logging.Error("Sending blob to device store", map[string]interface{}{"blob": string(bPayload)})
	devicePayload, err := http.NewRequest("POST", fmt.Sprintf("%s/device-store/v0/devices", config.StoreURL), bytes.NewBuffer(bPayload))
	if err != nil {
		logging.Error("Failed to create request", map[string]interface{}{"error": err.Error()})
		return err
	}
	devicePayload.Header.Set("Bridge-Key", config.AdapterKey)
	resp, err := http.DefaultClient.Do(
		devicePayload,
	)
	if err != nil {
		logging.Error("Failed to http.Do request", map[string]interface{}{"error": err.Error()})
		return err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read response body on response", map[string]interface{}{"error": err.Error()})
		return err
	}

	logging.Info("Sent payload to device store", map[string]interface{}{"Response Code": resp.Status, "Response Body": string(respBody)})
	return nil
}

func pushDeviceGroupUpdate(config config.StoreEnrollmentConfig, groupUpdate sduptemplates.DeviceGroupUpdate) error {
	bPayload, err := json.Marshal(groupUpdate.DeviceStorePatch())
	if err != nil {
		logging.Error("Failed to marshal struct to JSON to send to device store", map[string]interface{}{
			"error": err.Error(),
		})
		return err
	}
	logging.Info("Sending blob to device store", map[string]interface{}{"blob": string(bPayload)})
	devicePayload, err := http.NewRequest("POST", fmt.Sprintf("%s/device-store/v0/groups", config.StoreURL), bytes.NewBuffer(bPayload))
	if err != nil {
		logging.Error("Failed to create request", map[string]interface{}{"error": err.Error()})
		return err
	}
	devicePayload.Header.Set("Bridge-Key", config.AdapterKey)
	resp, err := http.DefaultClient.Do(
		devicePayload,
	)
	if err != nil {
		logging.Error("Failed to http.Do request", map[string]interface{}{"error": err.Error()})
		return err
	}
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Error("Failed to read response body on response", map[string]interface{}{"error": err.Error()})
		return err
	}

	logging.Info("Sent payload to device store", map[string]interface{}{"Response Code": resp.Status, "Response Body": string(respBody)})
	return nil
}

func InitDeviceStoreUpdater(config config.StoreEnrollmentConfig, subscriptions subscription.Subscriptions) error {
	logging.Info("Starting device store updater")
	sub := subscriptions.Subscribe()
	defer subscriptions.UnSubscribe(sub)
	for update := range sub.Updates() {
		if dUpdate, err := update.GetDeviceUpdate(); err == nil {
			if err := pushDeviceUpdate(config, dUpdate); err != nil {
				logging.Error("Failed to send device group update", map[string]interface{}{"error": err.Error()})
			}

		} else if gUpdate, err := update.GetDeviceGroupUpdate(); err == nil {
			if err := pushDeviceGroupUpdate(config, gUpdate); err != nil {
				logging.Error("Failed to send device group update", map[string]interface{}{"error": err.Error()})
			}
		} else {
			logging.Error("Update did not evaluate to group or device update")
		}
	}
	return nil
}
