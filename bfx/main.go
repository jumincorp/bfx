package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"../config"
	"../export"
)

const programName = "bfx"

var (
	minerGpuHashRate = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "miner_gpu_hashrate",
			Help: "Current hash rate of a GPU.",
		},
		[]string{"namespace", "miner", "gpu", "symbol"},
	)

	cfg      *config.Config
	exporter export.Exporter
)

func init() {
	cfg = config.NewConfig(programName)
	exporter = export.NewPrometheus(cfg.Prometheus.Address())

	// Metrics have to be registered to be exposed:
	prometheus.MustRegister(minerGpuHashRate)

}

type rpcCommand struct {
	Command   string `json:"command"`
	Parameter string `json:"parameter"`
}

func sendCommand(command string) (net.Conn, error) {
	conn, err := net.Dial("tcp", cfg.Miner.Address())
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
		//log.Printf("r %v\n", r)

		r.export()

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
			time.Sleep(time.Second * cfg.QueryDelay())
		}
	}()
	exporter.Setup()
}
