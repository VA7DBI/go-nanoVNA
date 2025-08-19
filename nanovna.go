package nanovna

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tarm/serial"
)

// HardwareVariant represents different NanoVNA hardware versions
type HardwareVariant int

const (
	VariantUnknown HardwareVariant = iota
	VariantV1                      // Original NanoVNA v1
	VariantVH                      // NanoVNA-H (Hardware version)
	VariantV2                      // NanoVNA v2 (SAA2)
	VariantV2Plus                  // NanoVNA v2 Plus
	VariantV2Plus4                 // NanoVNA v2 Plus4
	VariantSAA2                    // Standalone SAA2
	VariantTinysa                  // TinySA variant
	VariantLiteVNA                 // LiteVNA variant
)

// String returns the string representation of the hardware variant
func (hv HardwareVariant) String() string {
	switch hv {
	case VariantV1:
		return "NanoVNA v1"
	case VariantVH:
		return "NanoVNA-H"
	case VariantV2:
		return "NanoVNA v2"
	case VariantV2Plus:
		return "NanoVNA v2 Plus"
	case VariantV2Plus4:
		return "NanoVNA v2 Plus4"
	case VariantSAA2:
		return "SAA2"
	case VariantTinysa:
		return "TinySA"
	case VariantLiteVNA:
		return "LiteVNA"
	default:
		return "Unknown"
	}
}

// HardwareInfo contains hardware-specific information and capabilities
type HardwareInfo struct {
	Variant        HardwareVariant
	FrequencyRange FrequencyRange
	MaxSweepPoints int
	SupportedPorts []string // S11, S21, S12, S22
	CommandSet     CommandSet
	Capabilities   HardwareCapabilities
}

// FrequencyRange defines the frequency range for a hardware variant
type FrequencyRange struct {
	MinHz float64
	MaxHz float64
}

// CommandSet defines the command set for different hardware variants
type CommandSet struct {
	SweepCommand    string
	FreqCommand     string
	DataCommand     string
	InfoCommand     string
	VersionCommand  string
	CalibrationSave string
	CalibrationLoad string
	PromptPattern   string
}

// HardwareCapabilities defines what each hardware variant can do
type HardwareCapabilities struct {
	HasS21           bool
	HasTimeDomain    bool
	HasCalibration   bool
	HasMultiplePorts bool
	HasGenerator     bool
	HasSpectrumMode  bool
}

