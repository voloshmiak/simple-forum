package repository

import (
	"database/sql"
	"forum-project/internal/models"
)

type UserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	user := new(models.User)

	err := u.conn.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) InsertUser(user *models.User) (int, error) {
	query := `INSERT INTO users (username, email, password_hash, created_at, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := u.conn.QueryRow(query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.Role).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u *UserRepository) GetUserByID(id int) (*models.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	user := new(models.User)

	err := u.conn.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Role,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}
