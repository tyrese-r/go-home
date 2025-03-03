package service

import (
	"fmt"

	"github.com/tyrese-r/go-home/internal/models"
	"github.com/tyrese-r/go-home/internal/repository"
)

// DeviceService handles business logic for devices
type DeviceService struct {
	repo repository.DeviceRepository
}

// NewDeviceService creates a new DeviceService
func NewDeviceService(repo repository.DeviceRepository) *DeviceService {
	return &DeviceService{repo: repo}
}

// CreateDevice creates a new device
func (s *DeviceService) CreateDevice(device *models.DeviceCreate) (int64, error) {
	return s.repo.Create(device)
}

// GetDeviceByID retrieves a device by its ID
func (s *DeviceService) GetDeviceByID(id int64) (*models.Device, error) {
	return s.repo.GetByID(id)
}

// GetAllDevices retrieves all devices
func (s *DeviceService) GetAllDevices() ([]*models.Device, error) {
	return s.repo.GetAll()
}

// UpdateDevice updates a device
func (s *DeviceService) UpdateDevice(id int64, device *models.DeviceUpdate) error {
	return s.repo.Update(id, device)
}

// DeleteDevice deletes a device
func (s *DeviceService) DeleteDevice(id int64) error {
	return s.repo.Delete(id)
}

// TriggerAlarm triggers an alarm on a device
func (s *DeviceService) TriggerAlarm(id int64, alarm *models.AlarmRequest) error {
	// First check if device exists
	device, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if device == nil {
		return fmt.Errorf("device not found with ID: %d", id)
	}

	// Format reason with alarm level and timestamp
	formattedReason := fmt.Sprintf("[%s] %s", alarm.Level, alarm.Reason)

	// Trigger the alarm
	return s.repo.TriggerAlarm(id, formattedReason)
}
