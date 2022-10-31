package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Kaese72/sdup-converter-hue/config"
	"github.com/Kaese72/sdup-converter-hue/devicestoreupdater"
	"github.com/Kaese72/sdup-converter-hue/sduphue"
	"github.com/Kaese72/sdup-lib/httpsdup"
	log "github.com/Kaese72/sdup-lib/logging"
	"github.com/Kaese72/sdup-lib/sdupcache"
	"github.com/spf13/viper"
)

func main() {
	myVip := viper.New()
	// We have elected to no use AutomaticEnv() because of https://github.com/spf13/viper/issues/584
	// myVip.AutomaticEnv()
	// Set replaces to allow keys like "database.mongodb.connection-string"
	myVip.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// # Configuration file configuration
	myVip.SetConfigName("config")
	myVip.AddConfigPath(".")
	myVip.AddConfigPath("99_local")
	myVip.AddConfigPath("/etc/sdup-converter-hue/")
	if err := myVip.ReadInConfig(); err != nil {
		log.Error(err.Error())
	}

	// # API configuration
	// Listen address
	myVip.BindEnv("http-server.address")
	myVip.SetDefault("http-server.address", "0.0.0.0")
	// Listen port
	myVip.BindEnv("http-server.port")
	myVip.SetDefault("http-server.port", 8080)

	// # Hue server configuration
	// URL to server, eg. http://192.168.100.102:80/
	myVip.BindEnv("hue.url")
	// preconfigured api key string
	myVip.BindEnv("hue.api-key")

	// # Logging
	myVip.BindEnv("debug-logging")
	myVip.SetDefault("debug-logging", false)

	// # Enroll Config
	// ENROLL_STORE
	myVip.BindEnv("enroll.store")
	// ENROLL_STORE_BRIDGE_ADDRESS
	myVip.BindEnv("enroll.bridge.address")
	// ENROLL_STORE_BRIDGE_PORT
	myVip.BindEnv("enroll.bridge.port")
	// We set the default port since we know what port the container will be listening on,
	// but we can not set a default on the address since we have not clue what IP it will have
	myVip.SetDefault("enroll.bridge.port", 8080)

	var conf config.Config
	err := myVip.Unmarshal(&conf)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := myVip.WriteConfigAs("./config.used.yaml"); err != nil {
		log.Error(err.Error())
		return
	}

	if err := conf.Validate(); err != nil {
		log.Error(err.Error())
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
	myBasePath := fmt.Sprintf("%s:%d", conf.HTTPServer.ListenAddress, conf.HTTPServer.ListenPort)
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
