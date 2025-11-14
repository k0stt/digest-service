package email

import (
	"digest-service/internal/models"
	"fmt"
	"log"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

type IMAPClient struct{}

func NewIMAPClient() *IMAPClient {
	return &IMAPClient{}
}

// TestConnection - тестирует подключение к IMAP
func (ic *IMAPClient) TestConnection(settings *models.DigestSettings) error {
	c, err := client.DialTLS(settings.IMAPServer, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to IMAP: %v", err)
	}
	defer c.Logout()

	if err := c.Login(settings.Email, settings.AppPassword); err != nil {
		return fmt.Errorf("failed to login: %v", err)
	}

	return nil
}

// FetchRecentEmails - получает recent письма
func (ic *IMAPClient) FetchRecentEmails(settings *models.DigestSettings, days int) ([]string, error) {
	log.Printf("Connecting to IMAP: %s", settings.IMAPServer)

	// Подключаемся к серверу
	c, err := client.DialTLS(settings.IMAPServer, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %v", err)
	}
	defer c.Logout()

	// Логинимся
	if err := c.Login(settings.Email, settings.AppPassword); err != nil {
		return nil, fmt.Errorf("failed to login: %v", err)
	}

	// Выбираем папку INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		return nil, fmt.Errorf("failed to select INBOX: %v", err)
	}

	// Ищем письма за последние N дней
	since := time.Now().AddDate(0, 0, -days)
	criteria := imap.NewSearchCriteria()
	criteria.Since = since
	criteria.WithoutFlags = []string{"\\Seen"} // Только непрочитанные

	ids, err := c.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("failed to search emails: %v", err)
	}

	if len(ids) == 0 {
		return []string{"No new emails found for digest"}, nil
	}

	// Получаем содержимое писем
	var messages []string
	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	messagesCh := make(chan *imap.Message, 10)
	done := make(chan error, 1)

	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messagesCh)
	}()

	for msg := range messagesCh {
		if msg.Envelope != nil {
			subject := msg.Envelope.Subject
			from := ""
			if len(msg.Envelope.From) > 0 {
				from = msg.Envelope.From[0].Address()
			}
			date := msg.Envelope.Date.Format("2006-01-02 15:04")

			message := fmt.Sprintf("📧 From: %s\n   Subject: %s\n   Date: %s", from, subject, date)
			messages = append(messages, message)
		}
	}

	if err := <-done; err != nil {
		return nil, fmt.Errorf("failed to fetch emails: %v", err)
	}

	log.Printf("Found %d emails for digest", len(messages))
	return messages, nil
}
