package handlers

import (
	"fmt"
	"github.com/tyrese-r/go-home/internal/validation"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tyrese-r/go-home/internal/models"
)

// DeviceServiceInterface defines the interface for the device service
type DeviceServiceInterface interface {
	CreateDevice(device *models.DeviceCreate) (int64, error)
	GetDeviceByID(id int64) (*models.Device, error)
	GetAllDevices() ([]*models.Device, error)
	UpdateDevice(id int64, device *models.DeviceUpdate) error
	DeleteDevice(id int64) error
	TriggerAlarm(id int64, alarm *models.AlarmRequest) error
}

// Handler handles HTTP requests
type Handler struct {
	deviceService DeviceServiceInterface
	router        *gin.Engine
	startTime     time.Time
}

// New creates a new Handler
func New(deviceService DeviceServiceInterface) *Handler {
	h := &Handler{
		deviceService: deviceService,
		router:        gin.Default(),
		startTime:     time.Now(),
	}

	// Set up routes
	h.setupRoutes()

	return h
}

// setupRoutes configures the HTTP routes
func (h *Handler) setupRoutes() {
	// Health check endpoint
	h.router.GET("/health", h.healthCheck)

	api := h.router.Group("/api")
	{
		devices := api.Group("/devices")
		{
			devices.GET("", h.getAllDevices)
			devices.GET("/:id", h.getDeviceByID)
			devices.POST("", h.createDevice)
			devices.PUT("/:id", h.updateDevice)
			devices.DELETE("/:id", h.deleteDevice)
			devices.POST("/:id/alarm", h.triggerDeviceAlarm)
		}
	}
}

// StartServer starts the HTTP server
func (h *Handler) StartServer(addr string) error {
	return h.router.Run(addr)
}

// getAllDevices handles GET /api/devices
func (h *Handler) getAllDevices(c *gin.Context) {
	devices, err := h.deviceService.GetAllDevices()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, devices)
}

// getDeviceByID handles GET /api/devices/:id
func (h *Handler) getDeviceByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	device, err := h.deviceService.GetDeviceByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if device == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "device not found"})
		return
	}

	c.JSON(http.StatusOK, device)
}

// createDevice handles POST /api/devices
func (h *Handler) createDevice(c *gin.Context) {
	var deviceCreate models.DeviceCreate
	if err := c.ShouldBindJSON(&deviceCreate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	validationSuccessful, validationErrors := validation.ValidateDeviceCreate(&deviceCreate)
	if !validationSuccessful {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
		return
	}

	id, err := h.deviceService.CreateDevice(&deviceCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

// updateDevice handles PUT /api/devices/:id
func (h *Handler) updateDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	var deviceUpdate models.DeviceUpdate
	if err := c.ShouldBindJSON(&deviceUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = h.deviceService.UpdateDevice(id, &deviceUpdate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// deleteDevice handles DELETE /api/devices/:id
func (h *Handler) deleteDevice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	err = h.deviceService.DeleteDevice(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// healthCheck handles GET /health
func (h *Handler) healthCheck(c *gin.Context) {
	// Dummy request to check db status
	_, err := h.deviceService.GetAllDevices()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":  "error",
			"message": "Database connection failed",
			"error":   err.Error(),
		})
		return
	}

	// Calculate uptime
	uptime := time.Since(h.startTime).String()

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"message":  "Service is healthy",
		"uptime":   uptime,
		"database": "connected",
	})
}

// triggerDeviceAlarm handles POST /api/devices/:id/alarm
func (h *Handler) triggerDeviceAlarm(c *gin.Context) {
	// Parse device ID from URL
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid device ID"})
		return
	}

	// Parse request body
	var alarmRequest models.AlarmRequest
	if err := c.ShouldBindJSON(&alarmRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate alarm request
	validationSuccessful, validationErrors := validation.ValidateAlarmRequest(&alarmRequest)
	if !validationSuccessful {
		c.JSON(http.StatusBadRequest, gin.H{"errors": validationErrors})
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
