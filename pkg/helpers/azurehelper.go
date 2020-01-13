package helpers

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

type AzureHelper struct {

}

// generateAuthHeader generates the basic auth header.
// just base64-ing them etc.
func generateAuthHeader(username string, password string) string {
	s := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(s))
}

func Upload( config AppServiceConfig, deployZipBytes io.Reader ) error {
	authHeader := generateAuthHeader(config.Username, config.Password)

	urlTemplate := "https://%s.scm.azurewebsites.net/api/continuouswebjobs/%s"
	url := fmt.Sprintf(urlTemplate, config.AppServiceName, config.WebjobName)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, deployZipBytes)
	if err != nil {
		fmt.Printf("Error while uploading %s\n", err.Error())
		return err
	}

	req.Header.Add("Authorization", "Basic " +authHeader)
	req.Header.Add("Content-type", "application/zip")
	req.Header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", config.WebjobExeName))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("error on post %s\n", err.Error())
		panic(err)
	}
	fmt.Printf("return status code is %d\n", resp.StatusCode)
	return nil
}
