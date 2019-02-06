package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/jumincorp/constrictor"
	"github.com/jumincorp/micrometrics"
)

const (
	programName = "bfx"
)

var (
	label             = constrictor.StringVar("label", "l", "default", "Label to identify this miner's data")
	miner             = constrictor.AddressPortVar("miner", "m", ":4028", "Address:Port of the miner's RPC port")
	prometheusAddress = constrictor.AddressPortVar("prometheus", "p", ":40010", "Address:Port to expose to Prometheus")
	queryDelay        = constrictor.TimeDurationVar("time", "t", "30", "Delay between RPC calls to the miner")

	exporter micrometrics.Exporter
)

func init() {
	constrictor.App("bfx", "bfgminer metrics", "Export bfgminer metrics")

	log.Printf("miner %s prometheus %s\n", miner(), prometheusAddress())

	exporter = micrometrics.NewPrometheusExporter(prometheusAddress())
}

type rpcCommand struct {
	Command   string `json:"command"`
	Parameter string `json:"parameter"`
}

func sendCommand(command string) (net.Conn, error) {
	conn, err := net.Dial("tcp", miner())
	if err == nil {
		var msg []byte
		msg, err := json.Marshal(rpcCommand{Command: command})
		if err == nil {
			fmt.Fprintf(conn, string(msg))
		}
	}
	return conn, err
}

func gatherCommand(metricsList *[]micrometrics.Metric, command string) {
	conn, err := sendCommand(command)
	if err == nil {

		resp, _ := ioutil.ReadAll(conn)
		log.Printf("-------------------------------------\n")
		log.Printf(" %v\n", command)
		log.Printf("-------------------------------------\n")
		r := newResponse(command, resp)

		r.getMetrics(metricsList, programName, programName, label())

	} else {
		log.Printf("Error sending command to miner: %v\n", err)
	}
}

func gather() {
	var metricsList = make([]micrometrics.Metric, 0)

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

	exporter.Export(metricsList)
}

func main() {
	go func() {
		for {
			gather()
			time.Sleep(queryDelay())
		}
	}()
	exporter.Setup()
}
