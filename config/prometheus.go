package config

import "github.com/spf13/viper"

const (
	cfgPrometheusAddress = "Prometheus.Address"
)

// Prometheus represents the Prometheus section of the configuration
type PrometheusConfig struct {
}

func newPrometheusConfig() *PrometheusConfig {
	prometheusConfig := new(PrometheusConfig)
	viper.SetDefault(cfgPrometheusAddress, ":40010")

	return prometheusConfig
}

func (*PrometheusConfig) Address() string {
	return viper.Get(cfgPrometheusAddress).(string)
}
