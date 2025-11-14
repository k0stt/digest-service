package handler

import (
	"digest-service/internal/auth"
	"digest-service/internal/models"
	"digest-service/internal/repository"
	"encoding/json"
	"net/http"
)

type AuthHandler struct {
	repo *repository.PostgresRepository
}

func NewAuthHandler(repo *repository.PostgresRepository) *AuthHandler {
	return &AuthHandler{repo: repo}
}

// RegisterRequest - структура для регистрации
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// LoginRequest - структура для логина
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse - структура ответа с токеном
type AuthResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Проверяем что пользователь не существует
	existingUser, err := h.repo.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if existingUser != nil {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	// Хешируем пароль
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Создаем пользователя
	user := &models.User{
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Email:        req.Email,
	}

	if err := h.repo.CreateUser(user); err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Генерируем JWT токен
	token, err := auth.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Ищем пользователя
	user, err := h.repo.GetUserByUsername(req.Username)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}
	if user == nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Проверяем пароль
	if !auth.CheckPasswordHash(req.Password, user.PasswordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Генерируем JWT токен
	token, err := auth.GenerateJWT(user)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Возвращаем токен
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Token: token})
}
