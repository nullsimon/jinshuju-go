package main

import (
	"encoding/json"
	"io/ioutil"

	jinshuju "github.com/nullsimon/jinshuju-go"
	log "github.com/sirupsen/logrus"
)

func main() {

	var conf jinshuju.Conf
	raw, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(raw, &conf)

	client := jinshuju.NewClient(conf)
	form, err := client.GetFormFields("test")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Info(len(form))
	entries, err := client.GetFormEntries("test")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Info(len(entries))
}
