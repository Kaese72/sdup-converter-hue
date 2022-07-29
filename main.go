package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Kaese72/sdup-converter-hue/config"
	"github.com/Kaese72/sdup-converter-hue/devicestoreupdater"
	"github.com/Kaese72/sdup-converter-hue/sduphue"
	"github.com/Kaese72/sdup-lib/httpsdup"
	log "github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/sdupcache"
)

func main() {
	var conf config.Config

	if err := json.NewDecoder(os.Stdin).Decode(&conf); err != nil {
		log.Error(err.Error())
	}

	if err := conf.Validate(); err != nil {
		log.Error(err.Error())
		conf.PopulateExample()
		obj, err := json.Marshal(conf)
		if err != nil {
			log.Error(err.Error())
		}
		_, err = fmt.Fprintf(os.Stdout, "%s\n", obj)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	//Logger assumed initiated
	if conf.DebugLogging != nil {
		log.SetDebugLogging(*conf.DebugLogging)
	}
	myBasePath := fmt.Sprintf("%s:%d", conf.SDUP.ListenAddress, conf.SDUP.ListenPort)
	hueTarget := sduphue.InitSDUPHueTarget(conf.Hue.URL, conf.Hue.APIKey)
	sdupCache := sdupcache.NewSDUPCache(hueTarget)
	router, subscriptions := httpsdup.InitHTTPMux(sdupCache)
	log.Info("Starting HTTP Server")
	go func() {
		err := devicestoreupdater.InitDeviceStoreUpdater(conf.EnrollDeviceStore, subscriptions)
		if err != nil {
			log.Error("Failed to initiate device store updater", map[string]string{"error": err.Error()})
			return
		}
	}()
	if err := http.ListenAndServe(myBasePath, router); err != nil {
		log.Error(err.Error())
	}

}
