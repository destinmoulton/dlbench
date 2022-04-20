package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const settingsFile = "settings.json"

type JSONSections struct {
	Domains []JSONDomain `json:"domains"`
	Files   []JSONFile   `json:"files"`
}

type JSONDomain struct {
	Domain string `json:"domain"`
	Path   string `json:"path"`
	Host   string `json:"host"`
}

type JSONFile struct {
	Name string `json:"name"`
	Size string `json:"size"`
}

func main() {
	var settings JSONSections

	fmt.Println("test")
	jsonFile, err := os.Open(settingsFile)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &settings)

	for _, domain := range settings.Domains {
		for _, file := range settings.Files {
			var url = domain.Domain + domain.Path + file.Name

			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
		}
		fmt.Println("Domain:" + domain.Domain)
	}
}
