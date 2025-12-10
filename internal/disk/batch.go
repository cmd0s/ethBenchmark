package disk

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/vBenchmark/internal/types"
)

// BenchmarkBatch measures batch write performance
// This simulates LevelDB batch write patterns during block commitment
// Reference: geth/ethdb/leveldb/leveldb.go Write()
func BenchmarkBatch(testDir string, duration time.Duration, verbose bool) types.BatchResult {
	// Simulate LevelDB batch characteristics:
	// - WriteBuffer: ~64MB (cache/4)
	// - Typical batch: 1000-5000 key-value pairs
	const kvSize = 100      // Average KV pair size in bytes
	const batchSize = 2000  // KV pairs per batch

	testFile := filepath.Join(testDir, "ethbench_batch_test.dat")
	defer os.Remove(testFile)

	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC|os.O_SYNC, 0644)
	if err != nil {
		return types.BatchResult{Rating: "Error: " + err.Error()}
	}
	defer f.Close()

	var batchCount uint64
	var totalWritten uint64
	var totalLatency time.Duration

	// Pre-allocate batch buffer
	batchBuffer := make([]byte, batchSize*kvSize)

	start := time.Now()
	for time.Since(start) < duration {
		// Build batch in memory (simulates LevelDB batch accumulation)
		// Each KV pair: key (32 bytes) + value (68 bytes) = 100 bytes
		rand.Read(batchBuffer)

		// Write batch with fsync (simulates durable write)
		opStart := time.Now()
		n, err := f.Write(batchBuffer)
		// Force sync to disk
		f.Sync()
		opLatency := time.Since(opStart)

		if err == nil {
			totalWritten += uint64(n)
			totalLatency += opLatency
			batchCount++
		}
	}

	elapsed := time.Since(start)

	batchesPerSec := float64(batchCount) / elapsed.Seconds()
	throughputMBps := float64(totalWritten) / elapsed.Seconds() / (1024 * 1024)
	avgBatchLatencyMs := float64(totalLatency.Milliseconds()) / float64(batchCount)

	return types.BatchResult{
		BatchesPerSecond:  batchesPerSec,
		ThroughputMBps:    throughputMBps,
		AvgBatchLatencyMs: avgBatchLatencyMs,
		Duration:          elapsed,
		Rating:            rateBatch(throughputMBps),
	}
}

// rateBatch provides a rating based on batch write throughput
func rateBatch(throughputMBps float64) string {
	switch {
	case throughputMBps >= 100:
		return "Excellent"
	case throughputMBps >= 50:
		return "Good"
	case throughputMBps >= 25:
		return "Adequate"
	case throughputMBps >= 10:
		return "Marginal"
	default:
		return "Poor"
	}
}
