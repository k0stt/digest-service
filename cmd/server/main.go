package main

import (
	"digest-service/internal/cron"
	"digest-service/internal/handler"
	"digest-service/internal/middleware"
	"digest-service/internal/repository"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func main() {
	// Подключаемся к базе данных
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	repo, err := repository.NewPostgresRepository(dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer repo.Close()

	log.Println("Successfully connected to database!")

	// Запускаем планировщик
	scheduler := cron.NewScheduler(repo)
	defer scheduler.Stop()
	scheduler.Start()

	// Создаем хендлеры
	authHandler := handler.NewAuthHandler(repo)
	settingsHandler := handler.NewSettingsHandler(repo)
	emailHandler := handler.NewEmailHandler(repo)

	// Настраиваем роутер
	r := mux.NewRouter()

	// Публичные эндпоинты (без аутентификации)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/api/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/api/login", authHandler.Login).Methods("POST")

	// Защищенные эндпоинты (требуют аутентификации)
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	protected.HandleFunc("/test-email", emailHandler.TestConnection).Methods("POST")
	protected.HandleFunc("/send-test-digest", emailHandler.SendTestDigest).Methods("POST")

	protected.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("This is a protected endpoint!"))
	}).Methods("GET")

	// Эндпоинты для настроек
	protected.HandleFunc("/settings", settingsHandler.GetSettings).Methods("GET")
	protected.HandleFunc("/settings", settingsHandler.SaveSettings).Methods("POST")

	// Тест email подключения
	protected.HandleFunc("/test-email", emailHandler.TestConnection).Methods("POST")

	// Обслуживаем статические файлы фронтенда
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
