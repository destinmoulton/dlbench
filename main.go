package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const settingsFile = "settings.json"
const downloadsDir = "downloads"

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
		fmt.Println("Testing Domain: " + domain.Domain + " on " + domain.Host)
		for _, file := range settings.Files {
			var url = domain.Domain + domain.Path + file.Name
			fmt.Println("Downloading " + file.Name + " from " + domain.Domain)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				fmt.Errorf("Error. Domain returned a %s status.", resp.Status)
			} else {
				fmt.Printf("Successfully downloaded.")
			}

			var domainDir = cleanDomainDir(domain.Domain)
			createDownloadFolder(domainDir)
			var outpath = filepath.Join(downloadsDir, domainDir, file.Name)
			outfile, err := os.Create(outpath)
			if err != nil {
				fmt.Println(err)
			}
			defer outfile.Close()

			_, err = io.Copy(outfile, resp.Body)
			if err != nil {
				fmt.Println(err)
			}

		}
	}
}

// Clean the domain name for use as a directory
func cleanDomainDir(domain string) string {
	res := domain
	res = strings.TrimPrefix(res, "http://")
	return strings.TrimPrefix(res, "https://")
}

// Create a download folder
func createDownloadFolder(subdir string) {
	dir := filepath.Join(downloadsDir, subdir)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Unable to create the directory %s", dir)
		os.Exit(-1)
	}
}
