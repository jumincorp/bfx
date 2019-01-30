package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

const (
	cfgQueryDelay = "time"
	cfgID         = "id"
	cfgPrometheus = "prometheus"
	cfgMiner      = "miner"
)

func readConfig(name string) {
	//cfg := new(config)

	viper.SetConfigName(name)
	//viper.AddConfigPath(fmt.Sprintf("/etc/"))
	viper.AddConfigPath(fmt.Sprintf("."))

	viper.SetEnvPrefix(name)
	viper.AutomaticEnv()

	viper.SetDefault(cfgID, "default")
	viper.SetDefault(cfgMiner, ":4028")
	viper.SetDefault(cfgPrometheus, ":40010")
	viper.SetDefault(cfgQueryDelay, 15)

	err := viper.ReadInConfig()
	if err != nil { // Handle errors reading the config file
		fmt.Printf("error reading config file: %s", err)
	}

	log.Printf("cfg id %v\n", id())
	log.Printf("cfg miner %v\n", miner())
	log.Printf("cfg prometheus %v\n", prometheus())
	log.Printf("cfg time %v\n", queryDelay())
}

func queryDelay() time.Duration {
	if delay, ok := viper.Get(cfgQueryDelay).(int); ok {
		return time.Duration(time.Duration(delay) * time.Second)
	}
	return viper.GetDuration(cfgQueryDelay)
}

func miner() string {
	return viper.GetString(cfgMiner)
}

func id() string {
	return viper.GetString(cfgID)
}

func prometheus() string {
	return viper.GetString(cfgPrometheus)
}
