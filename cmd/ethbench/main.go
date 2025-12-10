// Package main provides the CLI entry point for ethbench
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/vBenchmark/internal/benchmark"
	"github.com/vBenchmark/internal/report"
	"github.com/vBenchmark/internal/system"
)

const (
	version = "0.1.0"
	banner  = `
 _____ _   _     ____                  _
| ____| |_| |__ | __ )  ___ _ __   ___| |__
|  _| | __| '_ \|  _ \ / _ \ '_ \ / __| '_ \
| |___| |_| | | | |_) |  __/ | | | (__| | | |
|_____|\__|_| |_|____/ \___|_| |_|\___|_| |_|

Ethereum Node Benchmark Tool v%s
Target: Raspberry Pi 5 / ARM64 Linux
`
)

func main() {
	// Get executable directory for default paths
	execPath, err := os.Executable()
	if err != nil {
		execPath = "."
	}
	execDir := filepath.Dir(execPath)

	// Parse command line arguments
	testDir := flag.String("test-dir", execDir, "Directory for disk I/O tests")
	outputDir := flag.String("output", execDir, "Directory for JSON output file")
	quick := flag.Bool("quick", false, "Quick mode: ~1 minute benchmark")
	verbose := flag.Bool("verbose", false, "Show detailed progress")
	showHelp := flag.Bool("help", false, "Show help message")

	flag.Parse()

	if *showHelp {
		printHelp()
		return
	}

	// Print banner
	fmt.Printf(banner, version)
	fmt.Println()

	// Detect system information
	fmt.Println("Detecting system information...")
	sysInfo, err := system.Detect()
	if err != nil {
		fmt.Printf("Warning: Could not detect all system info: %v\n", err)
	}

	// Print system info summary
	fmt.Printf("  System: %s %s (%s)\n", sysInfo.OS, sysInfo.OSVersion, sysInfo.Architecture)
	fmt.Printf("  CPU: %s (%d cores)\n", sysInfo.CPUModel, sysInfo.CPUCores)
	fmt.Printf("  RAM: %d MB\n", sysInfo.RAMTotalMB)
	fmt.Printf("  Storage: %s\n", sysInfo.DiskModel)
	fmt.Printf("  Serial: %s\n", sysInfo.SerialNumber)
	fmt.Println()

	// Check prerequisites
	fmt.Printf("Testing write access to %s...\n", *testDir)
	if err := system.CheckPrerequisites(*testDir); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("  OK")
	fmt.Println()

	// Configure benchmark
	var config *benchmark.Config
	if *quick {
		config = benchmark.QuickConfig()
		fmt.Println("Quick mode enabled - benchmark will take approximately 1 minute")
	} else {
		config = benchmark.DefaultConfig()
		fmt.Println("Full benchmark mode - this will take approximately 3 minutes")
	}
	config.TestDir = *testDir
	config.Verbose = *verbose

	fmt.Println()
	fmt.Println("Starting benchmarks...")
	fmt.Println()

	// Create and run benchmark
	runner := benchmark.NewRunner(config)
	results := runner.RunAll()

	// Generate report
	fmt.Println()
	fmt.Println("Generating report...")

	benchReport := report.NewReport(version, sysInfo, results, runner.Duration())

	// Print text report to terminal
	textOutput := report.FormatText(benchReport)
	fmt.Print(textOutput)

	// Save JSON report
	jsonPath, err := report.SaveJSON(benchReport, *outputDir)
	if err != nil {
		fmt.Printf("Warning: Could not save JSON report: %v\n", err)
	} else {
		fmt.Printf("\nJSON report saved to: %s\n", jsonPath)
	}
}

func printHelp() {
	fmt.Printf(banner, version)
	fmt.Println()
	fmt.Println("Usage: ethbench [options]")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -test-dir string    Directory for disk I/O tests (default: executable directory)")
	fmt.Println("  -output string      Directory for JSON output file (default: executable directory)")
	fmt.Println("  -quick              Quick mode: ~1 minute benchmark instead of 3 minutes")
	fmt.Println("  -verbose            Show detailed progress during benchmarks")
	fmt.Println("  -help               Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ethbench                        Run full benchmark")
	fmt.Println("  ethbench -test-dir /mnt/nvme    Use specific directory for disk tests")
	fmt.Println("  ethbench -quick                 Run quick 1-minute benchmark")
	fmt.Println("  ethbench -output /home/user     Save JSON to specific directory")
	fmt.Println()
	fmt.Println("System Requirements:")
	fmt.Println("  - sysbench (sudo apt install sysbench)")
	fmt.Println("  - fio (sudo apt install fio)")
	fmt.Println()
}
