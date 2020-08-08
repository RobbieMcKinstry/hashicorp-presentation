package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LifetimeKey = "LIFETIME"
	CPUKey      = "CPU"
	AddressKey  = "ADDRESSES"
)

func main() {
	// Fetch the duration for which this "noisy neighbor" will
	// steal CPU.
	var lifetimeStr = os.Getenv(LifetimeKey)
	// Fetch the amount of CPU this workload will steal.
	var cpuStr = os.Getenv(CPUKey)
	// Fetch the server addresses which this workload will steal from.
	var addressesStr = os.Getenv(AddressKey)

	if lifetimeStr == "" {
		log.Fatalf("Expected non-empty lifetime")
	}
	if cpuStr == "" {
		log.Fatalf("Expected non-empty CPU requirement")
	}
	if addressesStr == "" {
		log.Fatalf("Expected non-empty address list")
	}

	// var cpu = parseCPU(cpuStr)
	var lifetime = parseLifetime(lifetimeStr)
	var addresses = parseAddresses(addressesStr)

	// Now, ping each address and add this service as a neighbor.
	for _, addr := range addresses {
		var uri *url.URL
		var resp *http.Response
		var err error
		uri, err = url.Parse(addr)
		ExitOnError(err)
		uri.Path = "/neighbors/add"
		var params = uri.Query()
		params.Set("cpu", cpuStr)
		uri.RawQuery = params.Encode()
		resp, err = http.Get(uri.String())
		ExitOnError(err)
		if resp.StatusCode != 200 {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			ExitOnError(err)
			log.Fatal(body)
		}
	}
	// Now, this batch job sleeps for the specified duration.
	// The time slept represents the duration for which this process is working.
	// When we awake, we'll restore the CPU to the process we stole from.
	time.Sleep(lifetime)

	// Finally, ping each address and add this service as a neighbor.
	// TODO: Refactor this into a function which just takes the path.
	for _, addr := range addresses {
		var uri *url.URL
		var resp *http.Response
		var err error
		uri, err = url.Parse(addr)
		ExitOnError(err)
		uri.Path = "/neighbors/remove"
		var params = uri.Query()
		params.Set("cpu", cpuStr)
		uri.RawQuery = params.Encode()
		resp, err = http.Get(uri.String())
		ExitOnError(err)
		if resp.StatusCode != 200 {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			ExitOnError(err)
			log.Fatal(body)
		}
	}
}

// CPU requirements are presented as a number in the range [0…100]
// We simply need to parse this int to get the CPU.
func parseCPU(cpu string) uint64 {
	var i, err = strconv.ParseInt(cpu, 10, 64)
	ExitOnError(err)
	return uint64(i)
}

// A lifetime is an integer representing the amount of time this
// batch job takes to run in seconds. To return it as a time.Duration,
// we can tag the end of the string with an 's' and then parse the input
// using the standard library.
func parseLifetime(lifetime string) time.Duration {
	var inDurationFormat = fmt.Sprintf("%vs", lifetime)
	var duration, err = time.ParseDuration(inDurationFormat)
	ExitOnError(err)
	return duration
}

func parseAddresses(addresses string) []string {
	return strings.Split(addresses, ",")
}

func ExitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
