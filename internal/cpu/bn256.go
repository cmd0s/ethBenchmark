package cpu

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto/bn256/cloudflare"

	"github.com/vBenchmark/internal/types"
)

// BenchmarkBN256 measures BN256 elliptic curve operations
// These are used in EVM precompiled contracts for zkSNARK verification
// Reference: geth/core/vm/contracts.go (bn256Add, bn256ScalarMul, bn256Pairing)
func BenchmarkBN256(duration time.Duration, verbose bool) types.BN256Result {
	// Generate random test points
	_, g1a, err := bn256.RandomG1(rand.Reader)
	if err != nil {
		return types.BN256Result{Rating: "Error"}
	}
	_, g1b, _ := bn256.RandomG1(rand.Reader)
	_, g2a, _ := bn256.RandomG2(rand.Reader)

	// Generate random scalar for multiplication
	scalar := make([]byte, 32)
	rand.Read(scalar)
	scalarInt := new(big.Int).SetBytes(scalar)

	// Phase 1: G1 point addition (precompile 0x06)
	addDuration := duration * 3 / 10
	var addCount uint64
	start := time.Now()

	for time.Since(start) < addDuration {
		result := new(bn256.G1)
		result.Add(g1a, g1b)
		addCount++
	}
	addElapsed := time.Since(start)
	addRate := float64(addCount) / addElapsed.Seconds()

	// Phase 2: G1 scalar multiplication (precompile 0x07)
	mulDuration := duration * 3 / 10
	var mulCount uint64
	start = time.Now()

	for time.Since(start) < mulDuration {
		result := new(bn256.G1)
		result.ScalarMult(g1a, scalarInt)
		mulCount++
	}
	mulElapsed := time.Since(start)
	mulRate := float64(mulCount) / mulElapsed.Seconds()

	// Phase 3: Pairing operations (precompile 0x08)
	// This is the most expensive operation, used in zkSNARK verification
	pairDuration := duration * 4 / 10
	var pairCount uint64
	start = time.Now()

	for time.Since(start) < pairDuration {
		bn256.Pair(g1a, g2a)
		pairCount++
	}
	pairElapsed := time.Since(start)
	pairRate := float64(pairCount) / pairElapsed.Seconds()

	totalDuration := addElapsed + mulElapsed + pairElapsed

	return types.BN256Result{
		G1AddsPerSecond:       addRate,
		G1ScalarMulsPerSecond: mulRate,
		PairingsPerSecond:     pairRate,
		Duration:              totalDuration,
		Rating:                rateBN256(pairRate),
	}
}

// rateBN256 provides a rating based on pairing operations per second
func rateBN256(pairRate float64) string {
	switch {
	case pairRate >= 100:
		return "Excellent"
	case pairRate >= 50:
		return "Good"
	case pairRate >= 25:
		return "Adequate"
	case pairRate >= 10:
		return "Marginal"
	default:
		return "Poor"
	}
}
