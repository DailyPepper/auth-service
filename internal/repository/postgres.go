package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/DailyPepper/auth-service/internal/models"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, id int64) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	UpdateLastLogin(ctx context.Context, userID int64, loginTime time.Time) error
	Close() error
}

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(dbURL string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (first_name, surname, birthday, email, phone, password_hash, is_active, is_verified, role, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	err := r.db.QueryRowContext(ctx, query,
		user.FirstName,
		user.Surname,
		user.Birthday,
		user.Email,
		user.Phone,
		user.Password,
		user.IsActive,
		user.IsVerified,
		user.Role,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	return errors.Wrap(err, "failed to create user")
}

func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT id, first_name, surname, birthday, email, phone, password_hash, 
		       is_active, is_verified, last_login, role, created_at, updated_at
		FROM users WHERE email = $1
	`

	var user models.User
	var phone sql.NullString
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.FirstName,
		&user.Surname,
		&user.Birthday,
		&user.Email,
		&phone,
		&user.Password,
		&user.IsActive,
		&user.IsVerified,
		&lastLogin,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by email")
	}

	// Обработка nullable полей
	if phone.Valid {
		user.Phone = &phone.String
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

func (r *PostgresRepository) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
	query := `
		SELECT id, first_name, surname, birthday, email, phone, password_hash, 
		       is_active, is_verified, last_login, role, created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	var phone sql.NullString
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.FirstName,
		&user.Surname,
		&user.Birthday,
		&user.Email,
		&phone,
		&user.Password,
		&user.IsActive,
		&user.IsVerified,
		&lastLogin,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by ID")
	}

	// Обработка nullable полей
	if phone.Valid {
		user.Phone = &phone.String
	}
	if lastLogin.Valid {
		user.LastLogin = &lastLogin.Time
	}

	return &user, nil
}

func (r *PostgresRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users 
		SET first_name = $1, surname = $2, birthday = $3, email = $4, phone = $5,
		    is_active = $6, is_verified = $7, role = $8, updated_at = $9
		WHERE id = $10
	`

	_, err := r.db.ExecContext(ctx, query,
		user.FirstName,
		user.Surname,
		user.Birthday,
		user.Email,
		user.Phone,
		user.IsActive,
		user.IsVerified,
		user.Role,
		user.UpdatedAt,
		user.ID,
	)

	return errors.Wrap(err, "failed to update user")
}

func (r *PostgresRepository) UpdateLastLogin(ctx context.Context, userID int64, loginTime time.Time) error {
	query := `UPDATE users SET last_login = $1, updated_at = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, loginTime, time.Now(), userID)
	return errors.Wrap(err, "failed to update last login")
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}
