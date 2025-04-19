package activity

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) UploadActivityPhoto(c *gin.Context) {
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Photo file required"})
		return
	}

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	savePath := "./uploads/" + filename

	if err := c.SaveUploadedFile(file, savePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	fileURL := fmt.Sprintf("http://localhost:8080/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}

func (h *Handler) PostActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userIDInterface, _ := c.Get("id")
	userID := userIDInterface.(uint)

	startTime, _ := time.Parse(time.RFC3339, req.StartTime)
	endTime, _ := time.Parse(time.RFC3339, req.EndTime)
	date, _ := time.Parse("2006-01-02", req.Date)

	activity := Activity{
		UserID:       userID,
		Title:        req.Title,
		Caption:      req.Caption,
		Distance:     req.Distance,
		Pace:         req.Pace,
		StepsCount:   req.StepsCount,
		StartTime:    startTime,
		EndTime:      endTime,
		Location:     req.Location,
		ActivityType: req.ActivityType,
		Date:         date,
	}

	// Save activity
	if err := h.service.CreateActivity(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	// Save pictures
	for _, url := range req.PictureURLs {
		pic := Picture {
			ActivityID: activity.ID,
			URL: url,
		}
		if err := h.service.CreatePicture(&pic); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save picture"})
			return
		}
	}

	c.JSON(http.StatusCreated, activity)
}

func (h *Handler) GetActivityByID(c *gin.Context) {
	idParam := c.Param("id")

	activityID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	activity, err := h.service.GetByID(uint(activityID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "activity not found"})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
	}

	c.JSON(http.StatusOK, activity)
}