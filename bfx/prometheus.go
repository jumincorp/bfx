package main

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

var (
	mutex            = &sync.Mutex{}
	formattedMetrics = make([]string, 0)
)

type prometheusExporter struct {
	exporter
	address string
}

func newPrometheusExporter(address string) *prometheusExporter {
	p := new(prometheusExporter)
	p.address = address
	return p
}

func httpHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		defer mutex.Unlock()

		for _, metric := range formattedMetrics {
			w.Write([]byte(metric))
			w.Write([]byte("\n"))
		}
	}
}

func (p *prometheusExporter) setup() {
	http.Handle("/metrics", httpHandler())
	log.Fatal(http.ListenAndServe(p.address, nil))
}

func formatHeaders(headers map[string]string) string {
	// Sort Header keys
	var sb strings.Builder

	keys := make([]string, len(headers))
	i := 0
	for k := range headers {
		keys[i] = k
		i++
	}
	sort.Strings(keys)

	sb.WriteString("{")

	for k := range keys {
		if k != 0 {
			sb.WriteString(",")
		}
		sb.WriteString(keys[k])
		sb.WriteString("=\"")
		sb.WriteString(headers[keys[k]])
		sb.WriteString("\"")
	}

	sb.WriteString("}")
	return sb.String()
}

func (p *prometheusExporter) export(list []metrics) error {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	var numberOfMetrics int
	for _, m := range list {
		numberOfMetrics += len(m.values)
	}

	formattedMetrics = make([]string, numberOfMetrics)

	var metricIndex int
	for _, m := range list {
		formattedHeaders := formatHeaders(m.headers)

		for k, v := range m.values {
			formattedMetrics[metricIndex] = fmt.Sprintf("%s%s %f", k, formattedHeaders, v)
			metricIndex++
		}
	}
	return err
}
