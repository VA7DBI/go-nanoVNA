package nanovna

import (
	"errors"
	"testing"
)

// MockSerialPort implements the SerialPort interface for testing
type MockSerialPort struct {
	WriteBuffer []byte
	ReadBuffer  []byte
	Closed      bool
	ReadIndex   int
}

func (m *MockSerialPort) Write(p []byte) (int, error) {
	m.WriteBuffer = append(m.WriteBuffer, p...)
	return len(p), nil
}

func (m *MockSerialPort) Read(p []byte) (int, error) {
	if m.ReadIndex >= len(m.ReadBuffer) {
		return 0, errors.New("timeout")
	}
	n := copy(p, m.ReadBuffer[m.ReadIndex:])
	m.ReadIndex += n
	return n, nil
}

func (m *MockSerialPort) Close() error {
	m.Closed = true
	return nil
}

func TestDevice_IsPortSupported_TableDriven(t *testing.T) {
	dev := &Device{
		hardwareInfo: getHardwareInfo(VariantV2Plus4),
	}
	tests := []struct {
		port   string
		expect bool
	}{
		{"S21", true},
		{"S11", true},
		{"S12", true},
		{"S22", true},
		{"FOO", false},
		{"", false},
	}
	for _, tc := range tests {
		got := dev.IsPortSupported(tc.port)
		if got != tc.expect {
			t.Errorf("IsPortSupported(%q) = %v, want %v", tc.port, got, tc.expect)
		}
	}
}

func TestDevice_GetPortConfigAndDetails(t *testing.T) {
	// Try to open a real port (COM3), fallback to mock if it fails
	dev, err := Open("COM3")
	if err != nil {
		t.Log("Falling back to mock serial port for port config test")
		mock := &MockSerialPort{}
		dev, _ = Open("COM1", mock)
	}
	config := dev.GetPortConfig()
	if config == nil {
		t.Error("Expected port config to be set (real or mock)")
	}
	details := dev.GetPortDetails()
	if details == "" {
		t.Error("Expected port details string")
	}
}

func TestDevice_GetHardwareInfoAndVariant(t *testing.T) {
	dev := &Device{variant: VariantV2, hardwareInfo: getHardwareInfo(VariantV2)}
	if dev.GetHardwareVariant() != VariantV2 {
		t.Error("GetHardwareVariant mismatch")
	}
	info := dev.GetHardwareInfo()
	if info.Variant != VariantV2 {
		t.Error("GetHardwareInfo mismatch")
	}
}

func TestDevice_GetFrequencyRangeAndSweepPoints(t *testing.T) {
	dev := &Device{hardwareInfo: getHardwareInfo(VariantV2)}
	fr := dev.GetFrequencyRange()
	if fr.MinHz <= 0 || fr.MaxHz <= 0 {
		t.Error("Invalid frequency range")
	}
	if dev.GetMaxSweepPoints() <= 0 {
		t.Error("Invalid max sweep points")
	}
}

func TestDevice_GetSupportedPortsAndCapabilities(t *testing.T) {
	dev := &Device{hardwareInfo: getHardwareInfo(VariantV2Plus4)}
	ports := dev.GetSupportedPorts()
	if len(ports) == 0 {
		t.Error("Expected supported ports")
	}
	caps := dev.GetCapabilities()
	if !caps.HasS21 {
		t.Error("Expected HasS21 true")
	}
}

func TestDevice_SetPortHandleAndGetPortHandle(t *testing.T) {
	dev := &Device{}
	mock := &MockSerialPort{}
	dev.SetPortHandle(mock)
	if dev.GetPortHandle() != mock {
		t.Error("SetPortHandle/GetPortHandle failed")
	}
}

func TestDevice_SetSweepConfig_Validation(t *testing.T) {
	dev := &Device{hardwareInfo: getHardwareInfo(VariantV1), variant: VariantV1}
	err := dev.SetSweepConfig(10, 900000000, 101) // too low start
	if err == nil {
		t.Error("Expected error for startHz too low")
	}
	err = dev.SetSweepConfig(50000, 900000001, 101) // too high stop
	if err == nil {
		t.Error("Expected error for stopHz too high")
	}
	err = dev.SetSweepConfig(50000, 900000000, 102) // too many points
	if err == nil {
		t.Error("Expected error for too many points")
	}
}

func TestDevice_GetVersion_Default(t *testing.T) {
	dev := &Device{version: "v2"}
	if dev.GetVersion() != "v2" {
		t.Error("GetVersion did not return expected value")
	}
}

func TestDevice_CalibrationStubs(t *testing.T) {
	dev := &Device{}
	_, err := dev.GetCalibration()
	if err != nil {
		t.Error("GetCalibration should not error (stub)")
	}
	err = dev.SetCalibration(CalibrationData{})
	if err != nil {
		t.Error("SetCalibration should not error (stub)")
	}
	err = dev.SaveCalibration(1)
	if err != nil {
		t.Error("SaveCalibration should not error (stub)")
	}
	err = dev.LoadCalibration(1)
	if err != nil {
		t.Error("LoadCalibration should not error (stub)")
	}
}

func TestListDevices_NoDevices(t *testing.T) {
	// This test will always pass on non-Windows or without devices
	_, err := ListDevices()
	if err == nil {
		t.Log("ListDevices found ports (expected on Windows with devices attached)")
	}
}

func TestOpenWithMockSerialPort(t *testing.T) {
	mock := &MockSerialPort{}
	dev, err := Open("COM1", mock)
	if err != nil {
		t.Fatalf("Open failed: %v", err)
	}
	if dev.GetPortHandle() != mock {
		t.Error("Port handle was not set to mock")
	}
}

func TestDevice_Close(t *testing.T) {
	mock := &MockSerialPort{}
	dev, _ := Open("COM1", mock)
	err := dev.Close()
	if err != nil {
		t.Errorf("Close failed: %v", err)
	}
	if !mock.Closed {
		t.Error("MockSerialPort was not closed")
	}
}

func TestDevice_IsPortSupported(t *testing.T) {
	dev := &Device{
		hardwareInfo: getHardwareInfo(VariantV2Plus4),
	}
	if !dev.IsPortSupported("S21") {
		t.Error("S21 should be supported for V2Plus4")
	}
	if dev.IsPortSupported("FOO") {
		t.Error("FOO should not be supported")
	}
}
