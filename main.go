package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/miekg/dns"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var dnsVar string
var interval time.Duration

func main() {
	flag.StringVar(&dnsVar, "dnstargets", "", "comma separated list of dns servers which will be queried")
	flag.DurationVar(&interval, "interval", 5*time.Second, "interval between measurements. default: 5s")
	flag.Parse()

	dnsTargets := strings.Split(dnsVar, ",")
	if dnsVar == "" || len(dnsTargets) == 0 {
		log.Fatal("need at least one target via -dnstargets")
	}
	log.Printf("querying %d targets: %s", len(dnsTargets), dnsTargets)

	targetVec := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "dnsquery",
			Help: "dns query time in s",
			Objectives: map[float64]float64{
				0.5:  0.01,
				0.9:  0.01,
				0.99: 0.01,
			},
		},
		[]string{"target", "success"},
	)

	for _, target := range dnsTargets {
		go measure(targetVec, strings.TrimSpace(target))
	}

	reg := prometheus.NewRegistry()
	reg.Register(targetVec)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func measure(vec *prometheus.SummaryVec, target string) {
	client := dns.Client{
		ReadTimeout: time.Duration(time.Minute),
	}
	for {
		var msg dns.Msg
		msg.SetQuestion(dns.Fqdn("."), dns.TypeNS)
		msg.RecursionDesired = true

		addr := net.JoinHostPort(target, "53")
		_, rtt, err := client.Exchange(&msg, addr)

		success := "true"
		if err != nil {
			success = "false"
		}
		vec.WithLabelValues(target, success).Observe(rtt.Seconds())

		time.Sleep(interval)
	}
}
