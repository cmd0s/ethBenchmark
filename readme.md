# Ethereum Node Benchmark (ethbench)

A specialized benchmark tool for measuring hardware performance relevant to Ethereum node operations (Geth + Nimbus).

## Target Platform: Raspberry Pi 5

This benchmark is designed specifically for **Raspberry Pi 5** running as an Ethereum validator. The tests Ethereum node workloads and help you verify that your hardware configuration is suitable for running both execution (Geth) and consensus (Nimbus) clients.

## Features

- **CPU Benchmarks**: Keccak256 hashing, ECDSA/secp256k1 signatures, BLS12-381 operations (using gnark-crypto), BN256 pairing
- **Memory Benchmarks**: Merkle Patricia Trie simulation, object pool allocation, state cache patterns
- **Disk Benchmarks**: Sequential I/O, random 4K I/O (bypasses page cache), batch write simulation
- **Raspberry Pi 5 Detection**: Model, GPU firmware, bootloader version, kernel, CPU governor/frequency, core voltage
- **Ethereum-Focused**: Tests based on actual Geth and Nimbus operation patterns
- **Scoring System**: Hardware readiness verdict for running Ethereum nodes

## System Requirements

### Supported Operating Systems
- Ubuntu (24.04+ recommended)
- Raspberry Pi OS (64-bit)
- Armbian
- DietPi
- Debian (12+)

### Recommended Hardware (Raspberry Pi 5)

| Component | Minimum | Recommended |
|-----------|---------|-------------|
| Model | Raspberry Pi 5 4GB | Raspberry Pi 5 8GB |
| Storage | 256GB NVMe SSD | 1TB+ NVMe SSD |
| Cooling | Passive heatsink | Active cooler |
| Power | Official 27W PSU | Official 27W PSU |

### Optional Tools

For baseline comparison, you can install:

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install -y fio stress-ng
```

## Installation

### Pre-built Binary

Download the latest release for your platform:

```bash
# For Raspberry Pi 5 / ARM64 Linux
wget https://github.com/vBenchmark/releases/latest/download/ethbench-linux-arm64
chmod +x ethbench-linux-arm64
./ethbench-linux-arm64
```

### Build from Source

```bash
# Clone repository
git clone https://github.com/vBenchmark/vBenchmark.git
cd vBenchmark

# Build for current platform
make build

# Cross-compile for ARM64
make build-arm64

# Build for all platforms
make build-all
```

## Usage

```bash
ethbench [options]

Options:
  -test-dir string    Directory for disk I/O tests (default: executable directory)
  -output string      Directory for JSON output file (default: executable directory)
  -quick              Quick mode: ~1 minute benchmark instead of 3 minutes
  -verbose            Show detailed progress during benchmarks
  -help               Show this help message
```

### Examples

```bash
# Run full benchmark (3 minutes)
./ethbench

# Run with specific test directory for disk I/O
./ethbench -test-dir /mnt/nvme

# Quick benchmark mode
./ethbench -quick

# Save JSON output to specific directory
./ethbench -output /home/user/benchmarks
```

## Output

### Terminal Output
Human-readable report displayed in the terminal with:
- System information
- Individual benchmark results
- Overall score and verdict

### JSON Output
Automatically saved to: `ethbench-YYYY-MM-DD_HH-MM-SS.json`

Contains:
- Complete benchmark results
- System information including device serial number
- Timestamp and duration
- Scoring and recommendations

## Benchmark Details

### CPU Benchmarks (~60 seconds)

| Test | Duration | Ethereum Relevance |
|------|----------|-------------------|
| Keccak256 | 15s | State trie hashing, transaction hashing |
| ECDSA/secp256k1 | 20s | Transaction signature verification |
| BLS12-381 | 15s | Consensus layer signature verification |
| BN256 Pairing | 10s | zkSNARK precompile operations |

### Memory Benchmarks (~60 seconds)

| Test | Duration | Ethereum Relevance |
|------|----------|-------------------|
| Trie Operations | 25s | State storage insert/lookup/hash |
| Pool Allocation | 15s | EVM memory management patterns |
| State Cache | 20s | Account and storage caching |

### Disk Benchmarks (~60 seconds)

| Test | Duration | Ethereum Relevance |
|------|----------|-------------------|
| Sequential I/O | 20s | State sync, snapshot operations |
| Random 4K I/O | 25s | Trie node random access |
| Batch Writes | 15s | Block commitment patterns |

## Scoring System

- **80-100**: Ready - Hardware meets Ethereum node requirements
- **60-79**: Marginal - May struggle under high load
- **40-59**: Below Spec - Slow sync expected
- **0-39**: Unsuitable - Hardware upgrade recommended

Score weights:
- CPU: 40%
- Disk: 35%
- Memory: 25%

## License

GNU GENERAL PUBLIC LICENSE version 3
