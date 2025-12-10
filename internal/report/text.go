package report

import (
	"fmt"
	"strings"
)

// FormatText generates a human-readable text report
func FormatText(r *Report) string {
	var sb strings.Builder

	// Header
	sb.WriteString("\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n")
	sb.WriteString("                    Ethereum Node Benchmark Report\n")
	sb.WriteString(fmt.Sprintf("                    Generated: %s\n", r.Metadata.Timestamp.Format("2006-01-02 15:04:05")))
	sb.WriteString(strings.Repeat("=", 80) + "\n")

	// System Information
	sb.WriteString("\nSYSTEM INFORMATION\n")
	sb.WriteString(strings.Repeat("-", 40) + "\n")
	sb.WriteString(fmt.Sprintf("  Hostname:      %s\n", r.System.Hostname))
	sb.WriteString(fmt.Sprintf("  Serial:        %s\n", r.System.SerialNumber))
	sb.WriteString(fmt.Sprintf("  OS:            %s %s\n", r.System.OS, r.System.OSVersion))
	sb.WriteString(fmt.Sprintf("  Architecture:  %s\n", r.System.Architecture))
	sb.WriteString(fmt.Sprintf("  CPU:           %s (%d cores)\n", r.System.CPUModel, r.System.CPUCores))
	sb.WriteString(fmt.Sprintf("  RAM:           %d MB\n", r.System.RAMTotalMB))
	sb.WriteString(fmt.Sprintf("  Storage:       %s\n", r.System.DiskModel))

	// CPU Benchmarks
	sb.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	sb.WriteString("CPU BENCHMARKS (Execution Layer Critical)\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n")

	sb.WriteString("\nKeccak256 Hashing (state trie, tx hashing)\n")
	sb.WriteString(fmt.Sprintf("  Throughput:     %.2f hashes/sec\n", r.CPU.Keccak.HashesPerSecond))
	sb.WriteString(fmt.Sprintf("  Data Processed: %.2f MB\n", r.CPU.Keccak.DataProcessedMB))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.CPU.Keccak.Rating))

	sb.WriteString("\nECDSA/secp256k1 (transaction signatures)\n")
	sb.WriteString(fmt.Sprintf("  Sign:           %.2f sig/sec\n", r.CPU.ECDSA.SignaturesPerSecond))
	sb.WriteString(fmt.Sprintf("  Verify:         %.2f verify/sec\n", r.CPU.ECDSA.VerificationsPerSecond))
	sb.WriteString(fmt.Sprintf("  ECRECOVER:      %.2f recover/sec\n", r.CPU.ECDSA.RecoveriesPerSecond))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.CPU.ECDSA.Rating))

	sb.WriteString("\nBLS12-381 (consensus layer signatures)\n")
	sb.WriteString(fmt.Sprintf("  Sign:           %.2f sig/sec\n", r.CPU.BLS.SignaturesPerSecond))
	sb.WriteString(fmt.Sprintf("  Verify:         %.2f verify/sec\n", r.CPU.BLS.VerificationsPerSecond))
	sb.WriteString(fmt.Sprintf("  Aggregate:      %.2f agg/sec\n", r.CPU.BLS.AggregationsPerSecond))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.CPU.BLS.Rating))

	sb.WriteString("\nBN256 Pairing (zkSNARK precompiles)\n")
	sb.WriteString(fmt.Sprintf("  G1 Add:         %.2f ops/sec\n", r.CPU.BN256.G1AddsPerSecond))
	sb.WriteString(fmt.Sprintf("  G1 ScalarMul:   %.2f ops/sec\n", r.CPU.BN256.G1ScalarMulsPerSecond))
	sb.WriteString(fmt.Sprintf("  Pairing:        %.2f ops/sec\n", r.CPU.BN256.PairingsPerSecond))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.CPU.BN256.Rating))

	// Memory Benchmarks
	sb.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	sb.WriteString("MEMORY BENCHMARKS\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n")

	sb.WriteString("\nMerkle Patricia Trie (state storage)\n")
	sb.WriteString(fmt.Sprintf("  Insert:         %.2f ops/sec\n", r.Memory.Trie.InsertsPerSecond))
	sb.WriteString(fmt.Sprintf("  Lookup:         %.2f ops/sec\n", r.Memory.Trie.LookupsPerSecond))
	sb.WriteString(fmt.Sprintf("  Hash:           %.2f ops/sec\n", r.Memory.Trie.HashesPerSecond))
	sb.WriteString(fmt.Sprintf("  Peak Memory:    %.2f MB\n", r.Memory.Trie.PeakMemoryMB))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.Memory.Trie.Rating))

	sb.WriteString("\nObject Pool Allocation (EVM memory)\n")
	sb.WriteString(fmt.Sprintf("  Allocations:    %.2f alloc/sec\n", r.Memory.Pool.AllocationsPerSecond))
	sb.WriteString(fmt.Sprintf("  Reuses:         %.2f reuse/sec\n", r.Memory.Pool.ReusesPerSecond))
	sb.WriteString(fmt.Sprintf("  Memory Churn:   %.2f MB\n", r.Memory.Pool.MemoryChurnMB))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.Memory.Pool.Rating))

	sb.WriteString("\nState Cache (account/storage)\n")
	sb.WriteString(fmt.Sprintf("  Cache Hits:     %.2f ops/sec\n", r.Memory.StateCache.CacheHitsPerSecond))
	sb.WriteString(fmt.Sprintf("  Cache Misses:   %.2f ops/sec\n", r.Memory.StateCache.CacheMissesPerSecond))
	sb.WriteString(fmt.Sprintf("  Hit Ratio:      %.2f%%\n", r.Memory.StateCache.HitRatio*100))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.Memory.StateCache.Rating))

	// Disk Benchmarks
	sb.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	sb.WriteString("DISK I/O BENCHMARKS\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n")

	sb.WriteString("\nSequential I/O (state sync, snapshots)\n")
	sb.WriteString(fmt.Sprintf("  Write Speed:    %.2f MB/s\n", r.Disk.Sequential.WriteSpeedMBps))
	sb.WriteString(fmt.Sprintf("  Read Speed:     %.2f MB/s\n", r.Disk.Sequential.ReadSpeedMBps))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.Disk.Sequential.Rating))

	sb.WriteString("\nRandom 4K I/O (trie node access)\n")
	sb.WriteString(fmt.Sprintf("  Read IOPS:      %.0f\n", r.Disk.Random.ReadIOPS))
	sb.WriteString(fmt.Sprintf("  Write IOPS:     %.0f\n", r.Disk.Random.WriteIOPS))
	sb.WriteString(fmt.Sprintf("  Avg Latency:    %.2f us\n", r.Disk.Random.AvgLatencyUs))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.Disk.Random.Rating))

	sb.WriteString("\nBatch Write (block commitment)\n")
	sb.WriteString(fmt.Sprintf("  Batch Rate:     %.2f batch/sec\n", r.Disk.Batch.BatchesPerSecond))
	sb.WriteString(fmt.Sprintf("  Throughput:     %.2f MB/s\n", r.Disk.Batch.ThroughputMBps))
	sb.WriteString(fmt.Sprintf("  Avg Latency:    %.2f ms\n", r.Disk.Batch.AvgBatchLatencyMs))
	sb.WriteString(fmt.Sprintf("  Rating:         %s\n", r.Disk.Batch.Rating))

	// Summary
	sb.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	sb.WriteString("SUMMARY\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n")
	sb.WriteString(fmt.Sprintf("\n  CPU Score:      %d/100\n", r.Summary.CPUScore))
	sb.WriteString(fmt.Sprintf("  Memory Score:   %d/100\n", r.Summary.MemoryScore))
	sb.WriteString(fmt.Sprintf("  Disk Score:     %d/100\n", r.Summary.DiskScore))
	sb.WriteString(fmt.Sprintf("  ─────────────────────\n"))
	sb.WriteString(fmt.Sprintf("  Overall Score:  %d/100\n", r.Summary.TotalScore))

	// Verdict
	sb.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	sb.WriteString("VERDICT\n")
	sb.WriteString(strings.Repeat("=", 80) + "\n")
	sb.WriteString(fmt.Sprintf("\n  Overall Score:        %d/100\n", r.Verdict.OverallScore))
	sb.WriteString(fmt.Sprintf("\n  Execution Client:     %s\n", r.Verdict.ExecutionClient))
	sb.WriteString(fmt.Sprintf("  Consensus Client:     %s\n", r.Verdict.ConsensusClient))
	sb.WriteString("\nRecommendations:\n")
	for _, rec := range r.Verdict.Recommendations {
		sb.WriteString(fmt.Sprintf("  - %s\n", rec))
	}

	sb.WriteString("\n" + strings.Repeat("=", 80) + "\n")
	sb.WriteString(fmt.Sprintf("Benchmark completed in %.1f seconds\n", r.Metadata.DurationSeconds))
	sb.WriteString(strings.Repeat("=", 80) + "\n")

	return sb.String()
}
