package handler

import (
	"digest-service/internal/models"
	"digest-service/internal/repository"
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

type SettingsHandler struct {
	repo *repository.PostgresRepository
}

func NewSettingsHandler(repo *repository.PostgresRepository) *SettingsHandler {
	return &SettingsHandler{repo: repo}
}

// GetUserIDFromToken - извлекает ID пользователя из JWT токена
func GetUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	tokenString := authHeader[7:] // Убираем "Bearer "

	// Простой парсинг JWT (в реальном приложении нужно валидировать!)
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if userID, ok := claims["user_id"].(float64); ok {
			return int(userID), nil
		}
	}
	return 0, http.ErrMissingFile
}

func (h *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	settings, err := h.repo.GetSettings(userID)
	if err != nil {
		http.Error(w, "Error getting settings", http.StatusInternalServerError)
		return
	}

	// Не возвращаем пароль приложения в ответе
	settings.AppPassword = ""

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(settings)
}

func (h *SettingsHandler) SaveSettings(w http.ResponseWriter, r *http.Request) {
	userID, err := GetUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	var settings models.DigestSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	settings.UserID = userID

	if err := h.repo.SaveSettings(&settings); err != nil {
		http.Error(w, "Error saving settings", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Settings saved successfully"))
}
