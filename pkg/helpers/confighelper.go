package helpers


import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const (
	configFile string = "config.json"
	configDir string = ".webjobdeploy"
)

type AppServiceConfig struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	WebjobName      string `json:"webjobname"`
	WebjobExeName   string `json:"webjobexename"`
	AppServiceName  string `json:"appservicename"`
}

// Config has all the configs for the various app services.
// ideally it would be a map to an AppServiceConfig but unsure about JSON and maps.
// so will just leave as array. Yes, needs to traverse it, but hardly a perf concern.
type Config struct {
	AppServicesConfigs []AppServiceConfig  `json:"appserviceconfigs"`
}

func replaceIfNotEmpty( newText string, defaultText string ) string {
	if strings.TrimSpace(newText) != "" {
		return strings.TrimSpace(newText)
	} else {
		return defaultText
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}


// readConfig reads a JSON file and returns the appropriate config.
func ReadConfig(  ) (*Config,error) {

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath := home + "/" + configDir + "/" + configFile

	if !fileExists(configPath) {
		// file doesn't exist, create it.
		os.MkdirAll(home + "/" + configDir , 0700)
		_, _ = os.Create(configPath)
	}

	data, err := ioutil.ReadFile( configPath)
	if err != nil {
		log.Fatal(err)
	}
	config := Config{}
	json.Unmarshal(data, &config)

	return &config, nil
}

// readConfig reads a JSON file and returns the appropriate config.
func WriteConfig( config Config ) error {

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := home + "/" + configDir + "/" + configFile

	jsonBytes,err :=json.Marshal(config)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(configPath, jsonBytes, 0644)

  return err
}

// get config... will return a map of string to string, but will
// initially populate via the config file (~/.webjobdeploy/config.json) and will overload
// from the passed in params.
func GetConfig( username string, password string, appServiceName string, webjobExeName string, webjobName string) (*AppServiceConfig, error) {

	// load config
	fullConfig, err := ReadConfig()
	if err != nil {
		fmt.Printf("GetConfig error %s\n", err.Error())
		return nil, err
	}

	foundConfig := false
	var configToUse AppServiceConfig

	lowerAppServiceName := strings.ToLower(appServiceName)
	// find config for the appServiceName
	for _,c := range fullConfig.AppServicesConfigs {
		if strings.ToLower(c.AppServiceName) == lowerAppServiceName {
			foundConfig = true
			configToUse = c
		}
	}

	if !foundConfig {
		configToUse = AppServiceConfig{}
	}

	configToUse.WebjobName = replaceIfNotEmpty(webjobName,configToUse.WebjobName )
	configToUse.WebjobExeName = replaceIfNotEmpty(webjobExeName,configToUse.WebjobExeName)
	configToUse.AppServiceName = replaceIfNotEmpty(appServiceName,configToUse.AppServiceName)
	configToUse.Password = replaceIfNotEmpty(password,configToUse.Password)
	configToUse.Username = replaceIfNotEmpty(username,configToUse.Username)

	return &configToUse, nil
}

// storeConfig stores PARTS of the config to the config file.
// needs to replace
func StoreConfig( config AppServiceConfig) error {
	fullConfig, err := ReadConfig()
	if err != nil {
		return err
	}

	found := false
	// loop through to find appropriate appServicePlan in the full config.
	// replace if exists, add if new
	for i,c := range fullConfig.AppServicesConfigs {
		if c.AppServiceName == config.AppServiceName {
			fullConfig.AppServicesConfigs[i] = config
			found = true
		}
	}

	if !found {
		fullConfig.AppServicesConfigs = append(fullConfig.AppServicesConfigs, config)
	}

	// save...
	err = WriteConfig(*fullConfig)
	if err != nil {
		return err
	}

	return nil
}

// check if config is missing stuff.
func ValidConfig( config AppServiceConfig, zipFileName string, uploadPath string) bool {
	if config.Username == "" {
		fmt.Printf("Invalid Username")
		return false
	}

	if config.Password == "" {
		fmt.Printf("Invalid Password")
		return false
	}

	if config.AppServiceName == "" {
		fmt.Printf("Invalid AppServiceName")
		return false
	}

	if config.WebjobExeName == "" {
		fmt.Printf("Invalid WebjobExeName")
		return false
	}

	if config.WebjobName == "" {
		fmt.Printf("Invalid WebjobName")
		return false
	}

	if zipFileName == "" && uploadPath == "" {
		fmt.Printf("need a valid zipFileName or uploadPath")
		return false
	}
	return true
}
