# Go NanoVNA Communication Package Specification

*This documentation was created with the assistance of AI.*

## Overview

This package provides a Go API for communicating with and controlling a NanoVNA device over a serial (USB) connection. It enables device discovery, configuration, data acquisition, calibration management, and device information retrieval.

---

## Features

- Serial port device discovery and connection
- Command sending and response parsing (ASCII and binary)
- Sweep configuration (frequency range, number of points)
- Data acquisition (S-parameters, traces)
- Calibration management (load, save, apply)
- Device information retrieval (model, firmware, serial number)
- Error handling (timeouts, protocol errors, device not found)
- Extensible for additional models/firmware

---

## Package Workflow (Flowchart)

```mermaid
flowchart TD
  A[ListDevices()] --> B(Open())
  B --> C{Device Connected?}
  C -- Yes --> D[GetInfo()]
  D --> E[SetSweepConfig()]
  E --> F[RunSweep()]
  F --> G[GetCalibration() / SetCalibration()]
  G --> H[SaveCalibration() / LoadCalibration()]
  H --> I[Close()]
  C -- No --> Z[Error Handling]
```

## Public API

### Device Management

- `ListDevices() ([]string, error)`
  - Lists available serial ports with connected NanoVNA devices.
- `Open(port string) (*Device, error)`
  - Opens a connection to the NanoVNA on the specified port.
- `Close() error`
  - Closes the device connection.

### Device Information

- `GetInfo() (DeviceInfo, error)`
  - Retrieves device model, firmware version, and serial number.

### Sweep Configuration

- `SetSweepConfig(startHz, stopHz int, points int) error`
  - Configures sweep parameters.

### Data Acquisition

- `RunSweep() (SweepData, error)`
  - Triggers a sweep and returns measurement data.

### Calibration

- `GetCalibration() (CalibrationData, error)`
  - Retrieves current calibration data.
- `SetCalibration(cal CalibrationData) error`
  - Applies calibration data.
- `SaveCalibration(slot int) error`
  - Saves calibration data to device memory.
- `LoadCalibration(slot int) error`
  - Loads calibration data from device memory.

---

## Data Structures

```go

type Device struct {
  Port         string
  // (Unexported fields: portHandle, config, version, variant, hardwareInfo)
}

type DeviceInfo struct {
    Model     string
    Firmware  string
    SerialNum string
}

  Frequencies []float64
  S11         []complex128
  S21         []complex128
  // Additional traces as needed
}

type CalibrationData struct {
    // Calibration coefficients and metadata
}
```

---

## Communication Protocol

- Serial port (default 9600 baud for most models, configurable)
- ASCII and binary command support as per NanoVNA firmware
- Command/response with timeouts and retries
- Parsing of measurement and info responses

---

## Error Handling

- Custom error types:
  - `ErrDeviceNotFound`
  - `ErrCommunication`
  - `ErrTimeout`
  - `ErrProtocol`
- All public API functions return errors as appropriate

---

## Extensibility

- Support for additional NanoVNA models and firmware versions
- Optional: async/streaming data acquisition
- Optional: advanced calibration and trace management

---

## Example Usage

```go
devices, err := nanovna.ListDevices()
dev, err := nanovna.Open(devices[0])
info, err := dev.GetInfo()
err = dev.SetSweepConfig(1000000, 30000000, 101)
data, err := dev.RunSweep()
err = dev.Close()
```
