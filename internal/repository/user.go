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

	user := models.NewUser()

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
