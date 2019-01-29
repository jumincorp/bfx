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
	getMetrics(metricsList *[]metrics, prefix string, namespace string, id string)
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
	name = strings.Replace(name, "*", "", -1)
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
	case "devdetails":
		r = newDevdetailsResponse()
	case "procs":
		r = newProcsResponse()
	case "stats":
		r = newStatsResponse()
	case "pools":
		r = newPoolsResponse()
	case "coin":
		r = newCoinResponse()
	case "summary":
		r = newPoolsResponse()
	case "notify":
		r = newNotifyResponse()
	default:
		r = newDefaultResponse()
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

func (r *response) getMetrics(metricsList *[]metrics, prefix string, namespace string, id string) {

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
				metrics.headers[sanitizeName(name)] = fmt.Sprintf(format, val)
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
		*metricsList = append(*(metricsList), metrics)
	}
}

func (r *response) makeResponse() {
	r.headerNames = make(map[string]string)
}

type defaultResponse struct {
	response
}

func newDefaultResponse() *defaultResponse {
	r := new(defaultResponse)
	r.makeResponse()

	return r
}

type devdetailsResponse struct {
	response
}

func newDevdetailsResponse() *devdetailsResponse {
	r := new(devdetailsResponse)
	r.makeResponse()

	r.headerNames["Name"] = "%s"
	r.headerNames["ID"] = "%.0f"
	r.headerNames["DEVDETAILS"] = "%.0f"
	r.headerNames["Model"] = "%s"
	r.headerNames["Kernel"] = "%s"
	r.headerNames["Driver"] = "%s"

	return r
}

type devsResponse struct {
	response
}

func newDevsResponse() *devsResponse {
	r := new(devsResponse)
	r.makeResponse()

	r.headerNames["Name"] = "%s"
	r.headerNames["ID"] = "%.0f"

	return r
}

type statsResponse struct {
	response
}

func newStatsResponse() *statsResponse {
	r := new(statsResponse)
	r.makeResponse()

	r.headerNames["STATS"] = "%.0f"
	r.headerNames["ID"] = "%s"

	return r
}

type poolsResponse struct {
	response
}

func newPoolsResponse() *poolsResponse {
	r := new(poolsResponse)
	r.makeResponse()

	r.headerNames["POOL"] = "%.0f"
	r.headerNames["URL"] = "%s"

	return r
}

type coinResponse struct {
	response
}

func newCoinResponse() *coinResponse {
	r := new(coinResponse)
	r.makeResponse()

	r.headerNames["Hash Method"] = "%s"

	return r
}

type notifyResponse struct {
	response
}

func newNotifyResponse() *notifyResponse {
	r := new(notifyResponse)
	r.makeResponse()

	r.headerNames["Name"] = "%s"
	r.headerNames["ID"] = "%.0f"
	r.headerNames["NOTIFY"] = "%.0f"

	return r
}

type procsResponse struct {
	response
}

func newProcsResponse() *procsResponse {
	r := new(procsResponse)
	r.makeResponse()

	r.headerNames["Name"] = "%s"
	r.headerNames["ID"] = "%.0f"
	r.headerNames["PGA"] = "%.0f"

	return r
}
