package repository

import (
	"database/sql"
	"errors"
	"simple-forum/internal/model"
)

type UserRepository struct {
	conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{conn: conn}
}

func (u *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	user := new(model.User)

	err := u.conn.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}

func (u *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, email, password_hash, created_at, role FROM users WHERE username = $1`

	user := new(model.User)

	err := u.conn.QueryRow(query, username).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (u *UserRepository) InsertUser(user *model.User) (int, error) {
	query := `INSERT INTO users (username, email, password_hash, created_at, role) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	err := u.conn.QueryRow(query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.Role).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u *UserRepository) GetUserByID(id int) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	user := new(model.User)

	err := u.conn.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.Role,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}
