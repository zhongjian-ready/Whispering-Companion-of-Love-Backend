package handler

import (
	"miniapp-backend/internal/model"
	"miniapp-backend/internal/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type IntakeHandler struct {
	repo     *repository.IntakeRepository
	userRepo *repository.UserRepository
}

func NewIntakeHandler(repo *repository.IntakeRepository, userRepo *repository.UserRepository) *IntakeHandler {
	return &IntakeHandler{repo: repo, userRepo: userRepo}
}

type AddIntakeRequest struct {
	Amount int `json:"amount" binding:"required,min=1"`
}

func (h *IntakeHandler) AddIntake(c *gin.Context) {
	// In a real app, get UserID from context (JWT middleware)
	// userID := c.GetInt64("userID")
	userID := int64(1) // Mock user ID for now

	var req AddIntakeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	record := &model.IntakeRecord{
		UserID:     userID,
		Amount:     req.Amount,
		RecordedAt: now,
		Date:       now.Format("2006-01-02"),
	}

	if err := h.repo.Create(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success", "data": record})
}

func (h *IntakeHandler) GetToday(c *gin.Context) {
	userID := int64(1) // Mock
	today := time.Now().Format("2006-01-02")

	records, err := h.repo.GetByDate(userID, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch records"})
		return
	}

	total, err := h.repo.GetTotalIntakeByDate(userID, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch total"})
		return
	}

	// Get user goal
	user, err := h.userRepo.FindByID(userID)
	goal := 2000 // Default
	if err == nil && user.DailyGoal > 0 {
		goal = int(user.DailyGoal)
	}

	percentage := 0
	if goal > 0 {
		percentage = (total * 100) / goal
	}

	c.JSON(http.StatusOK, gin.H{
		"records":    records,
		"total":      total,
		"goal":       goal,
		"percentage": percentage,
	})
}

func (h *IntakeHandler) DeleteIntake(c *gin.Context) {
	userID := int64(1) // Mock
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.repo.Delete(id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Deleted"})
}

func (h *IntakeHandler) GetWeeklyStats(c *gin.Context) {
	userID := int64(1) // Mock
	today := time.Now().Format("2006-01-02")

	stats, err := h.repo.GetWeeklyStats(userID, today)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stats})
}


