package main

import (
	"fmt"
	"github.com/RobbieMcKinstry/hashicorp-presentation/client/cluster"
	"github.com/RobbieMcKinstry/hashicorp-presentation/client/display"
	"time"
)

type ServiceData struct {
	Address    string
	Stream     chan<- float64
	Alive      bool
	Initalized bool
}

type LoadBalancer struct {
	cluster          cluster.Cluster
	services         []*ServiceData
	load             uint64
	deployedServices int

	changeLoad func(uint64)
}

// Construct a new LoadBalancer.
func NewLoadBalancer(display *display.Display) *LoadBalancer {
	var lb = &LoadBalancer{
		cluster:          cluster.NewMockCluster(),
		deployedServices: 0,
		services: []*ServiceData{{
			Address:    "",
			Stream:     display.GetSparklineStream(0),
			Alive:      true,
			Initalized: false,
		}, {
			Address:    "",
			Stream:     display.GetSparklineStream(1),
			Alive:      true,
			Initalized: false,
		}, {
			Address:    "",
			Stream:     display.GetSparklineStream(2),
			Alive:      true,
			Initalized: false,
		}},
	}
	lb.changeLoad = func(load uint64) {
		display.SetLoad(load)
		lb.setLoad(load)
	}

	go lb.watchServices()

	return lb
}

func (lb *LoadBalancer) watchServices() {
	// Every second, iterate through each service with the load
	// they *would* receive if it were evenly dispersed.
	var ticker = time.Tick(1 * time.Second)
	for {
		select {
		case <-ticker:
			// Every second, we take the current current load and divide it up.
			var loadPerServer = lb.load / uint64(lb.countAlive())
			// Now, each server gets sent this load, and the throughput
			// is written to its output channel.
			for _, serv := range lb.services {
				if !serv.Initalized || !serv.Alive {
					continue
				}
				// Send it a request to validate if its alive or dead.
			}
		}
	}
}

func (lb *LoadBalancer) SetLoad(load uint64) {
	lb.changeLoad(load)
}

func (lb *LoadBalancer) setLoad(load uint64) {
	lb.load = load
}

// This function is called when the user enter "new service" into the
// textbox.
//
// 1. Check how many services we have. If we have 3 already, exit.
// 2. Next, we ask the Cluster for a new service.
// 3. Using the cluster response, we add this server address to our list
//    of addresses and ping them once a second.
func (lb *LoadBalancer) OnNewService(throughput, soft, hard uint64) {
	if lb.tooManyServices() {
		ExitOnError(fmt.Errorf("Cannot support more than 3 services."))
	}
	var resp = lb.cluster.NewService(throughput, soft, hard)
	var serviceIndex = lb.totalDeployed() + 1
	// Add this address to our address list.
	lb.services[serviceIndex].Address = resp.Address
	lb.services[serviceIndex].Initalized = true

	lb.deployedServices++
}

func (lb *LoadBalancer) tooManyServices() bool {
	return lb.totalDeployed() >= 3
}

func (lb *LoadBalancer) totalDeployed() int {
	return lb.deployedServices
}

func (lb *LoadBalancer) countAlive() int {
	var count int
	for _, serv := range lb.services {
		if serv.Alive {
			count++
		}
	}
	return count
}
