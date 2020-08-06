package main

import (
	"encoding/json"
	"fmt"
	"html"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	startingThroughput = 1000
	startingSoft       = 1500
	startingHard       = 2000
)

func main() {
	var service = NewSimulatedService(startingThroughput, startingSoft, startingHard)
	var port = ":8080"
	var server = &http.Server{
		Addr:           port,
		Handler:        service,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	fmt.Printf("Listening on port %v\n", port)
	server.ListenAndServe()
}

// A SimulatedService is an HTTP server which simulates HTTP traffic.
// It can:
// 1. Respond to liveliness checks.
// 2. Provive a simulated throughput metric.
// 3. Modify its throughput with requests from CPU-stealing noisy neighbors.
// HTTP API:
// GET  /healthz -> return the server status.
// GET  /metrics/throughput -> return the number of requests handled in the last second
// POST /server-state -> edit the max throughput, soft limit, or hard limit.
type SimulatedService struct {
	MaxThroughput,
	RequestSoftLimit,
	RequestHardLimit uint64
}

// NewSNewSimulatedService is the constructor for a SimulatedService.
func NewSimulatedService(maxThroughput, softLimit, hardLimit uint64) *SimulatedService {
	return &SimulatedService{
		MaxThroughput:    maxThroughput,
		RequestSoftLimit: softLimit,
		RequestHardLimit: hardLimit,
	}
}

// ServServeHTTP fulfills the http.Handler interface.
func (service *SimulatedService) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Switch on the URL:
	// Forward the handler func for each URL
	var path = html.EscapeString(req.URL.Path)
	switch {
	case strings.HasPrefix(path, "/healthz"):
		service.handleHealthCheck(w, req)
	case strings.HasPrefix(path, "/metrics/throughput"):
		service.handleThroughput(w, req)
	default:
		http.Error(w, "No such route.", http.StatusNotFound)
	}
}

func (service *SimulatedService) handleHealthCheck(w http.ResponseWriter, req *http.Request) {
	var load, err = service.getLoad(req)
	if err != nil {
		fmt.Fprintf(w, "Error parsing load param: %v", html.EscapeString(err.Error()))
		return
	}
	var encoder = json.NewEncoder(w)
	var alive = service.IsAlive(load)
	var responseBody = HealthCheckResponse{Alive: alive}
	err = encoder.Encode(responseBody)
	if err != nil {
		http.Error(w, "Error when writing response.", http.StatusInternalServerError)
	}
}

// getLoad returns the value of the load provided to this server within the
// last second. It's fetched from the request's URL parameters.
func (service *SimulatedService) getLoad(req *http.Request) (uint64, error) {
	// Fetch the simulated load.
	var load = req.FormValue("load")
	return strconv.ParseUint(load, 10, 64)
}

func (service *SimulatedService) handleThroughput(w http.ResponseWriter, req *http.Request) {
	// Are we changinge the throughput or requesting it?
	switch req.Method {
	case http.MethodGet:
		service.handleThroughputGET(w, req)
	case http.MethodPost:
		service.handleThroughputPOST(w, req)
	}
}

func (service *SimulatedService) handleThroughputGET(w http.ResponseWriter, req *http.Request) {
	var load, err = service.getLoad(req)
	if err != nil {
		fmt.Fprintf(w, "Error parsing load param: %v", html.EscapeString(err.Error()))
		return
	}
	// Now, reply with the throughput.
	var throughput = service.CalculateThroughput(load)
	var encoder = json.NewEncoder(w)
	var responseBody = ThroughputGETResponse{Throughput: throughput}
	err = encoder.Encode(responseBody)
	if err != nil {
		http.Error(w, "Error when writing response.", http.StatusInternalServerError)
	}
}

func (service *SimulatedService) CalculateThroughput(load uint64) uint64 {
	var throughput uint64
	if load < service.AvailableThroughput() {
		throughput = load
	} else if load <= service.ModifiedSoftLimit() {
		throughput = service.AvailableThroughput()
	} else if load <= service.ModifiedHardLimit() {
		throughput = service.degradedThroughput(load)
	} else {
		throughput = 0
	}

	return throughput
}

// degradedThroughput returns a value between 50% and 75% of the available throughput.
func (service *SimulatedService) degradedThroughput(load uint64) uint64 {
	var total = service.AvailableThroughput()
	var offset = total / 2
	var rngBound = total / 4
	var rng = rand.Int63n(int64(rngBound))
	return offset + uint64(rng)
}

func (service *SimulatedService) handleThroughputPOST(w http.ResponseWriter, req *http.Request) {
}

// IsAlive returns true if the server hasn't fallen over from too much load.
func (service *SimulatedService) IsAlive(load uint64) bool {
	return service.RequestHardLimit >= load
}

// AvailableCPU returns, as a fraction from 0 to 1, the amount of CPU
// available to this service. Noisy neighbors reduce the amount of CPU available.
func (service *SimulatedService) AvailableCPU() float64 {
	return 1.0
}

// AvailableThroughput returns the number of requests per second processable
// by this service. It's value is the max throughput modified by the available CPU.
// If noisy neighbors steal CPU, then the available CPU decreases.
func (service *SimulatedService) AvailableThroughput() uint64 {
	return service.scaleDown(service.MaxThroughput)
}

// ModifiedSoftLimit returns the new soft limit for this server once
// noisy neighbors have been accounted for. A service's soft limit is reduced
// by other services in the same pod interfering with it.
// If noisy neighbors steal CPU, then the available CPU decreases.
func (service *SimulatedService) ModifiedSoftLimit() uint64 {
	return service.scaleDown(service.RequestSoftLimit)
}

// MModifiedHardLimit returns the new hard limit for this server once
// noisy neighbors have been accounted for. A service's soft limit is reduced
// by other services in the same pod interfering with it.
// If noisy neighbors steal CPU, then the available CPU decreases.
func (service *SimulatedService) ModifiedHardLimit() uint64 {
	return service.scaleDown(service.RequestHardLimit)
}

// scaleDown takes the provided metric (throughput, soft limit, hard limit)
// and adjusts it to reflect the new limit provided by the noisy neighbor.
func (service *SimulatedService) scaleDown(metric uint64) uint64 {
	// Scale down the metric in proportion to CPU availability.
	var scaledMetric = float64(metric) * service.AvailableCPU()
	// Round, and then cast.
	return uint64(math.Round(scaledMetric))
}
