
# go-nanovna

<!-- Badges -->

[![Go](https://github.com/VA7DBI/go-nanovna/workflows/Go/badge.svg)](https://github.com/VA7DBI/go-nanovna/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/VA7DBI/go-nanovna.svg)](https://pkg.go.dev/github.com/VA7DBI/go-nanovna)
[![Go Report Card](https://goreportcard.com/badge/github.com/VA7DBI/go-nanovna)](https://goreportcard.com/report/github.com/VA7DBI/go-nanovna)
![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
[![Codecov](https://codecov.io/gh/VA7DBI/go-nanovna/branch/main/graph/badge.svg)](https://codecov.io/gh/VA7DBI/go-nanovna)

---

**Documentation:** See [docs/specification.md](docs/specification.md) for full API and protocol details, including diagrams and flowcharts.

---

A comprehensive Go package for controlling and communicating with NanoVNA (Vector Network Analyzer) devices.

## Features

- **Universal Hardware Support**: Supports all major NanoVNA hardware variants with automatic detection
- **Hardware Abstraction**: Unified interface that adapts to different hardware capabilities
- **Comprehensive Coverage**: Supports V1, VH, V2, V2+, V2+4, SAA2, TinySA, and LiteVNA variants
- **Smart Detection**: Automatic hardware variant detection and capability mapping
- **Frequency Validation**: Hardware-aware frequency range and sweep point validation
- **Debug Support**: Built-in serial debugging and hardware information display

## Supported Hardware Variants

| Variant | Frequency Range | Max Points | S-Parameters | Special Features |
|---------|----------------|------------|--------------|------------------|
| NanoVNA v1 | 50 kHz - 900 MHz | 101 | S11, S21 | Basic VNA |
| NanoVNA-H | 50 kHz - 1.5 GHz | 201 | S11, S21 | Time domain, Generator |
| NanoVNA v2 | 50 kHz - 3 GHz | 4000 | S11, S21 | High resolution, Spectrum |
| NanoVNA v2+ | 50 kHz - 6 GHz | 4000 | S11, S21 | Extended frequency |
| NanoVNA v2+4 | 50 kHz - 6 GHz | 4000 | S11, S21, S12, S22 | 4-port measurements |
| TinySA | 100 kHz - 960 MHz | 500 | S11 | Spectrum analyzer focus |
| LiteVNA | 100 kHz - 6.5 GHz | 4000 | S11, S21 | Extended range |

## Installation

```bash
go get github.com/VA7DBI/go-nanovna
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "github.com/VA7DBI/go-nanovna"
)

func main() {
    // Auto-detect and connect to NanoVNA
    device, err := nanovna.AutoDetect()
    if err != nil {
        log.Fatal("Failed to detect NanoVNA:", err)
    }
    defer device.Close()

    // Get hardware information
    info := device.GetHardwareInfo()
    fmt.Printf("Connected to: %s\n", info.Variant.String())
    fmt.Printf("Frequency Range: %.0f Hz - %.0f Hz\n", 
        info.FrequencyRange.MinHz, info.FrequencyRange.MaxHz)
    
    // Configure sweep
    err = device.SetSweepConfig(144000000, 146000000, 101)
    if err != nil {
        log.Fatal("Failed to configure sweep:", err)
    }
    
    // Run measurement
    data, err := device.RunSweep()
    if err != nil {
        log.Fatal("Failed to run sweep:", err)
    }
    
    fmt.Printf("Measured %d points\n", len(data.Frequencies))
}
```

## Manual Connection

```go
// Connect to specific port
device, err := nanovna.Open("COM3")
if err != nil {
    log.Fatal(err)
}

// Detect hardware variant
version, err := device.DetectVersion()
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Detected: %s\n", version)
```

## Hardware-Aware Programming

```go
// Check hardware capabilities
caps := device.GetCapabilities()
if caps.HasS21 {
    fmt.Println("S21 measurements supported")
}
if caps.HasTimeDomain {
    fmt.Println("Time domain analysis supported")
}

// Validate frequency range
freqRange := device.GetFrequencyRange()
if startFreq < freqRange.MinHz {
    log.Fatal("Start frequency too low for this hardware")
}
```

## Testing & Coverage

Run all tests:

```sh
go test ./...
```

Run with coverage (CLI):

```sh
go test -cover ./...
```

Generate a local HTML coverage report:

```sh
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

This will open a browser window with a detailed coverage report.

### Codecov badge

The badge above is powered by [Codecov](https://codecov.io/gh/VA7DBI/go-nanovna). To enable it for your fork or private repo:

1. Sign up at [codecov.io](https://codecov.io/) and enable your repository.
2. (If private) Add your Codecov token as a GitHub secret named `CODECOV_TOKEN`.
3. Push to main or open a pull request to trigger the workflow.

---

## API Reference

### Device Management

- AutoDetect() (*Device, error) - Auto-detect and connect to NanoVNA
- Open(port string) (*Device, error) - Connect to specific serial port
- OpenWithVariant(port, variant) - Force specific hardware variant
- ListDevices() ([]string, error) - List available serial ports

### Hardware Information

- GetHardwareVariant() HardwareVariant - Get detected hardware type
- GetHardwareInfo() HardwareInfo - Get complete hardware information
- GetFrequencyRange() FrequencyRange - Get supported frequency range
- GetCapabilities() HardwareCapabilities - Get hardware capabilities

### Measurements

- SetSweepConfig(start, stop, points int) error - Configure sweep parameters
- RunSweep() (SweepData, error) - Perform measurement sweep
- GetInfo() (DeviceInfo, error) - Get device information

### Data Structures

```go
type SweepData struct {
    Frequencies []float64      // Frequency points
    S11         []complex128   // S11 measurements
    S21         []complex128   // S21 measurements (if supported)
}

type HardwareCapabilities struct {
    HasS21           bool  // S21 transmission measurements
    HasTimeDomain    bool  // Time domain analysis
    HasCalibration   bool  // Calibration support
    HasMultiplePorts bool  // 4-port measurements
    HasGenerator     bool  // Signal generator
    HasSpectrumMode  bool  // Spectrum analyzer mode
}
```

## License

MIT License - see LICENSE file for details.

## Contributing

Contributions welcome! Please submit issues and pull requests on GitHub.
