package models

// DeviceType represents the type of a smart home device
type DeviceType string

// Device type enum values
const (
	DeviceTypeCamera        DeviceType = "CAMERA"
	DeviceTypeThermostat    DeviceType = "THERMOSTAT"
	DeviceTypeSmokeDetector DeviceType = "SMOKE_DETECTOR"
	DeviceTypeMotionSensor  DeviceType = "MOTION_SENSOR"
	DeviceTypeLock          DeviceType = "LOCK"
	DeviceTypeUnknown       DeviceType = "UNKNOWN"
)

// IsValid checks if the device type is a valid enum value
func (dt DeviceType) IsValid() bool {
	switch dt {
	case DeviceTypeCamera, DeviceTypeThermostat, DeviceTypeSmokeDetector,
		DeviceTypeMotionSensor, DeviceTypeLock, DeviceTypeUnknown:
		return true
	}
	return false
}

// String returns the string representation of the device type
func (dt DeviceType) String() string {
	return string(dt)
}

// DeviceTypeInfo represents information about a device type
type DeviceTypeInfo struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
}

// IsValidDeviceType checks if a given device type is valid
func IsValidDeviceType(dt DeviceType) bool {
	return dt.IsValid()
}

// GetAllDeviceTypes returns a list of all valid device types with their information
func GetAllDeviceTypes() []DeviceTypeInfo {
	return []DeviceTypeInfo{
		{ID: string(DeviceTypeCamera), DisplayName: "Camera", Description: "Smart security camera"},
		{ID: string(DeviceTypeThermostat), DisplayName: "Thermostat", Description: "Smart thermostat for controlling temperature"},
		{ID: string(DeviceTypeSmokeDetector), DisplayName: "Smoke Detector", Description: "Detects smoke and fire hazards"},
		{ID: string(DeviceTypeMotionSensor), DisplayName: "Motion Sensor", Description: "Detects movement in monitored areas"},
		{ID: string(DeviceTypeLock), DisplayName: "Lock", Description: "Smart lock with remote access capabilities"},
		{ID: string(DeviceTypeUnknown), DisplayName: "Unknown", Description: "Unknown device type"},
	}
}
