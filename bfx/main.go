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

		for _, metrics := range r.getMetrics(programName, programName, cfg.miner.id()) {
			export.export(metrics)
		}

		//for _, data := range r.data {
		//	log.Printf("data MHS rolling %v", data["MHS rolling"])
		//}

		//for i, device := range resp.DEVS {
		//log.Printf("%v Device %v %v Hashrate %v\n", i, device.Name, device.ID, device.MHS20S)

		//minerGpuHashRate.With(prometheus.Labels{
		//"namespace": programName,
		//"miner":     cfg.Miner.Program(),
		//"gpu":       fmt.Sprintf("GPU%d", device.ID),
		//"symbol":    cfg.Miner.Symbol(),
		//}).Set(device.MHS20S)
		//}

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
