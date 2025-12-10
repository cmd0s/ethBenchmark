// Package cpu provides CPU-intensive benchmarks for Ethereum operations
package cpu

import (
	"crypto/rand"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/vBenchmark/internal/types"
)

// hasherPool reuses Keccak256 hasher instances like Geth does
// Reference: geth/crypto/keccak.go
var hasherPool = sync.Pool{
	New: func() any {
		return sha3.NewLegacyKeccak256()
	},
}

// BenchmarkKeccak256 measures Keccak256 hashing performance
// This is critical for state trie operations and transaction hashing
func BenchmarkKeccak256(duration time.Duration, verbose bool) types.KeccakResult {
	// Input sizes matching Ethereum data patterns:
	// - 32 bytes: hash of hash (common in tries)
	// - 64 bytes: two concatenated hashes
	// - 128 bytes: typical small data
	// - 550 bytes: max fullNode encoding (see geth/trie/hasher.go line 41)
	inputSizes := []int{32, 64, 128, 550}

	// Pre-generate test data
	testData := make([][]byte, len(inputSizes))
	for i, size := range inputSizes {
		testData[i] = make([]byte, size)
		rand.Read(testData[i])
	}

	var totalHashes uint64
	var totalBytes uint64
	output := make([]byte, 32)

	start := time.Now()
	for time.Since(start) < duration {
		for i, data := range testData {
			// Get hasher from pool (like Geth does)
			hasher := hasherPool.Get().(sha3.ShakeHash)
			hasher.Reset()
			hasher.Write(data)
			hasher.Read(output)
			hasherPool.Put(hasher)

			totalHashes++
			totalBytes += uint64(inputSizes[i])
		}
	}

	elapsed := time.Since(start)
	hashesPerSec := float64(totalHashes) / elapsed.Seconds()
	dataMB := float64(totalBytes) / (1024 * 1024)

	return types.KeccakResult{
		HashesPerSecond: hashesPerSec,
		TotalHashes:     totalHashes,
		DataProcessedMB: dataMB,
		Duration:        elapsed,
		Rating:          rateKeccak(hashesPerSec),
	}
}

// rateKeccak provides a rating based on hashes per second
func rateKeccak(hps float64) string {
	switch {
	case hps >= 500000:
		return "Excellent"
	case hps >= 200000:
		return "Good"
	case hps >= 100000:
		return "Adequate"
	case hps >= 50000:
		return "Marginal"
	default:
		return "Poor"
	}
}
