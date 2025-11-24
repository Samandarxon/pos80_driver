package handlers

import (
	"fmt"
	"log"
	"net/http"
	"pos80/internal/audio"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
)

// 🎯 REQUEST/RESPONSE STRUCTURES
type AudioRequest struct {
	TicketID       string `json:"ticket_id" binding:"required"`
	RoomNumber     string `json:"room_number" binding:"required"`
	DepartmentName string `json:"department_name"`
	QueueNumber    string `json:"queue_number" binding:"required"`
	DoctorID       string `json:"doctor_id"`
}

type AudioResponse struct {
	Status    string                 `json:"status"`
	Message   string                 `json:"message,omitempty"`
	Error     string                 `json:"error,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Timestamp string                 `json:"timestamp"`
}

// 🎯 PERFORMANCE METRICS
var (
	requestCounter uint64
	successCounter uint64
	failureCounter uint64
	activeRequests int32
	startupTime    = time.Now()
)

// 🎯 AUDIO HANDLER WITH QUEUE SUPPORT
type AudioHandler struct {
	audioQueue *audio.AudioQueueService
}

// ⚠️ YANGI METOD: Allaqachon yaratilgan queue ni qabul qiladi
func NewAudioHandlerWithQueue(audioQueue *audio.AudioQueueService) *AudioHandler {
	handler := &AudioHandler{
		audioQueue: audioQueue,
	}

	log.Println("🚀 Audio Handler with Queue System ready!")
	return handler
}

// ⚡ QUEUE-BASED AUDIO ANNOUNCEMENT
func (h *AudioHandler) HandleAudioAnnouncement(c *gin.Context) {
	requestID := atomic.AddUint64(&requestCounter, 1)
	atomic.AddInt32(&activeRequests, 1)
	defer atomic.AddInt32(&activeRequests, -1)

	startTime := time.Now()

	var req AudioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		atomic.AddUint64(&failureCounter, 1)

		c.JSON(http.StatusBadRequest, AudioResponse{
			Status:    "error",
			Error:     "INVALID_REQUEST",
			Message:   "Noto'g'ri JSON: " + err.Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// ✅ Request validation
	if err := h.validateRequest(&req); err != nil {
		atomic.AddUint64(&failureCounter, 1)

		c.JSON(http.StatusBadRequest, AudioResponse{
			Status:    "error",
			Error:     "VALIDATION_ERROR",
			Message:   err.Error(),
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
		return
	}

	// 🚀 QUEUE GA QO'SHISH
	h.audioQueue.AddTask(req.QueueNumber, req.RoomNumber)

	// 📊 Queue status
	queueStatus := h.audioQueue.GetStatus()

	// ✅ Tezkor response
	responseTime := time.Since(startTime)
	c.JSON(http.StatusOK, AudioResponse{
		Status:  "success",
		Message: "Audio navbatga qo'shildi",
		Data: map[string]interface{}{
			"request_id":       requestID,
			"ticket_id":        req.TicketID,
			"queue_number":     req.QueueNumber,
			"room_number":      req.RoomNumber,
			"department_name":  req.DepartmentName,
			"doctor_id":        req.DoctorID,
			"response_time_ms": responseTime.Milliseconds(),
			"queue_position":   queueStatus["queue_length"],
			"active_workers":   queueStatus["worker_count"],
			"total_requests":   atomic.LoadUint64(&requestCounter),
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	})

	log.Printf("✅ [Req-%d] Navbatga qo'shildi: %s (navbat: %d)",
		requestID, req.QueueNumber, queueStatus["queue_length"])
}

// 🔧 REQUEST VALIDATION
func (h *AudioHandler) validateRequest(req *AudioRequest) error {
	if req.QueueNumber == "" {
		return fmt.Errorf("queue_number talab qilinadi")
	}
	if req.RoomNumber == "" {
		return fmt.Errorf("room_number talab qilinadi")
	}
	if req.TicketID == "" {
		return fmt.Errorf("ticket_id talab qilinadi")
	}
	return nil
}

// ⚡ HANDLE PLAY AUDIO - Deprecated
func (h *AudioHandler) HandlePlayAudio(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{
		"error": "Use /announcement endpoint instead",
	})
}

// 📊 QUEUE STATUS ENDPOINT
func (h *AudioHandler) HandleQueueStatus(c *gin.Context) {
	queueStatus := h.audioQueue.GetStatus()

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"queue": queueStatus,
			"metrics": gin.H{
				"total_requests":   atomic.LoadUint64(&requestCounter),
				"successful_plays": atomic.LoadUint64(&successCounter),
				"failed_plays":     atomic.LoadUint64(&failureCounter),
				"active_requests":  atomic.LoadInt32(&activeRequests),
			},
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// 🗑️ QUEUE NI TOZALASH
func (h *AudioHandler) HandleClearQueue(c *gin.Context) {
	clearedCount := h.audioQueue.ClearQueue()

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": fmt.Sprintf("Navbat tozalandi, %d task o'chirildi", clearedCount),
		"data": gin.H{
			"cleared_count": clearedCount,
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}

// 🏥 HEALTH CHECK
func (h *AudioHandler) HandleHealth(c *gin.Context) {
	queueStatus := h.audioQueue.GetStatus()

	healthStatus := "healthy"
	if !queueStatus["is_running"].(bool) {
		healthStatus = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    healthStatus,
		"service":   "audio-handler",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"uptime":    time.Since(startupTime).String(),
		"queue":     queueStatus,
	})
}
