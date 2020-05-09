package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sparrc/go-ping"
)

var appVersion = "v1.0"

var (
	addr    = flag.String("listen", ":8080", "The address to listen on for HTTP requests.")
	version = flag.Bool("version", false, "Show version")
)

var (
	lostRate = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_package_lost_rate",
			Help: "The package loss rate during ping",
		},
		[]string{"server"},
	)
	delayTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_delay_time_ms",
			Help: "Round-trip delay time in milliseconds",
		},
		[]string{"server"},
	)
)

func startMonitor(servers []string) {
	for _, s := range servers {
		go startMonitor4Server(s)
	}
}

func startMonitor4Server(s string) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		pinger, err := ping.NewPinger(s)
		if err != nil {
			log.Printf("ping server %s error: %s\n", s, err.Error())
			lostRate.DeleteLabelValues(s)
			delayTime.DeleteLabelValues(s)
		} else {
			pinger.SetPrivileged(true)
			pinger.Count = 10
			pinger.Timeout = time.Second * 30
			pinger.Run()
			stat := pinger.Statistics()
			lostRate.WithLabelValues(s).Set(stat.PacketLoss)
			delayTime.WithLabelValues(s).Set(float64(stat.AvgRtt.Milliseconds()))
		}
		<-ticker.C
	}
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println("version", appVersion)
		os.Exit(1)
	}
	servers := flag.Args()
	if len(servers) == 0 {
		fmt.Fprintln(os.Stderr, "At least one remote server is required")
		os.Exit(1)
	}
	log.Printf("Monitor remote servers: %s\n", strings.Join(servers, ","))
	log.Printf("Metric HTTP server listen on %s\n", *addr)
	http.Handle("/metrics", promhttp.Handler())
	startMonitor(servers)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
