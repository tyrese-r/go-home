package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/tyrese-r/go-home/internal/models"
)

// DeviceRepositoryImpl handles database operations for devices
type DeviceRepositoryImpl struct {
	db *sql.DB
}

// NewDeviceRepository creates a new DeviceRepository
func NewDeviceRepository(db *sql.DB) DeviceRepository {
	return &DeviceRepositoryImpl{db: db}
}

// Create adds a new device to the database
// Parameterised
func (r *DeviceRepositoryImpl) Create(device *models.DeviceCreate) (int64, error) {
	query := `INSERT INTO devices (name, description, device_type, owned_by) VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, device.Name, device.Description, device.DeviceType, device.OwnedBy)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetByID retrieves a device by its ID
func (r *DeviceRepositoryImpl) GetByID(id int64) (*models.Device, error) {
	query := `SELECT id, name, description, device_type, owned_by, is_online, last_alarm_reason, last_alarm_time, created_at, updated_at FROM devices WHERE id = ?`

	var device models.Device
	var lastAlarmTime, createdAt, updatedAt string

	err := r.db.QueryRow(query, id).Scan(
		&device.ID,
		&device.Name,
		&device.Description,
		&device.DeviceType,
		&device.OwnedBy,
		&device.IsOnline,
		&device.LastAlarmReason,
		&lastAlarmTime,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, err
	}

	// Parse time strings
	device.LastAlarmTime, _ = time.Parse(time.RFC3339, lastAlarmTime)
	device.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	device.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

	return &device, nil
}

// GetAll retrieves all devices
func (r *DeviceRepositoryImpl) GetAll() ([]*models.Device, error) {
	query := `SELECT id, name, description, device_type, owned_by, is_online, last_alarm_reason, last_alarm_time, created_at, updated_at FROM devices ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}()
	var devices []*models.Device

	for rows.Next() {
		var device models.Device
		var lastAlarmTime, createdAt, updatedAt string

		if err := rows.Scan(
			&device.ID,
			&device.Name,
			&device.Description,
			&device.DeviceType,
			&device.OwnedBy,
			&device.IsOnline,
			&device.LastAlarmReason,
			&lastAlarmTime,
			&createdAt,
			&updatedAt,
		); err != nil {
			return nil, err
		}

		// Parse time strings
		device.LastAlarmTime, _ = time.Parse(time.RFC3339, lastAlarmTime)
		device.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		device.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)

		devices = append(devices, &device)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return devices, nil
}

// Update updates a device in the database
func (r *DeviceRepositoryImpl) Update(id int64, device *models.DeviceUpdate) error {
	// First, get the current device data
	currentDevice, err := r.GetByID(id)
	if err != nil {
		return err
	}
	if currentDevice == nil {
		return sql.ErrNoRows
	}

	// Apply updates to fields that are present
	name := currentDevice.Name
	description := currentDevice.Description
	isOnline := currentDevice.IsOnline
	ownedBy := currentDevice.OwnedBy
	deviceType := currentDevice.DeviceType
	lastAlarmReason := currentDevice.LastAlarmReason

	if device.Name != nil {
		name = *device.Name
	}
	if device.Description != nil {
		description = *device.Description
	}
	if device.IsOnline != nil {
		isOnline = *device.IsOnline
	}
	if device.OwnedBy != nil {
		ownedBy = *device.OwnedBy
	}
	if device.DeviceType != nil {
		deviceType = *device.DeviceType
	}
	if device.LastAlarmReason != nil {
		lastAlarmReason = *device.LastAlarmReason
	}

	query := `UPDATE devices SET name = ?, description = ?, device_type = ?, is_online = ?, owned_by = ?, last_alarm_reason = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err = r.db.Exec(query, name, description, deviceType, isOnline, ownedBy, lastAlarmReason, id)
	return err
}

// Delete removes a device from the database
func (r *DeviceRepositoryImpl) Delete(id int64) error {
	query := `DELETE FROM devices WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// TriggerAlarm updates a device's alarm information
func (r *DeviceRepositoryImpl) TriggerAlarm(id int64, reason string) error {
	query := `UPDATE devices SET last_alarm_reason = ?, last_alarm_time = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := r.db.Exec(query, reason, id)
	return err
}
