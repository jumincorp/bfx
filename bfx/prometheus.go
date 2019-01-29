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
	address string
	exporter
}

func init() {
	//for _, c := range collectors {
	//prometheus.MustRegister(c)
	//}

	// // Metrics have to be registered to be exposed:
	//prometheus.MustRegister(minerGpuHashRate)

}

func newPrometheusExporter(address string) *prometheusExporter {
	p := new(prometheusExporter)
	p.address = address
	return p
}

func handleWithCare() http.HandlerFunc {
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
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	//http.Handle("/metrics", promhttp.Handler())
	http.Handle("/metrics", handleWithCare())
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
		//log.Printf("header %s %s\n", keys[k], headers[keys[k]])
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
			//log.Printf("metric %s\n", formattedMetrics[metricIndex])
			metricIndex++
		}
	}
	return err
}
