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
	Url     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
}

type DomainRequest struct {
	RequestMap map[string]Request `json:"requests"`
	Variables  map[string]string  `json:"variables"`
}

type Config map[string]DomainRequest

// configFilePath is the path to the config file
var configDir = "/curle"
var configFilePath = "/curle/domain.json"

// tartget is the url to fetch
var target = flag.String("target", "", "the target url")

func IsVariable(value string) bool {
	return value[0] == '$' && value[1] == '{' && value[len(value)-1] == '}'
}

func GetRequest(target string, config Config) {
	targetParts := strings.Split(target, ".")
	url := config[targetParts[0]].RequestMap[targetParts[1]].Url
	request, err := http.NewRequest("GET", url, nil)

	for k, v := range config[targetParts[0]].RequestMap[targetParts[1]].Headers {
		if IsVariable(v) {
			variable := v[1 : len(v)-1]
			v = config[targetParts[0]].Variables[variable]
			request.Header.Add(k, v)
		} else {
			request.Header.Add(k, v)
		}
	}

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
	var config Config
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

func validateTarget(target string, config Config) bool {
	targetParts := strings.Split(target, ".")
	if len(targetParts) != 2 {
		return false
	}
	website, ok := config[targetParts[0]]
	if !ok {
		return false
	}

	if _, ok := website.RequestMap[targetParts[1]]; ok {
		return true
	}
	return false
}

func init() {
  flag.Parse()
  userConfigDir, err := os.UserConfigDir()
  if err != nil {
    panic(err)
  }
  configFilePath = userConfigDir + configFilePath
  configDir = userConfigDir + configDir
  
  // check config file
  // if _, err := os.Stat(configDir); os.IsNotExist(err){
  //   err = os.Mkdir(configDir, 755)
  //   if err != nil {
  //     panic(err)
  //   }
  // }

  // crate configFile if not exists
  if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
    if err != nil {
      fmt.Println("The `~/.config/curle/domain.json` file does not exist. Please create it")
      os.Exit(1)
    }
  }
}

func main() {
	config := getConfig(configFilePath)
    fmt.Println("asdadasdd")
	fmt.Println(config)
	ok := validateTarget(*target, config)
	if !ok {
		fmt.Println("Invalid target")
		os.Exit(1)
	}
	GetRequest(*target, config)
}
