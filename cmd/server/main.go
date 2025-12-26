package main

import (
	"fmt"
	"log"
	"miniapp-backend/internal/config"
	"miniapp-backend/internal/handler"
	"miniapp-backend/internal/model"
	"miniapp-backend/internal/repository"
	"miniapp-backend/internal/router"
	"miniapp-backend/internal/service"
	"miniapp-backend/pkg/database"
	"miniapp-backend/pkg/wechat"
)

func main() {
	log.Println("Starting server v1.1 with StatusPhoto support...")

	// 1. Load Configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}


	// 2. Initialize Database
	log.Printf("Connecting to database at %s:%s...", cfg.Database.Host, cfg.Database.Port)
	db, err := database.NewMySQLDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// AutoMigrate
	if err := db.AutoMigrate(
		&model.User{},
		&model.IntakeRecord{},
		&model.Achievement{},
		&model.UserAchievement{},
		&model.Order{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// 3. Initialize Services
	wechatSvc := wechat.NewWeChatService(cfg.WeChat)

	// 4. Initialize Repository & Handler
	userRepo := repository.NewUserRepository(db)
	intakeRepo := repository.NewIntakeRepository(db)
	achievementRepo := repository.NewAchievementRepository(db)
	orderRepo := repository.NewOrderRepository(db)

	paymentSvc := service.NewPaymentService(cfg.WeChat, orderRepo, userRepo)

	userHandler := handler.NewUserHandler(userRepo, wechatSvc)
	intakeHandler := handler.NewIntakeHandler(intakeRepo, userRepo)
	achievementHandler := handler.NewAchievementHandler(achievementRepo)
	paymentHandler := handler.NewPaymentHandler(paymentSvc)

	// Initialize Reminder Service
	reminderSvc := service.NewReminderService(userRepo, intakeRepo, wechatSvc, cfg)
	reminderSvc.Start()

	// 5. Setup Router
	r := router.SetupRouter(cfg, db, userHandler, intakeHandler, achievementHandler, paymentHandler)

	// 6. Run Server
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
