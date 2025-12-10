// Package types provides common benchmark result types
package types

import (
	"time"
)

// Results holds all benchmark results
type Results struct {
	CPU    CPUResults    `json:"cpu"`
	Memory MemoryResults `json:"memory"`
	Disk   DiskResults   `json:"disk"`
}

// CPUResults contains all CPU benchmark results
type CPUResults struct {
	Keccak KeccakResult `json:"keccak"`
	ECDSA  ECDSAResult  `json:"ecdsa"`
	BLS    BLSResult    `json:"bls"`
	BN256  BN256Result  `json:"bn256"`
}

// KeccakResult holds Keccak256 benchmark results
type KeccakResult struct {
	HashesPerSecond float64       `json:"hashes_per_second"`
	TotalHashes     uint64        `json:"total_hashes"`
	DataProcessedMB float64       `json:"data_processed_mb"`
	Duration        time.Duration `json:"duration_ns"`
	Rating          string        `json:"rating"`
}

// ECDSAResult holds ECDSA/secp256k1 benchmark results
type ECDSAResult struct {
	SignaturesPerSecond    float64       `json:"signatures_per_second"`
	VerificationsPerSecond float64       `json:"verifications_per_second"`
	RecoveriesPerSecond    float64       `json:"recoveries_per_second"`
	Duration               time.Duration `json:"duration_ns"`
	Rating                 string        `json:"rating"`
}

// BLSResult holds BLS12-381 benchmark results
type BLSResult struct {
	SignaturesPerSecond    float64       `json:"signatures_per_second"`
	VerificationsPerSecond float64       `json:"verifications_per_second"`
	AggregationsPerSecond  float64       `json:"aggregations_per_second"`
	Duration               time.Duration `json:"duration_ns"`
	Rating                 string        `json:"rating"`
}

// BN256Result holds BN256 pairing benchmark results
type BN256Result struct {
	G1AddsPerSecond       float64       `json:"g1_adds_per_second"`
	G1ScalarMulsPerSecond float64       `json:"g1_scalar_muls_per_second"`
	PairingsPerSecond     float64       `json:"pairings_per_second"`
	Duration              time.Duration `json:"duration_ns"`
	Rating                string        `json:"rating"`
}

// MemoryResults contains all memory benchmark results
type MemoryResults struct {
	Trie       TrieResult       `json:"trie"`
	Pool       PoolResult       `json:"pool"`
	StateCache StateCacheResult `json:"state_cache"`
}

// TrieResult holds Merkle Patricia Trie benchmark results
type TrieResult struct {
	InsertsPerSecond float64       `json:"inserts_per_second"`
	LookupsPerSecond float64       `json:"lookups_per_second"`
	HashesPerSecond  float64       `json:"hashes_per_second"`
	PeakMemoryMB     float64       `json:"peak_memory_mb"`
	Duration         time.Duration `json:"duration_ns"`
	Rating           string        `json:"rating"`
}

// PoolResult holds object pool benchmark results
type PoolResult struct {
	AllocationsPerSecond float64       `json:"allocations_per_second"`
	ReusesPerSecond      float64       `json:"reuses_per_second"`
	MemoryChurnMB        float64       `json:"memory_churn_mb"`
	Duration             time.Duration `json:"duration_ns"`
	Rating               string        `json:"rating"`
}

// StateCacheResult holds state cache benchmark results
type StateCacheResult struct {
	CacheHitsPerSecond   float64       `json:"cache_hits_per_second"`
	CacheMissesPerSecond float64       `json:"cache_misses_per_second"`
	HitRatio             float64       `json:"hit_ratio"`
	ThroughputMBPerSec   float64       `json:"throughput_mb_per_sec"`
	Duration             time.Duration `json:"duration_ns"`
	Rating               string        `json:"rating"`
}

// DiskResults contains all disk benchmark results
type DiskResults struct {
	Sequential SequentialResult `json:"sequential"`
	Random     RandomResult     `json:"random"`
	Batch      BatchResult      `json:"batch"`
}

// SequentialResult holds sequential I/O benchmark results
type SequentialResult struct {
	WriteSpeedMBps float64       `json:"write_speed_mbps"`
	ReadSpeedMBps  float64       `json:"read_speed_mbps"`
	Duration       time.Duration `json:"duration_ns"`
	Rating         string        `json:"rating"`
}

// RandomResult holds random I/O benchmark results
type RandomResult struct {
	ReadIOPS     float64       `json:"read_iops"`
	WriteIOPS    float64       `json:"write_iops"`
	AvgLatencyUs float64       `json:"avg_latency_us"`
	Duration     time.Duration `json:"duration_ns"`
	Rating       string        `json:"rating"`
}

// BatchResult holds batch write benchmark results
type BatchResult struct {
	BatchesPerSecond  float64       `json:"batches_per_second"`
	ThroughputMBps    float64       `json:"throughput_mbps"`
	AvgBatchLatencyMs float64       `json:"avg_batch_latency_ms"`
	Duration          time.Duration `json:"duration_ns"`
	Rating            string        `json:"rating"`
}
