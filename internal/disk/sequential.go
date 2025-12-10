// Package disk provides disk I/O benchmarks for Ethereum operations
package disk

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/vBenchmark/internal/types"
)

// BenchmarkSequential measures sequential I/O performance
// This simulates state sync and snapshot operations
func BenchmarkSequential(testDir string, duration time.Duration, verbose bool) types.SequentialResult {
	// Block sizes matching Ethereum data patterns:
	// - 128KB: LevelDB SST file writes
	// - 1MB: State snapshot chunks
	blockSizes := []int{128 * 1024, 1024 * 1024}

	testFile := filepath.Join(testDir, "ethbench_seq_test.dat")
	defer os.Remove(testFile)

	// Phase 1: Sequential writes
	writeDuration := duration / 2
	var totalWritten uint64
	writeStart := time.Now()

	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return types.SequentialResult{Rating: "Error: " + err.Error()}
	}

	for time.Since(writeStart) < writeDuration {
		for _, blockSize := range blockSizes {
			data := make([]byte, blockSize)
			rand.Read(data)
			n, err := f.Write(data)
			if err != nil {
				break
			}
			totalWritten += uint64(n)
		}
	}
	f.Sync()
	f.Close()

	writeElapsed := time.Since(writeStart)
	writeSpeed := float64(totalWritten) / writeElapsed.Seconds() / (1024 * 1024)

	// Phase 2: Sequential reads
	readDuration := duration / 2
	var totalRead uint64
	readStart := time.Now()

	f, err = os.Open(testFile)
	if err != nil {
		return types.SequentialResult{
			WriteSpeedMBps: writeSpeed,
			Rating:         "Error: " + err.Error(),
		}
	}

	buffer := make([]byte, 1024*1024) // 1MB read buffer
	for time.Since(readStart) < readDuration {
		n, err := f.Read(buffer)
		if err != nil {
			// Loop back to start of file
			f.Seek(0, 0)
			continue
		}
		totalRead += uint64(n)
	}
	f.Close()

	readElapsed := time.Since(readStart)
	readSpeed := float64(totalRead) / readElapsed.Seconds() / (1024 * 1024)

	totalDuration := writeElapsed + readElapsed

	return types.SequentialResult{
		WriteSpeedMBps: writeSpeed,
		ReadSpeedMBps:  readSpeed,
		Duration:       totalDuration,
		Rating:         rateSequential(writeSpeed, readSpeed),
	}
}

// rateSequential provides a rating based on sequential I/O speeds
func rateSequential(writeSpeed, readSpeed float64) string {
	// Weight write speed slightly higher for Ethereum workloads
	avgSpeed := writeSpeed*0.6 + readSpeed*0.4

	switch {
	case avgSpeed >= 400:
		return "Excellent"
	case avgSpeed >= 200:
		return "Good"
	case avgSpeed >= 100:
		return "Adequate"
	case avgSpeed >= 50:
		return "Marginal"
	default:
		return "Poor"
	}
}
