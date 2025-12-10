// Package memory provides memory-intensive benchmarks for Ethereum operations
package memory

import (
	"crypto/rand"
	"runtime"
	"sync"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/vBenchmark/internal/types"
)

// hasher simulates Geth's hasher structure
// Reference: geth/trie/hasher.go
type hasher struct {
	tmp []byte
	sha sha3.ShakeHash
}

// hasherPool simulates Geth's hasher pooling pattern
var trieHasherPool = sync.Pool{
	New: func() any {
		return &hasher{
			tmp: make([]byte, 0, 550), // Same size as Geth
			sha: sha3.NewLegacyKeccak256().(sha3.ShakeHash),
		}
	},
}

// simulatedNode represents a trie node for benchmarking
// Reference: geth/trie/node.go
type simulatedNode struct {
	hash     [32]byte
	children [17]*simulatedNode // 16 children + value (fullNode pattern)
	key      []byte
	value    []byte
	dirty    bool
}

// BenchmarkTrie measures Merkle Patricia Trie operations
// This simulates state storage patterns in Geth
// Reference: geth/trie/trie.go
func BenchmarkTrie(duration time.Duration, verbose bool) types.TrieResult {
	nodes := make(map[[20]byte]*simulatedNode)
	nodeKeys := make([][20]byte, 0, 10000)

	var memBefore, memAfter runtime.MemStats
	runtime.ReadMemStats(&memBefore)

	// Phase 1: Trie insertions (simulates state updates during block processing)
	insertDuration := duration * 2 / 5
	var insertCount uint64
	start := time.Now()

	for time.Since(start) < insertDuration {
		// Simulate account address (20 bytes) -> account data
		var key [20]byte
		rand.Read(key[:])

		value := make([]byte, 100) // Typical account RLP size
		rand.Read(value)

		node := &simulatedNode{
			key:   key[:],
			value: value,
			dirty: true,
		}

		// Simulate trie path traversal and node hashing
		// Reference: geth/trie/trie.go insert() function
		h := trieHasherPool.Get().(*hasher)
		h.sha.Reset()
		h.sha.Write(key[:])
		h.sha.Write(value)
		h.sha.Read(node.hash[:])
		trieHasherPool.Put(h)

		nodes[key] = node
		nodeKeys = append(nodeKeys, key)
		insertCount++
	}
	insertElapsed := time.Since(start)
	insertRate := float64(insertCount) / insertElapsed.Seconds()

	// Phase 2: Trie lookups (simulates state reads during EVM execution)
	lookupDuration := duration * 2 / 5
	var lookupCount uint64
	start = time.Now()

	if len(nodeKeys) > 0 {
		for time.Since(start) < lookupDuration {
			// Random access pattern (simulates SLOAD operations)
			idx := int(lookupCount) % len(nodeKeys)
			key := nodeKeys[idx]
			_ = nodes[key]
			lookupCount++
		}
	}
	lookupElapsed := time.Since(start)
	lookupRate := float64(lookupCount) / lookupElapsed.Seconds()

	// Phase 3: Root hash computation (simulates block commitment)
	// Reference: geth/trie/trie.go hashRoot()
	hashDuration := duration / 5
	var hashCount uint64
	start = time.Now()

	for time.Since(start) < hashDuration {
		// Simulate parallel hashing like Geth when unhashed >= 100
		h := trieHasherPool.Get().(*hasher)
		for _, node := range nodes {
			if node.dirty {
				h.sha.Reset()
				h.sha.Write(node.hash[:])
				// Simulate hashing children
				for _, child := range node.children {
					if child != nil {
						h.sha.Write(child.hash[:])
					}
				}
			}
		}
		trieHasherPool.Put(h)
		hashCount++
	}
	hashElapsed := time.Since(start)
	hashRate := float64(hashCount) / hashElapsed.Seconds()

	runtime.ReadMemStats(&memAfter)
	peakMemMB := float64(memAfter.Alloc-memBefore.Alloc) / (1024 * 1024)
	if peakMemMB < 0 {
		peakMemMB = float64(memAfter.Alloc) / (1024 * 1024)
	}

	totalDuration := insertElapsed + lookupElapsed + hashElapsed

	return types.TrieResult{
		InsertsPerSecond: insertRate,
		LookupsPerSecond: lookupRate,
		HashesPerSecond:  hashRate,
		PeakMemoryMB:     peakMemMB,
		Duration:         totalDuration,
		Rating:           rateTrie(insertRate, lookupRate),
	}
}

// rateTrie provides a rating based on insert and lookup rates
func rateTrie(insertRate, lookupRate float64) string {
	// Weight lookups higher as they're more common
	score := insertRate*0.4 + lookupRate*0.001*0.6 // Scale lookup rate down

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
