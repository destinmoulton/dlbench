package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const settingsFile = "settings.json"

var settings JSONSettings

type JSONSettings struct {
	Config  JSONConfig   `json:"config"`
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

type JSONConfig struct {
	Rounds         int    `json:"rounds"`
	DownloadFolder string `json:"download_folder"`
	CSVFile        string `json:"csv_file"`
}

func main() {
	populateSettings(&settings)
	csvFile, csvWriter := createCSVFile()

	for _, domain := range settings.Domains {
		fmt.Println("Testing Domain: " + domain.Domain + " on " + domain.Host)

		for _, file := range settings.Files {
			var totalSeconds float64
			var lastSize int64
			cnt := 0

			// Run for the number of rounds
			for i := 0; i < settings.Config.Rounds; i++ {
				var url = domain.Domain + domain.Path + file.Name
				fmt.Println("Downloading " + file.Name + " from " + domain.Domain)

				startBench := time.Now()
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
				}
				endBench := time.Now()
				elapsed := endBench.Sub(startBench)
				totalSeconds += elapsed.Seconds()

				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Printf("Error. Domain returned a %s status.\n", resp.Status)
				} else {
					fmt.Printf("Successfully downloaded in %fs.\n", elapsed.Seconds())
				}

				domainDir := cleanDomainDir(domain.Domain)
				createDownloadFolder(domainDir)
				outpath := filepath.Join(settings.Config.DownloadFolder, domainDir, file.Name)
				outfile, err := os.Create(outpath)
				if err != nil {
					fmt.Println(err)
				}
				defer outfile.Close()

				_, err = io.Copy(outfile, resp.Body)
				if err != nil {
					fmt.Println(err)
				}

				info, err := os.Stat(outpath)
				if err != nil {
					fmt.Printf("Unable to stat %s\n", outpath)
					os.Exit(-1)
				}

				lastSize = info.Size()

				cnt++
			}
			kb := float64(lastSize) / (1 << 10)
			avgTime := totalSeconds / float64(cnt)
			avgRate := float64(kb) / float64(avgTime)
			row := []string{
				domain.Domain,
				file.Name,
				strconv.Itoa(cnt),
				strconv.FormatFloat(kb, 'f', 2, 64) + "kB",
				strconv.FormatFloat(avgTime, 'f', 2, 64) + "s",
				strconv.FormatFloat(avgRate, 'f', 2, 64) + "kB/s",
				domain.Host,
			}
			writeCSVRow(csvWriter, row)
		}
	}
	csvWriter.Flush()
	csvFile.Close()
}

// Clean the domain name for use as a directory
func cleanDomainDir(domain string) string {
	res := domain
	res = strings.TrimPrefix(res, "http://")
	return strings.TrimPrefix(res, "https://")
}

// Create a download folder
func createDownloadFolder(subdir string) {
	dir := filepath.Join(settings.Config.DownloadFolder, subdir)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		fmt.Printf("Unable to create the directory %s\n", dir)
		os.Exit(-1)
	}
}

func createCSVFile() (*os.File, *csv.Writer) {
	csvFile, err := os.Create(settings.Config.CSVFile)

	if err != nil {
		fmt.Printf("Unable to create the CSV file %s\n", settings.Config.CSVFile)
		os.Exit(-1)
	}
	csvWriter := csv.NewWriter(csvFile)

	// Write the header row
	header := []string{
		"Domain",
		"File",
		"Rounds",
		"File Size (kB)",
		"Avg. Time (s)",
		"Avg. Rate (kB/s)",
		"Host",
	}
	writeCSVRow(csvWriter, header)
	return csvFile, csvWriter
}

func writeCSVRow(csvWriter *csv.Writer, row []string) {
	err := csvWriter.Write(row)
	if err != nil {
		fmt.Printf("Unable to write %v to CSV\n", row)
		fmt.Println(err)
		os.Exit(-1)
	}
}

// Populate the settings struct from the settings json file
func populateSettings(settings *JSONSettings) {
	jsonFile, err := os.Open(settingsFile)

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &settings)
}
