package disk

import (
	"crypto/rand"
	mathrand "math/rand"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/vBenchmark/internal/types"
)

// BenchmarkRandom measures random 4K I/O performance
// This simulates trie node lookups during EVM execution
// Reference: geth/trie/trie.go resolveAndTrack()
func BenchmarkRandom(testDir string, duration time.Duration, verbose bool) types.RandomResult {
	const blockSize = 4096                 // 4KB - typical trie node size
	const fileSize = 1024 * 1024 * 1024    // 1GB test file - larger than typical cache

	testFile := filepath.Join(testDir, "ethbench_random_test.dat")
	defer os.Remove(testFile)

	// Create and populate test file
	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return types.RandomResult{Rating: "Error: " + err.Error()}
	}

	// Pre-allocate the file
	if err := f.Truncate(fileSize); err != nil {
		f.Close()
		return types.RandomResult{Rating: "Error: " + err.Error()}
	}

	// Fill with random data at intervals to ensure file is actually allocated
	data := make([]byte, blockSize)
	for offset := int64(0); offset < fileSize; offset += 4 * 1024 * 1024 { // Every 4MB
		rand.Read(data)
		f.WriteAt(data, offset)
	}
	f.Sync()

	numBlocks := fileSize / blockSize
	rng := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	// Drop page cache before reading
	fd := int(f.Fd())
	syscall.Syscall6(syscall.SYS_FADVISE64, uintptr(fd), 0, uintptr(fileSize), uintptr(4), 0, 0) // POSIX_FADV_DONTNEED = 4

	// Phase 1: Random reads (simulates trie lookups)
	readDuration := duration * 3 / 5
	var readOps uint64
	var totalReadLatency time.Duration

	readStart := time.Now()
	for time.Since(readStart) < readDuration {
		// Truly random offset within file
		blockNum := rng.Int63n(int64(numBlocks))
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

	// Phase 2: Random writes with sync (simulates dirty node flushes)
	writeDuration := duration * 2 / 5
	var writeOps uint64
	var totalWriteLatency time.Duration

	writeStart := time.Now()
	for time.Since(writeStart) < writeDuration {
		// Truly random offset within file
		blockNum := rng.Int63n(int64(numBlocks))
		offset := blockNum * blockSize

		rand.Read(data)

		opStart := time.Now()
		_, err := f.WriteAt(data, offset)
		// Sync periodically to measure real write latency (every 100 ops)
		if writeOps%100 == 99 {
			f.Sync()
		}
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
