package main

import (
	"fmt"
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
	fmt.Println("test")
	jsonFile, err := os.Open(settingsFile)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
}
