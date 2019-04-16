package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
)

type gatewayMainRecord struct {
	Level               int    `json:"level"`
	IsRequestSuccess    string `json:"isRequestSuccess"`
	Time                string `json:"time"`
	OperationName       string `json:"operationName"`
	Category            string `json:"category"`
	DurationMs          int    `json:"durationMs"`
	CallerIPAddress     string `json:"callerIpAddress"`
	CorrelationID       string `json:"correlationId"`
	Location            string `json:"location"`
	ResourceID          string `json:"resourceId"`
	BackendMethod       string `json:"backendMethod"`
	BackendURL          string `json:"backendUrl"`
	Method              string `json:"method"`
	URL                 string `json:"url"`
	BackendResponseCode int    `json:"backendResponseCode"`
	ResponseCode        int    `json:"responseCode"`
	ResponseSize        int    `json:"responseSize"`
	Cache               string `json:"cache"`
	BackendTime         int    `json:"backendTime"`
	RequestSize         int    `json:"requestSize"`
	APIID               string `json:"apiId"`
	OperationID         string `json:"operationId"`
	ApimSubscriptionID  string `json:"apimSubscriptionId"`
	ClientProtocol      string `json:"clientProtocol"`
	BackendProtocol     string `json:"backendProtocol"`
	APIRevision         string `json:"apiRevision"`
	Source              string `json:"source"`
	Reason              string `json:"reason"`
	Message             string `json:"message"`
	Section             string `json:"section"`
}

func serializeGatewayRecord(data []byte, conf appdConfig) int {
	//var record []gatewayMainRecord
	var d interface{}
	if err := json.Unmarshal([]byte(data), &d); err != nil {
		log.Fatal(err)
	}

	m := d.(map[string]interface{})
	j := m["records"].([]interface{})
	record := make([]gatewayMainRecord, len(j))

	for i := 0; i < len(j); i++ {
		/* Main Record Map */
		if j[i] != nil {
			k := j[i].(map[string]interface{})

			/* Set the record */
			record[i].IsRequestSuccess = fmt.Sprintf("%v", k["isRequestSuccess"])
			leveStr := fmt.Sprint(k["level"])
			record[i].Level, _ = strconv.Atoi(leveStr)
			record[i].Time = fmt.Sprintf("%v", k["time"])
			record[i].OperationName = fmt.Sprintf("%v", k["operationName"])
			record[i].Location = fmt.Sprintf("%v", k["location"])
			record[i].Category = fmt.Sprintf("%v", k["category"])
			record[i].CorrelationID = fmt.Sprintf("%v", k["correlationId"])
			record[i].CallerIPAddress = fmt.Sprintf("%v", k["callerIpAddress"])
			durationStr := fmt.Sprint(k["durationMs"])
			record[i].DurationMs, _ = strconv.Atoi(durationStr)
			record[i].ResourceID = fmt.Sprintf("%v", k["resourceId"])
			record[i].Time = fmt.Sprintf("%v", k["time"])
			/* Set the properties */
			if k["properties"] != nil {
				l := k["properties"].(map[string]interface{})
				record[i].Method = fmt.Sprintf("%v", l["method"])
				record[i].BackendMethod = fmt.Sprintf("%v", l["backendMethod"])
				record[i].BackendURL = fmt.Sprintf("%v", l["backendMethod"])
				record[i].URL = fmt.Sprintf("%v", l["backendMethod"])
				backendRespStr := fmt.Sprint(l["backendResponseCode"])
				record[i].BackendResponseCode, _ = strconv.Atoi(backendRespStr)
				responseCodeStr := fmt.Sprint(l["responseCode"])
				record[i].ResponseCode, _ = strconv.Atoi(responseCodeStr)
				responseSizeStr := fmt.Sprint(l["responseSize"])
				record[i].ResponseSize, _ = strconv.Atoi(responseSizeStr)
				record[i].Cache = fmt.Sprintf("%v", l["backendMethod"])
				backendTimeStr := fmt.Sprint(l["backendTime"])
				record[i].BackendTime, _ = strconv.Atoi(backendTimeStr)
				requestSizeStr := fmt.Sprint(l["requestSize"])
				record[i].RequestSize, _ = strconv.Atoi(requestSizeStr)
				record[i].APIID = fmt.Sprintf("%v", l["apiId"])
				record[i].OperationID = fmt.Sprintf("%v", l["operationId"])
				record[i].ApimSubscriptionID = fmt.Sprintf("%v", l["apimSubscriptionId"])
				record[i].ClientProtocol = fmt.Sprintf("%v", l["clientProtocol"])
				record[i].BackendProtocol = fmt.Sprintf("%v", l["backendProtocol"])
				record[i].APIRevision = fmt.Sprintf("%v", l["apiRevision"])
				/* Last error, only if the request failed */
				if l["lastError"] != nil {
					/* Last Error Map, if it exists */
					e := l["lastError"].(map[string]interface{})
					record[i].Source = fmt.Sprintf("%v", e["source"])
					record[i].Reason = fmt.Sprintf("%v", e["reason"])
					record[i].Message = fmt.Sprintf("%v", e["message"])
					record[i].Section = fmt.Sprintf("%v", e["section"])
				}
			} // End of Property check
		} // End of For Loop
	} // End of record Map check

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
		url := conf.AnalyticsEp + publishURL + conf.AnalyticsGatewaySchema
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
