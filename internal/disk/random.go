package disk

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/vBenchmark/internal/types"
)

// BenchmarkRandom measures random 4K I/O performance
// This simulates trie node lookups during EVM execution
// Reference: geth/trie/trie.go resolveAndTrack()
func BenchmarkRandom(testDir string, duration time.Duration, verbose bool) types.RandomResult {
	const blockSize = 4096           // 4KB - typical trie node size
	const fileSize = 256 * 1024 * 1024 // 256MB test file

	testFile := filepath.Join(testDir, "ethbench_random_test.dat")
	defer os.Remove(testFile)

	// Create test file with random data
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return types.RandomResult{Rating: "Error: " + err.Error()}
	}

	// Pre-allocate the file
	if err := f.Truncate(fileSize); err != nil {
		f.Close()
		return types.RandomResult{Rating: "Error: " + err.Error()}
	}

	// Fill with some random data at regular intervals
	data := make([]byte, blockSize)
	for offset := int64(0); offset < fileSize; offset += 1024 * 1024 {
		rand.Read(data)
		f.WriteAt(data, offset)
	}
	f.Sync()

	numBlocks := fileSize / blockSize

	// Phase 1: Random reads (simulates trie lookups)
	readDuration := duration * 3 / 5
	var readOps uint64
	var totalReadLatency time.Duration

	readStart := time.Now()
	for time.Since(readStart) < readDuration {
		// Random offset within file
		blockNum := int64(readOps) % int64(numBlocks)
		offset := blockNum * blockSize

		opStart := time.Now()
		_, err := f.ReadAt(data, offset)
		totalReadLatency += time.Since(opStart)

		if err == nil {
			readOps++
		}
	}
	readElapsed := time.Since(readStart)
	readIOPS := float64(readOps) / readElapsed.Seconds()

	// Phase 2: Random writes (simulates dirty node flushes)
	writeDuration := duration * 2 / 5
	var writeOps uint64
	var totalWriteLatency time.Duration

	writeStart := time.Now()
	for time.Since(writeStart) < writeDuration {
		// Random offset within file
		blockNum := int64(writeOps) % int64(numBlocks)
		offset := blockNum * blockSize

		rand.Read(data)

		opStart := time.Now()
		_, err := f.WriteAt(data, offset)
		totalWriteLatency += time.Since(opStart)

		if err == nil {
			writeOps++
		}
	}
	f.Sync()
	f.Close()

	writeElapsed := time.Since(writeStart)
	writeIOPS := float64(writeOps) / writeElapsed.Seconds()

	// Calculate average latency across all operations
	totalOps := readOps + writeOps
	totalLatency := totalReadLatency + totalWriteLatency
	avgLatencyUs := float64(totalLatency.Microseconds()) / float64(totalOps)

	totalDuration := readElapsed + writeElapsed

	return types.RandomResult{
		ReadIOPS:     readIOPS,
		WriteIOPS:    writeIOPS,
		AvgLatencyUs: avgLatencyUs,
		Duration:     totalDuration,
		Rating:       rateRandom(readIOPS, writeIOPS),
	}
}

// rateRandom provides a rating based on random I/O performance
func rateRandom(readIOPS, writeIOPS float64) string {
	// Read IOPS are more important for Ethereum workloads
	score := readIOPS*0.7 + writeIOPS*0.3

	switch {
	case score >= 50000:
		return "Excellent"
	case score >= 20000:
		return "Good"
	case score >= 10000:
		return "Adequate"
	case score >= 5000:
		return "Marginal"
	default:
		return "Poor"
	}
}
