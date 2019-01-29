package main

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	cfgQueryDelay = "QueryDelay"
)

type config struct {
	miner      *minerConfig
	prometheus *prometheusConfig
}

func newConfig(name string) *config {
	cfg := new(config)

	viper.SetConfigName(name)
	viper.AddConfigPath(fmt.Sprintf("/etc/"))

	cfg.prometheus = newPrometheusConfig()
	cfg.miner = newMinerConfig()

	viper.SetDefault(cfgQueryDelay, 15)

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Printf("error reading config file: %s", err)
	}
	return cfg
}

// QueryDelay returns the time we wait to interrogate the miner again
func (cfg *config) queryDelay() time.Duration {
	return time.Duration(viper.Get(cfgQueryDelay).(int))
}

const (
	cfgMinerAddress = "Miner.Address"
	cfgMinerID      = "Miner.Id"
)

// Miner represents the Miner section of the configuration
type minerConfig struct {
}

func newMinerConfig() *minerConfig {
	minerConfig := new(minerConfig)

	viper.SetDefault(cfgMinerAddress, ":4028")
	viper.SetDefault(cfgMinerID, "default")

	return minerConfig
}

func (*minerConfig) address() string {
	return viper.Get(cfgMinerAddress).(string)
}

func (*minerConfig) id() string {
	return viper.Get(cfgMinerID).(string)
}

const (
	cfgPrometheusAddress = "Prometheus.Address"
)

// Prometheus represents the Prometheus section of the configuration
type prometheusConfig struct {
}

func newPrometheusConfig() *prometheusConfig {
	prometheusConfig := new(prometheusConfig)
	viper.SetDefault(cfgPrometheusAddress, ":40010")

	return prometheusConfig
}

func (*prometheusConfig) address() string {
	return viper.Get(cfgPrometheusAddress).(string)
}
