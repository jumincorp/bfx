# bfx
[![Build Status](https://travis-ci.org/JUMINCORP/bfx.svg?branch=master)](https://travis-ci.org/JUMINCORP/bfx)

bfx exports bfgminer metrics to Prometheus.

## Building

If you have Go installed:
`go get github.com/JUMINCORP/bfx` should get the application to your $GOBIN directory.


## Usage

Sample usage: 

`bfx --label test_data --miner localhost:4028 --prometheus :50020 --time 5`

This would get the data from the bfgminer RPC port on localhost at port 4028, and send expose it to Prometheus on port 50020 (localhost implied). Data collection will be done every 5 seconds and every metric would have the label miner='test_data'

It is also possible to send options to the program using environment variables. In this case the equivalent configuration would be

```
export BFX_LABEL="test_data"
export BFX_MINER="localhost:4028"
export BFX_PROMETHEUS=":50020"
export BFX_TIME="5"
bfx
```

## Sample Output
[Sample Output](https://github.com/JUMINCORP/bfx/wiki/Sample-output) from a 4-GPU system
