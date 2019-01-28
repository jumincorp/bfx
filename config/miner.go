package config

import "github.com/spf13/viper"

const (
	cfgMinerAddress = "Miner.Address"
	cfgMinerID      = "Miner.Id"
)

// Miner represents the Miner section of the configuration
type MinerConfig struct {
}

func newMinerConfig() *MinerConfig {
	minerConfig := new(MinerConfig)

	viper.SetDefault(cfgMinerAddress, ":4028")
	viper.SetDefault(cfgMinerID, "default")

	return minerConfig
}

func (*MinerConfig) Address() string {
	return viper.Get(cfgMinerAddress).(string)
}

func (*MinerConfig) ID() string {
	return viper.Get(cfgMinerID).(string)
}
