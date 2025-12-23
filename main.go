package main

import (
	"fmt"
	"log"
	"miniapp-backend/internal/config"
	"miniapp-backend/internal/handler"
	"miniapp-backend/internal/model"
	"miniapp-backend/internal/repository"
	"miniapp-backend/internal/router"
	"miniapp-backend/pkg/cache"
	"miniapp-backend/pkg/database"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2. Initialize Database
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// AutoMigrate
	if err := db.AutoMigrate(
		&model.User{},
		&model.IntakeRecord{},
		&model.Achievement{},
		&model.UserAchievement{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 3. Initialize Redis
	rdb, err := cache.NewRedisClient(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	// 4. Initialize Repository & Handler
	userRepo := repository.NewUserRepository(db)
	intakeRepo := repository.NewIntakeRepository(db)
	achievementRepo := repository.NewAchievementRepository(db)

	userHandler := handler.NewUserHandler(userRepo)
	intakeHandler := handler.NewIntakeHandler(intakeRepo, userRepo)
	achievementHandler := handler.NewAchievementHandler(achievementRepo)

	// 5. Setup Router
	r := router.SetupRouter(cfg, db, rdb, userHandler, intakeHandler, achievementHandler)

	// 6. Run Server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
