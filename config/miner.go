package config

import "github.com/spf13/viper"

const (
	cfgMinerAddress = "Miner.Address"
	cfgMinerID      = "Miner.Id"
)

// Miner represents the Miner section of the configuration
type Miner struct {
}

func newMiner() *Miner {
	miner := new(Miner)

	viper.SetDefault(cfgMinerAddress, ":4028")
	viper.SetDefault(cfgMinerID, "default")

	return miner
}

func (miner *Miner) Address() string {
	return viper.Get(cfgMinerAddress).(string)
}

func (miner *Miner) ID() string {
	return viper.Get(cfgMinerID).(string)
}
