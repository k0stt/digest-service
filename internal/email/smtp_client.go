package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type SMTPClient struct{}

func NewSMTPClient() *SMTPClient {
	return &SMTPClient{}
}

// SendDigest - отправляет дайджест на email
func (sc *SMTPClient) SendDigest(toEmail, fromEmail, appPassword string, digestContent []string) error {
	// Для Gmail SMTP
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Аутентификация
	auth := smtp.PlainAuth("", fromEmail, appPassword, smtpHost)

	// Подготавливаем письмо с правильными переносами строк
	subject := "Ваш ежедневный дайджест задач"
	body := sc.buildHTMLDigest(digestContent)

	// Правильно формируем MIME сообщение
	msg := []byte(
		"To: " + toEmail + "\r\n" +
			"From: Digest Service <" + fromEmail + ">\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
			"\r\n" +
			body + "\r\n")

	// Отправляем письмо
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, fromEmail, []string{toEmail}, msg)
	if err != nil {
		// Пробуем с TLS
		return sc.sendWithTLS(smtpHost, smtpPort, fromEmail, appPassword, toEmail, msg)
	}

	log.Printf("Digest sent successfully to %s", toEmail)
	return nil
}

// sendWithTLS - отправка с TLS
func (sc *SMTPClient) sendWithTLS(host, port, from, password, to string, msg []byte) error {
	tlsconfig := &tls.Config{
		ServerName: host,
	}

	conn, err := tls.Dial("tcp", host+":"+port, tlsconfig)
	if err != nil {
		return fmt.Errorf("failed to dial: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	auth := smtp.PlainAuth("", from, password, host)
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to auth: %v", err)
	}

	if err = client.Mail(from); err != nil {
		return fmt.Errorf("failed to set from: %v", err)
	}

	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("failed to set rcpt: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to get data: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("failed to write: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("failed to close: %v", err)
	}

	return client.Quit()
}

// buildHTMLDigest - формирует красивый HTML дайджест
func (sc *SMTPClient) buildHTMLDigest(messages []string) string {
	var sb strings.Builder

	sb.WriteString(`<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ежедневный дайджест</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            line-height: 1.6; 
            color: #333; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            margin: 0;
            padding: 20px;
            min-height: 100vh;
        }
        .container {
            max-width: 800px;
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 40px rgba(0,0,0,0.1);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 40px 30px;
            text-align: center;
        }
        .header h1 {
            font-size: 2.5em;
            margin-bottom: 10px;
            font-weight: 300;
        }
        .header p {
            font-size: 1.2em;
            opacity: 0.9;
        }
        .content {
            padding: 40px 30px;
        }
        .stats {
            background: #f8f9fa;
            padding: 20px;
            border-radius: 15px;
            margin-bottom: 30px;
            text-align: center;
            border-left: 5px solid #667eea;
        }
        .stats h3 {
            color: #667eea;
            margin-bottom: 10px;
        }
        .email-card {
            background: white;
            border: 1px solid #e9ecef;
            border-radius: 15px;
            padding: 25px;
            margin-bottom: 20px;
            transition: all 0.3s ease;
            box-shadow: 0 5px 15px rgba(0,0,0,0.08);
        }
        .email-card:hover {
            transform: translateY(-5px);
            box-shadow: 0 15px 30px rgba(0,0,0,0.15);
            border-left: 5px solid #667eea;
        }
        .email-number {
            display: inline-block;
            background: #667eea;
            color: white;
            width: 30px;
            height: 30px;
            border-radius: 50%;
            text-align: center;
            line-height: 30px;
            font-weight: bold;
            margin-right: 15px;
        }
        .email-meta {
            color: #6c757d;
            font-size: 0.9em;
            margin-bottom: 10px;
            display: flex;
            justify-content: space-between;
            flex-wrap: wrap;
        }
        .email-from {
            font-weight: 600;
            color: #495057;
        }
        .email-date {
            color: #868e96;
        }
        .email-subject {
            font-size: 1.3em;
            font-weight: 600;
            color: #212529;
            margin: 10px 0;
            line-height: 1.4;
        }
        .no-emails {
            text-align: center;
            padding: 60px 30px;
            background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
            color: white;
            border-radius: 15px;
        }
        .no-emails h2 {
            font-size: 2em;
            margin-bottom: 15px;
        }
        .footer {
            background: #f8f9fa;
            padding: 30px;
            text-align: center;
            border-top: 1px solid #e9ecef;
        }
        .footer p {
            color: #6c757d;
            margin-bottom: 5px;
        }
        .badge {
            display: inline-block;
            padding: 5px 12px;
            background: #667eea;
            color: white;
            border-radius: 20px;
            font-size: 0.8em;
            margin-left: 10px;
        }
        @media (max-width: 600px) {
            .container {
                margin: 10px;
                border-radius: 15px;
            }
            .header h1 {
                font-size: 2em;
            }
            .email-meta {
                flex-direction: column;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>📊 Ваш ежедневный дайджест</h1>
            <p>Все важные письма за последние 24 часа</p>
        </div>
        
        <div class="content">`)

	if len(messages) == 0 {
		sb.WriteString(`
            <div class="no-emails">
                <h2>🎉 Отличный день!</h2>
                <p>Новых писем не обнаружено. Можете сосредоточиться на текущих задачах!</p>
            </div>`)
	} else {
		sb.WriteString(`
            <div class="stats">
                <h3>📨 Обнаружено писем: <span class="badge">` + strconv.Itoa(len(messages)) + `</span></h3>
                <p>За период: ` + time.Now().AddDate(0, 0, -1).Format("02.01.2006") + ` - ` + time.Now().Format("02.01.2006") + `</p>
            </div>
            <div class="emails-list">`)

		for i, msg := range messages {
			// Парсим информацию из сообщения
			lines := strings.Split(msg, "\n")
			var from, subject, date string

			for _, line := range lines {
				if strings.HasPrefix(line, "📧 From: ") {
					from = strings.TrimPrefix(line, "📧 From: ")
				} else if strings.HasPrefix(line, "   Subject: ") {
					subject = strings.TrimPrefix(line, "   Subject: ")
				} else if strings.HasPrefix(line, "   Date: ") {
					date = strings.TrimPrefix(line, "   Date: ")
				}
			}

			// Обрезаем длинные темы
			if len(subject) > 100 {
				subject = subject[:100] + "..."
			}

			sb.WriteString(`
                <div class="email-card">
                    <div class="email-meta">
                        <span class="email-number">` + strconv.Itoa(i+1) + `</span>
                        <span class="email-from">` + from + `</span>
                        <span class="email-date">` + date + `</span>
                    </div>
                    <div class="email-subject">` + subject + `</div>
                </div>`)
		}
		sb.WriteString(`</div>`)
	}

	sb.WriteString(`
        </div>
        
        <div class="footer">
            <p><strong>Digest Service</strong> - ваш помощник в управлении email</p>
            <p><small>Это письмо сгенерировано автоматически • ` + time.Now().Format("02.01.2006 15:04") + `</small></p>
            <p><small><a href="#" style="color: #667eea; text-decoration: none;">Настроить рассылку</a></small></p>
        </div>
    </div>
</body>
</html>`)

	return sb.String()
}
