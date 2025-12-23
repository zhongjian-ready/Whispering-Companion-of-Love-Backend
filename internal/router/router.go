package router

import (
	"fmt"
	"miniapp-backend/internal/config"
	"miniapp-backend/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func SetupRouter(cfg *config.Config, db *gorm.DB, rdb *redis.Client, userHandler *handler.UserHandler, intakeHandler *handler.IntakeHandler, achievementHandler *handler.AchievementHandler) *gin.Engine {
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	// Middleware to log WeChat Call ID
	r.Use(func(c *gin.Context) {
		callID := c.GetHeader("x-wx-call-id")
		if callID == "" {
			callID = c.GetHeader("x-request-id")
		}
		if callID != "" {
			fmt.Printf("[WeChat] CallID: %s | Path: %s | Method: %s\n", callID, c.Request.URL.Path, c.Request.Method)
			// Set it to context if handlers need it
			c.Set("callID", callID)
		}
		c.Next()
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// User routes
	r.GET("/users", userHandler.GetUsers)
	r.GET("/settings", userHandler.GetSettings)
	r.PUT("/settings", userHandler.UpdateSettings)

	// Intake routes
	r.POST("/intake", intakeHandler.AddIntake)
	r.GET("/intake/today", intakeHandler.GetToday)
	r.DELETE("/intake/:id", intakeHandler.DeleteIntake)
	r.GET("/intake/stats/weekly", intakeHandler.GetWeeklyStats)

	// Achievement routes
	r.GET("/achievements", achievementHandler.GetAchievements)

	return r
}
