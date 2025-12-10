package cpu

import (
	"math/big"
	"time"

	bls12381 "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"

	"github.com/vBenchmark/internal/types"
)

// BenchmarkBLS measures BLS12-381 operations performance
// This tests the actual cryptographic operations used in Ethereum consensus layer
// Reference: nimbus/beacon_chain/spec/crypto.nim, geth uses gnark-crypto
//
// BLS operations in consensus:
// - G1 scalar multiplication (signature generation)
// - Pairing operations (signature verification)
// - G2 point addition (signature aggregation)
func BenchmarkBLS(duration time.Duration, verbose bool) types.BLSResult {
	// Get generator points
	_, _, g1Gen, g2Gen := bls12381.Generators()

	// Phase 1: G1 scalar multiplication (simulates signature generation)
	// BLS signing involves multiplying the hash-to-curve point by secret key
	signDuration := duration / 4
	var signCount uint64
	start := time.Now()

	var scalar fr.Element
	var result bls12381.G1Affine

	for time.Since(start) < signDuration {
		// Generate random scalar (simulates secret key)
		scalar.SetRandom()
		// G1 scalar multiplication (core signing operation)
		result.ScalarMultiplication(&g1Gen, scalar.BigInt(new(big.Int)))
		signCount++
	}
	signElapsed := time.Since(start)
	signRate := float64(signCount) / signElapsed.Seconds()

	// Phase 2: Pairing operations (simulates signature verification)
	// BLS verify: e(sig, g2) == e(H(m), pk) requires pairing computation
	verifyDuration := duration / 4
	var verifyCount uint64
	start = time.Now()

	// Prepare points for pairing
	g1Points := []bls12381.G1Affine{g1Gen}
	g2Points := []bls12381.G2Affine{g2Gen}

	for time.Since(start) < verifyDuration {
		// Pairing operation (core verification)
		_, err := bls12381.Pair(g1Points, g2Points)
		if err == nil {
			verifyCount++
		}
	}
	verifyElapsed := time.Since(start)
	verifyRate := float64(verifyCount) / verifyElapsed.Seconds()

	// Phase 3: G2 point addition (simulates signature aggregation)
	// Aggregating multiple signatures involves G2 point additions
	aggDuration := duration / 4
	var aggCount uint64
	start = time.Now()

	var g2Jac bls12381.G2Jac
	g2Jac.FromAffine(&g2Gen)

	for time.Since(start) < aggDuration {
		// Simulate aggregating 64 signatures (typical committee size)
		var aggResult bls12381.G2Jac
		for i := 0; i < 64; i++ {
			aggResult.AddAssign(&g2Jac)
		}
		aggCount++
	}
	aggElapsed := time.Since(start)
	aggRate := float64(aggCount) / aggElapsed.Seconds()

	// Phase 4: Multi-pairing (simulates batch verification)
	// FastAggregateVerify uses multi-pairing for efficiency
	batchDuration := duration / 4
	var batchCount uint64
	start = time.Now()

	// Prepare multiple points for batch pairing (simulates 4 signature verification)
	multiG1 := []bls12381.G1Affine{g1Gen, g1Gen, g1Gen, g1Gen}
	multiG2 := []bls12381.G2Affine{g2Gen, g2Gen, g2Gen, g2Gen}

	for time.Since(start) < batchDuration {
		// Multi-pairing (batch verification)
		_, err := bls12381.Pair(multiG1, multiG2)
		if err == nil {
			batchCount++
		}
	}
	batchElapsed := time.Since(start)

	totalDuration := signElapsed + verifyElapsed + aggElapsed + batchElapsed

	return types.BLSResult{
		SignaturesPerSecond:    signRate,
		VerificationsPerSecond: verifyRate,
		AggregationsPerSecond:  aggRate,
		Duration:               totalDuration,
		Rating:                 rateBLS(verifyRate),
	}
}

// rateBLS provides a rating based on verification rate
// Thresholds calibrated for actual BLS12-381 pairing operations
func rateBLS(verifyRate float64) string {
	switch {
	case verifyRate >= 500:
		return "Excellent"
	case verifyRate >= 200:
		return "Good"
	case verifyRate >= 100:
		return "Adequate"
	case verifyRate >= 50:
		return "Marginal"
	default:
		return "Poor"
	}
}
