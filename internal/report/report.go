// Package report provides benchmark report generation
package report

import (
	"time"

	"github.com/vBenchmark/internal/system"
	"github.com/vBenchmark/internal/types"
)

// Report contains the complete benchmark report
type Report struct {
	Metadata Metadata          `json:"metadata"`
	System   *system.Info      `json:"system"`
	CPU      types.CPUResults    `json:"cpu"`
	Memory   types.MemoryResults `json:"memory"`
	Disk     types.DiskResults   `json:"disk"`
	Summary  Summary           `json:"summary"`
	Verdict  Verdict           `json:"verdict"`
}

// Metadata contains report metadata
type Metadata struct {
	Version         string    `json:"version"`
	Timestamp       time.Time `json:"timestamp"`
	DurationSeconds float64   `json:"duration_seconds"`
}

// Summary contains score summaries for each category
type Summary struct {
	CPUScore    int `json:"cpu_score"`
	MemoryScore int `json:"memory_score"`
	DiskScore   int `json:"disk_score"`
	TotalScore  int `json:"total_score"`
}

// Verdict contains the final hardware assessment
type Verdict struct {
	OverallScore      int      `json:"overall_score"`
	ExecutionClient   string   `json:"execution_client"`
	ConsensusClient   string   `json:"consensus_client"`
	Recommendations   []string `json:"recommendations"`
}

// NewReport creates a new benchmark report
func NewReport(version string, sysInfo *system.Info, results *types.Results, duration time.Duration) *Report {
	report := &Report{
		Metadata: Metadata{
			Version:         version,
			Timestamp:       time.Now(),
			DurationSeconds: duration.Seconds(),
		},
		System: sysInfo,
		CPU:    results.CPU,
		Memory: results.Memory,
		Disk:   results.Disk,
	}

	// Calculate scores
	report.Summary = calculateSummary(results)
	report.Verdict = determineVerdict(report.Summary.TotalScore, results)

	return report
}

// calculateSummary calculates scores for each category
func calculateSummary(results *types.Results) Summary {
	cpuScore := calculateCPUScore(&results.CPU)
	memoryScore := calculateMemoryScore(&results.Memory)
	diskScore := calculateDiskScore(&results.Disk)

	// Weighted total: CPU 40%, Disk 35%, Memory 25%
	totalScore := int(float64(cpuScore)*0.40 + float64(diskScore)*0.35 + float64(memoryScore)*0.25)

	return Summary{
		CPUScore:    cpuScore,
		MemoryScore: memoryScore,
		DiskScore:   diskScore,
		TotalScore:  totalScore,
	}
}

// calculateCPUScore scores CPU benchmark results (0-100)
func calculateCPUScore(cpu *types.CPUResults) int {
	var score float64

	// Keccak256 scoring (25% weight)
	keccakScore := scoreMetric(cpu.Keccak.HashesPerSecond, 50000, 100000, 200000, 500000)
	score += keccakScore * 0.25

	// ECDSA scoring (35% weight) - uses verification rate
	ecdsaScore := scoreMetric(cpu.ECDSA.VerificationsPerSecond, 250, 500, 1000, 2000)
	score += ecdsaScore * 0.35

	// BLS scoring (25% weight)
	blsScore := scoreMetric(cpu.BLS.VerificationsPerSecond, 50, 100, 200, 500)
	score += blsScore * 0.25

	// BN256 scoring (15% weight)
	bn256Score := scoreMetric(cpu.BN256.PairingsPerSecond, 10, 25, 50, 100)
	score += bn256Score * 0.15

	return int(score)
}

// calculateMemoryScore scores memory benchmark results (0-100)
func calculateMemoryScore(mem *types.MemoryResults) int {
	var score float64

	// Trie operations scoring (40% weight)
	trieScore := scoreMetric(mem.Trie.InsertsPerSecond, 5000, 10000, 20000, 50000)
	score += trieScore * 0.40

	// Pool operations scoring (30% weight)
	poolOps := mem.Pool.AllocationsPerSecond + mem.Pool.ReusesPerSecond
	poolScore := scoreMetric(poolOps, 50000, 100000, 200000, 500000)
	score += poolScore * 0.30

	// State cache scoring (30% weight)
	cacheScore := scoreMetric(mem.StateCache.CacheHitsPerSecond, 50000, 100000, 200000, 500000)
	score += cacheScore * 0.30

	return int(score)
}

