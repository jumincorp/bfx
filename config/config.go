package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

func init() {
	log.Printf("config init")
}

const (
	cfgPrometheusAddress = "Prometheus.Address"
)

// Prometheus represents the Prometheus section of the configuration
type Prometheus struct {
}

func newPrometheus() *Prometheus {
	prometheus := new(Prometheus)
	viper.SetDefault(cfgPrometheusAddress, ":40010")

	return prometheus
}

func (prometheus *Prometheus) Address() string {
	return viper.Get(cfgPrometheusAddress).(string)
}

const (
	cfgMinerAddress = "Miner.Address"
	cfgMinerID      = "Miner.Id"
	cfgMinerProgram = "Miner.Program"
	cfgMinerSymbol  = "Miner.Symbol"
)

// Miner represents the Miner section of the configuration
type Miner struct {
}

func newMiner() *Miner {
	miner := new(Miner)

	viper.SetDefault(cfgMinerAddress, ":4028")
	viper.SetDefault(cfgMinerID, "default")
	viper.SetDefault(cfgMinerProgram, "EQBminer")
	viper.SetDefault(cfgMinerSymbol, "EQB")

	return miner
}

func (miner *Miner) Address() string {
	return viper.Get(cfgMinerAddress).(string)
}

func (miner *Miner) ID() string {
	return viper.Get(cfgMinerID).(string)
}

func (miner *Miner) Program() string {
	return viper.Get(cfgMinerProgram).(string)
}

func (miner *Miner) Symbol() string {
	return viper.Get(cfgMinerSymbol).(string)
}

const (
	cfgQueryDelay = "QueryDelay"
)

// Config represents the configuration file. You should use NewConfig to create one.
type Config struct {
	name       string
	Miner      *Miner
	Prometheus *Prometheus
	queryDelay time.Duration
}

// NewConfig creates an instance of the configuration
func NewConfig(name string) *Config {
	cfg := new(Config)
	cfg.name = name

	cfg.Prometheus = newPrometheus()
	cfg.Miner = newMiner()

	viper.SetDefault(cfgQueryDelay, 15)

	viper.SetConfigName(name)
	viper.AddConfigPath(fmt.Sprintf("/etc/%v", name))

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Printf("error reading config file: %s", err)
	}
	return cfg
}

// Name returns the configuration name
func (cfg *Config) Name() string {
	return cfg.name
}

// QueryDelay returns the time we wait to interrogate the miner again
func (cfg *Config) QueryDelay() time.Duration {
	return time.Duration(viper.Get(cfgQueryDelay).(int))
}

//// Prometheus return the Prometheus configuration section
//func (cfg *Config) Prometheus() *Prometheus {
//return cfg.prometheus
//}

//// Miner return the Miner configuration section
//func (cfg *Config) Miner() *Miner {
//return cfg.miner
//}
