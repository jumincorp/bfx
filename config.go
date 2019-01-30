package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	programName = "bfx"

	cfgQueryDelay = "time"
	cfgLabel      = "label"
	cfgPrometheus = "prometheus"
	cfgMiner      = "miner"
)

var rootCmd = &cobra.Command{
	Use:   programName,
	Short: fmt.Sprintf("%v exports metrics from bfgminer", programName),
	Long:  `A simple way to get bfgminer information`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		run(args)
	},
}

func init() {
	cobra.OnInitialize(readConfig)

	rootCmd.PersistentFlags().StringP(cfgPrometheus, "p", ":40010", "Address:Port to expose to Prometheus")
	viper.BindPFlag(cfgPrometheus, rootCmd.PersistentFlags().Lookup(cfgPrometheus))

	rootCmd.PersistentFlags().StringP(cfgMiner, "m", ":4028", "Address:Port of the miner's RPC port")
	viper.BindPFlag(cfgMiner, rootCmd.PersistentFlags().Lookup(cfgMiner))

	rootCmd.PersistentFlags().StringP(cfgQueryDelay, "t", "30", "Delay between RPC calls to the miner")
	viper.BindPFlag(cfgQueryDelay, rootCmd.PersistentFlags().Lookup(cfgQueryDelay))

	rootCmd.PersistentFlags().StringP(cfgLabel, "l", "default", "Label to identify this miner's data")
	viper.BindPFlag(cfgLabel, rootCmd.PersistentFlags().Lookup(cfgLabel))
}

func readConfig() {
	//cfg := new(config)

	viper.SetConfigName(programName)
	//viper.AddConfigPath(fmt.Sprintf("/etc/"))
	viper.AddConfigPath(fmt.Sprintf("."))

	viper.SetEnvPrefix(programName)
	viper.AutomaticEnv()

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
	if i, err := strconv.ParseInt(viper.GetString(cfgMiner), 10, 64); err == nil {
		return fmt.Sprintf(":%d", i)
	}
	return viper.GetString(cfgMiner)
}

func id() string {
	return viper.GetString(cfgLabel)
}

func prometheus() string {
	if i, err := strconv.ParseInt(viper.GetString(cfgPrometheus), 10, 64); err == nil {
		return fmt.Sprintf(":%d", i)
	}
	return viper.GetString(cfgPrometheus)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