// getHardwareInfo returns hardware information for a given variant
func getHardwareInfo(variant HardwareVariant) HardwareInfo {
	switch variant {
	case VariantV1:
		return HardwareInfo{
			Variant:        VariantV1,
			FrequencyRange: FrequencyRange{MinHz: 50000, MaxHz: 900000000}, // 50kHz - 900MHz
			MaxSweepPoints: 101,
			SupportedPorts: []string{"S11", "S21"},
			CommandSet: CommandSet{
				SweepCommand:    "sweep %d %d %d",
				FreqCommand:     "frequencies",
				DataCommand:     "data %d",
				InfoCommand:     "info",
				VersionCommand:  "version",
				CalibrationSave: "save %d",
				CalibrationLoad: "recall %d",
				PromptPattern:   "ch>",
			},
			Capabilities: HardwareCapabilities{
				HasS21:           true,
				HasTimeDomain:    false,
				HasCalibration:   true,
				HasMultiplePorts: false,
				HasGenerator:     false,
				HasSpectrumMode:  false,
			},
		}
	case VariantVH:
		return HardwareInfo{
			Variant:        VariantVH,
			FrequencyRange: FrequencyRange{MinHz: 50000, MaxHz: 1500000000}, // 50kHz - 1.5GHz
			MaxSweepPoints: 201,
			SupportedPorts: []string{"S11", "S21"},
			CommandSet: CommandSet{
				SweepCommand:    "sweep %d %d %d",
				FreqCommand:     "frequencies",
				DataCommand:     "data %d",
				InfoCommand:     "info",
				VersionCommand:  "version",
				CalibrationSave: "save %d",
				CalibrationLoad: "recall %d",
				PromptPattern:   "ch>",
			},
			Capabilities: HardwareCapabilities{
				HasS21:           true,
				HasTimeDomain:    true,
				HasCalibration:   true,
				HasMultiplePorts: false,
				HasGenerator:     true,
				HasSpectrumMode:  false,
			},
		}
	case VariantV2:
		return HardwareInfo{
			Variant:        VariantV2,
			FrequencyRange: FrequencyRange{MinHz: 50000, MaxHz: 3000000000}, // 50kHz - 3GHz
			MaxSweepPoints: 4000,
			SupportedPorts: []string{"S11", "S21"},
			CommandSet: CommandSet{
				SweepCommand:    "sweep %d %d %d",
				FreqCommand:     "freq",
				DataCommand:     "data %d",
				InfoCommand:     "info",
				VersionCommand:  "version",
				CalibrationSave: "save %d",
				CalibrationLoad: "recall %d",
				PromptPattern:   "2>",
			},
			Capabilities: HardwareCapabilities{
				HasS21:           true,
				HasTimeDomain:    true,
				HasCalibration:   true,
				HasMultiplePorts: false,
				HasGenerator:     true,
				HasSpectrumMode:  true,
			},
		}
	case VariantV2Plus:
		return HardwareInfo{
			Variant:        VariantV2Plus,
			FrequencyRange: FrequencyRange{MinHz: 50000, MaxHz: 6000000000}, // 50kHz - 6GHz
			MaxSweepPoints: 4000,
			SupportedPorts: []string{"S11", "S21"},
			CommandSet: CommandSet{
				SweepCommand:    "sweep %d %d %d",
				FreqCommand:     "freq",
				DataCommand:     "data %d",
				InfoCommand:     "info",
				VersionCommand:  "version",
				CalibrationSave: "save %d",
				CalibrationLoad: "recall %d",
				PromptPattern:   "2>",
			},
			Capabilities: HardwareCapabilities{
				HasS21:           true,
				HasTimeDomain:    true,
				HasCalibration:   true,
				HasMultiplePorts: false,
				HasGenerator:     true,
				HasSpectrumMode:  true,
			},
		}
	case VariantV2Plus4:
		return HardwareInfo{
			Variant:        VariantV2Plus4,
			FrequencyRange: FrequencyRange{MinHz: 50000, MaxHz: 6000000000}, // 50kHz - 6GHz
			MaxSweepPoints: 4000,
			SupportedPorts: []string{"S11", "S21", "S12", "S22"},
			CommandSet: CommandSet{
				SweepCommand:    "sweep %d %d %d",
				FreqCommand:     "freq",
				DataCommand:     "data %d",
				InfoCommand:     "info",
				VersionCommand:  "version",
				CalibrationSave: "save %d",
				CalibrationLoad: "recall %d",
				PromptPattern:   "2>",
			},
			Capabilities: HardwareCapabilities{
				HasS21:           true,
				HasTimeDomain:    true,
				HasCalibration:   true,
				HasMultiplePorts: true,
				HasGenerator:     true,
				HasSpectrumMode:  true,
			},
		}
	case VariantTinysa:
		return HardwareInfo{
			Variant:        VariantTinysa,
			FrequencyRange: FrequencyRange{MinHz: 100000, MaxHz: 960000000}, // 100kHz - 960MHz
			MaxSweepPoints: 500,
			SupportedPorts: []string{"S11"},
			CommandSet: CommandSet{
				SweepCommand:    "sweep %d %d %d",
				FreqCommand:     "frequencies",
				DataCommand:     "data %d",
				InfoCommand:     "info",
				VersionCommand:  "version",
				CalibrationSave: "save %d",
				CalibrationLoad: "recall %d",
				PromptPattern:   "ch>",
			},
			Capabilities: HardwareCapabilities{
				HasS21:           false,
				HasTimeDomain:    false,
				HasCalibration:   true,
				HasMultiplePorts: false,
				HasGenerator:     true,
				HasSpectrumMode:  true,
			},
		}
	default:
		// Default/unknown hardware - use conservative settings
		return HardwareInfo{
			Variant:        VariantUnknown,
			FrequencyRange: FrequencyRange{MinHz: 50000, MaxHz: 900000000},
			MaxSweepPoints: 101,
			SupportedPorts: []string{"S11"},
			CommandSet: CommandSet{
				SweepCommand:    "sweep %d %d %d",
				FreqCommand:     "frequencies",
				DataCommand:     "data %d",
				InfoCommand:     "info",
				VersionCommand:  "version",
				CalibrationSave: "save %d",
				CalibrationLoad: "recall %d",
				PromptPattern:   "ch>",
			},
			Capabilities: HardwareCapabilities{
				HasS21:           false,
				HasTimeDomain:    false,
				HasCalibration:   true,
				HasMultiplePorts: false,
				HasGenerator:     false,
				HasSpectrumMode:  false,
			},
		}
	}
}

