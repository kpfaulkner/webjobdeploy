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

	username := flag.String("username", "", "Azure Username")
	password := flag.String("password", "", "Azure Password")
	zipFileName := flag.String("zipfilename", "", "Zip Filename")
	uploadPath := flag.String("uploadpath", "", "directory path to upload (if not using zip file)")
	appServiceName:= flag.String("appServiceName", "", "App Service name")
	webjobName:= flag.String("webjobName", "", "Webjob name")
	webjobExeName:= flag.String("webjobExeName", "", "Webjob executable filename. eg, mywebjob.exe")

	help := flag.Bool("help", false, "help me Obi-Wan")

  flag.Parse()

  if *help {
    fmt.Printf("help me now.....\n")
    return
  }
	if *username == "" || *password == "" || *appServiceName == "" || *webjobName == "" || *webjobExeName == ""{
		fmt.Printf("Please check params\n")
		return
	}

  zipFilePath := ""
  var err error

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
	if *username != "" && *password != "" && *appServiceName != "" && *webjobName != "" && *webjobExeName != ""{
		helpers.Upload(*username, *password, *appServiceName, *webjobName, *webjobExeName,  bufReader)
	}


}
