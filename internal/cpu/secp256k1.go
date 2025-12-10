package cpu

import (
	"crypto/ecdsa"
	"crypto/rand"
	"time"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/vBenchmark/internal/types"
)

// BenchmarkECDSA measures ECDSA/secp256k1 performance
// This is critical for transaction signature verification
// Reference: geth/crypto/crypto.go, geth/crypto/signature_cgo.go
func BenchmarkECDSA(duration time.Duration, verbose bool) types.ECDSAResult {
	// Generate test key pair
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return types.ECDSAResult{Rating: "Error"}
	}
	publicKey := privateKey.Public().(*ecdsa.PublicKey)
	pubKeyBytes := crypto.FromECDSAPub(publicKey)

	// Test message (typical transaction hash - 32 bytes)
	message := make([]byte, 32)
	rand.Read(message)

	// Phase 1: Signature generation
	signDuration := duration / 3
	var signCount uint64
	start := time.Now()

	for time.Since(start) < signDuration {
		_, err := crypto.Sign(message, privateKey)
		if err == nil {
			signCount++
		}
	}
	signElapsed := time.Since(start)
	signRate := float64(signCount) / signElapsed.Seconds()

	// Pre-generate signature for verification tests
	signature, _ := crypto.Sign(message, privateKey)

	// Phase 2: Signature verification (64-byte R||S format)
	verifyDuration := duration / 3
	var verifyCount uint64
	start = time.Now()

	for time.Since(start) < verifyDuration {
		// VerifySignature expects 64-byte signature (R||S without recovery byte)
		if crypto.VerifySignature(pubKeyBytes, message, signature[:64]) {
			verifyCount++
		}
	}
	verifyElapsed := time.Since(start)
	verifyRate := float64(verifyCount) / verifyElapsed.Seconds()

	// Phase 3: Public key recovery (ECRECOVER)
	// This is used in EVM precompiled contract 0x01
	recoverDuration := duration / 3
	var recoverCount uint64
	start = time.Now()

	for time.Since(start) < recoverDuration {
		_, err := crypto.Ecrecover(message, signature)
		if err == nil {
			recoverCount++
		}
	}
	recoverElapsed := time.Since(start)
	recoverRate := float64(recoverCount) / recoverElapsed.Seconds()

	totalDuration := signElapsed + verifyElapsed + recoverElapsed

	return types.ECDSAResult{
		SignaturesPerSecond:    signRate,
		VerificationsPerSecond: verifyRate,
		RecoveriesPerSecond:    recoverRate,
		Duration:               totalDuration,
		Rating:                 rateECDSA(verifyRate, recoverRate),
	}
}

// rateECDSA provides a rating based on verification and recovery rates
func rateECDSA(verifyRate, recoverRate float64) string {
	// Verification is more common, so weight it higher
	score := verifyRate*0.6 + recoverRate*0.4

	switch {
	case score >= 2000:
		return "Excellent"
	case score >= 1000:
		return "Good"
	case score >= 500:
		return "Adequate"
	case score >= 250:
		return "Marginal"
	default:
		return "Poor"
	}
}