// PortConfig holds serial port configuration details for debugging
type PortConfig struct {
	Name        string
	Baud        int
	ReadTimeout time.Duration
	Size        byte
	Parity      serial.Parity
	StopBits    serial.StopBits
}

// SerialPort is the interface for serial port operations (exported for debug wrapping)
type SerialPort interface {
	Write([]byte) (int, error)
	Read([]byte) (int, error)
	Close() error
}

type Device struct {
	Port         string
	portHandle   SerialPort
	config       *PortConfig     // Store configuration for debugging
	version      string          // Store detected version string (v1, vh, v2, etc.)
	variant      HardwareVariant // Store hardware variant enum
	hardwareInfo HardwareInfo    // Store hardware capabilities and info
}

// SetPortHandle allows replacing the underlying serial port (for debug wrapping)
func (d *Device) SetPortHandle(sp SerialPort) {
	d.portHandle = sp
}

// GetPortHandle returns the underlying serial port (for debug wrapping)
func (d *Device) GetPortHandle() SerialPort {
	return d.portHandle
}

// GetPortConfig returns the port configuration details (for debugging)
func (d *Device) GetPortConfig() *PortConfig {
	return d.config
}

// GetPortDetails returns detailed port information as a formatted string
func (d *Device) GetPortDetails() string {
	if d.config == nil {
		return fmt.Sprintf("Port: %s (config not available)", d.Port)
	}
	return fmt.Sprintf("Port: %s, Baud: %d, ReadTimeout: %v, Size: %d, Parity: %v, StopBits: %v",
		d.config.Name, d.config.Baud, d.config.ReadTimeout,
		d.config.Size, d.config.Parity, d.config.StopBits)
}

// GetHardwareVariant returns the detected hardware variant
func (d *Device) GetHardwareVariant() HardwareVariant {
	return d.variant
}

// GetHardwareInfo returns the hardware information and capabilities
func (d *Device) GetHardwareInfo() HardwareInfo {
	return d.hardwareInfo
}

// GetFrequencyRange returns the supported frequency range for this hardware
func (d *Device) GetFrequencyRange() FrequencyRange {
	return d.hardwareInfo.FrequencyRange
}

// GetMaxSweepPoints returns the maximum number of sweep points supported
func (d *Device) GetMaxSweepPoints() int {
	return d.hardwareInfo.MaxSweepPoints
}

// GetSupportedPorts returns the supported S-parameter ports
func (d *Device) GetSupportedPorts() []string {
	return d.hardwareInfo.SupportedPorts
}

// GetCapabilities returns the hardware capabilities
func (d *Device) GetCapabilities() HardwareCapabilities {
	return d.hardwareInfo.Capabilities
}

// IsPortSupported checks if a specific S-parameter port is supported
func (d *Device) IsPortSupported(port string) bool {
	for _, p := range d.hardwareInfo.SupportedPorts {
		if p == port {
			return true
		}
	}
	return false
}

type DeviceInfo struct {
	Model     string
	Firmware  string
	SerialNum string
}

type SweepData struct {
	Frequencies []float64
	S11         []complex128
	S21         []complex128
}

type CalibrationData struct {
	// TODO: define calibration fields
}

