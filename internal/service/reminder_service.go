package service

import (
	"fmt"
	"log"
	"miniapp-backend/internal/config"
	"miniapp-backend/internal/repository"
	"miniapp-backend/pkg/wechat"
	"time"
)

type ReminderService struct {
	userRepo   *repository.UserRepository
	intakeRepo *repository.IntakeRepository
	wechatSvc  *wechat.WeChatService
	cfg        *config.Config
}

func NewReminderService(userRepo *repository.UserRepository, intakeRepo *repository.IntakeRepository, wechatSvc *wechat.WeChatService, cfg *config.Config) *ReminderService {
	return &ReminderService{
		userRepo:   userRepo,
		intakeRepo: intakeRepo,
		wechatSvc:  wechatSvc,
		cfg:        cfg,
	}
}

func (s *ReminderService) Start() {
	log.Println("Reminder service started. Checking every 1 minute.")

	// Start a ticker immediately
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			s.SendReminders()
		}
	}()
}

func (s *ReminderService) SendReminders() {
	log.Println("Checking for reminders...")
	users, err := s.userRepo.FindUsersWithRemindersEnabled()
	if err != nil {
		log.Printf("Failed to fetch users for reminders: %v", err)
		return
	}

	now := time.Now()
	currentTimeStr := now.Format("15:04")
	todayStr := now.Format("2006-01-02")

	for _, user := range users {
		if user.OpenID == nil || *user.OpenID == "" {
			continue
		}

		// Check time range
		if currentTimeStr < user.ReminderStartTime || currentTimeStr > user.ReminderEndTime {
			continue
		}

		// Get intake data
		totalIntake, err := s.intakeRepo.GetTotalIntakeByDate(user.ID, todayStr)
		if err != nil {
			log.Printf("Failed to get total intake for user %d: %v", user.ID, err)
			continue
		}

		records, err := s.intakeRepo.GetByDate(user.ID, todayStr)
		if err != nil {
			log.Printf("Failed to get records for user %d: %v", user.ID, err)
			continue
		}

		lastIntakeTimeStr := "暂未饮水"
		if len(records) > 0 {
			lastIntake := records[0]
			// Calculate duration since last intake
			duration := now.Sub(lastIntake.RecordedAt)
			hours := int(duration.Hours())
			minutes := int(duration.Minutes()) % 60
			if hours > 0 {
				lastIntakeTimeStr = fmt.Sprintf("%d小时前", hours)
			} else {
				lastIntakeTimeStr = fmt.Sprintf("%d分钟前", minutes)
			}
		}

		// Send message
		req := &wechat.SubscribeMessageRequest{
			ToUser:     *user.OpenID,
			TemplateID: s.cfg.WeChat.TemplateID,
			Data: map[string]wechat.MessageData{
				"thing1":            {Value: "该喝水啦，保持身体水分充足哦！"},
				"character_string2": {Value: fmt.Sprintf("%d", user.DailyGoal)},
				"character_string3": {Value: fmt.Sprintf("%d", totalIntake)},
				"phrase4":           {Value: lastIntakeTimeStr},
				"thing5":            {Value: fmt.Sprintf("%d", user.DailyGoal)},
			},
		}

		if err := s.wechatSvc.SendSubscribeMessage(req); err != nil {
			log.Printf("Failed to send reminder to user %d: %v", user.ID, err)
		} else {
			log.Printf("Sent reminder to user %d", user.ID)
		}
	}
}
