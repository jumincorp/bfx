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

type metric struct {
	labels map[string]string
	name   string
	value  float64
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

func gatherCommand(metricsList *[]metric, command string) {
	conn, err := sendCommand(command)
	if err == nil {

		resp, _ := ioutil.ReadAll(conn)
		log.Printf("-------------------------------------\n")
		log.Printf(" %v\n", command)
		log.Printf("-------------------------------------\n")
		r := newResponse(command, resp)

		r.getMetrics(metricsList, programName, programName, cfg.miner.id())

	} else {
		log.Printf("Error sending command to miner: %v\n", err)
	}
}

func gather() {
	var metricsList = make([]metric, 0)

	gatherCommand(&metricsList, "devs")
	gatherCommand(&metricsList, "devdetails")
	gatherCommand(&metricsList, "summary")
	gatherCommand(&metricsList, "pools")
	gatherCommand(&metricsList, "stats")
	gatherCommand(&metricsList, "coin")
	gatherCommand(&metricsList, "procs")
	gatherCommand(&metricsList, "notify")

	//gatherCommand(&metricsList, "version")
	//gatherCommand(&metricsList, "config")

	export.export(metricsList)
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