// ListDevices lists available NanoVNA serial ports (Windows only, stub).
func ListDevices() ([]string, error) {
	var ports []string
	for i := 1; i <= 20; i++ {
		portName := fmt.Sprintf("COM%d", i)
		f, err := os.Open("//./" + portName)
		if err == nil {
			ports = append(ports, portName)
			f.Close()
		}
	}
	if len(ports) == 0 {
		return nil, errors.New("no serial ports found")
	}
	return ports, nil
}

// AutoDetect attempts to find and connect to a NanoVNA device automatically
func AutoDetect() (*Device, error) {
	ports, err := ListDevices()
	if err != nil {
		return nil, fmt.Errorf("failed to list serial ports: %v", err)
	}

	for _, port := range ports {
		device, err := Open(port)
		if err != nil {
			continue // Try next port
		}

		// Try to detect the version/hardware
		_, err = device.DetectVersion()
		if err != nil {
			device.Close()
			continue // Try next port
		}

		// Successfully detected a NanoVNA
		return device, nil
	}

	return nil, errors.New("no NanoVNA devices found on any serial port")
}

// OpenWithVariant opens a device and forces a specific hardware variant
// (useful for testing or when auto-detection fails)
func OpenWithVariant(port string, variant HardwareVariant) (*Device, error) {
	device, err := Open(port)
	if err != nil {
		return nil, err
	}

	// Override the detected variant
	device.variant = variant
	device.hardwareInfo = getHardwareInfo(variant)

	// Set version string based on variant
	switch variant {
	case VariantV1:
		device.version = "v1"
	case VariantVH:
		device.version = "vh"
	case VariantV2, VariantV2Plus, VariantV2Plus4, VariantSAA2:
		device.version = "v2"
	case VariantTinysa:
		device.version = "tinysa"
	case VariantLiteVNA:
		device.version = "litevna"
	default:
		device.version = "unknown"
	}

	return device, nil
}

// Open connects to a NanoVNA on the specified serial port. Optionally accepts a custom SerialPort for debug/testing.
func Open(port string, custom ...SerialPort) (*Device, error) {
	device := &Device{Port: port}

	if len(custom) > 0 && custom[0] != nil {
		device.portHandle = custom[0]
	} else {
		// Set a 5-second read timeout to prevent hanging
		c := &serial.Config{
			Name:        port,
			Baud:        9600,
			ReadTimeout: time.Second * 5,
			Size:        8,
			Parity:      serial.ParityNone,
			StopBits:    serial.Stop1,
		}
		s, err := serial.OpenPort(c)
		if err != nil {
			return nil, err
		}

		// Store configuration for debugging
		device.config = &PortConfig{
			Name:        c.Name,
			Baud:        c.Baud,
			ReadTimeout: c.ReadTimeout,
			Size:        c.Size,
			Parity:      c.Parity,
			StopBits:    c.StopBits,
		}

		device.portHandle = s
	}

	// Initialize with unknown hardware until detection
	device.variant = VariantUnknown
	device.hardwareInfo = getHardwareInfo(VariantUnknown)

	return device, nil
}

