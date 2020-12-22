package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Kaese72/sdup-hue/config"
	"github.com/Kaese72/sdup-hue/sdupserver"
	"github.com/amimof/huego"
	"github.com/gorilla/mux"
)

func main() {
	var conf config.Config

	if err := json.NewDecoder(os.Stdin).Decode(&conf); err != nil {
		log.Fatal(err)
	}

	if err := conf.Validate(); err != nil {
		obj, err2 := json.Marshal(config.NewExampleConfig())
		if err2 != nil {
			log.Fatal(err2)
		}
		_, err2 = fmt.Fprintf(os.Stdout, "%s\n", obj)
		if err2 != nil {
			log.Fatal(err2)
		}
		log.Fatal(err)
	}

	bridge := huego.New(conf.Hue.URL, conf.Hue.APIKey)
	SDUPServer, _ := sdupserver.NewSDUPServer(bridge, false, true)
	router := mux.NewRouter()
	router.HandleFunc("/discovery", SDUPServer.Discovery)
	router.HandleFunc("/subscribe", SDUPServer.Subscriptions.Subscribe)
	router.HandleFunc("/capability/{device_id}/{attribute_key}/{capability_key}", SDUPServer.CapabilityTrigger).Methods("POST")
	router.PathPrefix("/ui/").Handler(http.StripPrefix("/ui/", http.FileServer(http.Dir("./ui/"))))

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", conf.SDUP.ListenAddress, conf.SDUP.ListenPort), router); err != nil {
		log.Fatal(err)
	}
}
