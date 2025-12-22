package handler

import (
	"miniapp-backend/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, count, err := h.repo.FindAllBasicInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"users": users,
	})
}

type UpdateSettingsRequest struct {
	DailyGoal         *int64  `json:"daily_goal"`
	ReminderEnabled   *bool   `json:"reminder_enabled"`
	ReminderInterval  *int64  `json:"reminder_interval"`
	ReminderStartTime *string `json:"reminder_start_time"`
	ReminderEndTime   *string `json:"reminder_end_time"`
}

func (h *UserHandler) UpdateSettings(c *gin.Context) {
	userIDStr := c.Query("user_id")
	var userID int64 = 1
	if userIDStr != "" {
		var err error
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
	}

	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.DailyGoal != nil {
		updates["daily_goal"] = *req.DailyGoal
	}
	if req.ReminderEnabled != nil {
		updates["reminder_enabled"] = *req.ReminderEnabled
	}
	if req.ReminderInterval != nil {
		updates["reminder_interval"] = *req.ReminderInterval
	}
	if req.ReminderStartTime != nil {
		updates["reminder_start_time"] = *req.ReminderStartTime
	}
	if req.ReminderEndTime != nil {
		updates["reminder_end_time"] = *req.ReminderEndTime
	}

	if err := h.repo.UpdateSettings(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}

func (h *UserHandler) GetSettings(c *gin.Context) {
	userIDStr := c.Query("user_id")
	var userID int64 = 1
	if userIDStr != "" {
		var err error
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
	}

	user, err := h.repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"daily_goal":          user.DailyGoal,
		"reminder_enabled":    user.ReminderEnabled,
		"reminder_interval":   user.ReminderInterval,
		"reminder_start_time": user.ReminderStartTime,
		"reminder_end_time":   user.ReminderEndTime,
		"quick_add_presets":   user.QuickAddPresets,
	})
}
