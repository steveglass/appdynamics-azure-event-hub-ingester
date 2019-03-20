package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

type appdConfig struct {
	Controller              string `yaml:"controller"`
	Port                    int    `yaml:"port"`
	SSL                     string `yaml:"ssl"`
	AnalyticsEp             string `yaml:"analyticsEndPoint"`
	GlobalName              string `yaml:"globalAccountName"`
	Key                     string `yaml:"analyticsKey"`
	AnalyticsSchema         string `yaml:"analyticsSchema"`
	AzureSubscriptionID     string `yaml:"azureSubscriptionId"`
	AzureEventHubNameSpace  string `yaml:"azureEventHubNameSpace"`
	AzureEventHubConnString string `yaml:"azureEventHubConnString"`
	AzureClientID           string `yaml:"azureClientID"`
	AzureClientSecret       string `yaml:"azureClientSecret"`
	AzureTenantID           string `yaml:"azureTenantID"`
	AzureHubName            string `yaml:"azureHubName"`
	AzureResourceGroup      string `yaml:"azureResourceGroup"`
}

func initConfig() appdConfig {
	// Eventually this needs to support more than a single Event Hub
	var appd appdConfig

	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, &appd)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return appd
}

func checkSchema(schema string, conf appdConfig) bool {
	var exists bool
	url := conf.AnalyticsEp + schemaURL + schema

	if debug {
		fmt.Println("createSchema: [URL] " + url)
	}

	// Does schema exists?
	response := doRequest(url, conf.GlobalName, conf.Key, nil, "GET")
	if debug {
		fmt.Printf("checkSchema [HTTP Response] %d\n", response)
	}
	if response != 200 {
		exists = false
	} else {
		exists = true
	}

	return exists
}

func createSchema(schema string, conf appdConfig) {
	url := conf.AnalyticsEp + schemaURL + schema
	if debug {
		fmt.Println("createSchema: [URL] " + url)
	}
	schemaDef := map[string]map[string]string{"schema": {"level": "integer", "isRequestSuccess": "string", "time": "string", "operationName": "string", "category": "string", "durationMs": "integer", "callerIpAddress": "string", "correlationId": "string", "location": "string", "resourceId": "string", "backendMethod": "string", "backendUrl": "string", "method": "string", "url": "string", "backendResponseCode": "integer", "responseCode": "integer", "responseSize": "integer", "cache": "string", "backendTime": "integer", "requestSize": "integer", "apiId": "string", "operationId": "string", "apimSubscriptionId": "string", "clientProtocol": "string", "backendProtocol": "string", "apiRevision": "string", "source": "string", "reason": "string", "message": "string", "section": "string"}}
	response := doRequest(url, conf.GlobalName, conf.Key, schemaDef, "POST")
	if response != 201 {
		fmt.Printf("ERROR - Failed to create schema [response] %d\n", response)
	} else {
		fmt.Printf("Successfully created schema [response] %d\n", response)
	}
}
