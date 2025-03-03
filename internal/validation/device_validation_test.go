package validation

import (
	"testing"

	"github.com/tyrese-r/go-home/internal/models"
)

func TestIsValidDeviceName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty", "", false},
		{"Valid alphanumeric", "Device123", true},
		{"Valid letters only", "DeviceName", true},
		{"Valid numbers only", "123456", true},
		{"Valid min length", "a", true},
		{"Valid max length", generateString(MaxDeviceNameLength, 'a'), true},
		{"Invalid exceeds max length", generateString(MaxDeviceNameLength+1, 'a'), false},
		{"Invalid with space", "Device Name", false},
		{"Invalid with hyphen", "Device-123", false},
		{"Invalid with underscore", "Device_123", false},
		{"Invalid with special chars", "Device@123", false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidDeviceName(tc.input)
			if result != tc.expected {
				t.Errorf("IsValidDeviceName(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestIsValidOwner(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty", "", false},
		{"Valid", "owner123", true},
		{"Valid min length", "a", true},
		{"Valid max length", generateString(MaxOwnerLength, 'a'), true},
		{"Invalid exceeds max length", generateString(MaxOwnerLength+1, 'a'), false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsValidOwner(tc.input)
			if result != tc.expected {
				t.Errorf("IsValidOwner(%q) = %v; expected %v", tc.input, result, tc.expected)
			}
		})
	}
}

func TestValidateDeviceCreate(t *testing.T) {
	// Setup valid device type
	validDeviceType := models.DeviceTypeCamera
	invalidDeviceType := models.DeviceType("INVALID_TYPE")

	tests := []struct {
		name         string
		deviceCreate models.DeviceCreate
		expectValid  bool
		expectErrors []string
	}{
		{
			name: "Valid device",
			deviceCreate: models.DeviceCreate{
				Name:        "Device123",
				Description: "Test device",
				DeviceType:  validDeviceType,
				OwnedBy:     "owner1",
			},
			expectValid:  true,
			expectErrors: nil,
		},
		{
			name: "Invalid name with special character",
			deviceCreate: models.DeviceCreate{
				Name:        "Device@123",
				Description: "Test device",
				DeviceType:  validDeviceType,
				OwnedBy:     "owner1",
			},
			expectValid:  false,
			expectErrors: []string{"name"},
		},
		{
			name: "Invalid device type",
			deviceCreate: models.DeviceCreate{
				Name:        "Device123",
				Description: "Test device",
				DeviceType:  invalidDeviceType,
				OwnedBy:     "owner1",
			},
			expectValid:  false,
			expectErrors: []string{"device_type"},
		},
		{
			name: "Invalid owner (empty)",
			deviceCreate: models.DeviceCreate{
				Name:        "Device123",
				Description: "Test device",
				DeviceType:  validDeviceType,
				OwnedBy:     "",
			},
			expectValid:  false,
			expectErrors: []string{"owned_by"},
		},
		{
			name: "Invalid description (too long)",
			deviceCreate: models.DeviceCreate{
				Name:        "Device123",
				Description: generateString(MaxDescriptionLength+1, 'a'),
				DeviceType:  validDeviceType,
				OwnedBy:     "owner1",
			},
			expectValid:  false,
			expectErrors: []string{"description"},
		},
		{
			name: "Multiple validation errors",
			deviceCreate: models.DeviceCreate{
				Name:        "Device@123",
				Description: generateString(MaxDescriptionLength+1, 'a'),
				DeviceType:  invalidDeviceType,
				OwnedBy:     "",
			},
			expectValid:  false,
			expectErrors: []string{"name", "device_type", "owned_by", "description"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, errors := ValidateDeviceCreate(&tc.deviceCreate)

			if valid != tc.expectValid {
				t.Errorf("ValidateDeviceCreate() valid = %v, expected %v", valid, tc.expectValid)
			}

			if !tc.expectValid {
				for _, field := range tc.expectErrors {
					if _, exists := errors[field]; !exists {
						t.Errorf("Expected error for field %q but none was found", field)
					}
				}

				if len(errors) != len(tc.expectErrors) {
					t.Errorf("Got %d errors, expected %d", len(errors), len(tc.expectErrors))
				}
			}
		})
	}
}

