package digest

import (
	"digest-service/internal/email"
	"digest-service/internal/models"
	"digest-service/internal/repository"
	"log"
)

type DigestService struct {
	repo       *repository.PostgresRepository
	imapClient *email.IMAPClient
	smtpClient *email.SMTPClient
}

func NewDigestService(repo *repository.PostgresRepository) *DigestService {
	return &DigestService{
		repo:       repo,
		imapClient: email.NewIMAPClient(),
		smtpClient: email.NewSMTPClient(),
	}
}

// GenerateAndSendDigest - генерирует и отправляет дайджест для пользователя
func (ds *DigestService) GenerateAndSendDigest(userID int) error {
	// Получаем настройки пользователя
	settings, err := ds.repo.GetSettings(userID)
	if err != nil {
		return err
	}

	// Проверяем что настройки заполнены
	if settings.Email == "" || settings.AppPassword == "" {
		log.Printf("User %d: email settings not configured", userID)
		return nil
	}

	log.Printf("Generating digest for user %d (%s)", userID, settings.Email)

	// Получаем письма за последний день
	messages, err := ds.imapClient.FetchRecentEmails(settings, 1)
	if err != nil {
		return err
	}

	// Отправляем дайджест
	if err := ds.smtpClient.SendDigest(settings.Email, settings.Email, settings.AppPassword, messages); err != nil {
		return err
	}

	log.Printf("Digest successfully sent to %s", settings.Email)
	return nil
}

// TestEmailConnection - тестирует подключение к почте
func (ds *DigestService) TestEmailConnection(settings *models.DigestSettings) error {
	return ds.imapClient.TestConnection(settings)
}
