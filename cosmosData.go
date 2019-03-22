package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type cosmosMainRecord struct {
	Time           string  `json:"time"`
	ResourceID     string  `json:"resourceId"`
	Category       string  `json:"category"`
	OperationName  string  `json:"operationName"`
	ActivityID     string  `json:"activityId"`
	OpCode         string  `json:"opCode"`
	ErrorCode      string  `json:"errorCode"`
	Duration       int     `json:"duration"`
	RequestCharge  float64 `json:"requestCharge"`
	DatabaseName   string  `json:"databaseName"`
	CollectionName string  `json:"collectionName"`
	RetryCount     string  `json:"retryCount"`
}

func serializeCosmosRecord(data []byte, conf appdConfig) int {
	//var record []gatewayMainRecord
	var d interface{}
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		log.Fatal(err)
	}

	m := d.(map[string]interface{})
	j := m["records"].([]interface{})
	record := make([]cosmosMainRecord, len(j))

	for i := 0; i < len(j); i++ {
		/* Main Record Map */
		k := j[i].(map[string]interface{})
		/* Properties Map */
		l := k["properties"].(map[string]interface{})

		/* Set the record */
		record[i].Time = fmt.Sprintf("%v", k["time"])
		record[i].ResourceID = fmt.Sprintf("%v", k["resourceId"])
		record[i].Category = fmt.Sprintf("%v", k["category"])
		record[i].OperationName = fmt.Sprintf("%v", k["operationName"])
		record[i].ActivityID = fmt.Sprintf("%v", l["activityId"])
		record[i].OpCode = fmt.Sprintf("%v", l["opCode"])
		record[i].ErrorCode = fmt.Sprintf("%v", l["errorCode"])
		durationStr := fmt.Sprint(l["duration"])
		record[i].Duration, _ = strconv.Atoi(durationStr)
		requestChargeStr := fmt.Sprint(l["requestCharge"])
		record[i].RequestCharge, _ = strconv.ParseFloat(requestChargeStr, 64)
		record[i].DatabaseName = fmt.Sprintf("%v", l["databaseName"])
		record[i].CollectionName = fmt.Sprintf("%v", l["collectionName"])
		record[i].RetryCount = fmt.Sprintf("%v", l["retryCount"])
	}

	// Loop through each record in the batch and send to analytics
	for i := 0; i < len(record); i++ {
		// AppD API expects the JSON to be wrapped in []'s, and json.Marshal does not add them
		var a []byte
		a = append(a, '[')
		b, err := json.Marshal(record[i])
		b = append(b, ']')
		a = append(a, b...)
		if err != nil {
			fmt.Printf("ERROR marshalling JSON from record: %s\n", err)
		}

		// Push data to analytics
		url := conf.AnalyticsEp + publishURL + conf.AnalyticsCosmosSchema
		response := doRequestJSON(url, conf.GlobalName, conf.Key, a, "POST")
		if response != 200 {
			fmt.Printf("ERROR pushing analytics to AppD [message] %d\n", response)
		}
		if debug {
			fmt.Printf("Analytics [response] %d [record]: %+v\n", response, record)
		}
	}

	return len(record)
}
