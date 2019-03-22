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
	AnalyticsGatewaySchema  string `yaml:"analyticsGatewaySchema"`
	AnalyticsCosmosSchema   string `yaml:"analyticsCosmosSchema"`
	AzureSubscriptionID     string `yaml:"azureSubscriptionId"`
	AzureEventHubNameSpace  string `yaml:"azureEventHubNameSpace"`
	AzureEventHubConnString string `yaml:"azureEventHubConnString"`
	AzureClientID           string `yaml:"azureClientID"`
	AzureClientSecret       string `yaml:"azureClientSecret"`
	AzureTenantID           string `yaml:"azureTenantID"`
	AzureGatewayHubName     string `yaml:"azureGatewayHubName"`
	AzureCosmosHubName      string `yaml:"azureCosmosHubName"`
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
		fmt.Printf("checkSchema [name] %s [HTTP Response] %d\n", schema, response)
	}
	if response != 200 {
		exists = false
	} else {
		exists = true
	}

	return exists
}

func createSchema(schema string, conf appdConfig, deftype int) {
	var schemaDef map[string]map[string]string
	url := conf.AnalyticsEp + schemaURL + schema
	if debug {
		fmt.Println("createSchema: [URL] " + url)
	}
	switch deftype {
	case typeGateway:
		schemaDef = map[string]map[string]string{"schema": {"level": "integer", "isRequestSuccess": "string", "time": "string", "operationName": "string", "category": "string", "durationMs": "integer", "callerIpAddress": "string", "correlationId": "string", "location": "string", "resourceId": "string", "backendMethod": "string", "backendUrl": "string", "method": "string", "url": "string", "backendResponseCode": "integer", "responseCode": "integer", "responseSize": "integer", "cache": "string", "backendTime": "integer", "requestSize": "integer", "apiId": "string", "operationId": "string", "apimSubscriptionId": "string", "clientProtocol": "string", "backendProtocol": "string", "apiRevision": "string", "source": "string", "reason": "string", "message": "string", "section": "string"}}
	case typeCosmos:
		schemaDef = map[string]map[string]string{"schema": {"time": "string", "resourceId": "string", "category": "string", "operationName": "string", "activityId": "string", "opCode": "string", "errorCode": "string", "duration": "integer", "requestCharge": "float", "databaseName": "string", "collectionName": "string", "retryCount": "string"}}
	default:
		fmt.Printf("ERROR - Unknown definition type for [schema] %s [type] %d\n", schema, deftype)
	}

	response := doRequest(url, conf.GlobalName, conf.Key, schemaDef, "POST")
	if response != 201 {
		fmt.Printf("ERROR - Failed to create [schema] %s [response] %d\n", schema, response)
	} else {
		fmt.Printf("Successfully created [schema] %s [response] %d\n", schema, response)
	}
}
