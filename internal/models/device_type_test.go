package models

import (
	"testing"
)

func TestDeviceType_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeviceType
		expected bool
	}{
		{"Valid - Camera", DeviceTypeCamera, true},
		{"Valid - Thermostat", DeviceTypeThermostat, true},
		{"Valid - SmokeDetector", DeviceTypeSmokeDetector, true},
		{"Valid - MotionSensor", DeviceTypeMotionSensor, true},
		{"Valid - Lock", DeviceTypeLock, true},
		{"Valid - Unknown", DeviceTypeUnknown, true},
		{"Invalid - Empty", DeviceType(""), false},
		{"Invalid - Random string", DeviceType("INVALID_TYPE"), false},
		{"Invalid - Lowercase", DeviceType("camera"), false},
		{"Invalid - Mixed case", DeviceType("Camera"), false},
		{"Invalid - OldType", DeviceType("LIGHT_BULB"), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.dt.IsValid()
			if result != tc.expected {
				t.Errorf("DeviceType(%q).IsValid() = %v; expected %v", tc.dt, result, tc.expected)
			}
		})
	}
}

func TestDeviceType_String(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeviceType
		expected string
	}{
		{"Camera", DeviceTypeCamera, "CAMERA"},
		{"Thermostat", DeviceTypeThermostat, "THERMOSTAT"},
		{"SmokeDetector", DeviceTypeSmokeDetector, "SMOKE_DETECTOR"},
		{"MotionSensor", DeviceTypeMotionSensor, "MOTION_SENSOR"},
		{"Lock", DeviceTypeLock, "LOCK"},
		{"Unknown", DeviceTypeUnknown, "UNKNOWN"},
		{"Empty string", DeviceType(""), ""},
		{"Custom string", DeviceType("CUSTOM"), "CUSTOM"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.dt.String()
			if result != tc.expected {
				t.Errorf("DeviceType(%q).String() = %q; expected %q", tc.dt, result, tc.expected)
			}
		})
	}
}

func TestIsValidDeviceType(t *testing.T) {
	tests := []struct {
		name     string
		dt       DeviceType
		expected bool
	}{
		{"Valid - Camera", DeviceTypeCamera, true},
		{"Valid - Thermostat", DeviceTypeThermostat, true},
		{"Valid - SmokeDetector", DeviceTypeSmokeDetector, true},
		{"Valid - MotionSensor", DeviceTypeMotionSensor, true},
		{"Valid - Lock", DeviceTypeLock, true},
		{"Valid - Unknown", DeviceTypeUnknown, true},
		{"Invalid - Empty", DeviceType(""), false},
		{"Invalid - Random string", DeviceType("INVALID_TYPE"), false},
		{"Invalid - Lowercase", DeviceType("camera"), false},
		{"Invalid - Mixed case", DeviceType("Camera"), false},
		{"Invalid - OldType", DeviceType("LIGHT_BULB"), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidDeviceType(tc.dt)
			if result != tc.expected {
				t.Errorf("IsValidDeviceType(%q) = %v; expected %v", tc.dt, result, tc.expected)
			}
		})
	}
}

func TestGetAllDeviceTypes(t *testing.T) {
	deviceTypes := GetAllDeviceTypes()

	// Check that we have the expected number of device types
	expectedCount := 6 // matches the count in the implementation
	if len(deviceTypes) != expectedCount {
		t.Errorf("GetAllDeviceTypes() returned %d device types; expected %d", len(deviceTypes), expectedCount)
	}

	// Test that each device type has the expected properties
	expectedDeviceTypes := map[string]struct{}{
		string(DeviceTypeCamera):        {},
		string(DeviceTypeThermostat):    {},
		string(DeviceTypeSmokeDetector): {},
		string(DeviceTypeMotionSensor):  {},
		string(DeviceTypeLock):          {},
		string(DeviceTypeUnknown):       {},
	}

	// Verify each returned device type is in our expected map
	for _, dt := range deviceTypes {
		if _, exists := expectedDeviceTypes[dt.ID]; !exists {
			t.Errorf("Unexpected device type ID: %s", dt.ID)
		}

		// Check that the device type has non-empty display name and description
		if dt.DisplayName == "" {
			t.Errorf("Device type %s has empty display name", dt.ID)
		}

		if dt.Description == "" {
			t.Errorf("Device type %s has empty description", dt.ID)
		}
	}

	// Test for specific device types and their expected values
	expectedValues := map[string]struct {
		displayName string
		description string
	}{
		string(DeviceTypeCamera):        {"Camera", "Smart security camera"},
		string(DeviceTypeLock):          {"Lock", "Smart lock with remote access capabilities"},
		string(DeviceTypeSmokeDetector): {"Smoke Detector", "Detects smoke and fire hazards"},
		string(DeviceTypeUnknown):       {"Unknown", "Unknown device type"},
	}

	for _, dt := range deviceTypes {
		if expected, exists := expectedValues[dt.ID]; exists {
			if dt.DisplayName != expected.displayName {
				t.Errorf("Device type %s has display name %q; expected %q",
					dt.ID, dt.DisplayName, expected.displayName)
			}

			if dt.Description != expected.description {
				t.Errorf("Device type %s has description %q; expected %q",
					dt.ID, dt.Description, expected.description)
			}
		}
	}
}