// calculateDiskScore scores disk benchmark results (0-100)
func calculateDiskScore(disk *types.DiskResults) int {
	var score float64

	// Sequential I/O scoring (30% weight)
	seqAvg := (disk.Sequential.WriteSpeedMBps + disk.Sequential.ReadSpeedMBps) / 2
	seqScore := scoreMetric(seqAvg, 50, 100, 200, 400)
	score += seqScore * 0.30

	// Random I/O scoring (45% weight) - most important for Ethereum
	randomAvg := (disk.Random.ReadIOPS + disk.Random.WriteIOPS) / 2
	randomScore := scoreMetric(randomAvg, 5000, 10000, 20000, 50000)
	score += randomScore * 0.45

	// Batch write scoring (25% weight)
	batchScore := scoreMetric(disk.Batch.ThroughputMBps, 10, 25, 50, 100)
	score += batchScore * 0.25

	return int(score)
}

// scoreMetric converts a metric value to a 0-100 score
func scoreMetric(value, poor, marginal, good, excellent float64) float64 {
	switch {
	case value >= excellent:
		return 100
	case value >= good:
		return 75 + 25*(value-good)/(excellent-good)
	case value >= marginal:
		return 50 + 25*(value-marginal)/(good-marginal)
	case value >= poor:
		return 25 + 25*(value-poor)/(marginal-poor)
	default:
		return 25 * value / poor
	}
}

// determineVerdict determines hardware readiness for Ethereum nodes
func determineVerdict(score int, results *types.Results) Verdict {
	verdict := Verdict{
		OverallScore:    score,
		Recommendations: make([]string, 0),
	}

	// Determine client readiness
	switch {
	case score >= 80:
		verdict.ExecutionClient = "Ready"
		verdict.ConsensusClient = "Ready"
		verdict.Recommendations = append(verdict.Recommendations,
			"Your hardware meets Ethereum node requirements.",
			"Both Geth and Nimbus should run well on this system.",
		)
	case score >= 60:
		verdict.ExecutionClient = "Marginal"
		verdict.ConsensusClient = "Ready"
		verdict.Recommendations = append(verdict.Recommendations,
			"Consensus client (Nimbus) should work well.",
			"Execution client (Geth) may struggle during high network activity.",
			"Consider using checkpoint sync to reduce initial sync time.",
		)
	case score >= 40:
		verdict.ExecutionClient = "Marginal"
		verdict.ConsensusClient = "Marginal"
		verdict.Recommendations = append(verdict.Recommendations,
			"Hardware is below recommended specifications.",
			"Initial sync will be slow (potentially weeks).",
			"Consider using an external execution client RPC.",
		)
	default:
		verdict.ExecutionClient = "Unsuitable"
		verdict.ConsensusClient = "Marginal"
		verdict.Recommendations = append(verdict.Recommendations,
			"Hardware does not meet minimum requirements for execution client.",
			"Consider upgrading to NVMe storage.",
			"A more powerful single-board computer is recommended.",
		)
	}

	// Add specific recommendations based on weak areas
	if results.Disk.Random.ReadIOPS < 10000 {
		verdict.Recommendations = append(verdict.Recommendations,
			"Random I/O performance is low. NVMe SSD strongly recommended.",
		)
	}
	if results.CPU.ECDSA.VerificationsPerSecond < 500 {
		verdict.Recommendations = append(verdict.Recommendations,
			"ECDSA verification is slow. This may cause transaction validation delays.",
		)
	}
	if results.CPU.BLS.VerificationsPerSecond < 100 {
		verdict.Recommendations = append(verdict.Recommendations,
			"BLS signature verification is slow. Consensus layer may lag.",
		)
	}

	return verdict
}
