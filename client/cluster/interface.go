package cluster

import (
	"math/rand"
	"time"
)

type ServiceResponse struct {
	Address string
	// ThroughputStream <-chan uint64
}

type Cluster interface {
	NewBatch(addr string, cpu uint64)
	NewService(throughput, softLimit, hardLimit uint64) ServiceResponse
}

func NewMockCluster() Cluster {
	return new(mockCluster)
}

type mockCluster struct{ count int }

func (mock *mockCluster) NewBatch(addr string, cpu uint64) {
	// Pass
}

func (mock *mockCluster) NewService(throughput, softLimit, hardLimit uint64) ServiceResponse {
	switch mock.count {
	case 0:
		mock.count++
		return ServiceResponse{
			Address: "localhost:1234",
			// ThroughputStream: mock.mockThroughputStream(),
		}
	case 1:
		mock.count++
		return ServiceResponse{
			Address: "localhost:5678",
			//ThroughputStream: mock.mockThroughputStream(),
		}
	default:
		mock.count++
		return ServiceResponse{
			Address: "localhost:9999",
			// ThroughputStream: mock.mockThroughputStream(),
		}
	}
}

func (mock mockCluster) mockThroughputStream() <-chan uint64 {
	var stream = make(chan uint64, 100)
	go func() {
		// Every second, write a random value to the stream.
		for {
			var val = uint64(rand.Int63n(100))
			stream <- val
			time.Sleep(1 * time.Second)
		}
	}()

	return stream
}
