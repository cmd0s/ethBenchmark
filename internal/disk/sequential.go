// Package disk provides disk I/O benchmarks for Ethereum operations
package disk

import (
	"crypto/rand"
	"os"
	"path/filepath"
	"syscall"
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

	// Phase 1: Sequential writes with sync
	writeDuration := duration / 2
	var totalWritten uint64
	writeStart := time.Now()

	f, err := os.OpenFile(testFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return types.SequentialResult{Rating: "Error: " + err.Error()}
	}

	// Pre-allocate buffer to avoid GC during benchmark
	buffer := make([]byte, 1024*1024)
	rand.Read(buffer)

	for time.Since(writeStart) < writeDuration {
		for _, blockSize := range blockSizes {
			data := buffer[:blockSize]
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

	// Phase 2: Sequential reads - bypass page cache
	readDuration := duration / 2
	var totalRead uint64

	f, err = os.OpenFile(testFile, os.O_RDONLY, 0)
	if err != nil {
		return types.SequentialResult{
			WriteSpeedMBps: writeSpeed,
			Rating:         "Error: " + err.Error(),
		}
	}

	// Drop page cache for this file using fadvise
	fd := int(f.Fd())
	fileInfo, _ := f.Stat()
	fileSize := fileInfo.Size()
	syscall.Syscall6(syscall.SYS_FADVISE64, uintptr(fd), 0, uintptr(fileSize), uintptr(4), 0, 0) // POSIX_FADV_DONTNEED = 4

	readStart := time.Now()
	readBuffer := make([]byte, 1024*1024) // 1MB read buffer

	for time.Since(readStart) < readDuration {
		n, err := f.Read(readBuffer)
		if err != nil {
			// Loop back to start of file, drop cache again
			f.Seek(0, 0)
			syscall.Syscall6(syscall.SYS_FADVISE64, uintptr(fd), 0, uintptr(fileSize), uintptr(4), 0, 0)
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
