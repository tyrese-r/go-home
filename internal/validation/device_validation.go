package validation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/tyrese-r/go-home/internal/models"
)

// Validation constants
const (
	MaxDeviceNameLength      = 100
	MinDeviceNameLength      = 1
	MaxOwnerLength           = 50
	MinOwnerLength           = 1
	MaxDescriptionLength     = 500
	MaxLastAlarmReasonLength = 200
	MinAlarmReasonLength     = 1
)

// Regex patterns
var (
	// Matches alphanumeric characters only (A-Z, a-z, 0-9)
	alphanumericPattern = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
)

// ValidationErrors holds validation error messages for each field
type ValidationErrors map[string]string

// IsValidDeviceName checks if the device name meets criteria
func IsValidDeviceName(name string) bool {
	// Check length
	if len(name) < MinDeviceNameLength || len(name) > MaxDeviceNameLength {
		return false
	}

	return alphanumericPattern.MatchString(name)
}

// IsValidOwner checks if the owner field is valid
func IsValidOwner(owner string) bool {
	return len(owner) >= MinOwnerLength && len(owner) <= MaxOwnerLength
}

// ValidateDeviceCreate performs all validations on device creation data
func ValidateDeviceCreate(device *models.DeviceCreate) (bool, ValidationErrors) {
	errors := make(ValidationErrors)

	if !IsValidDeviceName(device.Name) {
		errors["name"] = fmt.Sprintf("must be between %d-%d characters and contain only alphanumeric characters (A-Z, a-z, 0-9)",
			MinDeviceNameLength, MaxDeviceNameLength)
	}

	if !models.IsValidDeviceType(device.DeviceType) {
		// Get all valid types for the error message
		allTypes := models.GetAllDeviceTypes()
		typeNames := make([]string, 0, len(allTypes))
		for _, t := range allTypes {
			typeNames = append(typeNames, t.ID)
		}

		errors["device_type"] = fmt.Sprintf("must be one of: %s", strings.Join(typeNames, ", "))
	}

	if !IsValidOwner(device.OwnedBy) {
		errors["owned_by"] = fmt.Sprintf("must be between %d-%d characters",
			MinOwnerLength, MaxOwnerLength)
	}

	if len(device.Description) > MaxDescriptionLength {
		errors["description"] = fmt.Sprintf("must not exceed %d characters", MaxDescriptionLength)
	}

	return len(errors) == 0, errors
}

// ValidateAlarmRequest performs all validations on device alarm trigger request
func ValidateAlarmRequest(alarm *models.AlarmRequest) (bool, ValidationErrors) {
	errors := make(ValidationErrors)

	// Validate reason
	if len(alarm.Reason) < MinAlarmReasonLength {
		errors["reason"] = "reason cannot be empty"
	} else if len(alarm.Reason) > MaxLastAlarmReasonLength {
		errors["reason"] = fmt.Sprintf("reason must not exceed %d characters", MaxLastAlarmReasonLength)
	}

	// Validate level
	if alarm.Level != "INFO" && alarm.Level != "WARNING" && alarm.Level != "CRITICAL" {
		errors["level"] = "level must be one of: INFO, WARNING, CRITICAL"
	}

	return len(errors) == 0, errors
}

// ValidateDeviceUpdate performs all validations on device update data
func ValidateDeviceUpdate(device *models.DeviceUpdate) (bool, ValidationErrors) {
	errors := make(ValidationErrors)

	if device.Name != nil && !IsValidDeviceName(*device.Name) {
		errors["name"] = fmt.Sprintf("must be between %d-%d characters and contain only alphanumeric characters (A-Z, a-z, 0-9)",
			MinDeviceNameLength, MaxDeviceNameLength)
	}

	if device.OwnedBy != nil && !IsValidOwner(*device.OwnedBy) {
		errors["owned_by"] = fmt.Sprintf("must be between %d-%d characters",
			MinOwnerLength, MaxOwnerLength)
	}

	if device.Description != nil && len(*device.Description) > MaxDescriptionLength {
		errors["description"] = fmt.Sprintf("must not exceed %d characters", MaxDescriptionLength)
	}

	if device.LastAlarmReason != nil && len(*device.LastAlarmReason) > MaxLastAlarmReasonLength {
		errors["last_alarm_reason"] = fmt.Sprintf("must not exceed %d characters", MaxLastAlarmReasonLength)
	}

	if device.DeviceType != nil && !models.IsValidDeviceType(*device.DeviceType) {
		// Get all valid types for the error message
		allTypes := models.GetAllDeviceTypes()
		typeNames := make([]string, 0, len(allTypes))
		for _, t := range allTypes {
			typeNames = append(typeNames, t.ID)
		}

		errors["device_type"] = fmt.Sprintf("must be one of: %s", strings.Join(typeNames, ", "))
	}

	return len(errors) == 0, errors
}
