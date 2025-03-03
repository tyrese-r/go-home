package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/tyrese-r/go-home/internal/models"
)

// MockDeviceRepo is a mock implementation of repository.DeviceRepository
type MockDeviceRepo struct {
	getByIDCalled      bool
	getByIDInput       int64
	getByIDOutput      *models.Device
	getByIDError       error
	triggerAlarmCalled bool
	triggerAlarmID     int64
	triggerAlarmReason string
	triggerAlarmError  error
}

// Implement the DeviceRepository interface methods
func (m *MockDeviceRepo) GetByID(id int64) (*models.Device, error) {
	m.getByIDCalled = true
	m.getByIDInput = id
	return m.getByIDOutput, m.getByIDError
}

func (m *MockDeviceRepo) TriggerAlarm(id int64, reason string) error {
	m.triggerAlarmCalled = true
	m.triggerAlarmID = id
	m.triggerAlarmReason = reason
	return m.triggerAlarmError
}

// Stub implementations of other repository methods
func (m *MockDeviceRepo) Create(*models.DeviceCreate) (int64, error) { return 0, nil }
func (m *MockDeviceRepo) GetAll() ([]*models.Device, error)          { return nil, nil }
func (m *MockDeviceRepo) Update(int64, *models.DeviceUpdate) error   { return nil }
func (m *MockDeviceRepo) Delete(int64) error                         { return nil }

func TestTriggerAlarm(t *testing.T) {
	tests := []struct {
		name                     string
		deviceID                 int64
		alarm                    *models.AlarmRequest
		mockGetByIDOutput        *models.Device
		mockGetByIDError         error
		mockTriggerAlarmError    error
		expectError              bool
		expectTriggerAlarmCalled bool
	}{
		{
			name:     "Successful alarm trigger",
			deviceID: 1,
			alarm: &models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "CRITICAL",
			},
			mockGetByIDOutput:        &models.Device{ID: 1, Name: "Test Device"},
			expectError:              false,
			expectTriggerAlarmCalled: true,
		},
		{
			name:     "Device not found",
			deviceID: 99,
			alarm: &models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "CRITICAL",
			},
			mockGetByIDOutput:        nil, // No device found
			expectError:              true,
			expectTriggerAlarmCalled: false,
		},
		{
			name:     "GetByID database error",
			deviceID: 1,
			alarm: &models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "CRITICAL",
			},
			mockGetByIDError:         errors.New("database error"),
			expectError:              true,
			expectTriggerAlarmCalled: false,
		},
		{
			name:     "TriggerAlarm database error",
			deviceID: 1,
			alarm: &models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "CRITICAL",
			},
			mockGetByIDOutput:        &models.Device{ID: 1, Name: "Test Device"},
			mockTriggerAlarmError:    errors.New("database error"),
			expectError:              true,
			expectTriggerAlarmCalled: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create the mock repository
			mockRepo := &MockDeviceRepo{
				getByIDOutput:     tc.mockGetByIDOutput,
				getByIDError:      tc.mockGetByIDError,
				triggerAlarmError: tc.mockTriggerAlarmError,
			}

			// Create service with the mock repository
			service := NewDeviceService(mockRepo)

			// Call the method being tested
			err := service.TriggerAlarm(tc.deviceID, tc.alarm)

			// Check error expectations
			if tc.expectError && err == nil {
				t.Errorf("Expected an error but got nil")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			// Check if GetByID was called with correct ID
			if !mockRepo.getByIDCalled {
				t.Errorf("Expected GetByID to be called")
			}
			if mockRepo.getByIDInput != tc.deviceID {
				t.Errorf("GetByID called with wrong ID, expected %d, got %d", tc.deviceID, mockRepo.getByIDInput)
			}

			// Check if TriggerAlarm was called when expected
			if tc.expectTriggerAlarmCalled && !mockRepo.triggerAlarmCalled {
				t.Errorf("Expected TriggerAlarm to be called but it wasn't")
			}
			if !tc.expectTriggerAlarmCalled && mockRepo.triggerAlarmCalled {
				t.Errorf("Expected TriggerAlarm not to be called but it was")
			}

			// If TriggerAlarm was called, check the parameters
			if mockRepo.triggerAlarmCalled {
				if mockRepo.triggerAlarmID != tc.deviceID {
					t.Errorf("TriggerAlarm called with wrong ID, expected %d, got %d", tc.deviceID, mockRepo.triggerAlarmID)
				}

				expectedReason := fmt.Sprintf("[%s] %s", tc.alarm.Level, tc.alarm.Reason)
				if mockRepo.triggerAlarmReason != expectedReason {
					t.Errorf("TriggerAlarm called with wrong reason, expected %q, got %q", expectedReason, mockRepo.triggerAlarmReason)
				}
			}
		})
	}
}
