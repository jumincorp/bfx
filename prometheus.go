package main

import (
	"log"
	"net/http"
	"sort"
	"strconv"
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

func formatMetric(m metric) string {
	var sb strings.Builder

	sb.WriteString(m.name)

	sortedLabels := make([]string, len(m.labels))
	i := 0
	for k := range m.labels {
		sortedLabels[i] = k
		i++
	}
	sort.Strings(sortedLabels)

	sb.WriteRune('{')
	for i := range sortedLabels {
		if i != 0 {
			sb.WriteRune(',')
		}
		sb.WriteString(sortedLabels[i])
		sb.WriteString("=\"")
		sb.WriteString(m.labels[sortedLabels[i]])
		sb.WriteRune('"')
	}
	sb.WriteString("} ")

	sb.WriteString(strconv.FormatFloat(m.value, 'f', -1, 64))

	return sb.String()
}

func (p *prometheusExporter) export(metrics []metric) error {
	var err error

	mutex.Lock()
	defer mutex.Unlock()

	formattedMetrics = make([]string, len(metrics))

	for i, m := range metrics {
		//formattedLabels := formatLabels(m.labels)
		formattedMetrics[i] = formatMetric(m)
	}
	return err
}
