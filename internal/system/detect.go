// Package system provides hardware and OS detection functionality
package system

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

// Info contains system hardware and OS information
type Info struct {
	Hostname     string `json:"hostname"`
	SerialNumber string `json:"serial_number"`
	OS           string `json:"os"`
	OSVersion    string `json:"os_version"`
	Architecture string `json:"architecture"`
	CPUModel     string `json:"cpu_model"`
	CPUCores     int    `json:"cpu_cores"`
	RAMTotalMB   int    `json:"ram_total_mb"`
	DiskModel    string `json:"disk_model"`

	// Raspberry Pi specific
	RPiModel          string `json:"rpi_model,omitempty"`
	KernelVersion     string `json:"kernel_version,omitempty"`
	GPUFirmware       string `json:"gpu_firmware,omitempty"`
	BootloaderVersion string `json:"bootloader_version,omitempty"`
	CPUGovernor       string `json:"cpu_governor,omitempty"`
	CPUFreqMHz        int    `json:"cpu_freq_mhz,omitempty"`
	CoreVoltage       string `json:"core_voltage,omitempty"`
}

// Detect gathers system information
func Detect() (*Info, error) {
	info := &Info{
		Architecture: runtime.GOARCH,
		CPUCores:     runtime.NumCPU(),
	}

	// Get hostname
	hostname, err := os.Hostname()
	if err == nil {
		info.Hostname = hostname
	}

	// Get OS info
	info.OS, info.OSVersion = detectOS()

	// Get serial number (Raspberry Pi specific)
	info.SerialNumber = detectSerialNumber()

	// Get CPU model
	info.CPUModel = detectCPUModel()

	// Get RAM total
	info.RAMTotalMB = detectRAM()

	// Get disk model
	info.DiskModel = detectDiskModel()

	// Raspberry Pi specific detection
	info.RPiModel = detectRPiModel()
	info.KernelVersion = detectKernelVersion()
	info.GPUFirmware = detectGPUFirmware()
	info.BootloaderVersion = detectBootloaderVersion()
	info.CPUGovernor = detectCPUGovernor()
	info.CPUFreqMHz = detectCPUFrequency()
	info.CoreVoltage = detectCoreVoltage()

	return info, nil
}

// detectOS reads /etc/os-release to determine OS name and version
func detectOS() (name, version string) {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return runtime.GOOS, ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "NAME=") {
			name = strings.Trim(strings.TrimPrefix(line, "NAME="), "\"")
		}
		if strings.HasPrefix(line, "VERSION_ID=") {
			version = strings.Trim(strings.TrimPrefix(line, "VERSION_ID="), "\"")
		}
	}

	if name == "" {
		name = runtime.GOOS
	}
	return name, version
}

// detectSerialNumber gets the device serial number
// Works for Raspberry Pi and similar ARM devices
func detectSerialNumber() string {
	// Try devicetree serial (Raspberry Pi 4/5)
	paths := []string{
		"/sys/firmware/devicetree/base/serial-number",
		"/proc/device-tree/serial-number",
	}

	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err == nil {
			serial := strings.TrimSpace(string(data))
			// Remove null bytes
			serial = strings.ReplaceAll(serial, "\x00", "")
			if serial != "" {
				return serial
			}
		}
	}

	// Try /proc/cpuinfo for Serial field
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Serial") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1])
			}
		}
	}

	// Fallback: try to get machine-id
	data, err := os.ReadFile("/etc/machine-id")
	if err == nil {
		return strings.TrimSpace(string(data))
	}

	return "unknown"
}

// detectCPUModel reads CPU model from /proc/cpuinfo
func detectCPUModel() string {
	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return "unknown"
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Try different CPU model fields
		for _, prefix := range []string{"model name", "Model", "Hardware", "CPU implementer"} {
			if strings.HasPrefix(line, prefix) {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) == 2 {
					model := strings.TrimSpace(parts[1])
					if model != "" {
						return model
					}
				}
			}
		}
	}

	// Fallback for ARM
	if runtime.GOARCH == "arm64" {
		return "ARM64 Processor"
	}
	return "unknown"
}