// SetSweepConfig configures sweep parameters (start, stop, points).
func (d *Device) SetSweepConfig(startHz, stopHz int, points int) error {
	// Validate frequency range against hardware capabilities
	if float64(startHz) < d.hardwareInfo.FrequencyRange.MinHz {
		return fmt.Errorf("start frequency %d Hz is below minimum %g Hz for %s",
			startHz, d.hardwareInfo.FrequencyRange.MinHz, d.variant.String())
	}
	if float64(stopHz) > d.hardwareInfo.FrequencyRange.MaxHz {
		return fmt.Errorf("stop frequency %d Hz is above maximum %g Hz for %s",
			stopHz, d.hardwareInfo.FrequencyRange.MaxHz, d.variant.String())
	}
	if points > d.hardwareInfo.MaxSweepPoints {
		return fmt.Errorf("requested %d points exceeds maximum %d for %s",
			points, d.hardwareInfo.MaxSweepPoints, d.variant.String())
	}

	// Use hardware-specific sweep command
	cmd := fmt.Sprintf(d.hardwareInfo.CommandSet.SweepCommand, startHz, stopHz, points)

	// For some hardware variants, we need to send additional commands
	switch d.variant {
	case VariantV2, VariantV2Plus, VariantV2Plus4:
		// V2 variants might need a different command sequence
		_, err := d.sendCommand(cmd)
		if err != nil {
			// Try alternative V2 command format
			altCmd := fmt.Sprintf("sweep start %d", startHz)
			if _, err2 := d.sendCommand(altCmd); err2 != nil {
				altCmd = fmt.Sprintf("sweep stop %d", stopHz)
				if _, err3 := d.sendCommand(altCmd); err3 != nil {
					altCmd = fmt.Sprintf("sweep points %d", points)
					if _, err4 := d.sendCommand(altCmd); err4 != nil {
						return fmt.Errorf("failed to set sweep config with any V2 command format: %v", err)
					}
				}
			}
		}
		return nil
	default:
		// Standard command for V1, VH, and other variants
		_, err := d.sendCommand(cmd)
		if err != nil {
			// Fallback to individual commands
			if _, err2 := d.sendCommand(fmt.Sprintf("start %d", startHz)); err2 != nil {
				if _, err3 := d.sendCommand(fmt.Sprintf("stop %d", stopHz)); err3 != nil {
					if _, err4 := d.sendCommand(fmt.Sprintf("points %d", points)); err4 != nil {
						return fmt.Errorf("failed to set sweep config: %v", err)
					}
				}
			}
		}
		return nil
	}
}

// RunSweep triggers a sweep and returns measurement data.
// Uses hardware-specific commands and handles different port configurations.
func (d *Device) RunSweep() (SweepData, error) {
	var data SweepData

	// Step 1: Get frequencies using hardware-specific command
	freqCmd := d.hardwareInfo.CommandSet.FreqCommand
	freqResp, err := d.sendCommand(freqCmd)
	if err != nil {
		return SweepData{}, fmt.Errorf("failed to get frequencies: %v", err)
	}

	// Parse frequencies
	freqLines := strings.Split(freqResp, "\n")
	for _, line := range freqLines {
		line = strings.TrimSpace(line)
		if line == "" || line == freqCmd ||
			strings.Contains(line, d.hardwareInfo.CommandSet.PromptPattern) ||
			strings.Contains(line, "?") {
			continue
		}

		freq, err := strconv.ParseFloat(line, 64)
		if err == nil {
			data.Frequencies = append(data.Frequencies, freq)
		}
	}

	// Step 2: Get S11 data (always available)
	s11Cmd := fmt.Sprintf(d.hardwareInfo.CommandSet.DataCommand, 0)
	s11Resp, err := d.sendCommand(s11Cmd)
	if err != nil {
		return SweepData{}, fmt.Errorf("failed to get S11 data: %v", err)
	}

	// Parse S11 data (complex numbers: real imaginary)
	s11Lines := strings.Split(s11Resp, "\n")
	for _, line := range s11Lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "data") ||
			strings.Contains(line, d.hardwareInfo.CommandSet.PromptPattern) ||
			strings.Contains(line, "?") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) >= 2 {
			real, err1 := strconv.ParseFloat(parts[0], 64)
			imag, err2 := strconv.ParseFloat(parts[1], 64)
			if err1 == nil && err2 == nil {
				data.S11 = append(data.S11, complex(real, imag))
			}
		}
	}

	// Step 3: Get S21 data if supported
	if d.hardwareInfo.Capabilities.HasS21 && d.IsPortSupported("S21") {
		s21Cmd := fmt.Sprintf(d.hardwareInfo.CommandSet.DataCommand, 1)
		s21Resp, err := d.sendCommand(s21Cmd)
		if err != nil {
			// S21 might not be available, create dummy data
			for range data.S11 {
				data.S21 = append(data.S21, complex(0, 0))
			}
		} else {
			// Parse S21 data
			s21Lines := strings.Split(s21Resp, "\n")
			for _, line := range s21Lines {
				line = strings.TrimSpace(line)
				if line == "" || strings.HasPrefix(line, "data") ||
					strings.Contains(line, d.hardwareInfo.CommandSet.PromptPattern) ||
					strings.Contains(line, "?") {
					continue
				}

				parts := strings.Fields(line)
				if len(parts) >= 2 {
					real, err1 := strconv.ParseFloat(parts[0], 64)
					imag, err2 := strconv.ParseFloat(parts[1], 64)
					if err1 == nil && err2 == nil {
						data.S21 = append(data.S21, complex(real, imag))
					}
				}
			}
		}
	}

	// Ensure S21 has same length as S11 (pad with zeros if needed)
	for len(data.S21) < len(data.S11) {
		data.S21 = append(data.S21, complex(0, 0))
	}

	// Validate we got some data
	if len(data.Frequencies) == 0 || len(data.S11) == 0 {
		return SweepData{}, fmt.Errorf("no valid measurement data received")
	}

	// Ensure all arrays have the same length
	minLen := len(data.Frequencies)
	if len(data.S11) < minLen {
		minLen = len(data.S11)
	}

	data.Frequencies = data.Frequencies[:minLen]
	data.S11 = data.S11[:minLen]
	if len(data.S21) >= minLen {
		data.S21 = data.S21[:minLen]
	} else {
		// Pad S21 with zeros if needed
		for len(data.S21) < minLen {
			data.S21 = append(data.S21, complex(0, 0))
		}
	}

	return data, nil
}

