package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	jinshuju "github.com/nullsimon/jinshuju-go"
)

func main() {

	var conf jinshuju.JinshujuConf
	raw, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(raw, &conf)

	client := jinshuju.NewClient(conf)
	form, err := client.GetFormFields("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(form))
	entries, err := client.GetFormEntries("test")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(entries))
}