func TestValidateDeviceUpdate(t *testing.T) {
	// Setup valid device type
	validDeviceType := models.DeviceTypeCamera
	invalidDeviceType := models.DeviceType("INVALID_TYPE")

	// Helper to create pointer values
	strPtr := func(s string) *string { return &s }
	boolPtr := func(b bool) *bool { return &b }
	deviceTypePtr := func(dt models.DeviceType) *models.DeviceType { return &dt }

	tests := []struct {
		name         string
		deviceUpdate models.DeviceUpdate
		expectValid  bool
		expectErrors []string
	}{
		{
			name: "Valid update (all fields)",
			deviceUpdate: models.DeviceUpdate{
				Name:            strPtr("Device123"),
				Description:     strPtr("Updated description"),
				IsOnline:        boolPtr(true),
				OwnedBy:         strPtr("newowner"),
				DeviceType:      deviceTypePtr(validDeviceType),
				LastAlarmReason: strPtr("Test alarm"),
			},
			expectValid:  true,
			expectErrors: nil,
		},
		{
			name: "Valid update (partial)",
			deviceUpdate: models.DeviceUpdate{
				Name:     strPtr("Device123"),
				IsOnline: boolPtr(false),
			},
			expectValid:  true,
			expectErrors: nil,
		},
		{
			name: "Invalid name",
			deviceUpdate: models.DeviceUpdate{
				Name: strPtr("Invalid Name!"),
			},
			expectValid:  false,
			expectErrors: []string{"name"},
		},
		{
			name: "Invalid owner",
			deviceUpdate: models.DeviceUpdate{
				OwnedBy: strPtr(generateString(MaxOwnerLength+1, 'a')),
			},
			expectValid:  false,
			expectErrors: []string{"owned_by"},
		},
		{
			name: "Invalid description",
			deviceUpdate: models.DeviceUpdate{
				Description: strPtr(generateString(MaxDescriptionLength+1, 'a')),
			},
			expectValid:  false,
			expectErrors: []string{"description"},
		},
		{
			name: "Invalid device type",
			deviceUpdate: models.DeviceUpdate{
				DeviceType: deviceTypePtr(invalidDeviceType),
			},
			expectValid:  false,
			expectErrors: []string{"device_type"},
		},
		{
			name: "Invalid last alarm reason",
			deviceUpdate: models.DeviceUpdate{
				LastAlarmReason: strPtr(generateString(MaxLastAlarmReasonLength+1, 'a')),
			},
			expectValid:  false,
			expectErrors: []string{"last_alarm_reason"},
		},
		{
			name: "Multiple validation errors",
			deviceUpdate: models.DeviceUpdate{
				Name:            strPtr("Invalid Name!"),
				Description:     strPtr(generateString(MaxDescriptionLength+1, 'a')),
				OwnedBy:         strPtr(generateString(MaxOwnerLength+1, 'a')),
				DeviceType:      deviceTypePtr(invalidDeviceType),
				LastAlarmReason: strPtr(generateString(MaxLastAlarmReasonLength+1, 'a')),
			},
			expectValid:  false,
			expectErrors: []string{"name", "description", "owned_by", "device_type", "last_alarm_reason"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, errors := ValidateDeviceUpdate(&tc.deviceUpdate)

			if valid != tc.expectValid {
				t.Errorf("ValidateDeviceUpdate() valid = %v, expected %v", valid, tc.expectValid)
			}

			if !tc.expectValid {
				for _, field := range tc.expectErrors {
					if _, exists := errors[field]; !exists {
						t.Errorf("Expected error for field %q but none was found", field)
					}
				}

				if len(errors) != len(tc.expectErrors) {
					t.Errorf("Got %d errors, expected %d", len(errors), len(tc.expectErrors))
				}
			}
		})
	}
}

// Helper function to generate strings of specified length
func generateString(length int, char rune) string {
	runes := make([]rune, length)
	for i := range runes {
		runes[i] = char
	}
	return string(runes)
}

func TestValidateAlarmRequest(t *testing.T) {
	tests := []struct {
		name         string
		alarmRequest models.AlarmRequest
		expectValid  bool
		expectErrors []string
	}{
		{
			name: "Valid alarm request",
			alarmRequest: models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "WARNING",
			},
			expectValid:  true,
			expectErrors: nil,
		},
		{
			name: "Empty reason",
			alarmRequest: models.AlarmRequest{
				Reason: "",
				Level:  "WARNING",
			},
			expectValid:  false,
			expectErrors: []string{"reason"},
		},
		{
			name: "Reason too long",
			alarmRequest: models.AlarmRequest{
				Reason: generateString(MaxLastAlarmReasonLength+1, 'a'),
				Level:  "CRITICAL",
			},
			expectValid:  false,
			expectErrors: []string{"reason"},
		},
		{
			name: "Invalid level",
			alarmRequest: models.AlarmRequest{
				Reason: "Motion detected",
				Level:  "LOW", // Not a valid level
			},
			expectValid:  false,
			expectErrors: []string{"level"},
		},
		{
			name: "Multiple validation errors",
			alarmRequest: models.AlarmRequest{
				Reason: "",
				Level:  "INVALID_LEVEL",
			},
			expectValid:  false,
			expectErrors: []string{"reason", "level"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			valid, errors := ValidateAlarmRequest(&tc.alarmRequest)

			if valid != tc.expectValid {
				t.Errorf("ValidateAlarmRequest() valid = %v, expected %v", valid, tc.expectValid)
			}

			if !tc.expectValid {
				for _, field := range tc.expectErrors {
					if _, exists := errors[field]; !exists {
						t.Errorf("Expected error for field %q but none was found", field)
					}
				}

				if len(errors) != len(tc.expectErrors) {
					t.Errorf("Got %d errors, expected %d", len(errors), len(tc.expectErrors))
				}
			}
		})
	}
}