// Close disconnects from the device.
func (d *Device) Close() error {
	if d.portHandle != nil {
		err := d.portHandle.Close()
		d.portHandle = nil
		return err
	}
	return nil
}

// sendCommand sends a command string to the NanoVNA and returns the response.
// Uses proper protocol based on detected version.
func (d *Device) sendCommand(cmd string) (string, error) {
	if d.portHandle == nil {
		return "", errors.New("device not open")
	}

	// Clear any existing data first
	buf := make([]byte, 1024)
	d.portHandle.Read(buf) // drain buffer

	// Send command with proper termination
	cmdBytes := []byte(cmd + "\r")
	_, err := d.portHandle.Write(cmdBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write command: %v", err)
	}

	// Small delay after sending to let device process
	time.Sleep(50 * time.Millisecond)

	// Read response with proper parsing for NanoVNA protocol
	var response strings.Builder
	maxAttempts := 10

	for attempts := 0; attempts < maxAttempts; attempts++ {
		n, err := d.portHandle.Read(buf)
		if err != nil {
			if strings.Contains(err.Error(), "timeout") && response.Len() > 0 {
				break // We got some data, timeout is OK
			}
			return response.String(), err
		}

		if n > 0 {
			chunk := string(buf[:n])
			response.WriteString(chunk)

			// Check if we've received the command prompt indicating end of response
			if strings.Contains(chunk, "ch>") {
				break
			}
		}

		// Small delay between reads
		time.Sleep(20 * time.Millisecond)
	}

	return response.String(), nil
}