// detectRAM reads total memory from /proc/meminfo
func detectRAM() int {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`MemTotal:\s+(\d+)\s+kB`)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			matches := re.FindStringSubmatch(line)
			if len(matches) == 2 {
				kb, err := strconv.Atoi(matches[1])
				if err == nil {
					return kb / 1024 // Convert to MB
				}
			}
		}
	}

	return 0
}

// detectDiskModel attempts to find the primary disk model
func detectDiskModel() string {
	// Look for NVMe devices first
	nvmeDevices, _ := filepath.Glob("/sys/block/nvme*")
	for _, dev := range nvmeDevices {
		modelPath := filepath.Join(dev, "device", "model")
		data, err := os.ReadFile(modelPath)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	// Look for SD cards (common on Raspberry Pi)
	sdDevices, _ := filepath.Glob("/sys/block/mmcblk*")
	for _, dev := range sdDevices {
		// Try to get card name
		namePath := filepath.Join(dev, "device", "name")
		data, err := os.ReadFile(namePath)
		if err == nil {
			return fmt.Sprintf("SD Card: %s", strings.TrimSpace(string(data)))
		}
	}

	// Look for SATA/SCSI devices
	sdaDevices, _ := filepath.Glob("/sys/block/sd*")
	for _, dev := range sdaDevices {
		modelPath := filepath.Join(dev, "device", "model")
		data, err := os.ReadFile(modelPath)
		if err == nil {
			return strings.TrimSpace(string(data))
		}
	}

	return "unknown"
}

// detectRPiModel reads Raspberry Pi model from device tree
func detectRPiModel() string {
	data, err := os.ReadFile("/proc/device-tree/model")
	if err != nil {
		return ""
	}
	model := strings.TrimSpace(string(data))
	model = strings.ReplaceAll(model, "\x00", "")
	return model
}

// detectKernelVersion reads kernel version
func detectKernelVersion() string {
	data, err := os.ReadFile("/proc/version")
	if err != nil {
		return ""
	}
	// Extract just the version part (e.g., "Linux version 6.1.0-rpi7-rpi-v8")
	parts := strings.Fields(string(data))
	if len(parts) >= 3 {
		return parts[2]
	}
	return strings.TrimSpace(string(data))
}

// detectGPUFirmware runs vcgencmd to get GPU firmware version
func detectGPUFirmware() string {
	cmd := exec.Command("vcgencmd", "version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	// Extract date from output (e.g., "Dec  5 2024 ...")
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "version") {
			return strings.TrimSpace(strings.TrimPrefix(line, "version"))
		}
		// First line usually contains date
		if len(line) > 0 && !strings.HasPrefix(line, "version") {
			return strings.TrimSpace(line)
		}
	}
	return strings.TrimSpace(string(output))
}

// detectBootloaderVersion runs vcgencmd to get bootloader version
func detectBootloaderVersion() string {
	cmd := exec.Command("vcgencmd", "bootloader_version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	// First line contains the date
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}
	return ""
}

// detectCPUGovernor reads current CPU scaling governor
func detectCPUGovernor() string {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_governor")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}

// detectCPUFrequency reads current CPU frequency in MHz
func detectCPUFrequency() int {
	data, err := os.ReadFile("/sys/devices/system/cpu/cpu0/cpufreq/scaling_cur_freq")
	if err != nil {
		return 0
	}
	// Value is in kHz, convert to MHz
	freqKHz, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0
	}
	return freqKHz / 1000
}

// detectCoreVoltage runs vcgencmd to get core voltage
func detectCoreVoltage() string {
	cmd := exec.Command("vcgencmd", "measure_volts", "core")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	// Output is like "volt=0.8700V"
	result := strings.TrimSpace(string(output))
	result = strings.TrimPrefix(result, "volt=")
	return result
}

// CheckPrerequisites verifies that required tools are available
func CheckPrerequisites(testDir string) error {
	// Check if test directory exists or can be created
	if err := os.MkdirAll(testDir, 0755); err != nil {
		return fmt.Errorf("cannot create test directory %s: %w", testDir, err)
	}

	// Check write permissions
	testFile := filepath.Join(testDir, ".ethbench_test")
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("cannot write to test directory %s: %w", testDir, err)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}
