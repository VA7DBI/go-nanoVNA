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
fmt.Printf("Max Sweep Points: %d\n", info.MaxSweepPoints)
fmt.Printf("Supported Ports: %v\n", info.SupportedPorts)

// Show capabilities
caps := info.Capabilities
fmt.Println("\nCapabilities:")
fmt.Printf("  S21 Transmission: %t\n", caps.HasS21)
fmt.Printf("  Time Domain: %t\n", caps.HasTimeDomain)
fmt.Printf("  Calibration: %t\n", caps.HasCalibration)
fmt.Printf("  Multiple Ports: %t\n", caps.HasMultiplePorts)
fmt.Printf("  Signal Generator: %t\n", caps.HasGenerator)
fmt.Printf("  Spectrum Mode: %t\n", caps.HasSpectrumMode)

// Configure sweep (2m amateur band)
startHz := 144000000  // 144 MHz
stopHz := 148000000   // 148 MHz
points := 101

err = device.SetSweepConfig(startHz, stopHz, points)
if err != nil {
log.Fatal("Failed to configure sweep:", err)
}

fmt.Printf("\nRunning sweep: %d - %d Hz, %d points\n", startHz, stopHz, points)

// Run measurement
data, err := device.RunSweep()
if err != nil {
log.Fatal("Failed to run sweep:", err)
}

fmt.Printf("Measured %d frequency points\n", len(data.Frequencies))

// Show first few measurement points
fmt.Println("\nFirst 5 measurement points:")
for i := 0; i < 5 && i < len(data.Frequencies); i++ {
freq := data.Frequencies[i]
s11 := data.S11[i]
fmt.Printf("  %.1f MHz: S11 = %.6f %+.6fi\n", 
freq/1e6, real(s11), imag(s11))
}
}
