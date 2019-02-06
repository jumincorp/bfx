package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/jumincorp/micrometrics"
)

type responseIface interface {
	append(map[string]interface{})
	getMetrics(metricsList *[]micrometrics.Metric, prefix string, namespace string, id string)
	setCommandName(string)
}

type response struct {
	command  string
	data     []map[string]interface{}
	labelFmt map[string]string
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

	delete(responseData, "STATUS")
	delete(responseData, "id")

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

	for _, respElement := range respList {
		r.append(respElement.(map[string]interface{}))
	}

	return r
}

func (r *response) getMetrics(metricsList *[]micrometrics.Metric, prefix string, namespace string, id string) {

	for _, element := range r.data {
		log.Printf("---")
		var m micrometrics.Metric

		m.Labels = make(map[string]string)
		m.Labels["namespace"] = namespace
		m.Labels["miner"] = id

		for name, val := range element {
			if format, isLabel := r.labelFmt[name]; isLabel {
				m.Labels[sanitizeName(name)] = fmt.Sprintf(format, val)
			}
		}

		for name, val := range element {
			if _, isLabel := r.labelFmt[name]; !isLabel {
				m.Name = r.metricName(prefix, name)
				switch casted := val.(type) {
				case int64:
					log.Printf("i! %s %d\n", name, casted)
				case float64:
					//log.Printf("f  %s >> %s %f\n", name, metricName, casted)
					m.Value = casted
				case string:
					log.Printf("s! %s %s\n", name, casted)
				default:
					log.Printf("?! %s %s\n", name, casted)
				}
				*metricsList = append(*(metricsList), m)
			}
		}
	}
}

func (r *response) makeResponse() {
	r.labelFmt = make(map[string]string)
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

	r.labelFmt["Name"] = "%s"
	r.labelFmt["ID"] = "%.0f"
	r.labelFmt["DEVDETAILS"] = "%.0f"
	r.labelFmt["Model"] = "%s"
	r.labelFmt["Kernel"] = "%s"
	r.labelFmt["Driver"] = "%s"

	return r
}

type devsResponse struct {
	response
}

func newDevsResponse() *devsResponse {
	r := new(devsResponse)
	r.makeResponse()

	r.labelFmt["Name"] = "%s"
	r.labelFmt["ID"] = "%.0f"

	return r
}

type statsResponse struct {
	response
}

func newStatsResponse() *statsResponse {
	r := new(statsResponse)
	r.makeResponse()

	r.labelFmt["STATS"] = "%.0f"
	r.labelFmt["ID"] = "%s"

	return r
}

type poolsResponse struct {
	response
}

func newPoolsResponse() *poolsResponse {
	r := new(poolsResponse)
	r.makeResponse()

	r.labelFmt["POOL"] = "%.0f"
	r.labelFmt["URL"] = "%s"

	return r
}

type coinResponse struct {
	response
}

func newCoinResponse() *coinResponse {
	r := new(coinResponse)
	r.makeResponse()

	r.labelFmt["Hash Method"] = "%s"

	return r
}

type notifyResponse struct {
	response
}

func newNotifyResponse() *notifyResponse {
	r := new(notifyResponse)
	r.makeResponse()

	r.labelFmt["Name"] = "%s"
	r.labelFmt["ID"] = "%.0f"
	r.labelFmt["NOTIFY"] = "%.0f"

	return r
}

type procsResponse struct {
	response
}

func newProcsResponse() *procsResponse {
	r := new(procsResponse)
	r.makeResponse()

	r.labelFmt["Name"] = "%s"
	r.labelFmt["ID"] = "%.0f"
	r.labelFmt["PGA"] = "%.0f"

	return r
}
