package sduphue

import (
	"crypto/tls"
	"net/http"
	"strings"
	"time"

	"github.com/Kaese72/device-store/ingestmodels"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/openhue/openhue-go"
)

type SDUPHueTarget struct {
	home            *openhue.Home
	host            string
	apiKey          string
	ignoreTLSErrors bool
	httpClient      *http.Client
}

// Initialize starts the updater by first fetching all devices and groups and create an initial state
// which is seent to the device store. After that it opens a eventstream to the Hue API and continously listens
// to changes which updates are then sent to the device store. If the connection to the eventstream is lost,
// it will reconnect and continue listening for changes. The channel returned from this function is never closed,
// as the updater is expected to run indefinitely until the adapter is stopped.
func (target SDUPHueTarget) Initialize() (chan adapter.Update, error) {
	channel := make(chan adapter.Update)
	go func() {
		if err := target.sendInitialState(channel); err != nil {
			log.Error("Error when creating initial state", map[string]interface{}{"error": err.Error()})
		}

		for {
			err := target.listenEventStream(channel)
			if err != nil {
				log.Error("Eventstream disconnected", map[string]interface{}{"error": err.Error()})
			}
			time.Sleep(2 * time.Second)
		}
	}()
	return channel, nil
}

func (target SDUPHueTarget) sendInitialState(updates chan adapter.Update) error {
	hueDevices, err := target.getAllDevices()
	if err != nil {
		return err
	}
	for _, newDevice := range hueDevices {
		updates <- adapter.Update{Device: &newDevice}
	}

	hueGroups, err := target.getAllGroups()
	if err != nil {
		return err
	}
	for _, newGroup := range hueGroups {
		updates <- adapter.Update{Group: &newGroup}
	}

	return nil
}

func (target SDUPHueTarget) DeviceTriggerCapability(deviceID string, capabilityKey string, argument ingestmodels.IngestDeviceCapabilityArgs) *adapter.AdapterError {
	capability, ok := capRegistry[capabilityKey]
	if !ok {
		// It might be worth looking into being able to differentiate between bridge not supporting and the capability truly not existing
		log.Debug("Could not find capability", map[string]interface{}{"device": string(deviceID), "capability": string(capabilityKey)})
		return &adapter.AdapterError{Code: 404, Message: "No such capability"}
	}

	if strings.TrimSpace(deviceID) == "" {
		return &adapter.AdapterError{Code: 404, Message: "No such device"}
	}

	return capability(target, deviceID, argument)
}

func (target SDUPHueTarget) GroupTriggerCapability(groupID string, capabilityKey string, argument ingestmodels.IngestGroupCapabilityArgs) *adapter.AdapterError {
	capability, ok := gCapRegistry[capabilityKey]
	if !ok {
		// It might be worth looking into being able to differentiate between bridge not supporting and the capability truly not existing
		log.Debug("Could not find capability", map[string]interface{}{"group": string(groupID), "capability": string(capabilityKey)})
		return &adapter.AdapterError{Code: 404, Message: "No such capability"}
	}

	if strings.TrimSpace(groupID) == "" {
		return &adapter.AdapterError{Code: 404, Message: "No such group"}
	}

	return capability(target, groupID, argument)
}

func InitSDUPHueTarget(host, APIKey string, ignoreTLSErrors bool) (SDUPHueTarget, error) {
	var httpClient *http.Client
	if ignoreTLSErrors {
		httpClient = &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	}
	home, err := openhue.NewHome(host, APIKey)
	if err != nil {
		return SDUPHueTarget{}, err
	}
	target := SDUPHueTarget{home: home, host: host, apiKey: APIKey, ignoreTLSErrors: ignoreTLSErrors, httpClient: httpClient}
	return target, nil
}
