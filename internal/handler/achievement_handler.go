package handler

import (
	"miniapp-backend/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AchievementHandler struct {
	repo *repository.AchievementRepository
}

func NewAchievementHandler(repo *repository.AchievementRepository) *AchievementHandler {
	return &AchievementHandler{repo: repo}
}

func (h *AchievementHandler) GetAchievements(c *gin.Context) {
	userIDStr := c.Query("user_id")
	userID, _ := strconv.ParseInt(userIDStr, 10, 64)

	all, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch achievements"})
		return
	}

	unlocked, err := h.repo.FindUserAchievements(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user achievements"})
		return
	}

	// Map unlocked IDs for easy lookup
	unlockedMap := make(map[int64]bool)
	for _, ua := range unlocked {
		unlockedMap[ua.AchievementID] = true
	}

	// Construct response with status
	var response []gin.H
	for _, a := range all {
		response = append(response, gin.H{
			"id":             a.ID,
			"name":           a.Name,
			"description":    a.Description,
			"icon_url":       a.IconURL,
			"condition_type": a.ConditionType,
			"condition_val":  a.ConditionVal,
			"is_unlocked":    unlockedMap[a.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{"achievements": response})
}
