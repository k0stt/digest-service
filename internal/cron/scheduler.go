package cron

import (
	"digest-service/internal/digest"
	"digest-service/internal/models"
	"digest-service/internal/repository"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

type Scheduler struct {
	repo          *repository.PostgresRepository
	digestService *digest.DigestService
	cron          *cron.Cron
}

func NewScheduler(repo *repository.PostgresRepository) *Scheduler {
	return &Scheduler{
		repo:          repo,
		digestService: digest.NewDigestService(repo),
		cron:          cron.New(),
	}
}

func (s *Scheduler) Start() {
	// Проверяем каждую минуту
	s.cron.AddFunc("* * * * *", s.checkScheduledDigests)
	s.cron.Start()

	log.Println("📅 Cron scheduler started - checking digests every minute")
}

func (s *Scheduler) Stop() {
	if s.cron != nil {
		s.cron.Stop()
		log.Println("📅 Cron scheduler stopped")
	}
}

func (s *Scheduler) checkScheduledDigests() {
	currentTime := time.Now().Format("15:04")
	log.Printf("⏰ Checking scheduled digests at %s", currentTime)

	// Получаем всех пользователей из БД
	users, err := s.repo.GetAllUsers()
	if err != nil {
		log.Printf("❌ Error getting users: %v", err)
		return
	}

	log.Printf("👥 Found %d users to check", len(users))

	for _, user := range users {
		s.checkUserDigest(user, currentTime)
	}
}

func (s *Scheduler) checkUserDigest(user *models.User, currentTime string) {
	// Получаем настройки пользователя
	settings, err := s.repo.GetSettings(user.ID)
	if err != nil {
		log.Printf("❌ Error getting settings for user %d: %v", user.ID, err)
		return
	}

	// Проверяем что настройки заполнены
	if settings.Email == "" || settings.AppPassword == "" {
		log.Printf("⚠️ User %d: email settings not configured", user.ID)
		return
	}

	// Проверяем время отправки (точное совпадение)
	if settings.Schedule == currentTime {
		log.Printf("🚀 Time to send digest for user %d (%s) at %s",
			user.ID, settings.Email, currentTime)

		// Отправляем дайджест
		if err := s.digestService.GenerateAndSendDigest(user.ID); err != nil {
			log.Printf("❌ Error sending digest for user %d: %v", user.ID, err)
		} else {
			log.Printf("✅ Digest sent successfully to %s", settings.Email)
		}
	}
}
