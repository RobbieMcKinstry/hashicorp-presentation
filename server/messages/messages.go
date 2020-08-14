package messages

// HealthCheckResponse tells the client if this server is alive or dead.
type HealthCheckResponse struct {
	Alive        bool   `json:"alive"`
	AvailableCPU string `json:"avaiable_cpu"`
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

// NeighborAddResponse is the JSON payload returned
// when a new noisy neighbor is added to the pod,
// resulting in a decrease in avaiable CPU.
type NeighborAddResponse struct {
	PreviousStolenCPU uint64
	StolenCPU         uint64
}

// NeighborRemoveResponse is the JSON payload returned
// when a noisy neighbor is removed from the pod,
// resulting in CPU being restored to this service.
type NeighborRemoveResponse struct {
	StolenCPU   uint64
	RestoredCPU uint64
}
