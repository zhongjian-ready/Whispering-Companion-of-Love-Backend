package router

import (
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
