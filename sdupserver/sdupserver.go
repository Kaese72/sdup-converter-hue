package sdupserver

import (
	"sync"
	"time"

	"github.com/Kaese72/sdup-hue/devicecontainer"
	"github.com/Kaese72/sdup-hue/log"
	"github.com/amimof/huego"
)

//SDUPServer contains everything for this sdup converter to function
type SDUPServer struct {
	DeviceContainer     devicecontainer.HueDeviceContainer
	Bridge              huego.Bridge
	Subscriptions       *Subscriptions
	lastDiscovered      time.Time
	deviceContainerLock sync.Mutex
	UseCache            bool
}

//NewSDUPServer creates and starts a SDUP service instance
func NewSDUPServer(bridge *huego.Bridge, useCache bool, refresher bool) (*SDUPServer, error) {
	deviceContainer := devicecontainer.NewHueDeviceContainer()
	SDUPServer := &SDUPServer{
		DeviceContainer: deviceContainer,
		Bridge:          *bridge,
		Subscriptions:   NewSubscriptions(),
	}
	if refresher {
		go func() {
			timer := time.NewTicker(2 * time.Second)
			defer timer.Stop()
			for {
				select {
				case <-timer.C:
					_, err := SDUPServer.discovery()
					if err != nil {
						log.Log(log.Error, err.Error(), nil)
					}
				}
			}
		}()
	}
	return SDUPServer, nil
}

//RecentlyRefreshed checks if the cached info should be used for a discovery request
func (sdup *SDUPServer) RecentlyRefreshed() bool {
	//If one minute has passed since the last discovery, it was no longer recently refreshed
	if sdup.UseCache {
		return time.Now().Before(sdup.lastDiscovered.Add(1 * time.Second))
	}
	return false
}
