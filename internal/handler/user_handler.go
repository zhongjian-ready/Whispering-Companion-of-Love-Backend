package handler

import (
	"miniapp-backend/internal/model"
	"miniapp-backend/internal/repository"
	"miniapp-backend/pkg/wechat"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo      *repository.UserRepository
	wechatSvc *wechat.WeChatService
}

func NewUserHandler(repo *repository.UserRepository, wechatSvc *wechat.WeChatService) *UserHandler {
	return &UserHandler{repo: repo, wechatSvc: wechatSvc}
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
	UserID            *int64  `json:"user_id"`
	DailyGoal         *int64  `json:"daily_goal"`
	ReminderEnabled   *bool   `json:"reminder_enabled"`
	ReminderInterval  *int64  `json:"reminder_interval"`
	ReminderStartTime *string `json:"reminder_start_time"`
	ReminderEndTime   *string `json:"reminder_end_time"`
}

func (h *UserHandler) UpdateSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID int64 = 1
	userIDStr := c.Query("user_id")
	if userIDStr != "" {
		var err error
		userID, err = strconv.ParseInt(userIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
			return
		}
	} else if req.UserID != nil {
		userID = *req.UserID
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

type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *UserHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. Call WeChat API to get OpenID
	wxResp, err := h.wechatSvc.Code2Session(req.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to login with WeChat: " + err.Error()})
		return
	}

	// 2. Find or Create User
	user, err := h.repo.FindByOpenID(wxResp.OpenID)
	if err != nil {
		// If not found, create new user
		if err.Error() == "record not found" {
			newUser := &model.User{
				OpenID:   &wxResp.OpenID,
				UnionID:  &wxResp.UnionID,
				Username: &wxResp.OpenID, // Use OpenID as default username
			}
			if err := h.repo.Create(newUser); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
			user = newUser
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// 3. Return User Info
	c.JSON(http.StatusOK, gin.H{
		"user_id": user.ID,
		"openid":  user.OpenID,
	})
}

type UpdateInfoRequest struct {
	UserID    int64  `json:"user_id" binding:"required"`
	Nickname  string `json:"nickname"`
	AvatarURL string `json:"avatar_url"`
	PhoneCode string `json:"phone_code"`
}

func (h *UserHandler) UpdateInfo(c *gin.Context) {
	var req UpdateInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := make(map[string]interface{})
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.AvatarURL != "" {
		updates["avatar_url"] = req.AvatarURL
	}

	if req.PhoneCode != "" {
		phoneInfo, err := h.wechatSvc.GetPhoneNumber(req.PhoneCode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get phone number: " + err.Error()})
			return
		}
		updates["phone"] = phoneInfo.PurePhoneNumber
	}

	if err := h.repo.UpdateSettings(req.UserID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user info"})
		return
	}

	response := gin.H{"message": "User info updated"}
	if val, ok := updates["phone"]; ok {
		response["phone"] = val
	}
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) GetInfo(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id"})
		return
	}

	user, err := h.repo.FindByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":    user.ID,
		"nickname":   user.Nickname,
		"avatar_url": user.AvatarURL,
		"openid":     user.OpenID,
		"phone":      user.Phone,
	})
}
