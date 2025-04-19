package activity

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service}
}

func (h *Handler) UploadPhoto(c *gin.Context) {
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

func (handler *Handler) PostActivity(c *gin.Context) {
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
	if err := handler.service.CreateActivity(&activity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create activity"})
		return
	}

	// Save pictures
	for _, url := range req.PictureURLs {
		pic := Picture {
			ActivityID: activity.ID,
			URL: url,
		}
		if err := handler.service.CreatePicture(&pic); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save picture"})
			return
		}
	}

	c.JSON(http.StatusCreated, activity)
}