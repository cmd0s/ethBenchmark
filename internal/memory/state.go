package memory

import (
	"crypto/rand"
	"time"

	"github.com/vBenchmark/internal/types"
)

// stateObject simulates Geth's state object caching
// Reference: geth/core/state/state_object.go
type stateObject struct {
	address        [20]byte
	data           []byte
	originStorage  map[[32]byte][32]byte // Original values
	dirtyStorage   map[[32]byte][32]byte // Modified values
	pendingStorage map[[32]byte][32]byte // Pending commit
	storageKeys    [][32]byte            // Keys for fast random access
}

// BenchmarkStateCache measures state access patterns
// This simulates account and storage caching in Geth
// Reference: geth/core/state/state_object.go
func BenchmarkStateCache(duration time.Duration, verbose bool) types.StateCacheResult {
	// Pre-populate cache with realistic state data
	// Simulating ~10000 accounts typical for a busy block
	cache := make(map[[20]byte]*stateObject)
	addresses := make([][20]byte, 0, 10000)

	for i := 0; i < 10000; i++ {
		var addr [20]byte
		rand.Read(addr[:])

		obj := &stateObject{
			address:        addr,
			data:           make([]byte, 100),
			originStorage:  make(map[[32]byte][32]byte),
			dirtyStorage:   make(map[[32]byte][32]byte),
			pendingStorage: make(map[[32]byte][32]byte),
			storageKeys:    make([][32]byte, 0, 50),
		}
		rand.Read(obj.data)

		// Pre-populate storage slots (typical contract state)
		for j := 0; j < 50; j++ {
			var key, val [32]byte
			rand.Read(key[:])
			rand.Read(val[:])
			obj.originStorage[key] = val
			obj.storageKeys = append(obj.storageKeys, key) // Store keys for this object
		}

		cache[addr] = obj
		addresses = append(addresses, addr)
	}

	var hits, misses uint64
	var totalBytes uint64

	start := time.Now()
	for time.Since(start) < duration {
		// 80% cache hits (typical during block processing)
		// This simulates the pattern where most accessed accounts are already cached
		opIndex := hits + misses
		if opIndex%5 < 4 { // 80% of the time
			// Cache hit path - access existing account
			idx := int(opIndex) % len(addresses)
			addr := addresses[idx]
			obj := cache[addr]

			// Use a key that belongs to THIS object
			keyIdx := int(opIndex) % len(obj.storageKeys)
			key := obj.storageKeys[keyIdx]

			// Check dirty first, then pending, then origin
			// This mirrors Geth's GetState() logic
			if _, ok := obj.dirtyStorage[key]; ok {
				hits++
				totalBytes += 32
			} else if _, ok := obj.pendingStorage[key]; ok {
				hits++
				totalBytes += 32
			} else if val, ok := obj.originStorage[key]; ok {
				// Simulate caching the read in dirty storage
				obj.dirtyStorage[key] = val
				hits++
				totalBytes += 32
			} else {
				// Should not happen with correct keys
				misses++
				totalBytes += 32
			}
		} else {
			// Cache miss - simulate new account access (20%)
			var newAddr [20]byte
			rand.Read(newAddr[:])
			_, exists := cache[newAddr]
			if !exists {
				misses++
			} else {
				hits++ // Rare case where random address matches
			}
			totalBytes += 100 // Account data size
		}
	}

	elapsed := time.Since(start)
	total := hits + misses
	hitRatio := float64(hits) / float64(total)

	return types.StateCacheResult{
		CacheHitsPerSecond:   float64(hits) / elapsed.Seconds(),
		CacheMissesPerSecond: float64(misses) / elapsed.Seconds(),
		HitRatio:             hitRatio,
		ThroughputMBPerSec:   float64(totalBytes) / elapsed.Seconds() / (1024 * 1024),
		Duration:             elapsed,
		Rating:               rateStateCache(float64(hits) / elapsed.Seconds()),
	}
}

// rateStateCache provides a rating based on cache hit rate
func rateStateCache(hitsPerSec float64) string {
	switch {
	case hitsPerSec >= 500000:
		return "Excellent"
	case hitsPerSec >= 200000:
		return "Good"
	case hitsPerSec >= 100000:
		return "Adequate"
	case hitsPerSec >= 50000:
		return "Marginal"
	default:
		return "Poor"
	}
}
