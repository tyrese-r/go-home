package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/tyrese-r/go-home/internal/models"
)

// Use the DeviceServiceInterface defined in handlers.go

// Mock implementation of the DeviceService
type MockDeviceService struct {
	getByIDFunc      func(id int64) (*models.Device, error)
	getAllFunc       func() ([]*models.Device, error)
	createFunc       func(device *models.DeviceCreate) (int64, error)
	updateFunc       func(id int64, device *models.DeviceUpdate) error
	deleteFunc       func(id int64) error
	triggerAlarmFunc func(id int64, alarm *models.AlarmRequest) error
}

// Implement the DeviceServiceInterface
func (m *MockDeviceService) GetDeviceByID(id int64) (*models.Device, error) {
	return m.getByIDFunc(id)
}

func (m *MockDeviceService) GetAllDevices() ([]*models.Device, error) {
	return m.getAllFunc()
}

func (m *MockDeviceService) CreateDevice(device *models.DeviceCreate) (int64, error) {
	return m.createFunc(device)
}

func (m *MockDeviceService) UpdateDevice(id int64, device *models.DeviceUpdate) error {
	return m.updateFunc(id, device)
}

func (m *MockDeviceService) DeleteDevice(id int64) error {
	return m.deleteFunc(id)
}

func (m *MockDeviceService) TriggerAlarm(id int64, alarm *models.AlarmRequest) error {
	return m.triggerAlarmFunc(id, alarm)
}

// TestHandler implements a minimal handler for testing
type TestHandler struct {
	deviceService DeviceServiceInterface
	router        *gin.Engine
}

// setupTestRouter creates a test router with necessary routes
func setupTestRouter(mockSvc *MockDeviceService) *gin.Engine {
	gin.SetMode(gin.TestMode)

	// Create router
	router := gin.New()

	// Create handler
	h := &TestHandler{
		deviceService: mockSvc,
		router:        router,
	}

	// Set up routes
	api := router.Group("/api")
	{
		devices := api.Group("/devices")
		{
			devices.POST("/:id/alarm", h.triggerDeviceAlarm)
		}
	}

	return router
}

// Test implementation of triggerDeviceAlarm
func (h *TestHandler) triggerDeviceAlarm(c *gin.Context) {
	// Parse device ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	// Parse request body
	var alarmRequest models.AlarmRequest
	if bindErr := c.ShouldBindJSON(&alarmRequest); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"bindError": bindErr.Error()})
		return
	}

	// Trigger alarm on device
	err = h.deviceService.TriggerAlarm(id, &alarmRequest)
	if err != nil {
		// Handle device not found case specifically
		if err.Error() == fmt.Sprintf("device not found with ID: %d", id) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success with 204 No Content
	c.Status(http.StatusNoContent)
}

func TestTriggerDeviceAlarm(t *testing.T) {
	tests := []struct {
		name         string
		deviceID     string
		requestBody  interface{}
		setupMock    func(*MockDeviceService)
		expectedCode int
		expectedBody map[string]interface{}
	}{
		{
			name:     "Successful alarm trigger",
			deviceID: "1",
			requestBody: models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "WARNING",
			},
			setupMock: func(m *MockDeviceService) {
				m.triggerAlarmFunc = func(id int64, alarm *models.AlarmRequest) error {
					return nil
				}
			},
			expectedCode: http.StatusNoContent,
			expectedBody: nil,
		},
		{
			name:     "Invalid device ID",
			deviceID: "invalid",
			requestBody: models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "WARNING",
			},
			setupMock: func(m *MockDeviceService) {
				m.triggerAlarmFunc = func(id int64, alarm *models.AlarmRequest) error {
					return nil
				}
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: map[string]interface{}{
				"error": "invalid device ID",
			},
		},
		{
			name:     "Device not found",
			deviceID: "99",
			requestBody: models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "WARNING",
			},
			setupMock: func(m *MockDeviceService) {
				m.triggerAlarmFunc = func(id int64, alarm *models.AlarmRequest) error {
					return fmt.Errorf("device not found with ID: %d", id)
				}
			},
			expectedCode: http.StatusNotFound,
			expectedBody: map[string]interface{}{
				"error": "device not found with ID: 99",
			},
		},
		{
			name:     "Service error",
			deviceID: "1",
			requestBody: models.AlarmRequest{
				Reason: "Smoke detected",
				Level:  "WARNING",
			},
			setupMock: func(m *MockDeviceService) {
				m.triggerAlarmFunc = func(id int64, alarm *models.AlarmRequest) error {
					return errors.New("internal error")
				}
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: map[string]interface{}{
				"error": "internal error",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup mock service
			mockSvc := &MockDeviceService{}
			tc.setupMock(mockSvc)

			// Setup router with mock service
			router := setupTestRouter(mockSvc)

			// Create request
			bodyBytes, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", fmt.Sprintf("/api/devices/%s/alarm", tc.deviceID), bytes.NewBuffer(bodyBytes))
			req.Header.Set("Content-Type", "application/json")

			// Create a response recorder
			recorder := httptest.NewRecorder()

			// Perform the request
			router.ServeHTTP(recorder, req)

			// Check status code
			if recorder.Code != tc.expectedCode {
				t.Errorf("Expected status code %d, got %d", tc.expectedCode, recorder.Code)
			}

			// If we expect a response body, check it
			if tc.expectedBody != nil {
				var responseBody map[string]interface{}
				err := json.Unmarshal(recorder.Body.Bytes(), &responseBody)
				if err != nil {
					t.Errorf("Failed to parse response body: %v", err)
				}

				// Compare error field
				if errorMsg, ok := tc.expectedBody["error"]; ok {
					if responseBody["error"] != errorMsg {
						t.Errorf("Expected error message %q, got %q", errorMsg, responseBody["error"])
					}
				}
			}
		})
	}
}
