package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ()

func doRequest(url string, account string, key string, body map[string]map[string]string, method string) int {
	var code int

	jsonValue, _ := json.Marshal(body)
	request, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonValue))
	if method == "POST" {
		request.Header.Set("Content-Type", "application/vnd.appd.events+json;v=2")
	}
	request.Header.Set("Accept", "application/vnd.appd.events+json;v=2")
	request.Header.Set("X-Events-API-AccountName", account)
	request.Header.Set("X-Events-API-Key", key)
	if debug {
		fmt.Printf("doRequest %s [jsonValue] %s X-Events-API-AccountName:%s X-Events-API-Key:%s\n", method, bytes.NewBuffer(jsonValue), account, key)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("doRequest %s [error] %s\n", method, err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		if debug {
			fmt.Printf("doRequest %s [response] %s\n", method, string(data))
		}
	}
	code = response.StatusCode
	defer response.Body.Close()
	return code
}

// Same as above really but instead of the body as map it's already json
func doRequestJSON(url string, account string, key string, body []byte, method string) int {
	var code int

	request, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	if method == "POST" {
		request.Header.Set("Content-Type", "application/vnd.appd.events+json;v=2")
	}
	request.Header.Set("Accept", "application/vnd.appd.events+json;v=2")
	request.Header.Set("X-Events-API-AccountName", account)
	request.Header.Set("X-Events-API-Key", key)
	if debug {
		fmt.Printf("doRequestJSON %s [jsonValue] %s X-Events-API-AccountName:%s X-Events-API-Key:%s\n", method, bytes.NewBuffer(body), account, key)
	}
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Printf("doRequestJSON %s [error] %s\n", method, err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		if debug {
			fmt.Printf("doRequestJSON %s [response] %s\n", method, string(data))
		}
	}
	code = response.StatusCode
	defer response.Body.Close()
	return code
}
