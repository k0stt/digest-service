package handler

import (
	"digest-service/internal/digest"
	"digest-service/internal/repository"
	"encoding/json"
	"net/http"
)

type EmailHandler struct {
	digestService *digest.DigestService
	repo          *repository.PostgresRepository
}

func NewEmailHandler(repo *repository.PostgresRepository) *EmailHandler {
	return &EmailHandler{
		digestService: digest.NewDigestService(repo),
		repo:          repo,
	}
}

func (h *EmailHandler) TestConnection(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Получаем настройки пользователя через репозиторий хендлера
	settings, err := h.repo.GetSettings(userID)
	if err != nil {
		http.Error(w, "Error getting settings", http.StatusInternalServerError)
		return
	}

	// Проверяем что настройки заполнены
	if settings.Email == "" || settings.AppPassword == "" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": "Email settings not configured",
		})
		return
	}

	// Тестируем подключение
	if err := h.digestService.TestEmailConnection(settings); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "✅ Email connection successful!",
	})
}

// SendTestDigest - отправляет тестовый дайджест
func (h *EmailHandler) SendTestDigest(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	if err := h.digestService.GenerateAndSendDigest(userID); err != nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "✅ Test digest sent successfully!",
	})
}
