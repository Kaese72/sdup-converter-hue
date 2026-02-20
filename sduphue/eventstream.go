package sduphue

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"

	"github.com/Kaese72/device-store/ingestmodels"
	log "github.com/Kaese72/huemie-lib/logging"
	"github.com/Kaese72/sdup-lib/adapter"
	"github.com/openhue/openhue-go"
)

var (
	seenLightBridgeIdentifiers = map[string]struct{}{}
	seenLightMutex             sync.Mutex
)

type HueEventStreamMessage struct {
	CreationTime string               `json:"creationtime"`
	Data         []HueEventStreamData `json:"data"`
	ID           string               `json:"id"`
	Type         string               `json:"type"`
}

type HueEventStreamData struct {
	ID      string                `json:"id"`
	IDV1    string                `json:"id_v1"`
	Type    string                `json:"type"`
	On      *HueEventStreamOn     `json:"on,omitempty"`
	Dimming *HueEventStreamDim    `json:"dimming,omitempty"`
	Status  *HueEventStreamStatus `json:"status,omitempty"`
}

type HueEventStreamOn struct {
	On bool `json:"on"`
}

type HueEventStreamDim struct {
	Brightness float64 `json:"brightness"`
}

type HueEventStreamStatus struct {
	Active string `json:"active"`
}

func (target SDUPHueTarget) listenEventStream(updates chan adapter.Update) error {
	_ = updates
	if target.home == nil {
		return errors.New("home not initialized")
	}

	eventURL, err := getEventStreamURL(target.host)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", eventURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "text/event-stream")
	req.Header.Set("hue-application-key", target.apiKey)

	client := target.httpClient
	if client == nil {
		client = &http.Client{}
		if target.ignoreTLSErrors {
			client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("eventstream connect failed: %s - %s", resp.Status, strings.TrimSpace(string(body)))
	}

	reader := bufio.NewReader(resp.Body)
	var dataLines []string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return err
			}
			return err
		}

		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			if len(dataLines) > 0 {
				payload := strings.Join(dataLines, "\n")
				dataLines = nil
				if err := target.handleEventStreamPayload(payload, updates); err != nil {
					log.Error("Failed to handle eventstream payload", map[string]interface{}{"error": err.Error()})
				}
			}
			continue
		}

		if after, ok := strings.CutPrefix(line, "data:"); ok {
			data := strings.TrimSpace(after)
			dataLines = append(dataLines, data)
		}
	}
}

func (target SDUPHueTarget) handleEventStreamPayload(payload string, updates chan adapter.Update) error {
	if payload == "" {
		return nil
	}

	var messages []HueEventStreamMessage
	if err := json.Unmarshal([]byte(payload), &messages); err != nil {
		var single HueEventStreamMessage
		if err := json.Unmarshal([]byte(payload), &single); err != nil {
			return err
		}
		messages = []HueEventStreamMessage{single}
	}

	for _, message := range messages {
		if !isUpdateMessageType(message.Type) {
			log.Debug("Skipping non-update message from eventstream", map[string]interface{}{"message_type": message.Type})
			continue
		}
		updatesByIdentifier := map[string]ingestmodels.IngestDevice{}
		newDevices := map[string]struct{}{}
		for _, data := range message.Data {
			if !isLightDataType(data.Type) {
				continue
			}
			bridgeIdentifier := strings.TrimSpace(data.ID)
			if bridgeIdentifier == "" {
				log.Debug("Skipping light update with missing id", map[string]interface{}{"id_v1": data.IDV1})
				continue
			}
			if !wasSeenBefore(bridgeIdentifier) {
				newDevices[bridgeIdentifier] = struct{}{}
				continue
			}
			if _, isNew := newDevices[bridgeIdentifier]; isNew {
				continue
			}

			device, ok := createLightDeviceUpdateFromEvent(data)
			if !ok {
				continue
			}
			existing, ok := updatesByIdentifier[bridgeIdentifier]
			if ok {
				existing.Attributes = mergeAttributes(existing.Attributes, device.Attributes)
				updatesByIdentifier[bridgeIdentifier] = existing
				continue
			}
			updatesByIdentifier[bridgeIdentifier] = device
		}

		for bridgeIdentifier := range newDevices {
			light, err := target.getLightByID(bridgeIdentifier)
			if err != nil {
				log.Error("Failed to fetch light from bridge", map[string]interface{}{"id": bridgeIdentifier, "error": err.Error()})
				continue
			}
			device := createLightDevice(*light)
			updates <- adapter.Update{Device: &device}
		}

		for _, device := range updatesByIdentifier {
			updates <- adapter.Update{Device: &device}
		}
	}

	return nil
}

func isUpdateMessageType(messageType string) bool {
	return strings.EqualFold(messageType, "update")
}

func isLightDataType(dataType string) bool {
	return strings.EqualFold(dataType, "light")
}

func createLightDeviceUpdateFromEvent(data HueEventStreamData) (ingestmodels.IngestDevice, bool) {
	bridgeIdentifier := strings.TrimSpace(data.ID)
	if bridgeIdentifier == "" {
		log.Debug("Skipping light update with missing id", map[string]interface{}{"id_v1": data.IDV1})
		return ingestmodels.IngestDevice{}, false
	}

	attributes := make([]ingestmodels.IngestAttribute, 0, 2)
	if data.On != nil {
		on := data.On.On
		attributes = append(attributes, ingestmodels.IngestAttribute{
			Name:    AttributeActive,
			Boolean: &on,
		})
	}
	if data.Dimming != nil {
		brightness := float32(data.Dimming.Brightness)
		attributes = append(attributes, ingestmodels.IngestAttribute{
			Name:    AttributeBrightness,
			Numeric: &brightness,
		})
	}
	if len(attributes) == 0 {
		return ingestmodels.IngestDevice{}, false
	}

	return ingestmodels.IngestDevice{
		BridgeIdentifier: bridgeIdentifier,
		Attributes:       attributes,
	}, true
}

func mergeAttributes(existing []ingestmodels.IngestAttribute, incoming []ingestmodels.IngestAttribute) []ingestmodels.IngestAttribute {
	if len(incoming) == 0 {
		return existing
	}
	if len(existing) == 0 {
		return incoming
	}
	indexByName := make(map[string]int, len(existing))
	for i, attr := range existing {
		indexByName[attr.Name] = i
	}
	for _, attr := range incoming {
		if idx, ok := indexByName[attr.Name]; ok {
			existing[idx] = attr
			continue
		}
		indexByName[attr.Name] = len(existing)
		existing = append(existing, attr)
	}
	return existing
}

func wasSeenBefore(bridgeIdentifier string) bool {
	seenLightMutex.Lock()
	defer seenLightMutex.Unlock()
	if _, ok := seenLightBridgeIdentifiers[bridgeIdentifier]; ok {
		return true
	}
	seenLightBridgeIdentifiers[bridgeIdentifier] = struct{}{}
	return false
}

func getEventStreamURL(host string) (string, error) {
	if strings.Contains(host, "://") {
		return "", fmt.Errorf("eventstream host must not include scheme: %s", host)
	}
	u := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/eventstream/clip/v2",
	}
	u.RawPath = path.Clean(u.Path)
	return u.String(), nil
}

func (target SDUPHueTarget) getLightByID(id string) (*openhue.LightGet, error) {
	lights, err := target.home.GetLights()
	if err != nil {
		return nil, err
	}
	if light, ok := lights[id]; ok {
		return &light, nil
	}
	return nil, fmt.Errorf("light not found: %s", id)
}
