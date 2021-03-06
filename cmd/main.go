package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/kpfaulkner/webjobdeploy/pkg/helpers"
	"os"
)


func main() {
	fmt.Printf("Starting....\n")
	var err error
	username := flag.String("username", "", "Azure Deployment Username")
	password := flag.String("password", "", "Azure Deployment Password")
	zipFileName := flag.String("zipfilename", "", "Zip Filename")
	uploadPath := flag.String("uploadpath", "", "directory path to upload (if not using zip file)")
	appServiceName:= flag.String("appServiceName", "", "App Service/Azure Function name")
	webjobName:= flag.String("webjobName", "", "Webjob name")
	webjobExeName:= flag.String("webjobExeName", "", "Webjob executable filename. eg, mywebjob.exe")
  deploy := flag.String("deploy", "webjob", "indicates if deploying webjob, azure function or app service. Values can be appservice, azurefunction or webjob")

	help := flag.Bool("help", false, "help me Obi-Wan")
	store := flag.Bool("store", false, "Store the username/password against the App Service Name so it can look it up later. Config is stored in HOME/.webjobdeploy/config.json")

  flag.Parse()

  if *help {
    fmt.Printf("help me now.....\n")
    return
  }

  // read config from ~/.webjobdeploy/config.json
  // use passed in params to override anything from config.
  config, err := helpers.GetConfig( *username, *password, *appServiceName, *webjobExeName, *webjobName)
	if err != nil {
		fmt.Printf("Cannot get config %s\n", err.Error())
		return
	}

  if !helpers.ValidConfig(*config, *zipFileName, *uploadPath, *deploy ) {
	  fmt.Printf("Please check params\n")
	  return
  }

	if *store {
		helpers.StoreConfig( *config )
		fmt.Printf("WARNING: The credentials are currently stored in plain text. If you think this was a mistake, delete the file ~/.webjobdeploy/config.json")
	}

	zipFilePath := ""

  // if upload path specified, that takes priority and we'll zip that up ready to use.
  if *uploadPath != "" {

	  zipFilePath, err = helpers.GenerateZipFile(*uploadPath)
  	if err != nil {
  		fmt.Printf("Error while generating zipfile %s\n", err.Error())
  		return
	  }

	  // cleanup.. but has issues on windows. Need to dig about more. TODO(kpfaulkner) make sure it gets removed.
	  defer os.Remove(zipFilePath)
	  fmt.Printf("generated temp zip file %s\nYou might need to remove manually\n", zipFilePath)

  } else {
		zipFilePath = *zipFileName
  }

	file, err := os.Open(zipFilePath)
	if err != nil {
		fmt.Printf("Error while uploading %s\n", err.Error())
		panic(err)
	}

	defer file.Close()
	bufReader := bufio.NewReader(file)

	switch *deploy {
	case "webjob":
		helpers.UploadWebjob(*config,  bufReader)
	case "azurefunction","appservice":
		helpers.UploadAppService(*config, bufReader)
	}


}
