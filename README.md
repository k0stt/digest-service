# Digest Service

Умный сервис для ежедневных дайджестов электронной почты.

## Возможности

- Автоматический сбор писем за последние 24 часа
- Планирование отправки дайджестов по расписанию  
- Безопасная аутентификация с JWT
- Docker-контейнеризация
- Адаптивный веб-интерфейс

## Технологии

- **Backend**: Go, Gorilla Mux, JWT, PostgreSQL
- **Frontend**: Vanilla JavaScript, HTML5, CSS3
- **Database**: PostgreSQL
- **Containerization**: Docker, Docker Compose
- **Email**: IMAP/SMTP интеграция

## Быстрый старт

### Требования
- Docker
- Docker Compose

### Запуск

```bash
# Клонируйте репозиторий
git clone https://github.com/k0stt/digest-service.git
cd digest-service

# Запустите приложение
docker-compose up --build

# Откройте в браузере
http://localhost:8080
