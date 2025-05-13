package repository

import (
	"userService/internal/domain"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) domain.UserRepository {
	return &userRepo{db}
}

func (r *userRepo) Create(u *domain.User) error {
	query := `INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id`
	return r.db.QueryRow(query, u.Username, u.Password).Scan(&u.ID)
}

func (r *userRepo) GetByUsername(username string) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE username=$1`, username)
	return &u, err
}

func (r *userRepo) GetByID(id int) (*domain.User, error) {
	var u domain.User
	err := r.db.Get(&u, `SELECT * FROM users WHERE id=$1`, id)
	return &u, err
}

func (r *userRepo) Update(u *domain.User) error {
	query := `UPDATE users SET username=$1 WHERE id=$2`
	_, err := r.db.Exec(query, u.Username, u.ID)
	return err
}
