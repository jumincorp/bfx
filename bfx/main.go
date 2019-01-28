package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"
)

const programName = "bfx"

var (
	cfg    *config
	export *prometheusExporter
)

func init() {
	cfg = newConfig(programName)
	export = newPrometheusExporter(cfg.prometheus.address())
}

type metrics struct {
	headers map[string]string
	values  map[string]float64
}

type rpcCommand struct {
	Command   string `json:"command"`
	Parameter string `json:"parameter"`
}

func sendCommand(command string) (net.Conn, error) {
	conn, err := net.Dial("tcp", cfg.miner.address())
	if err == nil {
		var msg []byte
		msg, err := json.Marshal(rpcCommand{Command: command})
		if err == nil {
			fmt.Fprintf(conn, string(msg))
		}
	}
	return conn, err
}

func gatherCommand(command string) {
	conn, err := sendCommand(command)
	if err == nil {

		resp, _ := ioutil.ReadAll(conn)
		log.Printf("-------------------------------------\n")
		log.Printf(" %v\n", command)
		log.Printf("-------------------------------------\n")
		r := newResponse(command, resp)

		metrics := r.getMetrics(programName, programName, cfg.miner.id())
		export.export(metrics)

	} else {
		log.Printf("Error sending command to miner: %v\n", err)
	}
}

func gather() {
	gatherCommand("devs")
}

func main() {
	go func() {
		for {
			gather()
			time.Sleep(time.Second * cfg.queryDelay())
		}
	}()
	export.setup()
}
