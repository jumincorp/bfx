package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type responseIface interface {
	append(map[string]interface{})
	getMetrics(prefix string, namespace string, id string) []metrics
	setCommandName(string)
}

type response struct {
	command     string
	data        []map[string]interface{}
	headerNames map[string]string
}

func (r *response) setCommandName(command string) {
	r.command = command
}

func (r *response) append(obj map[string]interface{}) {
	//log.Printf(":: OBJ :: %v\n", obj)
	r.data = append(r.data, obj)
}

func sanitizeName(name string) string {
	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, "%", "_Percent", -1)
	return name
}

func (r *response) metricName(prefix string, data string) string {
	return strings.Join([]string{prefix, r.command, sanitizeName(data)}, "_")
}

func newResponse(command string, responseBytes []byte) responseIface {

	var responseData = make(map[string]interface{})
	err := json.Unmarshal(bytes.Trim(responseBytes, "\x00"), &responseData)
	if err != nil {
		panic(err)
	}
	//fmt.Printf("JSON0: %v\n\n", data)

	delete(responseData, "STATUS")
	delete(responseData, "id")

	log.Printf("JSON: %v\n\n", responseData)

	// By now we should have only one element
	if len(responseData) != 1 {
		log.Printf("len(responseData) != 1 : %v\n", len(responseData))
		return nil
	}

	var r responseIface

	switch command {
	case "devs":
		r = newDevsResponse()
	}

	r.setCommandName(command)

	var respList []interface{}
	for _, wrapped := range responseData {
		respList = wrapped.([]interface{})
	}
	//log.Printf("respList: %v\n\n", respList)

	for _, respElement := range respList {
		r.append(respElement.(map[string]interface{}))
	}

	return r
}

func (r *response) getMetrics(prefix string, namespace string, id string) []metrics {

	var metricsArray = make([]metrics, len(r.data))

	for _, element := range r.data {
		var metrics metrics
		metrics.headers = make(map[string]string)
		metrics.headers["namespace"] = namespace
		metrics.headers["miner"] = id

		metrics.values = make(map[string]float64)

		log.Printf("---")
		for name, val := range element {

			if format, present := r.headerNames[name]; present {
				//log.Printf("header %s : %s\n", name, fmt.Sprintf(format, val))
				metrics.headers[name] = fmt.Sprintf(format, val)
			} else {
				metricName := r.metricName(prefix, name)
				switch casted := val.(type) {
				case int64:
					log.Printf("i! %s %d\n", name, casted)
				case float64:
					//log.Printf("f  %s >> %s %f\n", name, metricName, casted)
					metrics.values[metricName] = casted
				case string:
					log.Printf("s! %s %s\n", name, casted)
				default:
					log.Printf("?! %s %s\n", name, casted)
				}
			}
		}
		metricsArray = append(metricsArray, metrics)
	}
	return metricsArray
}

func (r *response) makeResponse() {
	r.headerNames = make(map[string]string)
}

type devsResponse struct {
	response
}

func newDevsResponse() *devsResponse {
	r := new(devsResponse)
	r.makeResponse()

	r.headerNames["Name"] = "%s"
	r.headerNames["ID"] = "%.0f"

	//r.response.data = make([]map[string]interface{}, 10)

	return r
}
