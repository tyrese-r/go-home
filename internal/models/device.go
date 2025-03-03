package models

import "time"

// Device database model
type Device struct {
	ID              int64      `json:"id"`
	OwnedBy         string     `json:"owned_by"`
	DeviceType      DeviceType `json:"device_type"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	IsOnline        bool       `json:"is_online"`
	LastAlarmTime   time.Time  `json:"last_alarm_time"`
	LastAlarmReason string     `json:"last_alarm_reason"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// API models

type DeviceCreate struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	DeviceType  DeviceType `json:"device_type" binding:"required"`
	OwnedBy     string     `json:"owned_by" binding:"required"`
}

type DeviceUpdate struct {
	Name            *string     `json:"name"`
	Description     *string     `json:"description"`
	IsOnline        *bool       `json:"is_online"`
	OwnedBy         *string     `json:"owned_by"`
	DeviceType      *DeviceType `json:"device_type"`
	LastAlarmReason *string     `json:"last_alarm_reason"`
}

// AlarmRequest represents a request to trigger a device alarm
type AlarmRequest struct {
	Reason string `json:"reason" binding:"required"`
	Level  string `json:"level" binding:"required"`
}
