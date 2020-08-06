package main

// HealthCheckResponse tells the client if this server is alive or dead.
type HealthCheckResponse struct {
	Alive bool `json:"alive"`
}

// ThroughputGETResponse returns the throughput at this point in time.
// A dead server returns a throughput of 0.
type ThroughputGETResponse struct {
	Throughput uint64 `json:"throughput"`
}

// ThroughputPOSTResponse returns the new state of the server
// after the throughput has been changed.
type ThroughputPOSTResponse struct {
	Throughput       uint64 `json:"throughput"`
	RequestSoftLimit uint64 `json:"soft_limit"`
	RequestHardLimit uint64 `json:"hard_limit"`
}