// GetInfo retrieves device information (model, firmware, serial number).
func (d *Device) GetInfo() (DeviceInfo, error) {
	// Use hardware-specific info command
	infoCmd := d.hardwareInfo.CommandSet.InfoCommand
	resp, err := d.sendCommand(infoCmd)
	if err != nil {
		return DeviceInfo{}, err
	}

	// Parse the response
	lines := strings.Split(resp, "\n")
	var info DeviceInfo

	// Set default model based on detected hardware
	info.Model = d.variant.String()

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || line == infoCmd ||
			strings.Contains(line, d.hardwareInfo.CommandSet.PromptPattern) {
			continue // Skip echo, empty lines, and command prompt
		}

		// Look for specific patterns based on hardware variant
		switch d.variant {
		case VariantV2, VariantV2Plus, VariantV2Plus4, VariantSAA2:
			// V2 variants have different info format
			if strings.Contains(strings.ToLower(line), "nanovna") ||
				strings.Contains(strings.ToLower(line), "saa2") {
				info.Model = line
			}
			if strings.Contains(strings.ToLower(line), "firmware") ||
				strings.Contains(strings.ToLower(line), "version") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					info.Firmware = strings.TrimSpace(parts[1])
				}
			}
		default:
			// V1, VH, and other variants
			// First non-empty line is usually the model/version info
			if info.Model == d.variant.String() {
				info.Model = line
			}

			// Look for serial number
			if strings.HasPrefix(strings.ToLower(line), "serial") {
				info.SerialNum = strings.TrimSpace(strings.TrimPrefix(line, "Serial:"))
			}

			// Look for firmware version patterns
			if strings.Contains(line, "v") && info.Firmware == "" {
				parts := strings.Fields(line)
				for _, p := range parts {
					if strings.HasPrefix(p, "v") && len(p) > 1 {
						info.Firmware = p
						break
					}
				}
			}
		}
	}

	// If we didn't get proper model info, use detected variant
	if info.Model == "" || info.Model == d.variant.String() {
		info.Model = fmt.Sprintf("%s (detected)", d.variant.String())
	}

	return info, nil
}

// DetectVersion detects the NanoVNA version by sending CR and analyzing the response
func (d *Device) DetectVersion() (string, error) {
	if d.portHandle == nil {
		return "", errors.New("device not open")
	}

	// Clear any existing data
	buf := make([]byte, 1024)
	d.portHandle.Read(buf) // drain buffer

	// Send carriage return to detect version
	_, err := d.portHandle.Write([]byte("\r"))
	if err != nil {
		return "", err
	}

	// Read response with timeout
	time.Sleep(50 * time.Millisecond)
	n, err := d.portHandle.Read(buf)
	if err != nil {
		return "", err
	}

	response := string(buf[:n])

	// Try to get more info to distinguish between variants
	info, _ := d.sendCommand("info")

	// Detect hardware variant based on response patterns and info
	if strings.HasPrefix(response, "ch> ") {
		d.version = "v1"
		d.variant = VariantV1
		if strings.Contains(strings.ToLower(info), "tinysa") {
			d.variant = VariantTinysa
		} else if strings.Contains(strings.ToLower(info), "litevna") {
			d.variant = VariantLiteVNA
		}
	} else if strings.HasPrefix(response, "\r\nch> ") || strings.HasPrefix(response, "\r\n?\r\nch> ") {
		d.version = "vh"
		d.variant = VariantVH
		// Check if it's actually a v1 with different prompt
		if strings.Contains(strings.ToLower(info), "nanovna v1") {
			d.variant = VariantV1
		}
	} else if strings.HasPrefix(response, "2") || strings.Contains(response, "2>") {
		d.version = "v2"
		d.variant = VariantV2
		// Distinguish between v2 variants based on info
		if strings.Contains(strings.ToLower(info), "plus4") {
			d.variant = VariantV2Plus4
		} else if strings.Contains(strings.ToLower(info), "plus") {
			d.variant = VariantV2Plus
		} else if strings.Contains(strings.ToLower(info), "saa2") {
			d.variant = VariantSAA2
		}
	} else {
		d.version = "unknown"
		d.variant = VariantUnknown
	}

	// Get hardware info for detected variant
	d.hardwareInfo = getHardwareInfo(d.variant)

	if d.variant == VariantUnknown {
		return "unknown", fmt.Errorf("unrecognized response: %q", response)
	}

	return d.version, nil
}

// GetVersion returns the detected NanoVNA version
func (d *Device) GetVersion() string {
	return d.version
}

// GetCalibration retrieves current calibration data.
func (d *Device) GetCalibration() (CalibrationData, error) {
	return CalibrationData{}, nil
}

// SetCalibration applies calibration data.
func (d *Device) SetCalibration(cal CalibrationData) error {
	return nil
}

// SaveCalibration saves calibration data to device memory.
func (d *Device) SaveCalibration(slot int) error {
	return nil
}

// LoadCalibration loads calibration data from device memory.
func (d *Device) LoadCalibration(slot int) error {
	return nil
}
