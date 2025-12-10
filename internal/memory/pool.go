package memory

import (
	"crypto/rand"
	"sync"
	"time"

	"github.com/vBenchmark/internal/types"
)

// memoryPool simulates EVM memory pool pattern
// Reference: geth/core/vm/memory.go
type memoryPool struct {
	pool sync.Pool
}

func newMemoryPool() *memoryPool {
	return &memoryPool{
		pool: sync.Pool{
			New: func() any {
				return make([]byte, 0, 4096)
			},
		},
	}
}

// stackPool simulates EVM stack pool
// Reference: geth/core/vm/stack.go
type stackPool struct {
	pool sync.Pool
}

func newStackPool() *stackPool {
	return &stackPool{
		pool: sync.Pool{
			New: func() any {
				// EVM stack: 1024 items of 32 bytes each
				return make([][32]byte, 0, 1024)
			},
		},
	}
}

// BenchmarkPool measures object pool allocation performance
// This simulates EVM memory management patterns
// Reference: geth/core/vm/memory.go, geth/core/vm/stack.go
func BenchmarkPool(duration time.Duration, verbose bool) types.PoolResult {
	memPool := newMemoryPool()
	stPool := newStackPool()

	var allocCount, reuseCount uint64
	var totalBytes uint64

	// Simulate EVM contract execution memory patterns
	start := time.Now()
	for time.Since(start) < duration {
		// Get memory from pool
		mem := memPool.pool.Get().([]byte)
		stack := stPool.pool.Get().([][32]byte)

		// Simulate memory expansion (like Memory.Resize)
		// Reference: geth/core/vm/memory.go Resize() lines 81-89
		// Target sizes: 1KB to 16KB (typical EVM memory usage)
		targetSize := 1024 + int(totalBytes%15360) // Deterministic but varied

		if cap(mem) < targetSize {
			mem = make([]byte, targetSize)
			allocCount++
		} else {
			mem = mem[:targetSize]
			reuseCount++
		}
		totalBytes += uint64(targetSize)

		// Simulate some memory operations (like MSTORE)
		if len(mem) >= 32 {
			for i := 0; i < len(mem)-32; i += 32 {
				rand.Read(mem[i : i+4]) // Partial fill to save time
			}
		}

		// Simulate stack operations
		stack = stack[:0]
		for i := 0; i < 16; i++ { // Typical stack depth during execution
			var item [32]byte
			stack = append(stack, item)
		}

		// Return to pool (like Memory.Free)
		// Reference: geth/core/vm/memory.go Free() lines 43-51
		const maxBufferSize = 16 << 10 // 16KB
		if cap(mem) <= maxBufferSize {
			clear(mem)
			memPool.pool.Put(mem[:0])
		}
		stPool.pool.Put(stack[:0])
	}

	elapsed := time.Since(start)
	totalOps := allocCount + reuseCount

	return types.PoolResult{
		AllocationsPerSecond: float64(allocCount) / elapsed.Seconds(),
		ReusesPerSecond:      float64(reuseCount) / elapsed.Seconds(),
		MemoryChurnMB:        float64(totalBytes) / (1024 * 1024),
		Duration:             elapsed,
		Rating:               ratePool(float64(totalOps) / elapsed.Seconds()),
	}
}

// ratePool provides a rating based on total operations per second
func ratePool(opsPerSec float64) string {
	switch {
	case opsPerSec >= 500000:
		return "Excellent"
	case opsPerSec >= 200000:
		return "Good"
	case opsPerSec >= 100000:
		return "Adequate"
	case opsPerSec >= 50000:
		return "Marginal"
	default:
		return "Poor"
	}
}
