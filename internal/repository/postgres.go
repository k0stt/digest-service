package repository

import (
	"database/sql"
	"digest-service/internal/models"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dataSourceName string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	// Проверяем подключение
	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Создаем таблицы если их нет
	if err := createTables(db); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func createTables(db *sql.DB) error {
	// Таблица пользователей
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username VARCHAR(50) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            email VARCHAR(100),
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `)
	if err != nil {
		return err
	}

	// Таблица настроек
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS digest_settings (
            user_id INTEGER PRIMARY KEY REFERENCES users(id),
            imap_server VARCHAR(100) DEFAULT 'imap.gmail.com:993',
            email VARCHAR(100),
            app_password VARCHAR(255),
            schedule VARCHAR(5) DEFAULT '09:00'
        )
    `)

	return err
}

// GetAllUsers - получает всех пользователей
func (r *PostgresRepository) GetAllUsers() ([]*models.User, error) {
	query := `SELECT id, username, email FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetSettings - получает настройки пользователя
func (r *PostgresRepository) GetSettings(userID int) (*models.DigestSettings, error) {
	settings := &models.DigestSettings{UserID: userID}
	query := `SELECT imap_server, email, app_password, schedule FROM digest_settings WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&settings.IMAPServer, &settings.Email, &settings.AppPassword, &settings.Schedule)

	if err == sql.ErrNoRows {
		// Если настроек нет - возвращаем настройки по умолчанию
		return &models.DigestSettings{
			UserID:     userID,
			IMAPServer: "imap.gmail.com:993",
			Schedule:   "09:00",
		}, nil
	}
	return settings, err
}

// SaveSettings - сохраняет настройки пользователя
func (r *PostgresRepository) SaveSettings(settings *models.DigestSettings) error {
	query := `
        INSERT INTO digest_settings (user_id, imap_server, email, app_password, schedule) 
        VALUES ($1, $2, $3, $4, $5)
        ON CONFLICT (user_id) 
        DO UPDATE SET imap_server = $2, email = $3, app_password = $4, schedule = $5
    `
	_, err := r.db.Exec(query,
		settings.UserID,
		settings.IMAPServer,
		settings.Email,
		settings.AppPassword,
		settings.Schedule)
	return err
}

// UserRepository методы
func (r *PostgresRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := r.db.QueryRow(query, user.Username, user.PasswordHash, user.Email).Scan(&user.ID, &user.CreatedAt)
	return err
}

func (r *PostgresRepository) GetUserByUsername(username string) (*models.User, error) {
	user := &models.User{}
	query := `SELECT id, username, password_hash, email, created_at FROM users WHERE username = $1`
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil // Пользователь не найден - это не ошибка
	}
	return user, err
}

func (r *PostgresRepository) GetDB() *sql.DB {
	return r.db
}

func (r *PostgresRepository) Close() error {
	if r.db != nil {
		return r.db.Close()
	}
	return nil
}
