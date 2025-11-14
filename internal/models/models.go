package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Пароль никогда не отдаем в JSON
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

type DigestSettings struct {
	UserID      int    `json:"user_id"`
	IMAPServer  string `json:"imap_server"`
	Email       string `json:"email"`
	AppPassword string `json:"app_password"`
	Schedule    string `json:"schedule"`
}
