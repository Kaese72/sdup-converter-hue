package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Kaese72/sdup-converter-hue/config"
	"github.com/Kaese72/sdup-converter-hue/sduphue"
	"github.com/Kaese72/sdup-lib/httpsdup"
	log "github.com/Kaese72/sdup-lib/logging"
)

func main() {
	var conf config.Config

	if err := json.NewDecoder(os.Stdin).Decode(&conf); err != nil {
		log.Error(err.Error())
	}

	if err := conf.Validate(); err != nil {
		obj, err := json.Marshal(config.NewExampleConfig())
		if err != nil {
			log.Error(err.Error())
		}
		_, err = fmt.Fprintf(os.Stdout, "%s\n", obj)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	SDUPServer := sduphue.InitSDUPHueTarget(conf.Hue.URL, conf.Hue.APIKey)
	router := httpsdup.InitHTTPMux(SDUPServer)
	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.SDUP.ListenAddress, conf.SDUP.ListenPort), router); err != nil {
		log.Error(err.Error())
	}
}
