package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	collectors = make(map[string](*prometheus.GaugeVec))

	minerGpuHashRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "miner_gpu_hashrate",
			Help: "Current hash rate of a GPU.",
		},
		[]string{"namespace", "miner", "gpu", "symbol"},
	)
)

type prometheusExporter struct {
	address string
	exporter
}

func init() {
	//collectors[clock] =
	//prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "amdgpu_clock", Help: "GPU Clock Rate in MHz"}, []string{"gpu", "name"})

	//collectors[power] =
	//prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "amdgpu_power", Help: "GPU Power Consumption in Watts"}, []string{"gpu", "name"})

	//collectors[temp] =
	//prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "amdgpu_temp", Help: "GPU Temperature in Celcius"}, []string{"gpu"})

	//collectors[load] =
	//prometheus.NewGaugeVec(prometheus.GaugeOpts{Name: "amdgpu_load", Help: "GPU Load Percentage"}, []string{"gpu"})

	for _, c := range collectors {
		prometheus.MustRegister(c)
	}

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(minerGpuHashRate)

}

func newPrometheusExporter(address string) *prometheusExporter {
	p := new(prometheusExporter)
	p.address = address
	return p
}

func (p *prometheusExporter) setup() {
	// The Handler function provides a default handler to expose metrics
	// via an HTTP server. "/metrics" is the usual endpoint for that.
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(p.address, nil))
}

func (p *prometheusExporter) export(m metrics) error {
	var err error

	for k, v := range m.headers {
		log.Printf("header %s %s\n", k, v)
	}

	for k, v := range m.values {
		log.Printf("value %s %f\n", k, v)
	}
	//_, err := strconv.ParseFloat(value, 64)
	//if err == nil {
	////switch ctype {
	////case clock, power:
	////collectors[ctype].With(prometheus.Labels{"gpu": gpu, "name": name}).Set(fValue)
	////default:
	////collectors[ctype].With(prometheus.Labels{"gpu": gpu}).Set(fValue)
	////}
	//}
	return err
}
