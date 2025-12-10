package benchmark

import (
	"fmt"
	"time"

	"github.com/vBenchmark/internal/cpu"
	"github.com/vBenchmark/internal/disk"
	"github.com/vBenchmark/internal/memory"
	"github.com/vBenchmark/internal/types"
)

// Runner orchestrates benchmark execution
type Runner struct {
	config    *Config
	StartTime time.Time
	verbose   bool
}

// NewRunner creates a new benchmark runner
func NewRunner(config *Config) *Runner {
	return &Runner{
		config:  config,
		verbose: config.Verbose,
	}
}

// RunAll executes all benchmarks and returns results
func (r *Runner) RunAll() *types.Results {
	r.StartTime = time.Now()
	results := &types.Results{}

	// Run CPU benchmarks
	r.log("Running CPU benchmarks...")
	results.CPU = r.runCPUBenchmarks()

	// Run Memory benchmarks
	r.log("Running Memory benchmarks...")
	results.Memory = r.runMemoryBenchmarks()

	// Run Disk benchmarks
	r.log("Running Disk benchmarks...")
	results.Disk = r.runDiskBenchmarks()

	return results
}

// runCPUBenchmarks executes all CPU benchmarks
func (r *Runner) runCPUBenchmarks() types.CPUResults {
	budget := r.config.GetCPUTimeBudget()
	results := types.CPUResults{}

	r.log("  [1/4] Keccak256 hashing...")
	results.Keccak = cpu.BenchmarkKeccak256(budget.Keccak256, r.verbose)

	r.log("  [2/4] ECDSA/secp256k1 signatures...")
	results.ECDSA = cpu.BenchmarkECDSA(budget.ECDSA, r.verbose)

	r.log("  [3/4] BLS12-381 operations...")
	results.BLS = cpu.BenchmarkBLS(budget.BLS, r.verbose)

	r.log("  [4/4] BN256 pairing...")
	results.BN256 = cpu.BenchmarkBN256(budget.BN256, r.verbose)

	return results
}

// runMemoryBenchmarks executes all memory benchmarks
func (r *Runner) runMemoryBenchmarks() types.MemoryResults {
	budget := r.config.GetMemoryTimeBudget()
	results := types.MemoryResults{}

	r.log("  [1/3] Merkle Patricia Trie simulation...")
	results.Trie = memory.BenchmarkTrie(budget.Trie, r.verbose)

	r.log("  [2/3] Object pool allocation...")
	results.Pool = memory.BenchmarkPool(budget.Pool, r.verbose)

	r.log("  [3/3] State cache operations...")
	results.StateCache = memory.BenchmarkStateCache(budget.StateCache, r.verbose)

	return results
}

// runDiskBenchmarks executes all disk benchmarks
func (r *Runner) runDiskBenchmarks() types.DiskResults {
	budget := r.config.GetDiskTimeBudget()
	results := types.DiskResults{}

	r.log("  [1/3] Sequential I/O...")
	results.Sequential = disk.BenchmarkSequential(r.config.TestDir, budget.Sequential, r.verbose)

	r.log("  [2/3] Random 4K I/O...")
	results.Random = disk.BenchmarkRandom(r.config.TestDir, budget.Random, r.verbose)

	r.log("  [3/3] Batch writes...")
	results.Batch = disk.BenchmarkBatch(r.config.TestDir, budget.Batch, r.verbose)

	return results
}

// log prints a message if verbose mode is enabled or always for progress
func (r *Runner) log(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Duration returns the total time elapsed since benchmark start
func (r *Runner) Duration() time.Duration {
	return time.Since(r.StartTime)
}
