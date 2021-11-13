package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Request struct {
	Url    string `json:"url"`
	Method string `json:"method"`
}

type Config map[string]map[string]Request

// configFilePath is the path to the config file
var configFilePath = flag.String("url", "./config.json", "url to fetch")

// tartget is the url to fetch
var target = flag.String("target", "", "the target url")

func GetRequest(target string, config Config) {
	targetParts := strings.Split(target, ".")
	url := config[targetParts[0]][targetParts[1]].Url
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(response.Header.Get("content-type"))
    responseType := response.Header.Get("content-type")
    if strings.Contains(responseType, "application/json") {
      err = json.Unmarshal(body, &body)
      var prettyJSON bytes.Buffer
      err = json.Indent(&prettyJSON, body, "", "  ")
      if err != nil {
          fmt.Println(err)
      }
      fmt.Println(string(prettyJSON.Bytes()))
    } else {
      fmt.Println(string(body))
    }
}

func getConfig(configFilePath string) Config {
	var config map[string]map[string]Request
	configFile, err := os.Open(configFilePath)
	if err != nil {
		fmt.Println(err)
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		fmt.Println(err)
	}
	return config
}

func validateTarget(target string, config map[string]map[string]Request) bool {
	targetParts := strings.Split(target, ".")
	if len(targetParts) != 2 {
		return false
	}
	website, ok := config[targetParts[0]]
	if ok {
		return true
	}

	if _, ok := website[targetParts[1]]; ok {
		return true
	}
	return false
}

func main() {
	flag.Parse()
	config := getConfig(*configFilePath)
	ok := validateTarget(*target, config)
	if !ok {
		fmt.Println("Invalid target")
		os.Exit(1)
	}
	GetRequest(*target, config)
}
