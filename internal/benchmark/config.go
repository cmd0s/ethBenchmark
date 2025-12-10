// Package benchmark provides the benchmark orchestration framework
package benchmark

import (
	"time"
)

// Config holds benchmark configuration
type Config struct {
	// Duration settings
	CPUDuration    time.Duration
	MemoryDuration time.Duration
	DiskDuration   time.Duration

	// Test directory for disk benchmarks
	TestDir string

	// Output settings
	Verbose bool
}

// DefaultConfig returns the default benchmark configuration
func DefaultConfig() *Config {
	return &Config{
		CPUDuration:    60 * time.Second,
		MemoryDuration: 60 * time.Second,
		DiskDuration:   60 * time.Second,
		TestDir:        ".",
		Verbose:        false,
	}
}

// QuickConfig returns a quick benchmark configuration (~1 minute total)
func QuickConfig() *Config {
	return &Config{
		CPUDuration:    20 * time.Second,
		MemoryDuration: 20 * time.Second,
		DiskDuration:   20 * time.Second,
		TestDir:        ".",
		Verbose:        false,
	}
}

// CPUTimeBudget returns time allocations for each CPU benchmark
type CPUTimeBudget struct {
	Keccak256 time.Duration
	ECDSA     time.Duration
	BLS       time.Duration
	BN256     time.Duration
}

// GetCPUTimeBudget calculates time budget for CPU benchmarks
func (c *Config) GetCPUTimeBudget() CPUTimeBudget {
	total := c.CPUDuration
	return CPUTimeBudget{
		Keccak256: total * 15 / 60, // 25%
		ECDSA:     total * 20 / 60, // 33%
		BLS:       total * 15 / 60, // 25%
		BN256:     total * 10 / 60, // 17%
	}
}

// MemoryTimeBudget returns time allocations for each memory benchmark
type MemoryTimeBudget struct {
	Trie       time.Duration
	Pool       time.Duration
	StateCache time.Duration
}

// GetMemoryTimeBudget calculates time budget for memory benchmarks
func (c *Config) GetMemoryTimeBudget() MemoryTimeBudget {
	total := c.MemoryDuration
	return MemoryTimeBudget{
		Trie:       total * 25 / 60, // 42%
		Pool:       total * 15 / 60, // 25%
		StateCache: total * 20 / 60, // 33%
	}
}

// DiskTimeBudget returns time allocations for each disk benchmark
type DiskTimeBudget struct {
	Sequential time.Duration
	Random     time.Duration
	Batch      time.Duration
}

// GetDiskTimeBudget calculates time budget for disk benchmarks
func (c *Config) GetDiskTimeBudget() DiskTimeBudget {
	total := c.DiskDuration
	return DiskTimeBudget{
		Sequential: total * 20 / 60, // 33%
		Random:     total * 25 / 60, // 42%
		Batch:      total * 15 / 60, // 25%
	}
}
