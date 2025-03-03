package repository

import "github.com/tyrese-r/go-home/internal/models"

// DeviceRepository defines the interface for device data operations
type DeviceRepository interface {
	Create(device *models.DeviceCreate) (int64, error)
	GetByID(id int64) (*models.Device, error)
	GetAll() ([]*models.Device, error)
	Update(id int64, device *models.DeviceUpdate) error
	Delete(id int64) error
	TriggerAlarm(id int64, reason string) error
}
