package main

import (
	"encoding/json"
	"fmt"
	"github.com/RobbieMcKinstry/hashicorp-presentation/server/messages"
	"net/http"
	"os"
	"time"
)

func main() {
	var lb = &LoadBalancer{}
	var systemLoad = 3000
	// Spin up 3 services.
	// Collect their addresses.

	// Start watching throughput
	go lb.watchThroughput(systemLoad)
	// Spin up one neighbor.
}

type Service struct {
	Address string
}

type LoadBalancer struct {
	services []*Service
}

func (lb *LoadBalancer) AddService() {
}

func (lb *LoadBalancer) watchThroughput(load int) {
	// Iterate through each server, retaining only the living ones.
	var ticker = time.Tick(1 * time.Second)
	var candidateServices = lb.services
	var disqualifiedServices = []*Service{}
	var candidateLoad = load / len(candidateServices)
	for range ticker {
		// If they fail the health check, they're removed this time, add back next round.
		var healthyServ, unhealthyServ = healthCheck(candidateLoad, candidateServices)
		var actualLoad = load / len(healthyServ)
		// Now, distribute this load across the healthy services.
		var stableServ, unstable = requestThroughput(actualLoad, healthyServ)

		if len(stableServ) == 0 {
			fmt.Println("Simulation Ended. System crash.")
			os.Exit(0)
		}
		candidateServices = append(stableServ, unhealthyServ...)
		candidateServices = append(candidateServices, disqualifiedServices...)
		disqualifiedServices = unstable
	}
}

func requestThroughput(load int, services []*Service) (stable, unstable []*Service) {
	for _, service := range services {
		var addr = fmt.Sprintf("http://%v/metrics/throughput?load=%v", service.Address, load)
		resp, _ := http.Get(addr)
		defer resp.Body.Close()
		var throughputResp = messages.ThroughputGETResponse{}
		json.NewDecoder(resp.Body).Decode(&throughputResp)
		if throughputResp.Throughput != 0 {
			stable = append(stable, service)
		} else {
			unstable = append(unstable, service)
		}
		fmt.Printf("Throughput for server %v: %v", service.Address, throughputResp.Throughput)
	}
	return stable, unstable
}

// Filter out any services which would be killed by this load.
func healthCheck(load int, instances []*Service) (healthy, unhealthy []*Service) {
	for _, serv := range instances {
		var addr = fmt.Sprintf("http://%v/healthz?load=%v", serv.Address, load)
		resp, err := http.Get(addr)
		if err != nil && resp.StatusCode == 200 {
			healthy = append(healthy, serv)
		} else {
			unhealthy = append(unhealthy, serv)
		}
	}
	return healthy, unhealthy
}
